package backend

import (
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/hashicorp/go-version"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// TODO Error handling
func CheckForVersionUpdate(page *rod.Page, currentVersion *version.Version) bool {
	page.MustNavigate(`https://github.com/GreenSoap/cfn-tracker/releases`).MustWaitLoad()
	el, _ := page.Element(`turbo-frame div.mr-md-0:nth-child(3) > a:nth-child(1)`)
	latestVerEl, _ := el.Text() //
	latestVersionText := strings.Split(latestVerEl, "v")[1]
	latestVersion, _ := version.NewVersion(latestVersionText)
	hasNewVersion := currentVersion.LessThan(latestVersion)
	if hasNewVersion {
		fmt.Println(`Has new version: `, latestVersionText)
	}
	return hasNewVersion
}

func SetupBrowser() *rod.Page {
	fmt.Println("Setting up browser")
	u := launcher.New().Leakless(false).Headless(true).MustLaunch()
	page := rod.New().ControlURL(u).MustConnect().MustPage("")
	router := page.HijackRequests()

	// Block all images, stylesheets, fonts and unessential scripts
	router.MustAdd("*", func(ctx *rod.Hijack) {
		if ctx.Request.Type() == proto.NetworkResourceTypeImage ||
			ctx.Request.Type() == proto.NetworkResourceTypeFont {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}

		if !strings.Contains(ctx.Request.URL().Hostname(), `steam`) &&
			ctx.Request.Type() == proto.NetworkResourceTypeStylesheet {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}
		/*


			// Only check for scripts on non-steam requests

				if !strings.Contains(ctx.Request.URL().Hostname(), `steam`) &&
					ctx.Request.Type() == proto.NetworkResourceTypeScript {
					ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
					return
				}*/

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()
	return page
}

func (a *App) Initialize(steamUsername string, steamPassword string, currentVersion *version.Version) int {
	if IsInitialized {
		return 1
	}

	page := SetupBrowser()
	loginStatus, page := Login(page, steamUsername, steamPassword)
	IsInitialized = (loginStatus == 1)
	runtime.EventsEmit(a.ctx, `initialized`, IsInitialized)

	if loginStatus == LoginError.returnCode {
		LogError(LoginError)
	} else if loginStatus == ProfileError.returnCode {
		LogError(ProfileError)
	} else if loginStatus == CaptchaError.returnCode {
		LogError(CaptchaError)
	}

	if IsInitialized == true {
		fmt.Println("Check for new version")
		hasNewVersion := CheckForVersionUpdate(page, currentVersion)
		runtime.EventsEmit(a.ctx, `version-update`, hasNewVersion)
	}

	PageInstance = page
	return loginStatus
}

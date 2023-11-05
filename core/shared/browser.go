package shared

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/hashicorp/go-version"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Browser struct {
	ctx          context.Context
	Page         *rod.Page
	HijackRouter *rod.HijackRouter
}

func NewBrowser(ctx context.Context, headless bool) (*Browser, error) {
	fmt.Println(`Setting up browser`)

	userHomeDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get cache dir for browser: %w", err)
	}
	userDataDir := filepath.Join(userHomeDir, "cfn-tracker")
	l := launcher.New()
	l.Set(flags.UserDataDir, userDataDir)
	l.RemoteDebuggingPort(6969)
	u, err := l.Leakless(false).Headless(headless).Launch()
	if err != nil {
		println("Failed to launch brosfsdfwsers", err)
		u = "127.0.0.1:6969"
	}

	// TODO: Connection to browser error handling
	browser := rod.New().ControlURL(u).MustConnect()
	var page *rod.Page
	if browser.MustPages().Empty() {
		page = browser.MustPage("")
	} else {
		page = browser.MustPages().First()
	}
	router := page.HijackRequests()

	// Block the browser from fetching unnecessary resources
	router.MustAdd(`*`, func(ctx *rod.Hijack) {
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

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	return &Browser{
		ctx:          ctx,
		Page:         page,
		HijackRouter: router,
	}, nil
}

// TODO: Error Handling
func (b *Browser) CheckForVersionUpdate(currentVersion *version.Version) {
	fmt.Println(`Check for new version`)
	b.Page.MustNavigate(`https://github.com/williamsjokvist/cfn-tracker/releases`).MustWaitLoad()
	el := b.Page.MustElement(`turbo-frame div.mr-md-0:nth-child(3) > a:nth-child(1)`)
	latestVerEl := el.MustText()
	latestVersionText := strings.Split(latestVerEl, `v`)[1]
	latestVersion, _ := version.NewVersion(latestVersionText)
	hasNewVersion := currentVersion.LessThan(latestVersion)
	if hasNewVersion {
		fmt.Println(`Has new version: `, latestVersionText)
		wails.EventsEmit(b.ctx, `version-update`, hasNewVersion)
	}
}

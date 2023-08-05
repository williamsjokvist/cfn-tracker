package shared

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"github.com/hashicorp/go-version"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
)

type Browser struct {
	ctx          context.Context
	Page         *rod.Page
	HijackRouter *rod.HijackRouter
}

func NewBrowser(ctx context.Context, headless bool) *Browser {
	fmt.Println(`Setting up browser`)
	u := launcher.New().Leakless(false).Headless(headless).MustLaunch()

	// TODO: Connection to browser error handling
	page := rod.New().ControlURL(u).MustConnect().MustPage(``)
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
	}
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

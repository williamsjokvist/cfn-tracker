package browser

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/hashicorp/go-version"
)

type Browser struct {
	Page         *rod.Page
	HijackRouter *rod.HijackRouter
}

func NewBrowser(headless bool) (*Browser, error) {
	log.Println(`Setting up browser`)

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
		return nil, fmt.Errorf("failed to launch temp browser: %w", err)
	}

	log.Println("Browser connecting to", u)
	browser := rod.New().ControlURL(u)
	err = browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to browser: %w", err)
	}

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
		Page:         page,
		HijackRouter: router,
	}, nil
}

func (b *Browser) GetLatestAppVersion() (*version.Version, error) {
	log.Println(`Check for new version`)
	err := b.Page.Navigate(`https://github.com/williamsjokvist/cfn-tracker/releases`)
	if err != nil {
		return nil, fmt.Errorf(`navigate to github: %w`, err)
	}
	err = b.Page.WaitLoad()
	if err != nil {
		return nil, fmt.Errorf(`wait for github to load: %w`, err)
	}
	versionElement, err := b.Page.Element(`turbo-frame div.mr-md-0:nth-child(3) > a:nth-child(1)`)
	if err != nil {
		return nil, fmt.Errorf(`get version element: %w`, err)
	}
	versionText, err := versionElement.Text()
	if err != nil {
		return nil, fmt.Errorf(`get version element text: %w`, err)
	}
	versionNumber := strings.Split(versionText, `v`)[1]
	latestVersion, err := version.NewVersion(versionNumber)
	if err != nil {
		return nil, fmt.Errorf(`parse version: %w`, err)
	}
	return latestVersion, nil
}

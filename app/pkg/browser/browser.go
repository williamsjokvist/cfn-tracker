package browser

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

type Browser struct {
	Page         *rod.Page
	HijackRouter *rod.HijackRouter
}

func setupHijack(page *rod.Page) *rod.HijackRouter {
	router := page.HijackRequests()
	router.MustAdd(`*`, hijackHandler)
	go router.Run()
	return router
}

func hijackHandler(ctx *rod.Hijack) {
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
}

func (b *Browser) NewTab() (*Browser, func(), error) {
	if b == nil || b.Page == nil {
		return nil, func() {}, errors.New("browser not initialized")
	}
	page, err := stealth.Page(b.Page.Browser())
	if err != nil {
		return nil, func() {}, fmt.Errorf("create stealth page: %w", err)
	}
	router := setupHijack(page)
	cleanup := func() {
		_ = router.Stop()
		_ = page.Close()
	}
	return &Browser{
		Page:         page,
		HijackRouter: router,
	}, cleanup, nil
}

func NewBrowser(headless bool) (*Browser, error) {
	slog.Debug("setting up browser")

	userHomeDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("get cache dir for browser: %w", err)
	}
	userDataDir := filepath.Join(userHomeDir, "cfn-tracker")
	l := launcher.New()
	l.Set(flags.UserDataDir, userDataDir)
	l.RemoteDebuggingPort(6969)
	u, err := l.Leakless(false).Headless(headless).Launch()
	if err != nil {
		return nil, fmt.Errorf("launch temp browser: %w", err)
	}

	slog.Debug("browser connecting to", slog.Any("url", u))
	browser := rod.New().ControlURL(u)
	err = browser.Connect()
	if err != nil {
		return nil, fmt.Errorf("connect to browser: %w", err)
	}
	page, err := stealth.Page(browser)
	if err != nil {
		return nil, fmt.Errorf("create stealth page: %w", err)
	}
	router := setupHijack(page)
	return &Browser{
		Page:         page,
		HijackRouter: router,
	}, nil
}

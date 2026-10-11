package browser

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
)

// browserCleanupTimeout bounds how long we wait for a closed browser to clean up.
// After that we move on, since a stuck authentication is worse than a leaked process.
const browserCleanupTimeout = 5 * time.Second

// ErrBrowserNotFound is returned instead of letting rod download Chromium.
var ErrBrowserNotFound = errors.New("no installed Chrome, Chromium or Edge found")

// Browser is the headless Chrome used to read buckler with the saved session.
// Logging in happens separately, in a plain Chrome window (see LaunchManualLogin).
type Browser struct {
	Page         *rod.Page
	HijackRouter *rod.HijackRouter

	mu       sync.Mutex
	launcher *launcher.Launcher
	rod      *rod.Browser
}

func NewBrowser() (*Browser, error) {
	slog.Debug("setting up browser")
	b := &Browser{}
	if err := b.launch(); err != nil {
		b.closeCurrent()
		return nil, err
	}
	return b, nil
}

func UserDataDir() (string, error) {
	userHomeDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("get cache dir for browser: %w", err)
	}
	return filepath.Join(userHomeDir, "cfn-tracker"), nil
}

func cookieKeyArgs() []string {
	args := make([]string, 0, len(cookieKeyFlags))
	for name, values := range cookieKeyFlags {
		arg := "--" + string(name)
		if len(values) > 0 {
			arg += "=" + strings.Join(values, ",")
		}
		args = append(args, arg)
	}
	return args
}

func (b *Browser) launch() error {
	userDataDir, err := UserDataDir()
	if err != nil {
		return err
	}
	path, found := launcher.LookPath()
	if !found {
		return ErrBrowserNotFound
	}
	// Without an explicit binary, rod downloads Chromium.
	l := launcher.New().Bin(path)
	l.Set(flags.UserDataDir, userDataDir)
	for name, values := range cookieKeyFlags {
		l.Set(name, values...)
	}
	l.Delete("enable-automation")
	l.Set("disable-blink-features", "AutomationControlled")
	l.RemoteDebuggingPort(6969)
	u, err := l.Leakless(false).Headless(true).Launch()
	if err != nil {
		return fmt.Errorf("launch temp browser: %w", err)
	}
	b.launcher = l

	slog.Debug("browser connecting to", slog.Any("url", u))
	b.rod = rod.New().NoDefaultDevice().ControlURL(u)
	err = b.rod.Connect()
	if err != nil {
		return fmt.Errorf("connect to browser: %w", err)
	}
	page, router, err := newPage(b.rod)
	if err != nil {
		return err
	}
	b.HijackRouter = router
	b.Page = page
	return nil
}

// NewTab gives work like the comparison its own page, so it doesn't navigate away
// from the page GetBattleLog polls with.
func (b *Browser) NewTab() (*Browser, func(), error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.rod == nil {
		return nil, func() {}, errors.New("browser not initialized")
	}
	page, router, err := newPage(b.rod)
	if err != nil {
		return nil, func() {}, err
	}
	cleanup := func() {
		_ = router.Stop()
		_ = page.Close()
	}
	return &Browser{Page: page, HijackRouter: router}, cleanup, nil
}

func newPage(browser *rod.Browser) (*rod.Page, *rod.HijackRouter, error) {
	page, err := browser.Page(proto.TargetCreateTarget{})
	if err != nil {
		return nil, nil, fmt.Errorf("open browser page: %w", err)
	}
	res, err := page.Eval(`() => navigator.userAgent`)
	if err != nil {
		_ = page.Close()
		return nil, nil, fmt.Errorf("read browser user agent: %w", err)
	}
	userAgent := res.Value.Str()
	if strings.Contains(userAgent, "HeadlessChrome") {
		userAgent = strings.ReplaceAll(userAgent, "HeadlessChrome", "Chrome")
		if err := page.SetUserAgent(&proto.NetworkSetUserAgentOverride{UserAgent: userAgent}); err != nil {
			_ = page.Close()
			return nil, nil, fmt.Errorf("set browser user agent: %w", err)
		}
	}

	router := page.HijackRequests()
	// Block the browser from fetching unnecessary resources
	router.MustAdd(`*`, func(ctx *rod.Hijack) {
		if shouldBlockRequest(ctx.Request.Type(), ctx.Request.URL().Hostname()) {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()
	return page, router, nil
}

func (b *Browser) closeCurrent() {
	if b.HijackRouter != nil {
		if err := b.HijackRouter.Stop(); err != nil {
			slog.Warn("failed to stop request hijacking", slog.Any("error", err))
		}
	}
	if b.rod != nil {
		if err := b.rod.Close(); err != nil {
			slog.Warn("failed to close browser", slog.Any("error", err))
		}
	}
	if b.launcher != nil {
		// Cleanup removes its configured UserDataDir. Clear that flag first so the
		// login profile survives — deleting it would force a new login on every
		// launch.
		b.launcher.Delete(flags.UserDataDir)

		// Cleanup blocks on <-l.exit, which is only closed once Chromium has
		// actually exited. A browser that fails to die would hang this call
		// forever while holding b.mu. Bound the wait: leaking a launcher is
		// recoverable, freezing authentication is not. A stale process only
		// means the profile lock is still held, which the retry loop in
		// Relaunch already handles.
		done := make(chan struct{})
		l := b.launcher
		go func() {
			l.Cleanup()
			close(done)
		}()
		select {
		case <-done:
		case <-time.After(browserCleanupTimeout):
			slog.Warn("browser cleanup timed out; continuing without it")
		}
	}
	b.Page = nil
	b.HijackRouter = nil
	b.rod = nil
	b.launcher = nil
}

// Relaunch replaces the current browser. It must only be called during
// authentication, before tracking starts; relaunching while polling would
// invalidate a Page being used by GetBattleLog.
func (b *Browser) Relaunch() error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.closeCurrent()
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		if attempt > 1 {
			time.Sleep(500 * time.Millisecond)
		}
		err = b.launch()
		if err == nil {
			return nil
		}
		b.closeCurrent()
	}
	return fmt.Errorf("relaunch browser: %w", err)
}

// Close stops the current browser and releases its launcher resources.
func (b *Browser) Close() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closeCurrent()
	return nil
}

func shouldBlockRequest(t proto.NetworkResourceType, hostname string) bool {
	return t == proto.NetworkResourceTypeImage ||
		t == proto.NetworkResourceTypeFont ||
		(t == proto.NetworkResourceTypeStylesheet && !strings.Contains(hostname, "steam"))
}

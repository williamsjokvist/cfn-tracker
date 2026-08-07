package browser

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

// browserCleanupTimeout は、閉じたブラウザの後始末を待つ上限。
// これを超えたら諦めて先へ進む（認証が固まるほうが害が大きい）。
const browserCleanupTimeout = 5 * time.Second

type Browser struct {
	Page         *rod.Page
	HijackRouter *rod.HijackRouter
	Headless     bool
	// PreferHeadless is the user-configured default mode. The browser returns to
	// this mode after it is temporarily made visible for image verification.
	PreferHeadless bool

	mu          sync.Mutex
	blockAssets atomic.Bool
	launcher    *launcher.Launcher
	rod         *rod.Browser
}

func NewBrowser(headless bool) (*Browser, error) {
	slog.Debug("setting up browser")
	b := &Browser{Headless: headless, PreferHeadless: headless}
	b.blockAssets.Store(true)
	if err := b.launch(headless); err != nil {
		b.closeCurrent()
		return nil, err
	}
	return b, nil
}

func (b *Browser) launch(headless bool) error {
	userHomeDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get cache dir for browser: %w", err)
	}
	userDataDir := filepath.Join(userHomeDir, "cfn-tracker")
	l := launcher.New()
	l.Set(flags.UserDataDir, userDataDir)
	l.RemoteDebuggingPort(6969)
	u, err := l.Leakless(false).Headless(headless).Launch()
	if err != nil {
		return fmt.Errorf("launch temp browser: %w", err)
	}
	b.launcher = l

	slog.Debug("browser connecting to", slog.Any("url", u))
	b.rod = rod.New().ControlURL(u)
	err = b.rod.Connect()
	if err != nil {
		return fmt.Errorf("connect to browser: %w", err)
	}
	page := stealth.MustPage(b.rod)

	router := page.HijackRequests()
	// Block the browser from fetching unnecessary resources
	router.MustAdd(`*`, func(ctx *rod.Hijack) {
		if shouldBlockRequest(ctx.Request.Type(), ctx.Request.URL().Hostname(), b.blockAssets.Load()) {
			ctx.Response.Fail(proto.NetworkErrorReasonBlockedByClient)
			return
		}

		ctx.ContinueRequest(&proto.FetchContinueRequest{})
	})

	go router.Run()

	b.Page = page
	b.HijackRouter = router
	b.Headless = headless
	return nil
}

func (b *Browser) closeCurrent() {
	if b.HijackRouter != nil {
		b.HijackRouter.Stop()
	}
	if b.rod != nil {
		if err := b.rod.Close(); err != nil {
			slog.Warn("failed to close browser", slog.Any("error", err))
		}
	}
	if b.launcher != nil {
		// Cleanup removes its configured UserDataDir. Clear that flag first so the
		// login profile survives — deleting it would force image verification on
		// every launch.
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

// Relaunch replaces the current browser in the requested mode. It must only be
// called during authentication, before tracking starts; relaunching while
// polling would invalidate a Page being used by GetBattleLog.
func (b *Browser) Relaunch(headless bool) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.closeCurrent()
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		if attempt > 1 {
			time.Sleep(500 * time.Millisecond)
		}
		err = b.launch(headless)
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

// SetAssetBlocking は画像・フォント・CSS のブロックを切り替える。
// 画像認証を人間が解く場面では false にする必要がある。
func (b *Browser) SetAssetBlocking(enabled bool) {
	b.blockAssets.Store(enabled)
}

func shouldBlockRequest(t proto.NetworkResourceType, hostname string, blocking bool) bool {
	if !blocking {
		return false
	}
	return t == proto.NetworkResourceTypeImage ||
		t == proto.NetworkResourceTypeFont ||
		(t == proto.NetworkResourceTypeStylesheet && !strings.Contains(hostname, "steam"))
}

package browser

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/launcher/flags"
	"github.com/go-rod/rod/lib/proto"
	"github.com/go-rod/stealth"
)

type Browser struct {
	Page         *rod.Page
	HijackRouter *rod.HijackRouter
	Headless     bool
	blockAssets  atomic.Bool
}

func NewBrowser(headless bool) (*Browser, error) {
	slog.Debug("setting up browser")
	b := &Browser{Headless: headless}
	b.blockAssets.Store(true)

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
	page := stealth.MustPage(browser)

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
	return b, nil
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

package browser

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/exec"
	"time"

	"github.com/go-rod/rod/lib/launcher"
)

var (
	// 起動直後に終了したかを見極める猶予。既存インスタンスへ引き渡された Chrome は
	// 1 秒以内に終了するので、これだけ生きていれば起動できたと見なしてよい。
	manualLoginStartupWindow = 3 * time.Second
	// プロファイルのロックが解放されるのを待つ間隔と、諦めるまでの試行回数。
	manualLoginRetryInterval = 2 * time.Second
	manualLoginStartAttempts = 4
)

// errManualLoginProfileLocked は、起動した Chrome が即座に終了したことを表す。
// UserDataDir は排他で、rod が残した Chromium などが掴んだままだと、新しいプロセスは
// 起動を既存インスタンスへ引き渡して自分は終了する。実機でアプリ終了後に Chromium が
// 7 プロセス残っていた例がある（2026-09-10）。
var errManualLoginProfileLocked = errors.New("manual login browser exited immediately; the user data directory may still be locked")

type manualLoginProcess interface {
	Wait() error
}

var findChromeForManualLogin = launcher.LookPath

var startManualLoginProcess = func(ctx context.Context, chromePath string, args ...string) (manualLoginProcess, error) {
	cmd := exec.CommandContext(ctx, chromePath, args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return cmd, nil
}

// LaunchManualLogin は CDP を使わない素の Chrome を起動し、ユーザーがログインを
// 完了してブラウザを閉じるまで待つ。rod 経由の Chrome では Cloudflare の検証を
// 通過できないため、ここでは launcher を使わず、デバッグポートも開かない。
func LaunchManualLogin(ctx context.Context, url string) error {
	chromePath, found := findChromeForManualLogin()
	if !found {
		return errors.New("Chrome executable not found")
	}
	dir, err := UserDataDir()
	if err != nil {
		return err
	}

	// 直前まで rod が使っていたプロファイルはすぐには解放されないことがある。
	// ロックは数秒で消えるのが普通なので、諦める前に間を置いて何度か試す。
	var lastErr error
	for attempt := 1; attempt <= manualLoginStartAttempts; attempt++ {
		if attempt > 1 {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(manualLoginRetryInterval):
			}
		}
		lastErr = runManualLoginBrowser(ctx, chromePath, dir, url)
		if !errors.Is(lastErr, errManualLoginProfileLocked) {
			return lastErr
		}
		slog.Info("manual login browser exited immediately; the profile is still locked, retrying",
			slog.Int("attempt", attempt),
			slog.Int("attempts", manualLoginStartAttempts))
	}
	return lastErr
}

// runManualLoginBrowser は素の Chrome を1回起動し、閉じられるまで待つ。
func runManualLoginBrowser(ctx context.Context, chromePath, dir, url string) error {
	process, err := startManualLoginProcess(ctx, chromePath, "--user-data-dir="+dir, url)
	if err != nil {
		return fmt.Errorf("start manual login browser: %w", err)
	}

	waitDone := make(chan error, 1)
	go func() {
		waitDone <- process.Wait()
	}()

	timer := time.NewTimer(manualLoginStartupWindow)
	defer timer.Stop()
	select {
	case waitErr := <-waitDone:
		if waitErr != nil {
			return fmt.Errorf("%w: %v", errManualLoginProfileLocked, waitErr)
		}
		return errManualLoginProfileLocked
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}

	select {
	case waitErr := <-waitDone:
		if waitErr != nil {
			return fmt.Errorf("wait for manual login browser: %w", waitErr)
		}
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

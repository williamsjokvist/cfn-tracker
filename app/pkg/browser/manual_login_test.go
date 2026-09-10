package browser_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
)

type fakeManualLoginProcess struct {
	wait func() error
}

func (p fakeManualLoginProcess) Wait() error {
	return p.wait()
}

func TestLaunchManualLoginReturnsWhenContextIsCancelled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	waiting := make(chan struct{})
	restore := browser.ReplaceManualLoginProcessForTest(
		func() (string, bool) { return "fake-chrome", true },
		func(context.Context, string, ...string) (browser.ManualLoginProcessForTest, error) {
			return fakeManualLoginProcess{wait: func() error {
				<-waiting
				return nil
			}}, nil
		},
	)
	t.Cleanup(func() {
		close(waiting)
		restore()
	})

	cancel()
	if err := browser.LaunchManualLogin(ctx, "https://example.com"); !errors.Is(err, context.Canceled) {
		t.Fatalf("LaunchManualLogin() error = %v, want context.Canceled", err)
	}
}

func TestLaunchManualLoginRejectsImmediateExit(t *testing.T) {
	t.Cleanup(browser.SetManualLoginTimingForTest(4, time.Millisecond, 10*time.Millisecond))

	restore := browser.ReplaceManualLoginProcessForTest(
		func() (string, bool) { return "fake-chrome", true },
		func(context.Context, string, ...string) (browser.ManualLoginProcessForTest, error) {
			return fakeManualLoginProcess{wait: func() error { return nil }}, nil
		},
	)
	t.Cleanup(restore)

	err := browser.LaunchManualLogin(context.Background(), "https://example.com")
	if err == nil || !strings.Contains(err.Error(), "user data directory may still be locked") {
		t.Fatalf("LaunchManualLogin() error = %v, want profile lock error", err)
	}
}

// プロファイルのロックは rod の Chromium が残っていると数秒握られたままになる。
// 1回目で諦めず、解放されてから起動できることを確かめる。
func TestLaunchManualLoginRetriesWhileProfileIsLocked(t *testing.T) {
	t.Cleanup(browser.SetManualLoginTimingForTest(4, time.Millisecond, 10*time.Millisecond))

	attempts := 0
	restore := browser.ReplaceManualLoginProcessForTest(
		func() (string, bool) { return "fake-chrome", true },
		func(context.Context, string, ...string) (browser.ManualLoginProcessForTest, error) {
			attempts++
			if attempts < 3 {
				// ロックされたまま。起動を既存インスタンスへ渡して即終了する。
				return fakeManualLoginProcess{wait: func() error { return nil }}, nil
			}
			// ロックが解放された。人間がログインを終えて閉じるまで生き続ける。
			return fakeManualLoginProcess{wait: func() error {
				time.Sleep(50 * time.Millisecond)
				return nil
			}}, nil
		},
	)
	t.Cleanup(restore)

	if err := browser.LaunchManualLogin(context.Background(), "https://example.com"); err != nil {
		t.Fatalf("LaunchManualLogin() error = %v, want nil", err)
	}
	if attempts != 3 {
		t.Fatalf("attempts = %d, want 3", attempts)
	}
}

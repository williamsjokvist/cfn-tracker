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

// A leftover rod Chromium can hold the profile lock for a few seconds.
// Verify we retry and launch once the lock is released.
func TestLaunchManualLoginRetriesWhileProfileIsLocked(t *testing.T) {
	t.Cleanup(browser.SetManualLoginTimingForTest(4, time.Millisecond, 10*time.Millisecond))

	attempts := 0
	restore := browser.ReplaceManualLoginProcessForTest(
		func() (string, bool) { return "fake-chrome", true },
		func(context.Context, string, ...string) (browser.ManualLoginProcessForTest, error) {
			attempts++
			if attempts < 3 {
				// Still locked: hands off to the existing instance and exits immediately.
				return fakeManualLoginProcess{wait: func() error { return nil }}, nil
			}
			// Lock released: stays alive until the user logs in and closes it.
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

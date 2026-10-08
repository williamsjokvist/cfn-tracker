package browser

import (
	"context"
	"time"
)

type ManualLoginProcessForTest interface {
	Wait() error
}

func ReplaceManualLoginProcessForTest(
	findChrome func() (string, bool),
	start func(context.Context, string, ...string) (ManualLoginProcessForTest, error),
) func() {
	originalFind := findChromeForManualLogin
	originalStart := startManualLoginProcess
	findChromeForManualLogin = findChrome
	startManualLoginProcess = func(ctx context.Context, path string, args ...string) (manualLoginProcess, error) {
		return start(ctx, path, args...)
	}
	return func() {
		findChromeForManualLogin = originalFind
		startManualLoginProcess = originalStart
	}
}

// SetManualLoginTimingForTest shortens the retry delays for tests.
func SetManualLoginTimingForTest(attempts int, interval, startupWindow time.Duration) func() {
	originalAttempts := manualLoginStartAttempts
	originalInterval := manualLoginRetryInterval
	originalWindow := manualLoginStartupWindow
	manualLoginStartAttempts = attempts
	manualLoginRetryInterval = interval
	manualLoginStartupWindow = startupWindow
	return func() {
		manualLoginStartAttempts = originalAttempts
		manualLoginRetryInterval = originalInterval
		manualLoginStartupWindow = originalWindow
	}
}

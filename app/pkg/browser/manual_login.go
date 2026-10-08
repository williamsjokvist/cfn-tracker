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
	// Grace period to detect an immediate exit. Chrome hands off to an existing
	// instance and exits within a second, so surviving this long means it started.
	manualLoginStartupWindow = 3 * time.Second
	// Delay between attempts while waiting for the profile lock, and how many attempts.
	manualLoginRetryInterval = 2 * time.Second
	manualLoginStartAttempts = 4
)

// errManualLoginProfileLocked means the launched Chrome exited immediately. The
// UserDataDir is exclusive, so if a leftover rod Chromium still holds it, the new
// process hands off to that instance and exits.
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

// LaunchManualLogin starts a plain Chrome without CDP and waits until the user has
// logged in and closed it. Chrome driven by rod can't pass Cloudflare's check, so no
// launcher or debugging port is used here.
func LaunchManualLogin(ctx context.Context, url string) error {
	chromePath, found := findChromeForManualLogin()
	if !found {
		return errors.New("Chrome executable not found")
	}
	dir, err := UserDataDir()
	if err != nil {
		return err
	}

	// The profile rod just used may not be released right away. The lock usually
	// clears within seconds, so retry a few times before giving up.
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

// runManualLoginBrowser starts plain Chrome once and waits for it to be closed.
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

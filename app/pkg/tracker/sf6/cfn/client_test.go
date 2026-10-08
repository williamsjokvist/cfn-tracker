package cfn

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type fakeAuthBrowser struct {
	sessions    []bool
	manualLogin func(ctx context.Context, url string) error
	relaunchErr error
	calls       []string
}

func (f *fakeAuthBrowser) HasBucklerSession(context.Context) bool {
	f.calls = append(f.calls, "session")
	if len(f.sessions) == 0 {
		return false
	}
	ok := f.sessions[0]
	f.sessions = f.sessions[1:]
	return ok
}

func (f *fakeAuthBrowser) Close() error {
	f.calls = append(f.calls, "close")
	return nil
}

func (f *fakeAuthBrowser) LaunchManualLogin(ctx context.Context, url string) error {
	f.calls = append(f.calls, "manual-login")
	if f.manualLogin == nil {
		return nil
	}
	return f.manualLogin(ctx, url)
}

func (f *fakeAuthBrowser) Relaunch() error {
	f.calls = append(f.calls, "relaunch")
	return f.relaunchErr
}

func authenticate(t *testing.T, auth authBrowser) []tracker.AuthStatus {
	t.Helper()
	statuses := make(chan tracker.AuthStatus, 10)
	client := &Client{auth: auth}
	client.Authenticate(context.Background(), statuses)
	close(statuses)
	var got []tracker.AuthStatus
	for s := range statuses {
		got = append(got, s)
	}
	return got
}

func assertCalls(t *testing.T, f *fakeAuthBrowser, want ...string) {
	t.Helper()
	if !slices.Equal(f.calls, want) {
		t.Fatalf("browser calls = %v, want %v", f.calls, want)
	}
}

func assertDone(t *testing.T, statuses []tracker.AuthStatus) {
	t.Helper()
	last := statuses[len(statuses)-1]
	if last.Err != nil || last.Progress != 100 {
		t.Fatalf("last status = %+v, want progress 100 without error", last)
	}
}

func assertFailed(t *testing.T, statuses []tracker.AuthStatus) {
	t.Helper()
	last := statuses[len(statuses)-1]
	if !errors.Is(last.Err, model.ErrAuthManualLoginFailed) {
		t.Fatalf("last status error = %v, want ErrAuthManualLoginFailed", last.Err)
	}
}

func TestAuthenticateReusesValidSession(t *testing.T) {
	f := &fakeAuthBrowser{sessions: []bool{true}}
	statuses := authenticate(t, f)

	assertCalls(t, f, "session")
	if len(statuses) != 1 {
		t.Fatalf("statuses = %+v, want only the done status", statuses)
	}
	assertDone(t, statuses)
}

func TestAuthenticateWaitsForManualLogin(t *testing.T) {
	var gotURL string
	var hasDeadline bool
	f := &fakeAuthBrowser{
		sessions: []bool{false, true},
		manualLogin: func(ctx context.Context, url string) error {
			gotURL = url
			_, hasDeadline = ctx.Deadline()
			return nil
		},
	}
	statuses := authenticate(t, f)

	assertCalls(t, f, "session", "close", "manual-login", "relaunch", "session")
	if gotURL != bucklerBaseURL+"/ja-jp" {
		t.Fatalf("manual login url = %q, want %q", gotURL, bucklerBaseURL+"/ja-jp")
	}
	if !hasDeadline {
		t.Fatal("manual login context has no deadline")
	}
	if len(statuses) != 2 {
		t.Fatalf("statuses = %+v, want relogin prompt then done", statuses)
	}
	if a := statuses[0].Action; a == nil || a.LocalizationKey != "authNeedRelogin" {
		t.Fatalf("first status = %+v, want authNeedRelogin action", statuses[0])
	}
	assertDone(t, statuses)
}

func TestAuthenticateChecksSessionWhenManualLoginErrors(t *testing.T) {
	// Chrome can exit with an error after the user has already logged in.
	f := &fakeAuthBrowser{
		sessions:    []bool{false, true},
		manualLogin: func(context.Context, string) error { return errors.New("exit status 1") },
	}
	assertDone(t, authenticate(t, f))
}

func TestAuthenticateFailsWhenUserDoesNotLogIn(t *testing.T) {
	f := &fakeAuthBrowser{sessions: []bool{false, false}}
	statuses := authenticate(t, f)

	assertCalls(t, f, "session", "close", "manual-login", "relaunch", "session")
	assertFailed(t, statuses)
}

func TestAuthenticateFailsWhenRelaunchFails(t *testing.T) {
	f := &fakeAuthBrowser{
		sessions:    []bool{false, true},
		relaunchErr: errors.New("profile locked"),
	}
	statuses := authenticate(t, f)

	assertCalls(t, f, "session", "close", "manual-login", "relaunch")
	assertFailed(t, statuses)
}

func TestAuthenticateWithoutBrowser(t *testing.T) {
	statuses := authenticate(t, nil)
	if len(statuses) != 1 || statuses[0].Err == nil {
		t.Fatalf("statuses = %+v, want a single error", statuses)
	}
}

func TestAuthenticateReturnsWhenCallerStopsListening(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	f := &fakeAuthBrowser{
		manualLogin: func(ctx context.Context, _ string) error {
			cancel()
			<-ctx.Done()
			return ctx.Err()
		},
	}
	// Unbuffered so the send after the manual login blocks until ctx is done.
	statuses := make(chan tracker.AuthStatus)
	done := make(chan struct{})
	go func() {
		(&Client{auth: f}).Authenticate(ctx, statuses)
		close(done)
	}()

	// The relogin prompt is sent before the manual login, so read it to get there.
	<-statuses
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("Authenticate blocked after the caller stopped listening")
	}
}

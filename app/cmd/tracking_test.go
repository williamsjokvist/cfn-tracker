package cmd_test

import (
	"context"
	"errors"
	"fmt"
	"net"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

func TestTrackingSelectGame(t *testing.T) {
	tests := []struct {
		name      string
		gameType  model.GameType
		expErrMsg string
	}{
		{name: "select tekken 8", gameType: model.GameTypeT8},
		{name: "select SF6", gameType: model.GameTypeSF6, expErrMsg: "unauthenticated: browser not initialized"},
		{name: "select game that doesn't exist", gameType: "undefined", expErrMsg: "select game: game does not exist"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := testSuite.trackingHandler.SelectGame(tt.gameType)
			if err == nil && tt.expErrMsg != "" {
				t.Errorf("expected error: %q, got: nil", tt.expErrMsg)
			}
			if err != nil && err.Error() != tt.expErrMsg {
				t.Errorf("unexpected error: got: %q, want: %q", err, tt.expErrMsg)
			}
		})
	}
}

type fakeTracker struct {
	mu              sync.Mutex
	polls           int
	poll            func(int) (*model.Match, error)
	pollWithContext func(context.Context, int) (*model.Match, error)
}

func (f *fakeTracker) GetUser(context.Context, string) (*model.User, error) {
	return &model.User{Code: fmt.Sprintf("test-%d", time.Now().UnixNano()), DisplayName: "test"}, nil
}
func (f *fakeTracker) Poll(ctx context.Context, _ *model.Session) (*model.Match, error) {
	f.mu.Lock()
	f.polls++
	polls := f.polls
	f.mu.Unlock()
	if f.pollWithContext != nil {
		return f.pollWithContext(ctx, polls)
	}
	return f.poll(polls)
}
func (f *fakeTracker) Authenticate(_ context.Context, _, _ string, statuses chan tracker.AuthStatus) {
	statuses <- tracker.AuthStatus{Progress: 100}
}

func runTracking(t *testing.T, fake *fakeTracker, emit func(string, ...interface{})) <-chan error {
	t.Helper()
	testSuite.trackingHandler.SetGameTracker(fake)
	testSuite.trackingHandler.SetEventEmitter(emit)
	done := make(chan error, 1)
	go func() { done <- testSuite.trackingHandler.StartTracking("test", false) }()
	return done
}

func TestPollTransientErrorRetries(t *testing.T) {
	matchSeen := make(chan struct{}, 1)
	fake := &fakeTracker{poll: func(n int) (*model.Match, error) {
		if n <= 3 {
			return nil, errors.New("temporary")
		}
		return &model.Match{ReplayID: "recovered"}, nil
	}}
	done := runTracking(t, fake, func(name string, data ...interface{}) {
		if name == "match" && len(data) > 0 {
			if match, ok := data[0].(model.Match); ok && match.ReplayID == "recovered" {
				select {
				case matchSeen <- struct{}{}:
				default:
				}
			}
		}
	})
	select {
	case <-matchSeen:
		testSuite.trackingHandler.StopTracking()
	case <-time.After(20 * time.Second):
		t.Fatal("match did not arrive after retries")
	}
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("StartTracking did not stop")
	}
}

func TestPollTimeoutRetriesWithoutTrackingErrorAndStopsOnParentCancel(t *testing.T) {
	restoreTimeout := cmd.SetPollTimeoutForTest(50 * time.Millisecond)
	t.Cleanup(restoreTimeout)

	pollDone := make(chan int, 4)
	fake := &fakeTracker{pollWithContext: func(ctx context.Context, n int) (*model.Match, error) {
		<-ctx.Done()
		if !errors.Is(ctx.Err(), context.DeadlineExceeded) && n < 2 {
			t.Errorf("Poll context error = %v, want deadline exceeded", ctx.Err())
		}
		pollDone <- n
		return nil, nil
	}}
	trackingError := make(chan struct{}, 1)
	done := runTracking(t, fake, func(name string, _ ...interface{}) {
		if name == "tracking-error" {
			trackingError <- struct{}{}
		}
	})

	select {
	case n := <-pollDone:
		if n != 1 {
			t.Fatalf("first completed Poll call = %d, want 1", n)
		}
	case <-time.After(time.Second):
		t.Fatal("first Poll context did not time out")
	}
	select {
	case err := <-done:
		t.Fatalf("tracking stopped after child timeout: %v", err)
	default:
	}
	testSuite.trackingHandler.ForcePoll()
	select {
	case n := <-pollDone:
		if n < 2 {
			t.Fatalf("completed Poll calls = %d, want at least 2", n)
		}
	case <-time.After(time.Second):
		t.Fatal("Poll was not retried after its context timed out")
	}
	select {
	case <-trackingError:
		t.Fatal("tracking-error was emitted for a Poll timeout")
	default:
	}

	testSuite.trackingHandler.StopTracking()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("poll loop did not return promptly after parent cancellation")
	}
}

func TestPollTransientErrorNeverGivesUp(t *testing.T) {
	retries := make(chan struct{}, 4)
	fake := &fakeTracker{poll: func(int) (*model.Match, error) { return nil, errors.New("temporary") }}
	done := runTracking(t, fake, func(name string, _ ...interface{}) {
		if name == "tracking-retrying" {
			retries <- struct{}{}
		}
	})
	for i := 0; i < 2; i++ {
		select {
		case <-retries:
		case <-time.After(8 * time.Second):
			t.Fatal("polling stopped retrying")
		}
	}
	select {
	case err := <-done:
		t.Fatalf("tracking ended unexpectedly: %v", err)
	default:
	}
	testSuite.trackingHandler.StopTracking()
	<-done
}

func TestPollFatalErrorEmitsAndStops(t *testing.T) {
	events := make(chan string, 4)
	fake := &fakeTracker{poll: func(int) (*model.Match, error) { return nil, &model.HTTPStatusError{Op: "poll", StatusCode: 404} }}
	done := runTracking(t, fake, func(name string, _ ...interface{}) { events <- name })
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("expected fatal error")
		}
	case <-time.After(time.Second):
		t.Fatal("StartTracking did not return")
	}
	seenError, seenStopped := false, false
	close(events)
	for event := range events {
		seenError = seenError || event == "tracking-error"
		seenStopped = seenStopped || event == "stopped-tracking"
	}
	if !seenError || !seenStopped {
		t.Fatalf("missing terminal events: error=%v stopped=%v", seenError, seenStopped)
	}
}

func TestStartTrackingAlwaysReturns(t *testing.T) {
	for _, err := range []error{&model.HTTPStatusError{Op: "poll", StatusCode: 404}, &model.HTTPStatusError{Op: "poll", StatusCode: 400}} {
		done := runTracking(t, &fakeTracker{poll: func(int) (*model.Match, error) { return nil, err }}, func(string, ...interface{}) {})
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatalf("StartTracking did not return for %v", err)
		}
	}
}

func TestStartTrackingNoGoroutineLeak(t *testing.T) {
	before := runtime.NumGoroutine()
	done := runTracking(t, &fakeTracker{poll: func(int) (*model.Match, error) { return nil, &model.HTTPStatusError{Op: "poll", StatusCode: 404} }}, func(string, ...interface{}) {})
	<-done
	time.Sleep(20 * time.Millisecond)
	if after := runtime.NumGoroutine(); after > before+1 {
		t.Fatalf("goroutines grew from %d to %d", before, after)
	}
}

func TestForcePollDoesNotBlockWhenStopped(t *testing.T) {
	done := make(chan struct{})
	go func() { testSuite.trackingHandler.ForcePoll(); close(done) }()
	select {
	case <-done:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("ForcePoll blocked")
	}
}

func TestClassifyPollError(t *testing.T) {
	parse := &model.ParseError{Op: "parse", Err: errors.New("bad")}
	timeout := &net.DNSError{IsTimeout: true}
	tests := []struct {
		err  error
		want model.ErrorClass
	}{
		{&model.HTTPStatusError{StatusCode: 401}, model.ClassAuth}, {&model.HTTPStatusError{StatusCode: 403}, model.ClassAuth}, {&model.HTTPStatusError{StatusCode: 404}, model.ClassFatal}, {&model.HTTPStatusError{StatusCode: 429}, model.ClassTransient}, {&model.HTTPStatusError{StatusCode: 503}, model.ClassTransient}, {&model.HTTPStatusError{StatusCode: 400}, model.ClassFatal}, {parse, model.ClassParse}, {timeout, model.ClassTransient}, {context.DeadlineExceeded, model.ClassTransient}, {errors.New("unknown"), model.ClassTransient},
	}
	for _, tt := range tests {
		if got := model.ClassifyPollError(tt.err); got != tt.want {
			t.Errorf("ClassifyPollError(%v)=%v, want %v", tt.err, got, tt.want)
		}
	}
}

func TestPollTimeoutErrorRetriesAutomatically(t *testing.T) {
	restore := cmd.SetPollTimeoutForTest(50 * time.Millisecond)
	t.Cleanup(restore)

	retries := make(chan struct{}, 4)
	trackingError := make(chan struct{}, 1)
	// 実際の GetBattleLog はタイムアウト時に DeadlineExceeded をラップして返す。
	fake := &fakeTracker{pollWithContext: func(ctx context.Context, _ int) (*model.Match, error) {
		<-ctx.Done()
		return nil, fmt.Errorf("cfn: get battle log: %w", ctx.Err())
	}}
	done := runTracking(t, fake, func(name string, data ...interface{}) {
		switch name {
		case "tracking-retrying":
			if len(data) > 0 {
				if status, ok := data[0].(cmd.RetryStatus); ok && status.Reason != "errNetwork" {
					t.Errorf("retry reason = %q, want errNetwork (transient)", status.Reason)
				}
			}
			select {
			case retries <- struct{}{}:
			default:
			}
		case "tracking-error":
			select {
			case trackingError <- struct{}{}:
			default:
			}
		}
	})

	// ForcePoll を使わず、バックオフ経由で自動的に再試行され続けることを確認する。
	for i := 0; i < 2; i++ {
		select {
		case <-retries:
		case <-time.After(8 * time.Second):
			t.Fatal("poll loop stopped retrying after timeout errors")
		}
	}
	select {
	case <-trackingError:
		t.Fatal("tracking-error was emitted for a Poll timeout")
	default:
	}
	select {
	case err := <-done:
		t.Fatalf("tracking ended unexpectedly: %v", err)
	default:
	}

	testSuite.trackingHandler.StopTracking()
	<-done
}

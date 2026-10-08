package cmd

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
)

type authTestCFNClient struct {
	authenticate func(context.Context, chan tracker.AuthStatus)
}

func (f *authTestCFNClient) GetBattleLog(context.Context, string) (*cfn.BattleLog, error) {
	return nil, errors.New("not used")
}

func (f *authTestCFNClient) Authenticate(ctx context.Context, statuses chan tracker.AuthStatus) {
	f.authenticate(ctx, statuses)
}

func authTestHandler(authenticate func(context.Context, chan tracker.AuthStatus)) *TrackingHandler {
	return NewTrackingHandler(nil, &authTestCFNClient{authenticate: authenticate}, nil, nil, nil, &config.BuildConfig{})
}

func TestSelectGameEmitsAuthActionRequired(t *testing.T) {
	handler := authTestHandler(func(_ context.Context, statuses chan tracker.AuthStatus) {
		statuses <- tracker.AuthStatus{Action: &tracker.AuthAction{LocalizationKey: "authSolveCaptcha", SecondsLeft: 42}}
		close(statuses)
	})
	var mu sync.Mutex
	events := make([]string, 0, 1)
	handler.SetEventEmitter(func(name string, _ ...interface{}) {
		mu.Lock()
		defer mu.Unlock()
		events = append(events, name)
	})

	if err := handler.SelectGame(model.GameTypeSF6); err != nil {
		t.Fatalf("SelectGame() error = %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if len(events) != 1 || events[0] != "auth-action-required" {
		t.Fatalf("events = %v, want only auth-action-required", events)
	}
}

func TestSelectGameReturnsOnAuthError(t *testing.T) {
	handler := authTestHandler(func(_ context.Context, statuses chan tracker.AuthStatus) {
		for i := 0; i < 3; i++ {
			statuses <- tracker.AuthStatus{Action: &tracker.AuthAction{LocalizationKey: "authSolveCaptcha", SecondsLeft: 3 - i}}
		}
		statuses <- tracker.AuthStatus{Err: errors.New("manual authentication failed")}
	})
	done := make(chan error, 1)
	go func() { done <- handler.SelectGame(model.GameTypeSF6) }()
	select {
	case err := <-done:
		if err == nil {
			t.Fatal("SelectGame() returned nil, want error")
		}
	case <-time.After(time.Second):
		t.Fatal("SelectGame() did not return after authentication error")
	}
}

func TestSelectGameDoesNotHangWhenChannelNeverCloses(t *testing.T) {
	originalTimeout := selectGameTimeout
	selectGameTimeout = 20 * time.Millisecond
	t.Cleanup(func() { selectGameTimeout = originalTimeout })
	handler := authTestHandler(func(context.Context, chan tracker.AuthStatus) {})
	done := make(chan error, 1)
	go func() { done <- handler.SelectGame(model.GameTypeSF6) }()
	select {
	case err := <-done:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("SelectGame() error = %v, want deadline exceeded", err)
		}
	case <-time.After(time.Second):
		t.Fatal("SelectGame() hung on a channel that never closed")
	}
}

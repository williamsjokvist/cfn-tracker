package cmd_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

func filterEvents(events []eventRecord, name string) []eventRecord {
	var result []eventRecord
	for _, e := range events {
		if e.Name == name {
			result = append(result, e)
		}
	}
	return result
}

func TestTrackingSelectGame(t *testing.T) {
	tests := []struct {
		name      string
		gameType  model.GameType
		expErrMsg string
	}{
		{
			name:     "select tekken 8",
			gameType: model.GameTypeT8,
		},
		{
			name:      "select SF6",
			gameType:  model.GameTypeSF6,
			expErrMsg: "unauthenticated: browser not initialized",
		},
		{
			name:      "select game that doesn't exist",
			gameType:  "undefined",
			expErrMsg: "select game: game does not exist",
		},
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

func TestTickSingleMatch(t *testing.T) {
	match1 := &model.Match{
		Victory: true, LP: 1500, MR: 0, Character: "Ryu",
		Wins: 1, Losses: 0, WinStreak: 1, WinRate: 100,
		Opponent: "Ken", ReplayID: "replay-001",
		Date: "2026-01-01", Time: "12:00",
	}
	handler, sqlDb, session, events := setupTrackingTest(t, tracker.NewMockPollResult(match1))
	match1.SessionId = session.Id
	match1.UserId = session.UserId
	ctx := context.Background()

	got, err := handler.Tick(ctx, session)
	if err != nil {
		t.Fatalf("Tick() error: %v", err)
	}
	if got == nil {
		t.Fatal("Tick() returned nil match, want non-nil")
	}

	if session.LP != 1500 {
		t.Errorf("session.LP = %d, want 1500", session.LP)
	}
	if len(session.Matches) != 1 {
		t.Errorf("len(session.Matches) = %d, want 1", len(session.Matches))
	}

	matches, err := sqlDb.GetMatches(ctx, session.Id, session.UserId, 0, 0)
	if err != nil {
		t.Fatalf("GetMatches: %v", err)
	}
	if len(matches) != 1 {
		t.Errorf("DB matches count = %d, want 1", len(matches))
	}
	if len(matches) > 0 && matches[0].LP != 1500 {
		t.Errorf("DB match LP = %d, want 1500", matches[0].LP)
	}

	matchEvents := filterEvents(*events, "match")
	if len(matchEvents) != 1 {
		t.Errorf("match events count = %d, want 1", len(matchEvents))
	}
}

func TestTickMultipleMatches(t *testing.T) {
	matchData := []*model.Match{
		{Victory: true, LP: 1500, Wins: 1, Losses: 0, WinStreak: 1, WinRate: 100, ReplayID: "r1", Date: "2026-01-01", Time: "12:00"},
		{Victory: true, LP: 1550, Wins: 2, Losses: 0, WinStreak: 2, WinRate: 100, ReplayID: "r2", Date: "2026-01-01", Time: "12:01"},
		{Victory: false, LP: 1500, Wins: 2, Losses: 1, WinStreak: 0, WinRate: 66, ReplayID: "r3", Date: "2026-01-01", Time: "12:02"},
	}
	handler, sqlDb, session, events := setupTrackingTest(t,
		tracker.NewMockPollResult(matchData[0]),
		tracker.NewMockPollResult(matchData[1]),
		tracker.NewMockPollResult(matchData[2]),
	)
	for _, m := range matchData {
		m.SessionId = session.Id
		m.UserId = session.UserId
	}
	ctx := context.Background()

	for i, want := range matchData {
		got, err := handler.Tick(ctx, session)
		if err != nil {
			t.Fatalf("Tick() %d error: %v", i, err)
		}
		if got == nil {
			t.Fatalf("Tick() %d returned nil", i)
		}
		if session.LP != want.LP {
			t.Errorf("after tick %d: session.LP = %d, want %d", i, session.LP, want.LP)
		}
		if len(session.Matches) != i+1 {
			t.Errorf("after tick %d: len(session.Matches) = %d, want %d", i, len(session.Matches), i+1)
		}
	}

	dbMatches, err := sqlDb.GetMatches(ctx, session.Id, session.UserId, 0, 0)
	if err != nil {
		t.Fatalf("GetMatches: %v", err)
	}
	if len(dbMatches) != 3 {
		t.Errorf("DB matches count = %d, want 3", len(dbMatches))
	}

	matchEvents := filterEvents(*events, "match")
	if len(matchEvents) != 3 {
		t.Errorf("match events = %d, want 3", len(matchEvents))
	}
}

func TestTickNoMatch(t *testing.T) {
	handler, sqlDb, session, events := setupTrackingTest(t, tracker.NewMockPollResult(nil))
	ctx := context.Background()
	originalLP := session.LP

	got, err := handler.Tick(ctx, session)
	if err != nil {
		t.Fatalf("Tick() error: %v", err)
	}
	if got != nil {
		t.Errorf("Tick() returned %v, want nil", got)
	}
	if session.LP != originalLP {
		t.Errorf("session.LP changed to %d, want %d (unchanged)", session.LP, originalLP)
	}
	if len(session.Matches) != 0 {
		t.Errorf("session.Matches has %d entries, want 0", len(session.Matches))
	}

	dbMatches, err := sqlDb.GetMatches(ctx, session.Id, session.UserId, 0, 0)
	if err != nil {
		t.Fatalf("GetMatches: %v", err)
	}
	if len(dbMatches) != 0 {
		t.Errorf("DB matches count = %d, want 0", len(dbMatches))
	}

	matchEvents := filterEvents(*events, "match")
	if len(matchEvents) != 0 {
		t.Errorf("match events = %d, want 0", len(matchEvents))
	}
}

func TestTickPollError(t *testing.T) {
	pollErr := fmt.Errorf("network timeout")
	handler, sqlDb, session, events := setupTrackingTest(t, tracker.NewMockPollError(pollErr))
	ctx := context.Background()

	got, err := handler.Tick(ctx, session)
	if err == nil {
		t.Fatal("Tick() returned nil error, want error")
	}
	if got != nil {
		t.Errorf("Tick() returned match %v, want nil on error", got)
	}
	if len(session.Matches) != 0 {
		t.Errorf("session.Matches has %d entries, want 0", len(session.Matches))
	}

	dbMatches, err := sqlDb.GetMatches(ctx, session.Id, session.UserId, 0, 0)
	if err != nil {
		t.Fatalf("GetMatches: %v", err)
	}
	if len(dbMatches) != 0 {
		t.Errorf("DB matches count = %d, want 0", len(dbMatches))
	}

	matchEvents := filterEvents(*events, "match")
	if len(matchEvents) != 0 {
		t.Errorf("match events = %d, want 0", len(matchEvents))
	}
}

func TestTickWinStreak(t *testing.T) {
	streakMatches := []*model.Match{
		{Victory: true, WinStreak: 1, Wins: 1, Losses: 0, LP: 1100, ReplayID: "r1", Date: "2026-01-01", Time: "12:00"},
		{Victory: true, WinStreak: 2, Wins: 2, Losses: 0, LP: 1200, ReplayID: "r2", Date: "2026-01-01", Time: "12:01"},
		{Victory: false, WinStreak: 0, Wins: 2, Losses: 1, LP: 1150, ReplayID: "r3", Date: "2026-01-01", Time: "12:02"},
		{Victory: true, WinStreak: 1, Wins: 3, Losses: 1, LP: 1250, ReplayID: "r4", Date: "2026-01-01", Time: "12:03"},
	}
	results := []tracker.MockPollResult{
		tracker.NewMockPollResult(streakMatches[0]),
		tracker.NewMockPollResult(streakMatches[1]),
		tracker.NewMockPollResult(streakMatches[2]),
		tracker.NewMockPollResult(streakMatches[3]),
	}
	expectedStreaks := []int{1, 2, 0, 1}

	handler, _, session, _ := setupTrackingTest(t, results...)
	for _, m := range streakMatches {
		m.SessionId = session.Id
		m.UserId = session.UserId
	}
	ctx := context.Background()

	for i, expected := range expectedStreaks {
		got, err := handler.Tick(ctx, session)
		if err != nil {
			t.Fatalf("Tick() %d error: %v", i, err)
		}
		if got == nil {
			t.Fatalf("Tick() %d returned nil", i)
		}
		if got.WinStreak != expected {
			t.Errorf("tick %d: WinStreak = %d, want %d", i, got.WinStreak, expected)
		}
		if session.Matches[0].WinStreak != expected {
			t.Errorf("tick %d: session.Matches[0].WinStreak = %d, want %d", i, session.Matches[0].WinStreak, expected)
		}
	}
}

func TestSetGameTracker(t *testing.T) {
	sqlDb, err := sql.NewStorage(true)
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	t.Cleanup(func() { sqlDb.DB().Close() })

	ctx := context.Background()
	user := &model.User{DisplayName: "Player", Code: "player-456"}
	if err := sqlDb.SaveUser(ctx, *user); err != nil {
		t.Fatalf("save user: %v", err)
	}
	session, err := sqlDb.CreateSession(ctx, user.Code)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	session.UserName = user.DisplayName

	match := &model.Match{Victory: true, LP: 2000, Wins: 1, ReplayID: "r1", Date: "2026-01-01", Time: "12:00"}
	match.SessionId = session.Id
	match.UserId = session.UserId
	mock := tracker.NewMockGameTracker(user, tracker.NewMockPollResult(match))

	cfg := &config.BuildConfig{}
	handler := cmd.NewTrackingHandler(nil, nil, sqlDb, nil, nil, cfg)
	handler.SetGameTracker(mock)
	handler.SetEventEmitter(func(eventName string, optionalData ...interface{}) {})

	got, err := handler.Tick(ctx, session)
	if err != nil {
		t.Fatalf("Tick() error: %v", err)
	}
	if got == nil {
		t.Fatal("Tick() returned nil, want match")
	}
	if got.LP != 2000 {
		t.Errorf("got.LP = %d, want 2000", got.LP)
	}
}

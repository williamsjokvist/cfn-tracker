package sql

import (
	"context"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestGetUserStatistics(t *testing.T) {
	ctx := context.Background()
	storage, err := NewInMemoryStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	for _, code := range []string{"123456789", "987654321"} {
		if err := storage.SaveUser(ctx, model.User{Code: code, DisplayName: code}); err != nil {
			t.Fatalf("SaveUser() error = %v", err)
		}
	}
	first, err := storage.CreateSession(ctx, "123456789", "sf6")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	second, err := storage.CreateSession(ctx, "123456789", "sf6")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	other, err := storage.CreateSession(ctx, "987654321", "sf6")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	matches := []model.Match{
		{UserId: "123456789", SessionId: first.Id, Character: "Ken", Victory: true, LP: 1000, MR: 1500, Date: "2026-01-01", Time: "10:00"},
		{UserId: "123456789", SessionId: first.Id, Character: "Ryu", Victory: false, LP: 500, Date: "2026-01-01", Time: "10:05"},
		{UserId: "123456789", SessionId: second.Id, Character: "Ken", Victory: false, LP: 1100, MR: 1520, Date: "2026-01-02", Time: "09:00"},
		{UserId: "123456789", SessionId: second.Id, Character: "Ken", Victory: true, LP: 1200, MR: 1540, Date: "2026-01-02", Time: "09:05"},
		{UserId: "987654321", SessionId: other.Id, Character: "Ken", Victory: true, LP: 9000, Date: "2026-01-03", Time: "09:00"},
	}
	for _, m := range matches {
		if err := storage.SaveMatch(ctx, m); err != nil {
			t.Fatalf("SaveMatch() error = %v", err)
		}
	}

	got, err := storage.GetUserStatistics(ctx, "123456789")
	if err != nil {
		t.Fatalf("GetUserStatistics() error = %v", err)
	}
	if got.Sessions != 2 || got.Wins != 2 || got.Losses != 2 || got.MainCharacter != "Ken" {
		t.Fatalf("GetUserStatistics() = %d sessions, %d-%d, main %q, want 2 sessions, 2-2, main Ken", got.Sessions, got.Wins, got.Losses, got.MainCharacter)
	}
	want := []model.CharacterStatistics{
		{Character: "Ken", Wins: 2, Losses: 1, LP: 1200, MR: 1540},
		{Character: "Ryu", Wins: 0, Losses: 1, LP: 500},
	}
	if len(got.Characters) != len(want) {
		t.Fatalf("GetUserStatistics() characters = %d, want %d", len(got.Characters), len(want))
	}
	for i, c := range got.Characters {
		if *c != want[i] {
			t.Errorf("GetUserStatistics() character %d = %+v, want %+v", i, *c, want[i])
		}
	}

	wantHistory := []model.HistoryMatch{
		{Date: "2026-01-01", Time: "10:00", Character: "Ken", Victory: true, LP: 1000, MR: 1500},
		{Date: "2026-01-01", Time: "10:05", Character: "Ryu", LP: 500},
		{Date: "2026-01-02", Time: "09:00", Character: "Ken", LP: 1100, MR: 1520},
		{Date: "2026-01-02", Time: "09:05", Character: "Ken", Victory: true, LP: 1200, MR: 1540},
	}
	if len(got.History) != len(wantHistory) {
		t.Fatalf("GetUserStatistics() history = %d, want %d", len(got.History), len(wantHistory))
	}
	for i, h := range got.History {
		if *h != wantHistory[i] {
			t.Errorf("GetUserStatistics() history %d = %+v, want %+v", i, *h, wantHistory[i])
		}
	}

	empty, err := storage.GetUserStatistics(ctx, "000000000")
	if err != nil {
		t.Fatalf("GetUserStatistics() of a user without matches error = %v", err)
	}
	if empty.Sessions != 0 || len(empty.Characters) != 0 || len(empty.History) != 0 || empty.MainCharacter != "" {
		t.Fatalf("GetUserStatistics() of a user without matches = %+v, want nothing", empty)
	}
}

package sql

import (
	"context"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	storage, err := NewStorage(true)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "123456789", DisplayName: "Kakeru"}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	created, err := storage.CreateSession(ctx, "123456789")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	got, err := storage.GetSession(ctx, created.Id)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if got.Id != created.Id || got.UserId != "123456789" || got.UserName != "Kakeru" || got.CreatedAt != created.CreatedAt {
		t.Fatalf("GetSession() = %+v, want session %d by Kakeru created at %s", got, created.Id, created.CreatedAt)
	}

	if _, err := storage.GetSession(ctx, created.Id+1); err == nil {
		t.Fatal("GetSession() of a missing session returned no error")
	}
}

func TestGetSessionsGains(t *testing.T) {
	ctx := context.Background()
	storage, err := NewStorage(true)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "111222333", DisplayName: "Tokido"}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	sesh, err := storage.CreateSession(ctx, "111222333")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	matches := []model.Match{
		{UserId: "111222333", SessionId: sesh.Id, Character: "Ken", Victory: true, LPGain: 100, MRGain: 10, Date: "2026-01-01", Time: "10:00"},
		{UserId: "111222333", SessionId: sesh.Id, Character: "Ken", Victory: false, LPGain: 80, MRGain: 2, Date: "2026-01-01", Time: "10:05"},
		{UserId: "111222333", SessionId: sesh.Id, Character: "Ryu", Victory: false, LPGain: -30, MRGain: -8, Date: "2026-01-01", Time: "10:10"},
	}
	for _, m := range matches {
		if err := storage.SaveMatch(ctx, m); err != nil {
			t.Fatalf("SaveMatch() error = %v", err)
		}
	}

	sessions, err := storage.GetSessions(ctx, "111222333", "", 0, 0)
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("GetSessions() returned %d sessions, want 1", len(sessions))
	}
	if got := sessions[0]; got.LPGain != 50 || got.MRGain != -6 || got.MatchesWon != 1 || got.MatchesLost != 2 {
		t.Fatalf("GetSessions() = %+v, want 1-2 with +50 LP and -6 MR", got)
	}
}

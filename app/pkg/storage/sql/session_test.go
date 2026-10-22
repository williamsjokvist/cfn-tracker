package sql

import (
	"context"
	"fmt"
	"maps"
	"path/filepath"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestGetSession(t *testing.T) {
	ctx := context.Background()
	storage, err := NewInMemoryStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "123456789", DisplayName: "Kakeru"}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	created, err := storage.CreateSession(ctx, "123456789", "sf6")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	got, err := storage.GetSession(ctx, created.Id)
	if err != nil {
		t.Fatalf("GetSession() error = %v", err)
	}
	if got.Id != created.Id || got.UserId != "123456789" || got.UserName != "Kakeru" || got.GameId != "sf6" || got.CreatedAt != created.CreatedAt {
		t.Fatalf("GetSession() = %+v, want sf6 session %d by Kakeru created at %s", got, created.Id, created.CreatedAt)
	}

	if _, err := storage.GetSession(ctx, created.Id+1); err == nil {
		t.Fatal("GetSession() of a missing session returned no error")
	}
}

func TestGetSessionsGains(t *testing.T) {
	ctx := context.Background()
	storage, err := NewInMemoryStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "111222333", DisplayName: "Tokido"}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	sesh, err := storage.CreateSession(ctx, "111222333", "sf6")
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

func TestMigrationGivesExistingSessionsAndMatchesAGameId(t *testing.T) {
	dataSource := filepath.Join(t.TempDir(), "cfn-tracker.db")
	beforeGameIds := 4
	if err := migrateSchemaAt(dataSource, &beforeGameIds); err != nil {
		t.Fatalf("migrate to before game IDs: %v", err)
	}
	db, err := sqlx.Open("sqlite", dataSource)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.Exec(`
		INSERT INTO users (code, display_name, game) VALUES ('3141592653', 'Kakeru', 'sf6'), ('4amNGJ6EdB6J', 'Tekk!', 't8');
		INSERT INTO sessions (id, user_id) VALUES (1, '3141592653'), (2, '4amNGJ6EdB6J');
		INSERT INTO matches (user_id, session_id, character, date, time) VALUES
			('3141592653', 1, 'Ken', '2026-01-01', '10:00'),
			('4amNGJ6EdB6J', 2, 'Jin', '2026-01-01', '10:00');
	`); err != nil {
		t.Fatalf("insert rows: %v", err)
	}

	if err := migrateSchemaAt(dataSource, nil); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	gameIds := map[string]string{}
	rows, err := db.Query(`
		SELECT 'session ' || id, game_id FROM sessions
		UNION ALL
		SELECT 'match ' || session_id, game_id FROM matches
	`)
	if err != nil {
		t.Fatalf("select game IDs: %v", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var row, gameId string
		if err := rows.Scan(&row, &gameId); err != nil {
			t.Fatalf("scan game ID: %v", err)
		}
		gameIds[row] = gameId
	}
	want := map[string]string{"session 1": "sf6", "match 1": "sf6", "session 2": "t8", "match 2": "t8"}
	if !maps.Equal(gameIds, want) {
		t.Fatalf("game IDs after migration = %v, want %v", gameIds, want)
	}
}

func TestGetLatestSessionOfGame(t *testing.T) {
	ctx := context.Background()
	storage, err := NewInMemoryStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "naggamurasang", DisplayName: "naggamurasang", Game: model.GameTypeFightcade}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	ids := map[string]uint16{}
	for i, gameId := range []string{"sfiii3nr1", "kof98"} {
		sesh, err := storage.CreateSession(ctx, "naggamurasang", gameId)
		if err != nil {
			t.Fatalf("CreateSession(%s) error = %v", gameId, err)
		}
		ids[gameId] = sesh.Id
		match := model.Match{UserId: "naggamurasang", SessionId: sesh.Id, GameId: gameId, Date: "2026-10-12", Time: fmt.Sprintf("20:0%d", i)}
		if err := storage.SaveMatch(ctx, match); err != nil {
			t.Fatalf("SaveMatch() error = %v", err)
		}
	}

	for gameId, want := range ids {
		got, err := storage.GetLatestSession(ctx, "naggamurasang", gameId)
		if err != nil || got == nil || got.Id != want || got.GameId != gameId || len(got.Matches) != 1 || got.Matches[0].GameId != gameId {
			t.Fatalf("GetLatestSession(%s) = %+v, %v, want session %d with its one match", gameId, got, err, want)
		}
	}
}

func TestGetSessionsOfUserInMonth(t *testing.T) {
	ctx := context.Background()
	storage, err := NewInMemoryStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	for _, code := range []string{"5550001", "5550002"} {
		if err := storage.SaveUser(ctx, model.User{Code: code, DisplayName: code}); err != nil {
			t.Fatalf("SaveUser() error = %v", err)
		}
		sesh, err := storage.CreateSession(ctx, code, "sf6")
		if err != nil {
			t.Fatalf("CreateSession() error = %v", err)
		}
		if err := storage.SaveMatch(ctx, model.Match{UserId: code, SessionId: sesh.Id, Date: "2026-10-12", Time: "20:00"}); err != nil {
			t.Fatalf("SaveMatch() error = %v", err)
		}
	}

	sessions, err := storage.GetSessions(ctx, "5550001", time.Now().Format("2006-01"), 0, 0)
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) != 1 || sessions[0].UserId != "5550001" {
		t.Fatalf("GetSessions() = %+v, want only user 5550001's session", sessions)
	}
}

func TestGetLatestSessionWithoutMatches(t *testing.T) {
	ctx := context.Background()
	storage, err := NewInMemoryStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "dakidmelo", DisplayName: "dakidmelo", Game: model.GameTypeFightcade}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	played, err := storage.CreateSession(ctx, "dakidmelo", "kf2k2pls")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	if err := storage.SaveMatch(ctx, model.Match{UserId: "dakidmelo", SessionId: played.Id, GameId: "kf2k2pls", Date: "2026-10-12", Time: "20:00"}); err != nil {
		t.Fatalf("SaveMatch() error = %v", err)
	}
	// Stopped before the first match was played
	empty, err := storage.CreateSession(ctx, "dakidmelo", "kf2k2pls")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	got, err := storage.GetLatestSession(ctx, "dakidmelo", "kf2k2pls")
	if err != nil || got == nil || got.Id != empty.Id || got.UserName != "dakidmelo" || got.GameId != "kf2k2pls" || len(got.Matches) != 0 {
		t.Fatalf("GetLatestSession() = %+v, %v, want empty session %d", got, err, empty.Id)
	}

	if got, err := storage.GetLatestSession(ctx, "dakidmelo", "sfiii3nr1"); err != nil || got != nil {
		t.Fatalf("GetLatestSession() of an unplayed game = %+v, %v, want none", got, err)
	}
}

func TestDeleteSessionIfEmpty(t *testing.T) {
	ctx := context.Background()
	storage, err := NewInMemoryStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "5550003", DisplayName: "5550003"}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	played, err := storage.CreateSession(ctx, "5550003", "sf6")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	if err := storage.SaveMatch(ctx, model.Match{UserId: "5550003", SessionId: played.Id, Date: "2026-10-12", Time: "20:00"}); err != nil {
		t.Fatalf("SaveMatch() error = %v", err)
	}
	empty, err := storage.CreateSession(ctx, "5550003", "sf6")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}

	if deleted, err := storage.DeleteSessionIfEmpty(ctx, played.Id); err != nil || deleted {
		t.Fatalf("DeleteSessionIfEmpty(played) = %v, %v, want it kept", deleted, err)
	}
	if deleted, err := storage.DeleteSessionIfEmpty(ctx, empty.Id); err != nil || !deleted {
		t.Fatalf("DeleteSessionIfEmpty(empty) = %v, %v, want it deleted", deleted, err)
	}
	if _, err := storage.GetSession(ctx, empty.Id); err == nil {
		t.Fatal("GetSession() found the deleted session")
	}
	if got, err := storage.GetLatestSession(ctx, "5550003", "sf6"); err != nil || got == nil || got.Id != played.Id {
		t.Fatalf("GetLatestSession() = %+v, %v, want played session %d", got, err, played.Id)
	}
}

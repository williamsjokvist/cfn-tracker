package sql

import (
	"context"
	"path/filepath"
	"slices"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func userCodes(users []*model.User) []string {
	codes := make([]string, 0, len(users))
	for _, u := range users {
		codes = append(codes, u.Code)
	}
	slices.Sort(codes)
	return codes
}

func TestGetUsersByGame(t *testing.T) {
	ctx := context.Background()
	storage, err := NewStorage(true)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	for _, user := range []model.User{
		{Code: "1111", DisplayName: "Kakeru", Game: model.GameTypeSF6},
		{Code: "2222", DisplayName: "Punk", Game: model.GameTypeSF6},
		{Code: "4amNGJ6EdB6J", DisplayName: "Tekk!", Game: model.GameTypeT8},
	} {
		if err := storage.SaveUser(ctx, user); err != nil {
			t.Fatalf("SaveUser(%s) error = %v", user.Code, err)
		}
	}

	tests := []struct {
		game model.GameType
		want []string
	}{
		{model.GameTypeSF6, []string{"1111", "2222"}},
		{model.GameTypeT8, []string{"4amNGJ6EdB6J"}},
	}
	for _, tt := range tests {
		users, err := storage.GetUsers(ctx, tt.game)
		if err != nil {
			t.Fatalf("GetUsers(%s) error = %v", tt.game, err)
		}
		if got := userCodes(users); !slices.Equal(got, tt.want) {
			t.Fatalf("GetUsers(%s) = %v, want %v", tt.game, got, tt.want)
		}
	}
}

func TestSaveUserMovesExistingUserToGame(t *testing.T) {
	ctx := context.Background()
	storage, err := NewStorage(true)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "9876543210", DisplayName: "Tekk!", Game: model.GameTypeSF6}); err != nil {
		t.Fatalf("SaveUser() error = %v", err)
	}
	if err := storage.SaveUser(ctx, model.User{Code: "9876543210", DisplayName: "Tekk!", Game: model.GameTypeT8}); err != nil {
		t.Fatalf("SaveUser() again error = %v", err)
	}

	user, err := storage.GetUserByCode(ctx, "9876543210")
	if err != nil {
		t.Fatalf("GetUserByCode() error = %v", err)
	}
	if user.Game != model.GameTypeT8 {
		t.Fatalf("user game = %q, want %q", user.Game, model.GameTypeT8)
	}
}

func TestMigrationGivesExistingUsersAGame(t *testing.T) {
	dataSource := filepath.Join(t.TempDir(), "cfn-tracker.db")
	beforeGames := 3
	if err := migrateSchemaAt(dataSource, &beforeGames); err != nil {
		t.Fatalf("migrate to before games: %v", err)
	}
	db, err := sqlx.Open("sqlite", dataSource)
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { _ = db.Close() })
	if _, err := db.Exec(`INSERT INTO users (code, display_name) VALUES ('3141592653', 'Kakeru'), ('4amNGJ6EdB6J', 'Tekk!')`); err != nil {
		t.Fatalf("insert users: %v", err)
	}

	if err := migrateSchemaAt(dataSource, nil); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	games := map[string]string{}
	rows, err := db.Query(`SELECT code, game FROM users`)
	if err != nil {
		t.Fatalf("select users: %v", err)
	}
	defer func() { _ = rows.Close() }()
	for rows.Next() {
		var code, game string
		if err := rows.Scan(&code, &game); err != nil {
			t.Fatalf("scan user: %v", err)
		}
		games[code] = game
	}
	if games["3141592653"] != "sf6" || games["4amNGJ6EdB6J"] != "t8" {
		t.Fatalf("games after migration = %v, want the CFN code on sf6 and the Polaris ID on t8", games)
	}
}

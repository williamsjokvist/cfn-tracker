package cmd_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
)

func TestExportUserMatches(t *testing.T) {
	ctx := context.Background()
	sqlDb, err := sql.NewStorage(true)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	for _, u := range []model.User{
		{Code: "123456789", DisplayName: "Kakeru", Game: model.GameTypeSF6},
		{Code: "987654321", DisplayName: "Other", Game: model.GameTypeSF6},
	} {
		if err := sqlDb.SaveUser(ctx, u); err != nil {
			t.Fatalf("SaveUser() error = %v", err)
		}
	}
	session, err := sqlDb.CreateSession(ctx, "123456789")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	other, err := sqlDb.CreateSession(ctx, "987654321")
	if err != nil {
		t.Fatalf("CreateSession() error = %v", err)
	}
	for _, m := range []model.Match{
		{UserId: "123456789", SessionId: session.Id, Character: "JP", Opponent: "Daigo", OpponentCharacter: "Ken", OpponentLeague: "Master", Victory: false, LP: 25000, LPGain: -20, MR: 1600, MRGain: -8, OpponentLP: 25000, OpponentMR: 1620, Date: "2026-01-01", Time: "10:05", ReplayID: "B"},
		{UserId: "123456789", SessionId: session.Id, Character: "JP", Opponent: "Tokido", OpponentCharacter: "Akuma", OpponentLeague: "Master", Victory: true, LP: 25000, LPGain: 0, MR: 1608, MRGain: 8, OpponentLP: 25000, OpponentMR: 1590, WinStreak: 1, Date: "2026-01-01", Time: "10:00", ReplayID: "A"},
		{UserId: "987654321", SessionId: other.Id, Character: "Ken", Victory: true, Date: "2026-01-01", Time: "11:00"},
	} {
		if err := sqlDb.SaveMatch(ctx, m); err != nil {
			t.Fatalf("SaveMatch() error = %v", err)
		}
	}

	path := filepath.Join(t.TempDir(), "matches.csv")
	var suggested string
	exporter := cmd.NewExporter()
	exporter.SetSaveFileDialog(func(defaultFilename string) (string, error) {
		suggested = defaultFilename
		return path, nil
	})
	handler := cmd.NewCommandHandler(sqlDb, nil, nil, &config.BuildConfig{}, nil, exporter)

	if err := handler.ExportUserMatches("123456789"); err != nil {
		t.Fatalf("ExportUserMatches() error = %v", err)
	}
	if suggested != "Kakeru-matches.csv" {
		t.Errorf("suggested filename = %q, want Kakeru-matches.csv", suggested)
	}

	exported, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read exported csv: %v", err)
	}
	got, err := csv.NewReader(bytes.NewReader(exported)).ReadAll()
	if err != nil {
		t.Fatalf("read exported csv: %v", err)
	}
	want := [][]string{
		{"session_id", "session_started", "date", "time", "character", "opponent", "opponent_character", "opponent_league", "victory", "lp", "lp_gain", "mr", "mr_gain", "opponent_lp", "opponent_mr", "win_streak", "replay_id"},
		{strconv.Itoa(int(session.Id)), session.CreatedAt, "2026-01-01", "10:00", "JP", "Tokido", "Akuma", "Master", "true", "25000", "0", "1608", "8", "25000", "1590", "1", "A"},
		{strconv.Itoa(int(session.Id)), session.CreatedAt, "2026-01-01", "10:05", "JP", "Daigo", "Ken", "Master", "false", "25000", "-20", "1600", "-8", "25000", "1620", "0", "B"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("exported csv = %v, want %v", got, want)
	}
}

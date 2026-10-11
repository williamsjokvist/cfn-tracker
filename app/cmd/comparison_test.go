package cmd_test

import (
	"context"
	"encoding/csv"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestExportSF6BattleStatsComparison(t *testing.T) {
	path := filepath.Join(t.TempDir(), "comparison.csv")
	var suggested string
	exporter := cmd.NewExporter()
	exporter.SetSaveFileDialog(func(defaultFilename string) (string, error) {
		suggested = defaultFilename
		return path, nil
	})
	handler := cmd.NewSF6ComparisonHandler(nil, exporter)

	report := model.SF6BattleStatsComparisonReport{
		UserName:      "Kakeru",
		CharacterName: "JP",
		Metrics: []model.SF6BattleStatsMetric{
			{Key: "rankMatchPlayCount", Unit: "matches", Current: 120, TopAvg: 100, TopMin: 80, TopMax: 130},
			{Key: "driveParry", Unit: "%", Current: 10.5, TopAvg: 12.25, TopMin: 9, TopMax: 15},
		},
	}
	if err := handler.ExportSF6BattleStatsComparison(report); err != nil {
		t.Fatalf("ExportSF6BattleStatsComparison() error = %v", err)
	}
	if suggested != "Kakeru-comparison.csv" {
		t.Errorf("suggested filename = %q, want Kakeru-comparison.csv", suggested)
	}

	exported, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read exported csv: %v", err)
	}
	got, err := csv.NewReader(strings.NewReader(string(exported))).ReadAll()
	if err != nil {
		t.Fatalf("read exported csv: %v", err)
	}
	want := [][]string{
		{"metric", "unit", "Kakeru", "top avg", "top min", "top max", "delta"},
		{"rankMatchPlayCount", "matches", "120", "100", "80", "130", "20"},
		{"driveParry", "%", "10.5", "12.25", "9", "15", "-1.75"},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("exported csv = %v, want %v", got, want)
	}
}

func TestExportCancelled(t *testing.T) {
	exporter := cmd.NewExporter()
	exporter.SetSaveFileDialog(func(string) (string, error) { return "", nil })
	handler := cmd.NewSF6ComparisonHandler(nil, exporter)
	if err := handler.ExportSF6BattleStatsComparison(model.SF6BattleStatsComparisonReport{}); err != nil {
		t.Fatalf("ExportSF6BattleStatsComparison() after cancelling = %v, want no error", err)
	}
}

func TestGetSF6BattleStatsComparisonCancelledBeforeStart(t *testing.T) {
	handler := cmd.NewSF6ComparisonHandler(nil, cmd.NewExporter())
	handler.CancelSF6BattleStatsComparison("run")

	if _, err := handler.GetSF6BattleStatsComparison("123456789", "run"); !errors.Is(err, context.Canceled) {
		t.Fatalf("GetSF6BattleStatsComparison() error = %v, want context.Canceled", err)
	}
	// The early cancel is used up, so the same id runs again. Without a browser
	// it fails at opening a tab instead.
	if _, err := handler.GetSF6BattleStatsComparison("123456789", "run"); errors.Is(err, context.Canceled) {
		t.Fatalf("GetSF6BattleStatsComparison() error = %v, want it to not be cancelled", err)
	}
}

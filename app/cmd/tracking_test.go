package cmd_test

import (
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

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

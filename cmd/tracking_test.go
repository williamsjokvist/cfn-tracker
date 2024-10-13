package cmd_test

import (
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestTrackingSelectGame(t *testing.T) {
	tEnv.trackingHandler.SelectGame(model.GameTypeT8)

	tests := []struct {
		name     string
		gameType model.GameType
		setup    func(gameType model.GameType)
		expErr   error
	}{}

	for _, tt := range tests {
		tt.setup(tt.gameType)
	}
}

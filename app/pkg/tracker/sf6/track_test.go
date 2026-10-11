package sf6

import (
	"context"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn/cfntest"
)

func TestPoll(t *testing.T) {
	tests := []struct {
		name      string
		match     cfntest.Match
		wantMatch bool
	}{
		{"win as player 1", cfntest.Match{Won: true, Opponent: "Punk"}, true},
		{"loss as player 2", cfntest.Match{Opponent: "Punk", AsP2: true}, true},
		{"replay older than 15 minutes", cfntest.Match{Won: true, Opponent: "Punk", At: time.Now().Add(-16 * time.Minute)}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := cfntest.NewClient("Kakeru", "Ryu", 10000, 1500)
			replay := client.Play(tt.match)
			match, err := NewSF6Tracker(client).Poll(context.Background(), &model.Session{UserId: "123456789"})
			if err != nil {
				t.Fatalf("Poll() error = %v", err)
			}
			if (match != nil) != tt.wantMatch {
				t.Fatalf("Poll() match = %+v, want match: %v", match, tt.wantMatch)
			}
			if match != nil && (match.Victory != tt.match.Won || match.Opponent != "Punk" || match.ReplayID != replay.ReplayID) {
				t.Fatalf("Poll() match = %+v, want victory %v against Punk in %s", match, tt.match.Won, replay.ReplayID)
			}
		})
	}
}

func TestPollMapsOpponent(t *testing.T) {
	client := cfntest.NewClient("Kakeru", "Ryu", 10000, 1500)
	client.Play(cfntest.Match{Won: true, Opponent: "Punk", OpponentCharacter: "Cammy", OpponentLP: 25000, OpponentMR: 1800, LPChange: 50, MRChange: 12, AsP2: true})

	match, err := NewSF6Tracker(client).Poll(context.Background(), &model.Session{UserId: "123456789", UserName: "Kakeru"})
	if err != nil || match == nil {
		t.Fatalf("Poll() = %+v, %v, want a match", match, err)
	}
	if match.Character != "Ryu" || match.LP != 10050 || match.MR != 1512 || match.UserName != "Kakeru" || match.UserId != "123456789" {
		t.Fatalf("Poll() player = %+v, want Kakeru on Ryu at 10050 LP and 1512 MR", match)
	}
	if match.OpponentCharacter != "Cammy" || match.OpponentLP != 25000 || match.OpponentMR != 1800 || match.OpponentLeague != "Master" {
		t.Fatalf("Poll() opponent = %+v, want Master Cammy", match)
	}
}

func TestSession(t *testing.T) {
	client := cfntest.NewClient("Kakeru", "Ryu", 10000, 0)
	client.SetPoints("Ken", 5000, 0)
	session := tracker.NewSession(NewSF6Tracker(client), &model.Session{UserId: "123456789", LP: 10000})
	tick := func() *model.Match {
		t.Helper()
		match, err := session.Tick(context.Background())
		if err != nil {
			t.Fatalf("Tick() error = %v", err)
		}
		return match
	}

	client.Play(cfntest.Match{Won: true, Opponent: "Punk", LPChange: 100})
	if m := tick(); m == nil || m.Wins != 1 || m.WinStreak != 1 || m.LPGain != 100 {
		t.Fatalf("first match = %+v, want a win gaining 100 LP", m)
	}
	if m := tick(); m != nil {
		t.Fatalf("polling again recorded %+v, want nothing new", m)
	}

	client.Play(cfntest.Match{Opponent: "Punk", LPChange: -80, AsP2: true})
	if m := tick(); m == nil || m.Wins != 1 || m.Losses != 1 || m.WinStreak != 0 || m.LPGain != 20 {
		t.Fatalf("second match = %+v, want a loss down to +20 LP", m)
	}

	client.Play(cfntest.Match{Won: true, Opponent: "Punk", Character: "Ken", LPChange: 30})
	if m := tick(); m == nil || m.Character != "Ken" || m.Wins != 1 || m.Losses != 0 || m.LPGain != 0 {
		t.Fatalf("Ken match = %+v, want Ken's first win with no LP baseline", m)
	}

	client.Play(cfntest.Match{Won: true, Opponent: "Punk", Character: "Ryu", LPChange: 50})
	if m := tick(); m == nil || m.Wins != 2 || m.Losses != 1 || m.LPGain != 70 {
		t.Fatalf("back on Ryu = %+v, want Ryu's stats to continue at +70 LP", m)
	}
}

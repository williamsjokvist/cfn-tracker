package tracker_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

type stats struct {
	Wins, Losses, WinStreak, WinRate, LPGain, MRGain int
}

func statsOf(m model.Match) stats {
	return stats{m.Wins, m.Losses, m.WinStreak, m.WinRate, m.LPGain, m.MRGain}
}

func win(character string, lp, mr int) model.Match {
	return model.Match{Character: character, Victory: true, LP: lp, MR: mr}
}

func loss(character string, lp, mr int) model.Match {
	return model.Match{Character: character, LP: lp, MR: mr}
}

func TestRecord(t *testing.T) {
	type step struct {
		outcome model.Match
		want    stats
	}
	tests := []struct {
		name    string
		session model.Session
		steps   []step
	}{
		{
			name:    "counts wins, losses and the streak",
			session: model.Session{LP: 10000, MR: 1500},
			steps: []step{
				{win("Ryu", 10100, 1510), stats{1, 0, 1, 100, 100, 10}},
				{win("Ryu", 10200, 1520), stats{2, 0, 2, 100, 200, 20}},
				{loss("Ryu", 10150, 1510), stats{2, 1, 0, 66, 150, 10}},
				{win("Ryu", 10250, 1520), stats{3, 1, 1, 75, 250, 20}},
			},
		},
		{
			name:    "counts each character separately",
			session: model.Session{LP: 10000},
			steps: []step{
				{win("Ryu", 10100, 0), stats{1, 0, 1, 100, 100, 0}},
				{win("Ryu", 10200, 0), stats{2, 0, 2, 100, 200, 0}},
				{loss("Ken", 5000, 0), stats{0, 1, 0, 0, 0, 0}},
				{win("Ken", 5100, 0), stats{1, 1, 1, 50, 100, 0}},
				{win("Ryu", 10300, 0), stats{3, 0, 3, 100, 300, 0}},
			},
		},
		{
			name:    "gains nothing from placement",
			session: model.Session{LP: -1},
			steps: []step{
				{win("Ryu", -1, 0), stats{1, 0, 1, 100, 0, 0}},
				{win("Ryu", 1200, 0), stats{2, 0, 2, 100, 0, 0}},
				{win("Ryu", 1300, 0), stats{3, 0, 3, 100, 100, 0}},
			},
		},
		{
			name: "continues a restored session",
			session: model.Session{LP: 10500, Matches: []*model.Match{
				{Character: "Ryu", LP: 10500, Wins: 4, Losses: 2, WinStreak: 3, LPGain: 500},
			}},
			steps: []step{
				{win("Ryu", 10600, 0), stats{5, 2, 4, 71, 600, 0}},
			},
		},
		{
			name:    "has no points to gain in Tekken",
			session: model.Session{},
			steps: []step{
				{win("Steve", 0, 0), stats{1, 0, 1, 100, 0, 0}},
				{loss("Steve", 0, 0), stats{1, 1, 0, 50, 0, 0}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			state := tt.session
			state.Id = 7
			session := tracker.NewSession(nil, &state)
			for i, step := range tt.steps {
				match := session.Record(step.outcome)
				if got := statsOf(match); got != step.want {
					t.Fatalf("match %d: stats = %+v, want %+v", i+1, got, step.want)
				}
				if match.SessionId != 7 {
					t.Fatalf("match %d: session id = %d, want 7", i+1, match.SessionId)
				}
				if *state.Matches[0] != match || state.LP != match.LP || state.MR != match.MR {
					t.Fatalf("match %d: session state not updated: %+v", i+1, state)
				}
			}
		})
	}
}

type scriptedGame struct {
	outcomes []*model.Match
	err      error
}

func (g *scriptedGame) GetUser(context.Context, string) (*model.User, error)  { return nil, nil }
func (g *scriptedGame) Authenticate(context.Context, chan tracker.AuthStatus) {}
func (g *scriptedGame) Poll(context.Context, *model.Session) (*model.Match, error) {
	if g.err != nil {
		return nil, g.err
	}
	outcome := g.outcomes[0]
	g.outcomes = g.outcomes[1:]
	return outcome, nil
}

func TestTick(t *testing.T) {
	game := &scriptedGame{outcomes: []*model.Match{nil, {Character: "Ryu", Victory: true}}}
	state := &model.Session{}
	session := tracker.NewSession(game, state)

	if match, err := session.Tick(context.Background()); match != nil || err != nil {
		t.Fatalf("Tick() = %+v, %v, want nothing new", match, err)
	}
	if len(state.Matches) != 0 {
		t.Fatalf("session has %d matches after an empty poll, want 0", len(state.Matches))
	}
	match, err := session.Tick(context.Background())
	if err != nil || match == nil || match.Wins != 1 {
		t.Fatalf("Tick() = %+v, %v, want the first win", match, err)
	}
	if len(state.Matches) != 1 {
		t.Fatalf("session has %d matches, want 1", len(state.Matches))
	}

	game.err = errors.New("offline")
	if _, err := session.Tick(context.Background()); !errors.Is(err, game.err) {
		t.Fatalf("Tick() error = %v, want %v", err, game.err)
	}
}

func TestCharacterSwitch(t *testing.T) {
	state := &model.Session{LP: 10000}
	session := tracker.NewSession(nil, state)
	characterStats := func(character string) stats {
		t.Helper()
		match, ok := session.Character(character)
		if !ok {
			t.Fatalf("no stats for %s", character)
		}
		return statsOf(match)
	}

	session.Record(win("Ryu", 10100, 0))
	session.Record(win("Ryu", 10200, 0))
	session.Record(loss("Ryu", 10150, 0))
	ryu := stats{2, 1, 0, 66, 150, 0}
	if got := characterStats("Ryu"); got != ryu {
		t.Fatalf("Ryu = %+v, want %+v", got, ryu)
	}

	if got := statsOf(session.Record(win("Ken", 5000, 0))); got != (stats{1, 0, 1, 100, 0, 0}) {
		t.Fatalf("Ken's first match = %+v, want a fresh start", got)
	}
	session.Record(win("Ken", 5100, 0))
	ken := stats{2, 0, 2, 100, 100, 0}
	if got := characterStats("Ken"); got != ken {
		t.Fatalf("Ken = %+v, want %+v", got, ken)
	}
	if got := characterStats("Ryu"); got != ryu {
		t.Fatalf("Ryu after switching = %+v, want %+v", got, ryu)
	}

	if got := statsOf(session.Record(win("Ryu", 10250, 0))); got != (stats{3, 1, 1, 75, 250, 0}) {
		t.Fatalf("back on Ryu = %+v, want his stats to continue", got)
	}
	if got := characterStats("Ken"); got != ken {
		t.Fatalf("Ken after switching back = %+v, want %+v", got, ken)
	}

	var played []string
	for _, match := range state.Matches {
		played = append(played, match.Character)
	}
	if want := "Ryu Ken Ken Ryu Ryu Ryu"; fmt.Sprint(played) != "["+want+"]" {
		t.Fatalf("session matches = %v, want [%s]", played, want)
	}
}

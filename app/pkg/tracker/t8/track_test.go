package t8

import (
	"context"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu/wavutest"
)

type fakeWavuClient struct {
	replay wavu.Replay
}

func (f fakeWavuClient) GetLastReplay(context.Context, string) (wavu.Replay, error) {
	return f.replay, nil
}

func (f fakeWavuClient) GetUserName(context.Context, string) (string, error) {
	return "", nil
}

func replayAt(at time.Time, opponent string) wavu.Replay {
	return wavu.Replay{
		BattleAt:    at.Unix(),
		P1Name:      "Tekk!",
		P1PolarisId: "4amNGJ6EdB6J",
		P1CharaId:   9,
		P2Name:      opponent,
		P2PolarisId: opponent + "-id",
		P2CharaId:   8,
		P2Rank:      33,
		Winner:      1,
	}
}

func recorded(replay wavu.Replay) *model.Match {
	at := time.Unix(replay.BattleAt, 0)
	return &model.Match{
		ReplayID: replay.ID(),
		Date:     at.Format("2006-01-02"),
		Time:     at.Format("15:04"),
		Wins:     1,
	}
}

func TestPoll(t *testing.T) {
	// Whole minutes, so replays can be placed in the same minute or the next.
	minute := time.Now().Add(-5 * time.Minute).Truncate(time.Minute)
	prev := replayAt(minute.Add(5*time.Second), "bilol")

	tests := []struct {
		name      string
		replay    wavu.Replay
		prev      *model.Match
		wantMatch bool
	}{
		{"no replay for the player", wavu.Replay{}, nil, false},
		{"first replay of the session", replayAt(minute, "bilol"), nil, true},
		{"replay older than 10 minutes", replayAt(time.Now().Add(-11*time.Minute), "bilol"), nil, false},
		{"replay just under 10 minutes old", replayAt(time.Now().Add(-9*time.Minute), "bilol"), nil, true},
		{"same replay as the last match", prev, recorded(prev), false},
		{"replay from before the last match", replayAt(minute.Add(-2*time.Minute), "lowhigh"), recorded(prev), false},
		{"new replay in the same minute as the last match", replayAt(minute.Add(50*time.Second), "lowhigh"), recorded(prev), true},
		{"new replay a minute after the last match", replayAt(minute.Add(time.Minute), "lowhigh"), recorded(prev), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			session := &model.Session{UserId: "4amNGJ6EdB6J"}
			if tt.prev != nil {
				session.Matches = []*model.Match{tt.prev}
			}
			match, err := NewT8Tracker(fakeWavuClient{tt.replay}).Poll(context.Background(), session)
			if err != nil {
				t.Fatalf("Poll() error = %v", err)
			}
			if (match != nil) != tt.wantMatch {
				t.Fatalf("Poll() match = %+v, want match: %v", match, tt.wantMatch)
			}
		})
	}
}

func TestPollMapsReplayFromPlayerSide(t *testing.T) {
	client := wavutest.NewClient("4amNGJ6EdB6J", "Tekk!")
	replay := client.Play(wavutest.Match{Won: true, Opponent: "bilol", Character: 9, OpponentCharacter: 8, OpponentRank: 33, AsP2: true})

	match, err := NewT8Tracker(client).Poll(context.Background(), &model.Session{UserId: "4amNGJ6EdB6J"})
	if err != nil || match == nil {
		t.Fatalf("Poll() = %+v, %v, want a match", match, err)
	}

	at := time.Unix(replay.BattleAt, 0)
	want := model.Match{
		UserName:          "Tekk!",
		UserId:            "4amNGJ6EdB6J",
		Opponent:          "bilol",
		Victory:           true,
		ReplayID:          replay.ID(),
		Character:         wavu.ConvCharaIdToName(9),
		OpponentCharacter: wavu.ConvCharaIdToName(8),
		OpponentLeague:    wavu.ConvRankToName(33),
		Date:              at.Format("2006-01-02"),
		Time:              at.Format("15:04"),
	}
	if *match != want {
		t.Fatalf("Poll() match = %+v\nwant %+v", *match, want)
	}
}

func TestSession(t *testing.T) {
	client := wavutest.NewClient("4amNGJ6EdB6J", "Tekk!")
	session := tracker.NewSession(NewT8Tracker(client), &model.Session{UserId: "4amNGJ6EdB6J"})
	tick := func() *model.Match {
		t.Helper()
		match, err := session.Tick(context.Background())
		if err != nil {
			t.Fatalf("Tick() error = %v", err)
		}
		return match
	}

	client.Play(wavutest.Match{Won: true, Opponent: "bilol", Character: 9})
	client.Play(wavutest.Match{Won: true, Opponent: "lowhigh", Character: 9, AsP2: true})
	// Only the latest replay is served, so the first match is never seen.
	if m := tick(); m == nil || m.Opponent != "lowhigh" || m.Wins != 1 || m.WinStreak != 1 {
		t.Fatalf("first match = %+v, want a win against lowhigh", m)
	}
	if m := tick(); m != nil {
		t.Fatalf("polling again recorded %+v, want nothing new", m)
	}

	client.Play(wavutest.Match{Opponent: "bilol", Character: 9})
	if m := tick(); m == nil || m.Wins != 1 || m.Losses != 1 || m.WinStreak != 0 || m.WinRate != 50 {
		t.Fatalf("second match = %+v, want a loss resetting the streak", m)
	}

	client.Play(wavutest.Match{Won: true, Opponent: "bilol", Character: 8})
	if m := tick(); m == nil || m.Character != "Kazuya" || m.Wins != 1 || m.Losses != 0 {
		t.Fatalf("Kazuya match = %+v, want Kazuya's first win", m)
	}
	if got := len(session.State().Matches); got != 3 {
		t.Fatalf("session has %d matches, want 3", got)
	}
}

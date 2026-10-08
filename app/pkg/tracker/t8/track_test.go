package t8

import (
	"context"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
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
	at := time.Now().Add(-time.Minute)
	replay := replayAt(at, "bilol")
	replay.P1PolarisId, replay.P2PolarisId = replay.P2PolarisId, replay.P1PolarisId
	replay.P1Name, replay.P2Name = replay.P2Name, replay.P1Name
	replay.P1CharaId, replay.P2CharaId = replay.P2CharaId, replay.P1CharaId
	replay.P1Rank = 33
	replay.Winner = 2

	session := &model.Session{Id: 7, UserId: "4amNGJ6EdB6J", Matches: []*model.Match{{Wins: 2, Losses: 1, WinStreak: 1}}}
	match, err := NewT8Tracker(fakeWavuClient{replay}).Poll(context.Background(), session)
	if err != nil || match == nil {
		t.Fatalf("Poll() = %+v, %v, want a match", match, err)
	}

	want := model.Match{
		SessionId:         7,
		UserName:          "Tekk!",
		UserId:            "4amNGJ6EdB6J",
		Opponent:          "bilol",
		Victory:           true,
		ReplayID:          replay.ID(),
		Wins:              3,
		Losses:            1,
		WinStreak:         2,
		WinRate:           75,
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

func TestPollResetsStreakOnLoss(t *testing.T) {
	replay := replayAt(time.Now().Add(-time.Minute), "bilol")
	replay.Winner = 2
	session := &model.Session{UserId: "4amNGJ6EdB6J", Matches: []*model.Match{{Wins: 2, WinStreak: 2}}}

	match, err := NewT8Tracker(fakeWavuClient{replay}).Poll(context.Background(), session)
	if err != nil || match == nil {
		t.Fatalf("Poll() = %+v, %v, want a match", match, err)
	}
	if match.Victory || match.Losses != 1 || match.WinStreak != 0 || match.WinRate != 66 {
		t.Fatalf("Poll() match = %+v, want a loss with the streak reset", match)
	}
}

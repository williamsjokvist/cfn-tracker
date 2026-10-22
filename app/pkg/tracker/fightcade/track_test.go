package fightcade_test

import (
	"context"
	"errors"
	"slices"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/fightcade"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/fightcade/fcapi"
)

type fakeClient struct {
	user    fcapi.User
	replays []fcapi.Replay
	games   map[string]string
}

func (c *fakeClient) Open(context.Context, func(string)) error { return nil }
func (c *fakeClient) Close() error                             { return nil }

func (c *fakeClient) GetUser(context.Context, string) (*fcapi.User, error) {
	return &c.user, nil
}

func (c *fakeClient) GetReplays(_ context.Context, _ string, limit int) ([]fcapi.Replay, error) {
	return c.replays[:min(limit, len(c.replays))], nil
}

func (c *fakeClient) GetGame(_ context.Context, gameId string) (*fcapi.Game, error) {
	name, ok := c.games[gameId]
	if !ok {
		return nil, errors.New("no such game")
	}
	return &fcapi.Game{GameId: gameId, Name: name}, nil
}

func ptr[T any](v T) *T { return &v }

var started = time.Date(2026, 10, 12, 20, 0, 0, 0, time.Local)

func set(quarkId, gameId string, at time.Time, me, them int) fcapi.Replay {
	return fcapi.Replay{
		QuarkId:  quarkId,
		GameId:   gameId,
		Date:     fcapi.Timestamp(at.UnixMilli()),
		Duration: 600,
		Ranked:   []byte("1"),
		Players: []fcapi.Player{
			{Name: "naggamurasang", Score: ptr(me), Rank: ptr(fcapi.Rank(4))},
			{Name: "Rival", Score: ptr(them), Rank: ptr(fcapi.Rank(6))},
		},
	}
}

func session(matches ...*model.Match) *model.Session {
	return &model.Session{
		UserId:    "naggamurasang",
		UserName:  "naggamurasang",
		GameId:    "sfiii3nr1",
		CreatedAt: started.Format(time.RFC3339),
		Matches:   matches,
	}
}

func TestPollRecordsLatestSetOfSessionGame(t *testing.T) {
	cancelled := set("cancelled", "sfiii3nr1", started.Add(40*time.Minute), 3, 0)
	cancelled.Ranked = []byte(`"cancelled"`)
	client := &fakeClient{replays: []fcapi.Replay{
		set("other-game", "kof98", started.Add(50*time.Minute), 3, 1),
		cancelled,
		set("latest", "sfiii3nr1", started.Add(30*time.Minute), 3, 1),
		set("earlier", "sfiii3nr1", started.Add(10*time.Minute), 0, 3),
		set("before-tracking", "sfiii3nr1", started.Add(-time.Hour), 3, 0),
	}}

	match, err := fightcade.NewFightcadeTracker(client).Poll(context.Background(), session())
	if err != nil {
		t.Fatalf("Poll() error = %v", err)
	}
	if match == nil || match.ReplayID != "latest" || !match.Victory || match.Opponent != "Rival" || match.OpponentLeague != "S" {
		t.Fatalf("Poll() = %+v, want a win over Rival (S) in the latest sfiii3nr1 set", match)
	}
	at := started.Add(30 * time.Minute)
	if match.Date != at.Format("2006-01-02") || match.Time != at.Format("15:04") {
		t.Fatalf("Poll() played at %s %s, want %s", match.Date, match.Time, at.Format("2006-01-02 15:04"))
	}
}

func TestPollCountsSetThatEndedAfterTrackingStarted(t *testing.T) {
	client := &fakeClient{replays: []fcapi.Replay{set("ongoing", "sfiii3nr1", started.Add(-5*time.Minute), 1, 3)}}

	match, err := fightcade.NewFightcadeTracker(client).Poll(context.Background(), session())
	if err != nil || match == nil || match.Victory {
		t.Fatalf("Poll() = %+v, %v, want the loss in the set that was still going when tracking started", match, err)
	}
}

func TestPollFindsUserAsSecondPlayer(t *testing.T) {
	replay := set("p2", "sfiii3nr1", started.Add(time.Minute), 3, 2)
	replay.Players[0], replay.Players[1] = replay.Players[1], replay.Players[0]
	replay.Players[1].Name = "NaggaMurasang"
	client := &fakeClient{replays: []fcapi.Replay{replay}}

	match, err := fightcade.NewFightcadeTracker(client).Poll(context.Background(), session())
	if err != nil || match == nil || !match.Victory || match.Opponent != "Rival" {
		t.Fatalf("Poll() = %+v, %v, want a win over Rival", match, err)
	}
}

func TestPollSkipsRecordedAndDrawnSets(t *testing.T) {
	tests := map[string]struct {
		replay  fcapi.Replay
		session *model.Session
	}{
		"recorded": {set("done", "sfiii3nr1", started.Add(time.Minute), 3, 0), session(&model.Match{ReplayID: "done"})},
		"drawn":    {set("draw", "sfiii3nr1", started.Add(time.Minute), 2, 2), session()},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			client := &fakeClient{replays: []fcapi.Replay{tt.replay}}
			match, err := fightcade.NewFightcadeTracker(client).Poll(context.Background(), tt.session)
			if err != nil || match != nil {
				t.Fatalf("Poll() = %+v, %v, want nothing to record", match, err)
			}
		})
	}
}

func TestGetPlayedGamesListsMostRecentFirst(t *testing.T) {
	client := &fakeClient{
		user: fcapi.User{Name: "naggamurasang", GameInfo: map[string]fcapi.GameInfo{
			"kof98":     {NumMatches: 40, LastMatch: 100},
			"sfiii3nr1": {NumMatches: 900, LastMatch: 300, Rank: ptr(fcapi.Rank(5))},
			"garou":     {NumMatches: 3, LastMatch: 200},
		}},
		games: map[string]string{"sfiii3nr1": "Street Fighter III: 3rd Strike", "kof98": "The King of Fighters '98"},
	}

	games, err := fightcade.NewFightcadeTracker(client).GetPlayedGames(context.Background(), "naggamurasang")
	if err != nil {
		t.Fatalf("GetPlayedGames() error = %v", err)
	}
	want := []model.PlayedGame{
		{Id: "sfiii3nr1", Name: "Street Fighter III: 3rd Strike", Matches: 900, Rank: "A"},
		{Id: "garou", Name: "garou", Matches: 3},
		{Id: "kof98", Name: "The King of Fighters '98", Matches: 40},
	}
	if !slices.Equal(games, want) {
		t.Fatalf("GetPlayedGames() = %+v, want %+v", games, want)
	}
}

func TestGetUserUsesFightcadeSpelling(t *testing.T) {
	client := &fakeClient{user: fcapi.User{Name: "naggamurasang"}}
	user, err := fightcade.NewFightcadeTracker(client).GetUser(context.Background(), " NaggaMurasang ")
	if err != nil || user.Code != "naggamurasang" || user.Game != model.GameTypeFightcade {
		t.Fatalf("GetUser() = %+v, %v, want the Fightcade user naggamurasang", user, err)
	}

	_, err = fightcade.NewFightcadeTracker(&fakeClient{}).GetUser(context.Background(), "nobody")
	if err == nil {
		t.Fatal("GetUser() of an unknown user returned no error")
	}
}

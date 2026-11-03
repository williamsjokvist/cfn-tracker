package fightcade

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/fightcade/fcapi"
)

const (
	// Enough to find the session's game among sets of other games played meanwhile.
	replaysPerPoll = 10
	// Each game's name is its own request, so only the most recent ones are listed.
	maxListedGames = 20
)

type FightcadeTracker struct {
	client fcapi.FightcadeClient
}

var (
	_ tracker.GameTracker = (*FightcadeTracker)(nil)
	_ tracker.GameLister  = (*FightcadeTracker)(nil)
)

func NewFightcadeTracker(client fcapi.FightcadeClient) *FightcadeTracker {
	return &FightcadeTracker{client}
}

func (t *FightcadeTracker) GetUser(ctx context.Context, username string) (*model.User, error) {
	user, err := t.client.GetUser(ctx, strings.TrimSpace(username))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	if user.Name == "" {
		return nil, errors.New("fightcade user not found")
	}
	// Replays name players as Fightcade spells them, whatever case was typed.
	return &model.User{DisplayName: user.Name, Code: user.Name, Game: model.GameTypeFightcade}, nil
}

func (t *FightcadeTracker) GetPlayedGames(ctx context.Context, username string) ([]model.PlayedGame, error) {
	user, err := t.client.GetUser(ctx, strings.TrimSpace(username))
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	ids := make([]string, 0, len(user.GameInfo))
	for id := range user.GameInfo {
		ids = append(ids, id)
	}
	slices.SortFunc(ids, func(a, b string) int {
		ga, gb := user.GameInfo[a], user.GameInfo[b]
		if c := cmp.Compare(gb.LastMatch, ga.LastMatch); c != 0 {
			return c
		}
		return cmp.Compare(gb.NumMatches, ga.NumMatches)
	})
	if len(ids) > maxListedGames {
		ids = ids[:maxListedGames]
	}
	games := make([]model.PlayedGame, 0, len(ids))
	for _, id := range ids {
		info := user.GameInfo[id]
		game := model.PlayedGame{Id: id, Name: id, Matches: info.NumMatches}
		if info.Rank != nil {
			game.Rank = info.Rank.String()
		}
		// The ID still tells the games apart, so a missing name isn't worth failing over.
		if details, err := t.client.GetGame(ctx, id); err == nil && details.Name != "" {
			game.Name = details.Name
		} else if err != nil {
			slog.Warn("fightcade: get game name", slog.String("game_id", id), slog.Any("error", err))
		}
		games = append(games, game)
	}
	return games, nil
}

func (t *FightcadeTracker) Poll(ctx context.Context, session *model.Session) (*model.Match, error) {
	replays, err := t.client.GetReplays(ctx, session.UserId, replaysPerPoll)
	if err != nil {
		return nil, fmt.Errorf("fightcade: get replays: %w", err)
	}
	startedAt, err := time.Parse(time.RFC3339, session.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("parse session start: %w", err)
	}
	replay, ok := latestReplay(replays, session.GameId, startedAt)
	if !ok || recorded(session, replay.QuarkId) {
		return nil, nil
	}
	me, opponent, ok := players(replay, session.UserId)
	if !ok {
		slog.Info("fightcade: replay has no opponent or score", slog.String("quark_id", replay.QuarkId))
		return nil, nil
	}
	playedAt := replay.Date.Time()
	match := &model.Match{
		UserId:   session.UserId,
		UserName: session.UserName,
		Opponent: opponent.Name,
		Victory:  *me.Score > *opponent.Score,
		ReplayID: replay.QuarkId,
		Date:     playedAt.Format("2006-01-02"),
		Time:     playedAt.Format("15:04"),
	}
	if opponent.Rank != nil {
		match.OpponentLeague = opponent.Rank.String()
	}
	return match, nil
}

// A replay is only listed once its sitting is over, so a sitting counts if it
// ended after tracking started, even if it began before.
func latestReplay(replays []fcapi.Replay, gameId string, startedAt time.Time) (fcapi.Replay, bool) {
	var latest fcapi.Replay
	found := false
	for _, r := range replays {
		if r.GameId != gameId || r.Cancelled() {
			continue
		}
		endedAt := r.Date.Time().Add(time.Duration(r.Duration * float64(time.Second)))
		if endedAt.Before(startedAt) {
			continue
		}
		if !found || r.Date > latest.Date {
			latest, found = r, true
		}
	}
	return latest, found
}

func recorded(session *model.Session, quarkId string) bool {
	return slices.ContainsFunc(session.Matches, func(m *model.Match) bool { return m.ReplayID == quarkId })
}

// A drawn set has no winner to record, so it's treated like one without scores.
func players(replay fcapi.Replay, username string) (me, opponent fcapi.Player, ok bool) {
	if len(replay.Players) != 2 {
		return me, opponent, false
	}
	me, opponent = replay.Players[0], replay.Players[1]
	if strings.EqualFold(opponent.Name, username) {
		me, opponent = opponent, me
	}
	if !strings.EqualFold(me.Name, username) || me.Score == nil || opponent.Score == nil || *me.Score == *opponent.Score {
		return me, opponent, false
	}
	return me, opponent, true
}

func (t *FightcadeTracker) Authenticate(ctx context.Context, statChan chan tracker.AuthStatus) {
	send(ctx, statChan, tracker.AuthStatus{Step: "authOpeningFightcade"})
	err := t.client.Open(ctx, func(step string) {
		send(ctx, statChan, tracker.AuthStatus{Step: step})
	})
	if err != nil {
		send(ctx, statChan, tracker.AuthStatus{Err: err})
		return
	}
	send(ctx, statChan, tracker.AuthStatus{Done: true})
}

func send(ctx context.Context, statChan chan tracker.AuthStatus, status tracker.AuthStatus) {
	select {
	case statChan <- status:
	case <-ctx.Done():
	}
}

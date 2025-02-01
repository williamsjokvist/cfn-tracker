package t8

import (
	"context"
	"fmt"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

type T8Tracker struct {
	wavuClient wavu.WavuClient
}

var _ tracker.GameTracker = (*T8Tracker)(nil)

func NewT8Tracker(wavuClient wavu.WavuClient) *T8Tracker {
	return &T8Tracker{
		wavuClient,
	}
}

func (t *T8Tracker) GetUser(ctx context.Context, polarisId string) (*model.User, error) {
	lastReplay, err := t.wavuClient.GetLastReplay(ctx, polarisId)
	if err != nil {
		return nil, fmt.Errorf("get last replay: %w", err)
	}
	name := polarisId
	if lastReplay != nil {
		if lastReplay.P1PolarisId == polarisId {
			name = lastReplay.P1Name
		} else {
			name = lastReplay.P2Name
		}
	}
	return &model.User{
		DisplayName: name,
		Code:        polarisId,
	}, nil
}

func (t *T8Tracker) Poll(ctx context.Context, cancel context.CancelFunc, session *model.Session, onNewMatch func(model.Match)) {
	lastReplay, err := t.wavuClient.GetLastReplay(ctx, session.UserId)
	if err != nil {
		cancel()
	}
	var prevMatch model.Match
	if len(session.Matches) > 0 {
		prevMatch = *session.Matches[0]
	}
	if lastReplay == nil || prevMatch.ReplayID == lastReplay.BattleId {
		return
	}
	battleAt := time.Unix(lastReplay.BattleAt, 0)
	if time.Since(battleAt).Hours() >= 24 {
		return
	}
	polarisId := lastReplay.P1PolarisId
	userName := lastReplay.P1Name
	character := lastReplay.P1CharaId
	opponentCharacter := lastReplay.P2CharaId
	victory := lastReplay.Winner == 1
	opponent := lastReplay.P2Name
	opponentRank := lastReplay.P2Rank

	if lastReplay.P2PolarisId == session.UserId {
		polarisId = lastReplay.P2PolarisId
		userName = lastReplay.P2Name
		character = lastReplay.P2CharaId
		opponentCharacter = lastReplay.P1CharaId
		victory = lastReplay.Winner == 2
		opponent = lastReplay.P1Name
		opponentRank = lastReplay.P1Rank
	}

	wins := prevMatch.Wins
	losses := prevMatch.Losses
	winStreak := prevMatch.WinStreak
	if victory {
		wins++
		winStreak++
	} else {
		losses++
		winStreak = 0
	}

	onNewMatch(model.Match{
		SessionId: session.Id,
		UserName:  userName,
		UserId:    polarisId,
		Opponent:  opponent,
		Victory:   victory,
		ReplayID:  lastReplay.BattleId,
		Wins:      wins,
		Losses:    losses,
		WinStreak: winStreak,
		WinRate: func() int {
			totalGames := wins + losses
			if totalGames == 0 {
				return 0
			}
			return int((float64(wins) / float64(totalGames)) * 100)
		}(),
		Character:         wavu.ConvCharaIdToName(character),
		OpponentCharacter: wavu.ConvCharaIdToName(opponentCharacter),
		OpponentLeague:    wavu.ConvRankToName(opponentRank),
		Date:              battleAt.Format("2006-01-02"),
		Time:              battleAt.Format("15:04"),
	})
}

func (t *T8Tracker) Authenticate(ctx context.Context, email string, password string, statChan chan tracker.AuthStatus) {
	statChan <- tracker.AuthStatus{Progress: 100, Err: nil}
}

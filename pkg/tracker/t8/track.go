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
	wm, err := t.wavuClient.GetLastReplay(ctx, session.UserId)
	if err != nil {
		cancel()
	}
	var prevMatch model.Match
	if len(session.Matches) > 0 {
		prevMatch = *session.Matches[0]
	}
	if wm == nil || prevMatch.ReplayID == wm.BattleId {
		return
	}
	battleAt := time.Unix(wm.BattleAt, 0)
	polarisId := wm.P1PolarisId
	userName := wm.P1Name
	character := wm.P1CharaId
	opponentCharacter := wm.P2CharaId
	victory := wm.Winner == 1
	opponent := wm.P2Name
	opponentRank := wm.P2Rank

	if wm.P2PolarisId == session.UserId {
		polarisId = wm.P2PolarisId
		userName = wm.P2Name
		character = wm.P2CharaId
		opponentCharacter = wm.P1CharaId
		victory = wm.Winner == 2
		opponent = wm.P1Name
		opponentRank = wm.P1Rank
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
		ReplayID:  wm.BattleId,
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

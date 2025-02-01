package sf6

import (
	"context"
	"fmt"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
)

type SF6Tracker struct {
	cfnClient cfn.CFNClient
}

var _ tracker.GameTracker = (*SF6Tracker)(nil)

func NewSF6Tracker(cfnClient cfn.CFNClient) *SF6Tracker {
	return &SF6Tracker{
		cfnClient,
	}
}

func (t *SF6Tracker) GetUser(ctx context.Context, userCode string) (*model.User, error) {
	bl, err := t.cfnClient.GetBattleLog(ctx, userCode)
	if err != nil {
		return nil, fmt.Errorf("fetch battle log: %w", err)
	}
	return &model.User{
		DisplayName: bl.GetCFN(),
		Code:        userCode,
		LP:          bl.GetLP(),
		MR:          bl.GetMR(),
	}, nil
}

func (t *SF6Tracker) Poll(ctx context.Context, cancel context.CancelFunc, session *model.Session, onNewMatch func(model.Match)) {
	bl, err := t.cfnClient.GetBattleLog(ctx, session.UserId)
	if err != nil {
		cancel()
	}
	if len(bl.ReplayList) == 0 {
		return
	}
	lastReplay := bl.ReplayList[0]
	var prevMatch model.Match
	if len(session.Matches) > 0 {
		prevMatch = getPreviousMatchForCharacter(session, bl.GetCharacter())
	}
	if session.LP == bl.GetLP() || prevMatch.ReplayID == lastReplay.ReplayID {
		return
	}

	battleAt := time.Unix(lastReplay.UploadedAt, 0)
	if time.Since(battleAt).Hours() >= 24 {
		return
	}

	var opponent cfn.PlayerInfo
	if bl.GetCFN() == lastReplay.Player1Info.Player.FighterID {
		opponent = lastReplay.Player2Info
	} else {
		opponent = lastReplay.Player1Info
	}
	victory := !isVictory(opponent.RoundResults)
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
		Character:         bl.GetCharacter(),
		LP:                bl.GetLP(),
		MR:                bl.GetMR(),
		Opponent:          opponent.Player.FighterID,
		OpponentCharacter: opponent.CharacterName,
		OpponentLP:        opponent.LeaguePoint,
		OpponentLeague:    getLeagueFromLP(opponent.LeaguePoint),
		OpponentMR:        opponent.MasterRating,
		Victory:           victory,
		Wins:              wins,
		Losses:            losses,
		WinStreak:         winStreak,
		Date:              time.Now().Format(`2006-01-02`),
		Time:              time.Now().Format(`15:04`),
		LPGain:            prevMatch.LPGain + bl.GetLP() - session.LP,
		MRGain:            prevMatch.MRGain + bl.GetMR() - session.MR,
		WinRate:           int((float64(wins) / float64(wins+losses)) * 100),
		UserId:            session.UserId,
		UserName:          session.UserName,
		SessionId:         session.Id,
		ReplayID:          lastReplay.ReplayID,
	})
}

func getPreviousMatchForCharacter(sesh *model.Session, character string) model.Match {
	for i := 0; i < len(sesh.Matches); i++ {
		match := sesh.Matches[i]
		if match.Character == character {
			return *match
		}
	}
	return model.Match{}
}

func isVictory(roundResults []int) bool {
	roundsPlayed := len(roundResults)
	losses := make([]int, 0, roundsPlayed)
	for _, result := range roundResults {
		if result == 0 {
			losses = append(losses, result)
		}
	}
	return (roundsPlayed == 3 && len(losses) == 1) || len(losses) == 0
}

func getLeagueFromLP(lp int) string {
	if lp >= 25000 {
		return `Master`
	} else if lp >= 20000 {
		return `Diamond`
	} else if lp >= 14000 {
		return `Platinum`
	} else if lp >= 9000 {
		return `Gold`
	} else if lp >= 5000 {
		return `Silver`
	} else if lp >= 3000 {
		return `Bronze`
	} else if lp >= 1000 {
		return `Iron`
	}

	return `Rookie`
}

func (t *SF6Tracker) Authenticate(ctx context.Context, email string, password string, statChan chan tracker.AuthStatus) {
	t.cfnClient.Authenticate(ctx, email, password, statChan)
}

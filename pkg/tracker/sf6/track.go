package sf6

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
)

type SF6Tracker struct {
	cfnClient *cfn.Client
	sqlDb     *sql.Storage
}

var _ tracker.GameTracker = (*SF6Tracker)(nil)

func NewSF6Tracker(browser *browser.Browser, sqlDb *sql.Storage) *SF6Tracker {
	return &SF6Tracker{
		cfnClient: cfn.NewCFNClient(browser),
		sqlDb:     sqlDb,
	}
}

// Start will update the tracking state when new matches are played.
func (t *SF6Tracker) Init(ctx context.Context, userCode string, restore bool) (*model.Session, error) {
	if restore {
		session, err := t.sqlDb.GetLatestSession(ctx, userCode)
		if err != nil {
			return nil, errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get last session: %w`, err))
		}
		_, err = t.sqlDb.GetUserByCode(ctx, userCode)
		if err != nil {
			return nil, errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get user: %w`, err))
		}
		return session, nil
	}

	bl, err := t.cfnClient.GetBattleLog(userCode)
	if err != nil {
		return nil, errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to fetch battle log: %w`, err))
	}

	err = t.sqlDb.SaveUser(ctx, model.User{
		DisplayName: bl.GetCFN(),
		Code:        userCode,
	})
	if err != nil {
		return nil, errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to save user: %w`, err))
	}
	session, err := t.sqlDb.CreateSession(ctx, userCode)
	if err != nil {
		return nil, errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to create session: %w`, err))
	}

	// set starting LP so we don't count the first polled match
	session.LP = bl.GetLP()
	session.MR = bl.GetMR()
	session.UserName = bl.GetCFN()
	return session, nil
}

func (t *SF6Tracker) Poll(ctx context.Context, cancel context.CancelFunc, session *model.Session, onNewMatch func(model.Match)) {
	bl, err := t.cfnClient.GetBattleLog(session.UserId)
	if err != nil {
		cancel()
	}
	if len(bl.ReplayList) == 0 {
		return
	}
	lastReplay := bl.ReplayList[0]
	var prevMatch model.Match
	if len(session.Matches) > 0 {
		prevMatch = *session.Matches[0]
	}
	if session.LP == bl.GetLP() || prevMatch.ReplayID == lastReplay.ReplayID {
		return
	}
	onNewMatch(getMatch(session, bl))
}

func getOpponentInfo(myCfn string, replay *cfn.Replay) cfn.PlayerInfo {
	if myCfn == replay.Player1Info.Player.FighterID {
		return replay.Player2Info
	} else {
		return replay.Player1Info
	}
}

func getMatch(sesh *model.Session, bl *cfn.BattleLog) model.Match {
	latestReplay := bl.ReplayList[0]
	opponent := getOpponentInfo(bl.GetCFN(), &latestReplay)
	victory := !isVictory(opponent.RoundResults)

	prevMatch := getPreviousMatchForCharacter(sesh, bl.GetCharacter())

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

	return model.Match{
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
		LPGain:            prevMatch.LPGain + bl.GetLP() - sesh.LP,
		MRGain:            prevMatch.MRGain + bl.GetMR() - sesh.MR,
		WinRate:           int((float64(wins) / float64(wins+losses)) * 100),
		UserId:            sesh.UserId,
		UserName:          sesh.UserName,
		SessionId:         sesh.Id,
		ReplayID:          latestReplay.ReplayID,
	}
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

func (t *SF6Tracker) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	t.cfnClient.Authenticate(email, password, statChan)
}

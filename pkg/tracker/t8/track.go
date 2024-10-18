package t8

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

type T8Tracker struct {
	wavuClient *wavu.Client
	sqlDb      *sql.Storage
	txtDb      *txt.Storage
}

var _ tracker.GameTracker = (*T8Tracker)(nil)

func NewT8Tracker(sqlDb *sql.Storage, txtDb *txt.Storage) *T8Tracker {
	return &T8Tracker{
		wavuClient: wavu.NewClient(),
		sqlDb:      sqlDb,
		txtDb:      txtDb,
	}
}

func (t *T8Tracker) Init(ctx context.Context, polarisId string, restore bool) (*model.Session, error) {
	if restore {
		session, err := t.sqlDb.GetLatestSession(ctx, polarisId)
		if err != nil {
			return nil, errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("get last session: %w", err))
		}
		return session, nil
	}

	user, err := t.sqlDb.GetUserByCode(ctx, polarisId)
	if err != nil && !errors.Is(err, sql.ErrUserNotFound) {
		return nil, errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("get user: %w", err))
	}

	if user == nil {
		if err := t.createUser(ctx, polarisId); err != nil {
			return nil, errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("create user: %w", err))
		}
	}

	session, err := t.sqlDb.CreateSession(ctx, polarisId)
	if err != nil {
		return nil, errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("create session: %w", err))
	}
	return session, nil
}

func (t *T8Tracker) createUser(ctx context.Context, polarisId string) error {
	lastReplay, err := t.wavuClient.GetLastReplay(polarisId)
	if err != nil {
		return fmt.Errorf("get last replay: %w", err)
	}
	name := polarisId
	if lastReplay != nil {
		if lastReplay.P1PolarisId == polarisId {
			name = lastReplay.P1Name
		} else {
			name = lastReplay.P2Name
		}
	}
	if err := t.sqlDb.SaveUser(ctx, model.User{DisplayName: name, Code: polarisId}); err != nil {
		return fmt.Errorf("save user: %w", err)
	}
	return nil
}

func (t *T8Tracker) Poll(ctx context.Context, cancel context.CancelFunc, session *model.Session, matchChan chan model.Match) {
	lastReplay, err := t.wavuClient.GetLastReplay(session.UserId)
	if err != nil {
		cancel()
	}
	var prevMatch *model.Match
	if len(session.Matches) > 0 {
		prevMatch = session.Matches[0]
	}
	if lastReplay == nil || (prevMatch != nil && prevMatch.ReplayID == lastReplay.BattleId) {
		return
	}
	p2 := lastReplay.P2PolarisId == session.UserId
	matchChan <- getMatch(lastReplay, prevMatch.Wins, prevMatch.Losses, prevMatch.WinStreak, session.Id, p2)
}

func getMatch(wm *wavu.Replay, wins, losses, winStreak int, sessionId uint16, p2 bool) model.Match {
	polarisId := wm.P1PolarisId
	userName := wm.P1Name
	character := wm.P1CharaId
	opponentCharacter := wm.P2CharaId
	victory := wm.Winner == 1
	opponent := wm.P2Name
	opponentRank := wm.P2Rank

	if p2 {
		polarisId = wm.P2PolarisId
		userName = wm.P2Name
		character = wm.P2CharaId
		opponentCharacter = wm.P1CharaId
		victory = wm.Winner == 2
		opponent = wm.P1Name
		opponentRank = wm.P1Rank
	}

	if victory {
		wins++
		winStreak++
	} else {
		losses++
		winStreak = 0
	}

	battleAt := time.Unix(wm.BattleAt, 0)
	return model.Match{
		SessionId: sessionId,
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
	}
}

func (t *T8Tracker) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	statChan <- tracker.AuthStatus{Progress: 100, Err: nil}
}

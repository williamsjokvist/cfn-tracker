package t8

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
	"github.com/williamsjokvist/cfn-tracker/pkg/utils"
)

type T8Tracker struct {
	cancel     context.CancelFunc
	wavuClient *wavu.Client
	sqlDb      *sql.Storage
	txtDb      *txt.Storage
}

var _ tracker.GameTracker = (*T8Tracker)(nil)

func NewT8Tracker(sqlDb *sql.Storage, txtDb *txt.Storage) *T8Tracker {
	return &T8Tracker{
		cancel:     func() {},
		sqlDb:      sqlDb,
		txtDb:      txtDb,
		wavuClient: wavu.NewClient(),
	}
}

func (t *T8Tracker) Start(ctx context.Context, polarisId string, restore bool, pollRate time.Duration) error {
	if restore {
		session, err := t.sqlDb.GetLatestSession(ctx, polarisId)
		if err != nil {
			return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("get last session: %w", err))
		}
		if len(session.Matches) > 0 {
			lastMatch := session.Matches[0]
			trackingState := model.ConvMatchToTrackingState(*lastMatch)
			wails.EventsEmit(ctx, "cfn-data", trackingState)
		}
		go t.poll(ctx, session, polarisId, pollRate)
		return nil
	}

	user, err := t.sqlDb.GetUserByCode(ctx, polarisId)
	if err != nil && !errors.Is(err, sql.ErrUserNotFound) {
		return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("get user: %w", err))
	}

	if user == nil {
		lastReplay, err := t.wavuClient.GetLastReplay(polarisId)
		if err != nil {
			return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("get last replay: %w", err))
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
			return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("save user: %w", err))
		}
	}

	session, err := t.sqlDb.CreateSession(ctx, polarisId)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("create session: %w", err))
	}
	go t.poll(ctx, session, polarisId, pollRate)
	return nil
}

func (t *T8Tracker) poll(ctx context.Context, session *model.Session, polarisId string, pollRate time.Duration) {
	// todo: more sophisticated poll rate
	ticker := time.NewTicker(pollRate)
	defer func() {
		ticker.Stop()
		wails.EventsEmit(ctx, "stopped-tracking")
	}()

	pollCtx, cancelFn := context.WithCancel(ctx)
	t.cancel = cancelFn

	i := 0
	for {
		select {
		case <-ticker.C:
			i++
			log.Println("polling", i)
			lastReplay, err := t.wavuClient.GetLastReplay(polarisId)
			if err != nil {
				return
			}
			var prevMatch *model.Match
			if len(session.Matches) > 0 {
				prevMatch = session.Matches[0]
			}
			if lastReplay == nil || (prevMatch != nil && prevMatch.ReplayID == lastReplay.BattleId) {
				continue
			}

			match := getMatch(lastReplay, prevMatch, lastReplay.P2PolarisId == polarisId)
			if match.SessionId == 0 {
				match.SessionId = session.Id
			}
			session.Matches = append([]*model.Match{&match}, session.Matches...)
			if err := t.sqlDb.SaveMatch(ctx, match); err != nil {
				return
			}

			trackingState := model.ConvMatchToTrackingState(match)
			t.txtDb.SaveTrackingState(&trackingState)
			wails.EventsEmit(ctx, "cfn-data", trackingState)
		case <-pollCtx.Done():
			return
		}
	}
}

func getMatch(wm *wavu.Replay, prevMatch *model.Match, p2 bool) model.Match {
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

	biota := utils.Biota(victory)
	wins := biota
	losses := (1 - biota)
	winStreak := biota
	sessionId := uint16(0)
	if prevMatch != nil {
		wins = prevMatch.Wins + biota
		losses = prevMatch.Losses + (1 - biota)
		winStreak = prevMatch.WinStreak*biota + biota
		sessionId = prevMatch.SessionId
	}

	battleAt := time.Unix(wm.BattleAt, 0)
	return model.Match{
		SessionId:         sessionId,
		UserName:          userName,
		UserId:            polarisId,
		Opponent:          opponent,
		Victory:           victory,
		ReplayID:          wm.BattleId,
		Wins:              wins,
		Losses:            losses,
		WinStreak:         winStreak,
		WinRate:           int((float64(wins) / float64(wins+losses)) * 100),
		Character:         wavu.ConvCharaIdToName(character),
		OpponentCharacter: wavu.ConvCharaIdToName(opponentCharacter),
		OpponentLeague:    wavu.ConvRankToName(opponentRank),
		Date:              battleAt.Format("2006-01-02"),
		Time:              battleAt.Format("15:04"),
	}
}

func (t *T8Tracker) Stop() {
	t.cancel()
}

func (t *T8Tracker) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	statChan <- tracker.AuthStatus{Progress: 100, Err: nil}
}

package t8

import (
	"context"
	"fmt"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"net/http"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
	"github.com/williamsjokvist/cfn-tracker/pkg/utils"
)

type T8Tracker struct {
	stopPolling context.CancelFunc
	wavuClient  *wavu.Client
	sqlDb       *sql.Storage
	txtDb       *txt.Storage
}

var _ tracker.GameTracker = (*T8Tracker)(nil)

func NewT8Tracker(sqlDb *sql.Storage, txtDb *txt.Storage) *T8Tracker {
	return &T8Tracker{
		stopPolling: func() {},
		sqlDb:       sqlDb,
		txtDb:       txtDb,
		wavuClient:  wavu.NewClient(),
	}
}

func (t *T8Tracker) Start(ctx context.Context, polarisId string, restore bool, pollRate time.Duration) error {
	var session *model.Session
	var user *model.User
	var err error
	if restore {
		session, err = t.sqlDb.GetLatestSession(ctx, polarisId)
		if err != nil {
			return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("get last session: %w", err))
		}
		user, err = t.sqlDb.GetUserByCode(ctx, polarisId)
		if err != nil {
			return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("get user: %w", err))
		}
		if len(session.Matches) > 0 {
			lastMatch := session.Matches[0]
			trackingState := model.ConvMatchToTrackingState(*lastMatch)
			wails.EventsEmit(ctx, "cfn-data", trackingState)
		}
		go t.poll(ctx, session, polarisId, pollRate)
		return nil
	}

	lastReplay, err := t.wavuClient.GetLastReplay(polarisId)
	if lastReplay.P1PolarisId == polarisId {
		user = &model.User{DisplayName: lastReplay.P1Name, Code: polarisId}
	} else {
		user = &model.User{DisplayName: lastReplay.P2Name, Code: polarisId}
	}
	err = t.sqlDb.SaveUser(ctx, *user)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("save user: %w", err))
	}
	session, err = t.sqlDb.CreateSession(ctx, polarisId)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("create session: %w", err))
	}
	go t.poll(ctx, session, polarisId, pollRate)
	return nil
}

func (t *T8Tracker) poll(ctx context.Context, session *model.Session, polarisId string, pollRate time.Duration) {
	pollCtx, cancelFn := context.WithCancel(ctx)
	t.stopPolling = cancelFn

	didStop := func() bool {
		return utils.SleepOrBreak(pollRate, func() bool {
			select {
			case <-pollCtx.Done():
				return true
			default:
				return false
			}
		})
	}
	i := 0
	for {
		i++
		log.Println("polling", i)
		if didStop() {
			wails.EventsEmit(ctx, "stopped-tracking")
			break
		}
		lastReplay, err := t.wavuClient.GetLastReplay(polarisId)
		if err != nil {
			wails.EventsEmit(ctx, "stopped-tracking")
			t.stopPolling()
			return
		}
		var prevMatch *model.Match
		if len(session.Matches) > 0 {
			prevMatch = session.Matches[0]
		}
		latestMatch := t.getLatestMatch(session.Id, lastReplay, prevMatch, lastReplay.P2PolarisId == polarisId)
		if lastReplay == nil || (len(session.Matches) > 0 && session.Matches[0].ReplayID == latestMatch.ReplayID) {
			continue
		}
		session.Matches = append([]*model.Match{&latestMatch}, session.Matches...)
		err = t.sqlDb.UpdateSession(ctx, session, latestMatch, session.Id)
		if err != nil {
			wails.EventsEmit(ctx, "stopped-tracking")
			t.stopPolling()
			return
		}

		trackingState := model.ConvMatchToTrackingState(latestMatch)
		t.txtDb.SaveTrackingState(&trackingState)
		wails.EventsEmit(ctx, "cfn-data", trackingState)
	}
}

func (t *T8Tracker) getLatestMatch(sessionId uint16, wm *wavu.Replay, prevMatch *model.Match, p2 bool) model.Match {
	polarisId := wm.P1PolarisId
	userName := wm.P1Name
	character := wm.P1CharaId
	opponentCharacter := wm.P2CharaId
	victory := wm.Winner == 1
	opponent := wm.P2Name
	opponentLeague := wm.P2Rank

	if p2 {
		polarisId = wm.P2PolarisId
		userName = wm.P2Name
		character = wm.P2CharaId
		opponentCharacter = wm.P1CharaId
		victory = wm.Winner == 2
		opponent = wm.P1Name
		opponentLeague = wm.P1Rank
	}

	biota := utils.Biota(victory)
	wins := biota
	losses := (1 - biota)
	winStreak := biota
	if prevMatch != nil {
		wins = prevMatch.Wins + biota
		losses = prevMatch.Losses + (1 - biota)
		winStreak = prevMatch.WinStreak*biota + biota
	}

	battleTime := time.Unix(wm.BattleAt, 0)
	return model.Match{
		UserName:          userName,
		UserId:            polarisId,
		Opponent:          opponent,
		Character:         wavu.ConvCharaIdToName(character),
		Victory:           victory,
		OpponentCharacter: wavu.ConvCharaIdToName(opponentCharacter),
		OpponentLeague:    wavu.ConvDanToRank(opponentLeague),
		Date:              battleTime.Format("2006-01-02"),
		Time:              battleTime.Format("15:04"),
		ReplayID:          wm.BattleId,
		Wins:              wins,
		Losses:            losses,
		WinStreak:         winStreak,
		WinRate:           int((float64(wins) / float64(wins+losses)) * 100),
		SessionId:         sessionId,
	}
}

func (t *T8Tracker) Stop() {
	t.stopPolling()
}

func (t *T8Tracker) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	statChan <- tracker.AuthStatus{Progress: 100, Err: nil}
}

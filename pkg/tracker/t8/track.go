package t8

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

type T8Tracker struct {
	cancel        context.CancelFunc
	forcePollChan chan os.Signal
	wavuClient    *wavu.Client
	sqlDb         *sql.Storage
	txtDb         *txt.Storage
}

var _ tracker.GameTracker = (*T8Tracker)(nil)

func NewT8Tracker(sqlDb *sql.Storage, txtDb *txt.Storage) *T8Tracker {
	var forcePollChan = make(chan os.Signal, 2)
	signal.Notify(forcePollChan, os.Interrupt)
	return &T8Tracker{
		cancel:        func() {},
		forcePollChan: forcePollChan,
		sqlDb:         sqlDb,
		txtDb:         txtDb,
		wavuClient:    wavu.NewClient(),
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

		go t.poll(ctx, session, pollRate)
		return nil
	}

	user, err := t.sqlDb.GetUserByCode(ctx, polarisId)
	if err != nil && !errors.Is(err, sql.ErrUserNotFound) {
		return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("get user: %w", err))
	}

	if user == nil {
		if err := t.createUser(ctx, polarisId); err != nil {
			return errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf("create user: %w", err))
		}
	}

	session, err := t.sqlDb.CreateSession(ctx, polarisId)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("create session: %w", err))
	}

	go t.poll(ctx, session, pollRate)
	return nil
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

func (t *T8Tracker) poll(ctx context.Context, session *model.Session, pollRate time.Duration) {
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
		case <-t.forcePollChan:
			i++
			log.Println("forced poll", i)
			t.pollFn(ctx, session)
		case <-ticker.C:
			i++
			log.Println("polling", i)
			t.pollFn(ctx, session)
		case <-pollCtx.Done():
			return
		}
	}
}

func (t *T8Tracker) pollFn(ctx context.Context, session *model.Session) {
	lastReplay, err := t.wavuClient.GetLastReplay(session.UserId)
	if err != nil {
		t.Stop()
	}
	var prevMatch *model.Match
	if len(session.Matches) > 0 {
		prevMatch = session.Matches[0]
	}
	if lastReplay == nil || (prevMatch != nil && prevMatch.ReplayID == lastReplay.BattleId) {
		return
	}
	match := getMatch(lastReplay, prevMatch, lastReplay.P2PolarisId == session.UserId)
	if match.SessionId == 0 {
		match.SessionId = session.Id
	}
	session.Matches = append([]*model.Match{&match}, session.Matches...)
	if err := t.sqlDb.SaveMatch(ctx, match); err != nil {
		t.Stop()
	}

	trackingState := model.ConvMatchToTrackingState(match)
	wails.EventsEmit(ctx, "cfn-data", trackingState)
	if err := t.txtDb.SaveTrackingState(&trackingState); err != nil {
		t.Stop()
	}
}

func (t *T8Tracker) ForcePoll() {
	t.forcePollChan <- os.Interrupt
}

func (t *T8Tracker) Stop() {
	t.cancel()
}

func (t *T8Tracker) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	statChan <- tracker.AuthStatus{Progress: 100, Err: nil}
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

	sessionId := uint16(0)
	wins := 0
	losses := 0
	winStreak := 0
	if prevMatch != nil {
		wins = prevMatch.Wins
		losses = prevMatch.Losses
		winStreak = prevMatch.WinStreak
		sessionId = prevMatch.SessionId
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

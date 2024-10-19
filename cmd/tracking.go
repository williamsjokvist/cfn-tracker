package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/nosql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8"
)

type TrackingHandler struct {
	gameTracker tracker.GameTracker

	browser *browser.Browser

	cancelPolling context.CancelFunc
	forcePollChan chan struct{}

	sqlDb   *sql.Storage
	nosqlDb *nosql.Storage
	txtDb   *txt.Storage

	eventEmitter EventEmitFn

	cfg *config.Config
}

var _ CmdHandler = (*TrackingHandler)(nil)

func NewTrackingHandler(browser *browser.Browser, sqlDb *sql.Storage, nosqlDb *nosql.Storage, txtDb *txt.Storage, cfg *config.Config) *TrackingHandler {
	return &TrackingHandler{
		sqlDb:   sqlDb,
		nosqlDb: nosqlDb,
		txtDb:   txtDb,
		browser: browser,
		cfg:     cfg,
	}
}

func (ch *TrackingHandler) SetEventEmitter(eventEmitter EventEmitFn) {
	ch.eventEmitter = eventEmitter
}

func (ch *TrackingHandler) StartTracking(userCode string, restore bool) {
	log.Printf(`Starting tracking for %s, restoring = %v`, userCode, restore)
	ticker := time.NewTicker(30 * time.Second)
	ctx, cancel := context.WithCancel(context.Background())
	ch.cancelPolling = cancel
	ch.forcePollChan = make(chan struct{})
	var matchChan = make(chan model.Match)

	defer func() {
		ticker.Stop()
		cancel()
		close(ch.forcePollChan)
		ch.forcePollChan = nil
		ch.eventEmitter("stopped-tracking")
	}()

	session, err := ch.gameTracker.Init(ctx, userCode, restore)
	if err != nil {
		return
	}

	if len(session.Matches) > 0 {
		ch.eventEmitter("match", *session.Matches[0])
	}

	go func() {
		log.Println("polling")
		ch.gameTracker.Poll(ctx, cancel, session, matchChan)
		for {
			select {
			case <-ch.forcePollChan:
				log.Println("forced poll")
				ch.gameTracker.Poll(ctx, cancel, session, matchChan)
			case <-ticker.C:
				log.Println("polling")
				ch.gameTracker.Poll(ctx, cancel, session, matchChan)
			case <-ctx.Done():
				close(matchChan)
				return
			}
		}
	}()

	for match := range matchChan {
		ch.eventEmitter("match", match)

		session.LP = match.LP
		session.MR = match.MR
		session.Matches = append([]*model.Match{&match}, session.Matches...)

		if err := ch.sqlDb.UpdateSession(ctx, session); err != nil {
			log.Println("failed to update session", err)
			break
		}
		if err := ch.sqlDb.SaveMatch(ctx, match); err != nil {
			log.Println("failed to save match to database", err)
			break
		}
		if err := ch.txtDb.SaveMatch(match); err != nil {
			log.Println("failed to save to text files:", err)
			break
		}
	}
}

func (ch *TrackingHandler) StopTracking() {
	ch.cancelPolling()
}

func (ch *TrackingHandler) SelectGame(game model.GameType) error {
	var username, password string

	switch game {
	case model.GameTypeT8:
		ch.gameTracker = t8.NewT8Tracker(ch.sqlDb, ch.txtDb)
	case model.GameTypeSF6:
		ch.gameTracker = sf6.NewSF6Tracker(ch.browser, ch.sqlDb, ch.txtDb)
		username = ch.cfg.CapIDEmail
		password = ch.cfg.CapIDPassword
	default:
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to select game`))
	}

	authChan := make(chan tracker.AuthStatus)
	go ch.gameTracker.Authenticate(username, password, authChan)
	for status := range authChan {
		if status.Err != nil {
			return errorsx.NewFormattedError(http.StatusUnauthorized, status.Err)
		}

		ch.eventEmitter("auth-progress", status.Progress)

		if status.Progress >= 100 {
			close(authChan)
			break
		}
	}
	return nil
}

func (ch *TrackingHandler) ForcePoll() {
	if ch.forcePollChan != nil {
		ch.forcePollChan <- struct{}{}
	}
}

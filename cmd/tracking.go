package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/nosql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6"
	_ "github.com/williamsjokvist/cfn-tracker/pkg/tracker/sfv"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8"
)

type TrackingHandler struct {
	ctx         context.Context
	gameTracker tracker.GameTracker

	browser *browser.Browser

	cancelPolling context.CancelFunc
	forcePollChan chan struct{}

	sqlDb   *sql.Storage
	nosqlDb *nosql.Storage
	txtDb   *txt.Storage

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

// The CommandHandler needs the wails runtime context in order to emit events
func (ch *TrackingHandler) SetContext(ctx context.Context) {
	ch.ctx = ctx
}

func (ch *TrackingHandler) StartTracking(userCode string, restore bool) {
	log.Printf(`Starting tracking for %s, restoring = %v`, userCode, restore)
	ticker := time.NewTicker(30 * time.Second)
	pollCtx, cancel := context.WithCancel(ch.ctx)
	ch.cancelPolling = cancel
	ch.forcePollChan = make(chan struct{})
	var matchChan = make(chan model.Match)

	defer func() {
		ticker.Stop()
		cancel()
		close(ch.forcePollChan)
		ch.forcePollChan = nil
		wailsRuntime.EventsEmit(ch.ctx, "stopped-tracking")
		log.Println("stopped polling")
	}()

	session, err := ch.gameTracker.Init(pollCtx, userCode, restore)
	if err != nil {
		return
	}

	go func() {
		log.Println("polling")
		ch.gameTracker.Poll(pollCtx, cancel, session, matchChan)
		for {
			select {
			case <-ch.forcePollChan:
				log.Println("forced poll")
				ch.gameTracker.Poll(pollCtx, cancel, session, matchChan)
			case <-ticker.C:
				log.Println("polling")
				ch.gameTracker.Poll(pollCtx, cancel, session, matchChan)
			case <-pollCtx.Done():
				close(matchChan)
				return
			}
		}
	}()

	for match := range matchChan {
		session.LP = match.LP
		session.MR = match.MR
		session.Matches = append([]*model.Match{&match}, session.Matches...)
		if err := ch.sqlDb.UpdateSession(ch.ctx, session); err != nil {
			log.Println("failed to update session", err)
			return
		}
		if err := ch.sqlDb.SaveMatch(ch.ctx, match); err != nil {
			log.Println("failed to save match", err)
			return
		}

		wailsRuntime.EventsEmit(ch.ctx, `cfn-data`, match)
		if err := ch.txtDb.SaveMatch(match); err != nil {
			log.Print("failed to save tracking state:", err)
			return
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
	case model.GameTypeSFV:
		// gameTracker = sfv.NewSFVTracker(ch.browser)
		// username = ch.cfg.SteamUsername
		// password = ch.cfg.SteamPassword
		fallthrough
	default:
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to select game`))
	}

	authChan := make(chan tracker.AuthStatus)
	go ch.gameTracker.Authenticate(username, password, authChan)
	for status := range authChan {
		if status.Err != nil {
			return errorsx.NewFormattedError(http.StatusUnauthorized, status.Err)
		}
		wailsRuntime.EventsEmit(ch.ctx, "auth-progress", status.Progress)

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

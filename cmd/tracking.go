package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	cfgDb "github.com/williamsjokvist/cfn-tracker/pkg/storage/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

type EventEmitFn func(eventName string, optionalData ...interface{})

type TrackingHandler struct {
	sqlDb   *sql.Storage
	nosqlDb *cfgDb.Storage
	txtDb   *txt.Storage

	wavuClient wavu.WavuClient
	cfnClient  cfn.CFNClient

	cfg        *config.BuildConfig
	matchChans []chan model.Match

	cancelPolling context.CancelFunc
	forcePollChan chan struct{}
	gameTracker   tracker.GameTracker
	eventEmitter  EventEmitFn
}

func NewTrackingHandler(
	wavuClient wavu.WavuClient,
	cfnClient cfn.CFNClient,
	sqlDb *sql.Storage,
	nosqlDb *cfgDb.Storage,
	txtDb *txt.Storage,
	cfg *config.BuildConfig,
	matchChans ...chan model.Match,
) *TrackingHandler {
	return &TrackingHandler{
		wavuClient: wavuClient,
		cfnClient:  cfnClient,
		sqlDb:      sqlDb,
		nosqlDb:    nosqlDb,
		txtDb:      txtDb,
		cfg:        cfg,
		matchChans: matchChans,
	}
}

func (ch *TrackingHandler) SetEventEmitter(eventEmitter EventEmitFn) {
	ch.eventEmitter = eventEmitter
}

func (ch *TrackingHandler) StartTracking(userCode string, restore bool) error {
	log.Printf(`Starting tracking for %s, restoring = %v`, userCode, restore)

	ctx, cancel := context.WithCancel(context.Background())
	ch.cancelPolling = cancel

	var session *model.Session
	if restore {
		sesh, err := ch.sqlDb.GetLatestSession(ctx, userCode)
		if err != nil {
			return model.WrapError(model.ErrGetLatestSession, err)
		}
		session = sesh
	} else {
		user, err := ch.sqlDb.GetUserByCode(ctx, userCode)
		if err != nil && !errors.Is(err, sql.ErrUserNotFound) {
			return model.WrapError(model.ErrGetUser, err)
		}

		if user == nil {
			usr, err := ch.gameTracker.GetUser(ctx, userCode)
			if err != nil {
				return model.WrapError(model.ErrGetUser, err)
			}
			if err := ch.sqlDb.SaveUser(ctx, *usr); err != nil {
				return model.WrapError(model.ErrSaveUser, err)
			}
		}

		sesh, err := ch.sqlDb.CreateSession(ctx, userCode)
		if err != nil {
			return model.WrapError(model.ErrCreateSession, err)
		}
		session = sesh
		// session.LP = bl.GetLP()
		// session.MR = bl.GetMR()
		// session.UserName = bl.GetCFN()
	}

	if session == nil {
		return model.ErrCreateSession
	}

	ticker := time.NewTicker(30 * time.Second)
	ch.forcePollChan = make(chan struct{})
	defer func() {
		ch.eventEmitter("stopped-tracking")
		ticker.Stop()
		cancel()
		close(ch.forcePollChan)
		ch.forcePollChan = nil
	}()

	matchChan := make(chan model.Match)

	onNewMatch := func(match model.Match) {
		matchChan <- match
		for _, mc := range ch.matchChans {
			if mc != nil {
				mc <- match
			}
		}
	}

	if len(session.Matches) > 0 {
		match := *session.Matches[0]
		ch.eventEmitter("match", match)
		for _, mc := range ch.matchChans {
			if mc != nil {
				mc <- match
			}
		}
	}

	go func() {
		log.Println("polling")
		ch.gameTracker.Poll(ctx, cancel, session, onNewMatch)
		for {
			select {
			case <-ch.forcePollChan:
				log.Println("forced poll")
				ch.gameTracker.Poll(ctx, cancel, session, onNewMatch)
			case <-ticker.C:
				log.Println("polling")
				ch.gameTracker.Poll(ctx, cancel, session, onNewMatch)
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
			log.Println("failed to update session:", err)
			break
		}
		if err := ch.sqlDb.SaveMatch(ctx, match); err != nil {
			log.Println("failed to save match to database:", err)
			break
		}
		if err := ch.txtDb.SaveMatch(match); err != nil {
			log.Println("failed to save to text files:", err)
			break
		}
	}
	return nil
}

func (ch *TrackingHandler) StopTracking() {
	ch.cancelPolling()
}

func (ch *TrackingHandler) SelectGame(game model.GameType) error {
	var username, password string
	switch game {
	case model.GameTypeT8:
		ch.gameTracker = t8.NewT8Tracker(ch.wavuClient)
	case model.GameTypeSF6:
		ch.gameTracker = sf6.NewSF6Tracker(ch.cfnClient)
		username = ch.cfg.CapIDEmail
		password = ch.cfg.CapIDPassword
	default:
		return model.WrapError(model.ErrSelectGame, fmt.Errorf("game does not exist"))
	}

	authChan := make(chan tracker.AuthStatus)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	go ch.gameTracker.Authenticate(ctx, username, password, authChan)
	for status := range authChan {
		if status.Err != nil {
			return model.WrapError(model.ErrAuth, status.Err)
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

package cmd

import (
	"context"
	"fmt"
	"log/slog"
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
	slog.Info("started tracking", slog.String("user_code", userCode), slog.Bool("restoring", restore))

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
		sesh, err := ch.sqlDb.CreateSession(ctx, userCode)
		if err != nil {
			return model.WrapError(model.ErrCreateSession, err)
		}
		session = sesh
	}
	if session == nil {
		return model.ErrCreateSession
	}

	user, err := ch.gameTracker.GetUser(ctx, userCode)
	if err != nil {
		return model.WrapError(model.ErrGetUser, err)
	}
	if err := ch.sqlDb.SaveUser(ctx, *user); err != nil {
		return model.WrapError(model.ErrSaveUser, err)
	}
	session.LP = user.LP
	session.MR = user.MR
	session.UserName = user.DisplayName

	ch.eventEmitter("match", model.Match{
		UserName:  session.UserName,
		LP:        session.LP,
		MR:        session.MR,
		SessionId: session.Id,
		UserId:    session.UserId,
	})

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

	onNewMatch := func(match *model.Match) {
		if match == nil {
			return
		}
		matchChan <- *match
		for _, mc := range ch.matchChans {
			if mc != nil {
				mc <- *match
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
		slog.Info("polling")
		match, err := ch.gameTracker.Poll(ctx, session)
		if err != nil {
			cancel()
			return
		}
		onNewMatch(match)
		for {
			select {
			case <-ch.forcePollChan:
				slog.Info("forced poll")
				match, err := ch.gameTracker.Poll(ctx, session)
				if err != nil {
					cancel()
					return
				}
				onNewMatch(match)
			case <-ticker.C:
				slog.Info("polling")
				match, err := ch.gameTracker.Poll(ctx, session)
				if err != nil {
					cancel()
					return
				}
				onNewMatch(match)
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
			slog.Error("update session:", slog.Any("error", err))
			break
		}
		if err := ch.sqlDb.SaveMatch(ctx, match); err != nil {
			slog.Error("save match to database", slog.Any("error", err))
			break
		}
		if err := ch.txtDb.SaveMatch(match); err != nil {
			slog.Error("save to text files:", slog.Any("error", err))
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

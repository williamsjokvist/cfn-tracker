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

func (t *TrackingHandler) SetEventEmitter(eventEmitter EventEmitFn) {
	t.eventEmitter = eventEmitter
}

func (t *TrackingHandler) CreateSession(userCode string, restore bool) (*model.Session, error) {
	ctx := context.Background()

	user, err := t.gameTracker.GetUser(ctx, userCode)
	if err != nil {
		return nil, model.WrapError(model.ErrGetUser, err)
	}
	if err := t.sqlDb.SaveUser(ctx, *user); err != nil {
		return nil, model.WrapError(model.ErrSaveUser, err)
	}

	var session *model.Session
	if restore {
		sesh, err := t.sqlDb.GetLatestSession(ctx, userCode)
		if err != nil {
			return nil, model.WrapError(model.ErrGetLatestSession, err)
		}
		session = sesh
	} else {
		sesh, err := t.sqlDb.CreateSession(ctx, userCode)
		if err != nil {
			return nil, model.WrapError(model.ErrCreateSession, err)
		}
		session = sesh
	}

	session.LP = user.LP
	session.MR = user.MR
	session.UserName = user.DisplayName
	if err := t.sqlDb.UpdateSession(ctx, session); err != nil {
		slog.Error("update session:", slog.Any("error", err))
	}
	return session, nil
}

func (t *TrackingHandler) StartTracking(sessionId string) error {
	slog.Info("started tracking", slog.String("session_id", sessionId))

	ctx, cancel := context.WithCancel(context.Background())
	t.cancelPolling = cancel

	session, err := t.sqlDb.GetSession(ctx, sessionId)
	if err != nil {
		return model.WrapError(model.ErrCreateSession, err)
	}

	// update the ui with the latest match
	// or with the user's base stats
	var initMatch *model.Match
	if len(session.Matches) > 0 {
		initMatch = session.Matches[0]
	} else {
		initMatch = &model.Match{
			UserName:  session.UserName,
			LP:        session.LP,
			MR:        session.MR,
			SessionId: session.Id,
			UserId:    session.UserId,
		}
	}
	t.onNewMatch(ctx, session, initMatch, true)

	ticker := time.NewTicker(30 * time.Second)
	t.forcePollChan = make(chan struct{})
	defer func() {
		t.eventEmitter("stopped-tracking")
		ticker.Stop()
		cancel()
		close(t.forcePollChan)
		t.forcePollChan = nil
	}()

	go func() {
		slog.Info("polling")
		match, err := t.gameTracker.Poll(ctx, session)
		if err != nil {
			cancel()
			return
		}
		t.onNewMatch(ctx, session, match, false)
		for {
			select {
			case <-t.forcePollChan:
				slog.Info("forced poll")
				match, err := t.gameTracker.Poll(ctx, session)
				if err != nil {
					cancel()
					return
				}
				t.onNewMatch(ctx, session, match, false)
			case <-ticker.C:
				slog.Info("polling")
				match, err := t.gameTracker.Poll(ctx, session)
				if err != nil {
					cancel()
					return
				}
				t.onNewMatch(ctx, session, match, false)
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (t *TrackingHandler) onNewMatch(ctx context.Context, session *model.Session, match *model.Match, dry bool) {
	if match == nil {
		return
	}
	t.eventEmitter("match", *match)
	for _, mc := range t.matchChans {
		if mc != nil {
			mc <- *match
		}
	}
	if dry {
		return
	}

	session.LP = match.LP
	session.MR = match.MR
	session.Matches = append([]*model.Match{match}, session.Matches...)

	if err := t.sqlDb.UpdateSession(ctx, session); err != nil {
		slog.Error("update session:", slog.Any("error", err))
		return
	}
	if err := t.sqlDb.SaveMatch(ctx, *match); err != nil {
		slog.Error("save match to database", slog.Any("error", err))
		return
	}
	if err := t.txtDb.SaveMatch(*match); err != nil {
		slog.Error("save to text files:", slog.Any("error", err))
		return
	}
}

func (t *TrackingHandler) StopTracking() {
	t.cancelPolling()
}

func (t *TrackingHandler) SelectGame(game model.GameType) error {
	var username, password string
	switch game {
	case model.GameTypeT8:
		t.gameTracker = t8.NewT8Tracker(t.wavuClient)
	case model.GameTypeSF6:
		t.gameTracker = sf6.NewSF6Tracker(t.cfnClient)
		username = t.cfg.CapIDEmail
		password = t.cfg.CapIDPassword
	default:
		return model.WrapError(model.ErrSelectGame, fmt.Errorf("game does not exist"))
	}

	authChan := make(chan tracker.AuthStatus)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	go t.gameTracker.Authenticate(ctx, username, password, authChan)
	for status := range authChan {
		if status.Err != nil {
			return model.WrapError(model.ErrAuth, status.Err)
		}

		t.eventEmitter("auth-progress", status.Progress)

		if status.Progress >= 100 {
			close(authChan)
			break
		}
	}
	return nil
}

func (t *TrackingHandler) ForcePoll() {
	if t.forcePollChan != nil {
		t.forcePollChan <- struct{}{}
	}
}

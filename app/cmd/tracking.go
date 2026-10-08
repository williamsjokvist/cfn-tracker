package cmd

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"sync"
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

const (
	pollInterval       = 30 * time.Second
	retryBaseDelay     = 2 * time.Second
	retryThrottleDelay = 5 * time.Second
	retryMaxDelay      = 60 * time.Second
	parseFailThreshold = 3
	parseRetryDelay    = 5 * time.Minute
	maxReauthAttempts  = 3
)

// Longer than the SF6 manual login timeout (10 minutes), so the user has time to log in.
var selectGameTimeout = 12 * time.Minute

// Upper bound for a single Poll. Longer than the inner battleLogTimeout (25s) so that
// fires first; this is a safety net so the poll loop always returns.
var pollTimeout = 40 * time.Second

type EventEmitFn func(eventName string, optionalData ...interface{})

type RetryStatus struct {
	Attempt       int    `json:"attempt"`
	NextRetryInMs int64  `json:"nextRetryInMs"`
	Reason        string `json:"reason"`
}

type TrackingHandler struct {
	sqlDb      *sql.Storage
	nosqlDb    *cfgDb.Storage
	txtDb      *txt.Storage
	wavuClient wavu.WavuClient
	cfnClient  cfn.CFNClient
	cfg        *config.BuildConfig
	matchChans []chan model.Match

	mu            sync.Mutex
	cancelPolling context.CancelFunc
	forcePollChan chan struct{}
	gameTracker   tracker.GameTracker
	eventEmitter  EventEmitFn
}

func NewTrackingHandler(wavuClient wavu.WavuClient, cfnClient cfn.CFNClient, sqlDb *sql.Storage, nosqlDb *cfgDb.Storage, txtDb *txt.Storage, cfg *config.BuildConfig, matchChans ...chan model.Match) *TrackingHandler {
	return &TrackingHandler{wavuClient: wavuClient, cfnClient: cfnClient, sqlDb: sqlDb, nosqlDb: nosqlDb, txtDb: txtDb, cfg: cfg, matchChans: matchChans}
}

func (ch *TrackingHandler) SetEventEmitter(eventEmitter EventEmitFn) { ch.eventEmitter = eventEmitter }

func (ch *TrackingHandler) emit(name string, data ...interface{}) {
	if ch.eventEmitter != nil {
		ch.eventEmitter(name, data...)
	}
}

func (ch *TrackingHandler) StartTracking(userCodeInput string, restore bool) error {
	slog.Info("started tracking", slog.String("user_code", userCodeInput), slog.Bool("restoring", restore))
	ctx, cancel := context.WithCancel(context.Background())
	forcePoll := make(chan struct{}, 1)
	ch.mu.Lock()
	ch.cancelPolling, ch.forcePollChan = cancel, forcePoll
	ch.mu.Unlock()
	defer func() {
		cancel()
		ch.mu.Lock()
		if ch.forcePollChan == forcePoll {
			ch.forcePollChan = nil
			ch.cancelPolling = nil
		}
		ch.mu.Unlock()
		ch.emit("stopped-tracking")
	}()

	if ch.gameTracker == nil {
		err := errors.New("game tracker not selected")
		ch.emit("tracking-error", model.FormatError(err))
		return err
	}
	user, err := ch.gameTracker.GetUser(ctx, userCodeInput)
	if err != nil {
		return ch.terminalError(model.WrapError(model.ErrGetUser, err))
	}
	if err := ch.sqlDb.SaveUser(ctx, *user); err != nil {
		return ch.terminalError(model.WrapError(model.ErrSaveUser, err))
	}
	var session *model.Session
	if restore {
		session, err = ch.sqlDb.GetLatestSession(ctx, user.Code)
	} else {
		session, err = ch.sqlDb.CreateSession(ctx, user.Code)
	}
	if err != nil {
		if restore {
			return ch.terminalError(model.WrapError(model.ErrGetLatestSession, err))
		}
		return ch.terminalError(model.WrapError(model.ErrCreateSession, err))
	}
	if session == nil {
		return ch.terminalError(model.ErrCreateSession)
	}
	session.LP, session.MR, session.UserName = user.LP, user.MR, user.DisplayName
	ch.emit("match", model.Match{UserName: session.UserName, LP: session.LP, MR: session.MR, SessionId: session.Id, UserId: session.UserId})

	matchChan := make(chan model.Match)
	pollErrChan := make(chan error, 1)
	go ch.poll(ctx, forcePoll, session, matchChan, pollErrChan)

	if len(session.Matches) > 0 {
		match := *session.Matches[0]
		ch.emit("match", match)
		if !ch.sendMatches(ctx, match) {
			return ctx.Err()
		}
	}
	for match := range matchChan {
		ch.emit("match", match)
		session.LP, session.MR = match.LP, match.MR
		session.Matches = append([]*model.Match{&match}, session.Matches...)
		if err := ch.sqlDb.UpdateSession(ctx, session); err != nil {
			return ch.storageError(cancel, "update session", err)
		}
		if err := ch.sqlDb.SaveMatch(ctx, match); err != nil {
			return ch.storageError(cancel, "save match to database", err)
		}
		if err := ch.txtDb.SaveMatch(match); err != nil {
			return ch.storageError(cancel, "save to text files", err)
		}
	}
	select {
	case err := <-pollErrChan:
		return err
	default:
		return nil
	}
}

func (ch *TrackingHandler) terminalError(err error) error {
	ch.emit("tracking-error", model.FormatError(err))
	return err
}

func (ch *TrackingHandler) storageError(cancel context.CancelFunc, op string, err error) error {
	wrapped := fmt.Errorf("%s: %w", op, err)
	slog.Error(op, slog.Any("error", err))
	ch.emit("tracking-error", model.FormatError(wrapped))
	cancel()
	return wrapped
}

func (ch *TrackingHandler) sendMatches(ctx context.Context, match model.Match) bool {
	for _, out := range ch.matchChans {
		if out != nil {
			if !sendMatch(ctx, out, match) {
				return false
			}
		}
	}
	return true
}

func sendMatch(ctx context.Context, out chan<- model.Match, match model.Match) bool {
	select {
	case out <- match:
		return true
	case <-ctx.Done():
		return false
	}
}

func (ch *TrackingHandler) poll(ctx context.Context, force <-chan struct{}, session *model.Session, matches chan<- model.Match, result chan<- error) {
	defer close(matches)
	ticker := time.NewTicker(pollInterval)
	defer ticker.Stop()
	attempt, parseFailures, authAttempts := 0, 0, 0
	wasFailing := false
	for {
		pollCtx, cancelPoll := context.WithTimeout(ctx, pollTimeout)
		started := time.Now()
		match, err := ch.gameTracker.Poll(pollCtx, session)
		cancelPoll()
		if err == nil {
			if match != nil {
				slog.Info("poll: new match", slog.String("replay_id", match.ReplayID), slog.Duration("took", time.Since(started)))
			} else {
				slog.Info("poll: no new match", slog.Duration("took", time.Since(started)))
			}
			if wasFailing {
				ch.emit("tracking-recovered")
			}
			attempt, parseFailures, authAttempts, wasFailing = 0, 0, 0, false
			if match != nil {
				if !sendMatch(ctx, matches, *match) || !ch.sendMatches(ctx, *match) {
					return
				}
			}
		} else {
			if ctx.Err() != nil {
				return
			}
			class := model.ClassifyPollError(err)
			attempt++
			wasFailing = true
			slog.Error("poll failed", slog.Any("error", err), slog.Int("class", int(class)), slog.Int("attempt", attempt))
			if class == model.ClassFatal {
				ch.emit("tracking-error", model.FormatError(err))
				select {
				case result <- err:
				case <-ctx.Done():
				}
				return
			}
			if class == model.ClassAuth {
				authAttempts++
				if ch.reauthenticate(ctx, attempt) {
					continue
				}
				if authAttempts >= maxReauthAttempts {
					ch.emit("tracking-error", model.FormatError(err))
					select {
					case result <- err:
					case <-ctx.Done():
					}
					return
				}
				delay := ch.retryDelay(err, attempt)
				ch.emit("tracking-retrying", RetryStatus{Attempt: attempt, NextRetryInMs: delay.Milliseconds(), Reason: "errAuth"})
				if !waitFor(ctx, delay) {
					return
				}
				continue
			}
			delay, reason := ch.retryDelay(err, attempt), "errNetwork"
			if class == model.ClassParse {
				parseFailures++
				if parseFailures > parseFailThreshold {
					delay, reason = parseRetryDelay, "errStructureChanged"
				}
			}
			var statusErr *model.HTTPStatusError
			if errors.As(err, &statusErr) && statusErr.StatusCode == 429 {
				reason = "errThrottled"
			}
			ch.emit("tracking-retrying", RetryStatus{Attempt: attempt, NextRetryInMs: delay.Milliseconds(), Reason: reason})
			if !waitFor(ctx, delay) {
				return
			}
			continue
		}
		select {
		case <-ticker.C:
		case <-force:
		case <-ctx.Done():
			return
		}
	}
}

func (ch *TrackingHandler) retryDelay(err error, attempt int) time.Duration {
	base := retryBaseDelay
	var statusErr *model.HTTPStatusError
	if errors.As(err, &statusErr) && statusErr.StatusCode == 429 {
		base = retryThrottleDelay
	}
	max := base
	for i := 1; i < attempt && max < retryMaxDelay; i++ {
		max *= 2
	}
	if max > retryMaxDelay {
		max = retryMaxDelay
	}
	if max <= 0 {
		return 0
	}
	return time.Duration(rand.Int63n(int64(max)))
}

func waitFor(ctx context.Context, delay time.Duration) bool {
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return true
	case <-ctx.Done():
		return false
	}
}

func (ch *TrackingHandler) reauthenticate(ctx context.Context, attempt int) bool {
	statuses := make(chan tracker.AuthStatus, 1)
	go ch.gameTracker.Authenticate(ctx, statuses)
	for {
		select {
		case status, ok := <-statuses:
			if !ok {
				return false
			}
			if status.Err != nil {
				return false
			}
			if status.Action != nil {
				// Show the login prompt in the tracking banner while the user logs in.
				ch.emit("tracking-retrying", RetryStatus{Attempt: attempt, Reason: status.Action.LocalizationKey})
				continue
			}
			if status.Progress >= 100 {
				return true
			}
		case <-ctx.Done():
			return false
		}
	}
}

func (ch *TrackingHandler) StopTracking() {
	ch.mu.Lock()
	cancel := ch.cancelPolling
	ch.mu.Unlock()
	if cancel != nil {
		cancel()
	}
}

func (ch *TrackingHandler) SelectGame(game model.GameType) error {
	switch game {
	case model.GameTypeT8:
		ch.gameTracker = t8.NewT8Tracker(ch.wavuClient)
	case model.GameTypeSF6:
		ch.gameTracker = sf6.NewSF6Tracker(ch.cfnClient)
	default:
		return model.WrapError(model.ErrSelectGame, fmt.Errorf("game does not exist"))
	}
	authChan := make(chan tracker.AuthStatus)
	ctx, cancel := context.WithTimeout(context.Background(), selectGameTimeout)
	defer cancel()
	go ch.gameTracker.Authenticate(ctx, authChan)
	for {
		select {
		case status, ok := <-authChan:
			if !ok {
				return nil
			}
			if status.Err != nil {
				// Return errors that already have a specific localization key as-is.
				// Wrapping them in ErrAuth would hide the actionable message behind
				// a generic "authentication failed".
				var localized *model.FGCTrackerError
				if errors.As(status.Err, &localized) {
					return localized
				}
				return model.WrapError(model.ErrAuth, status.Err)
			}
			if status.Action != nil {
				ch.emit("auth-action-required", *status.Action)
				continue
			}
			ch.emit("auth-progress", status.Progress)
			if status.Progress >= 100 {
				return nil
			}
		case <-ctx.Done():
			return model.WrapError(model.ErrAuth, ctx.Err())
		}
	}
}

func (ch *TrackingHandler) ForcePoll() {
	ch.mu.Lock()
	force, cancel := ch.forcePollChan, ch.cancelPolling
	ch.mu.Unlock()
	if force == nil || cancel == nil {
		return
	}
	select {
	case force <- struct{}{}:
	default:
	}
}

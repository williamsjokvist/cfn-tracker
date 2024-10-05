package t8

import (
	"context"
	"time"
	"log"
	"fmt"
	"net/http"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/utils"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

type T8Tracker struct {
	stopPolling context.CancelFunc
	state       map[string]*model.TrackingState
	sesh        *model.Session
	user        *model.User

	wavuClient *wavu.Client
	sqlDb      *sql.Storage
	txtDb      *txt.Storage
}

var _ tracker.GameTracker = (*T8Tracker)(nil)

func NewT8Tracker(sqlDb *sql.Storage, txtDb *txt.Storage) *T8Tracker {
	wavu := wavu.NewClient()
	return &T8Tracker{
		stopPolling: func() {},
		sqlDb:       sqlDb,
		txtDb:       txtDb,
		state:       make(map[string]*model.TrackingState, 4),
		wavuClient:  &wavu,
	}
}

func (t *T8Tracker) Start(ctx context.Context, polarisId string, restore bool, pollRate time.Duration) error {
	sesh, err := t.sqlDb.CreateSession(ctx, polarisId)
	if err != nil {
		return errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf("failed to create session: %w", err))
	}
	t.sesh = sesh
	t.user = &model.User{
		DisplayName: "Kazuya",
		Code:        polarisId,
	}
	pollCtx, cancelFn := context.WithCancel(ctx)
	t.stopPolling = cancelFn
	go t.poll(pollCtx, polarisId, pollRate)
	wails.EventsEmit(ctx, "cfn-data", model.TrackingState{
		CFN:       polarisId,
		Character: "Kazuya",
	})
	return nil
}

func (t *T8Tracker) poll(ctx context.Context, polarisId string, pollRate time.Duration) {
	i := 0
	didStop := func() bool {
		return utils.SleepOrBreak(pollRate, func() bool {
			select {
			case <-ctx.Done():
				return true
			default:
				return false
			}
		})
	}
	for {
		i++
		log.Println(`polling`, i)
		if didStop() {
			wails.EventsEmit(ctx, `stopped-tracking`)
			break
		}

		replays, err := t.wavuClient.GetReplays(polarisId)
		if err != nil {
			wails.EventsEmit(ctx, `stopped-tracking`)
			t.stopPolling()
			return
		}
		if len(replays) == 0 {
			continue
		}
		log.Println("replays", replays)

		latestReplay := replays[0]
		latestMatch := wavu.ConvWavuReplayToModelMatch(&latestReplay, latestReplay.P2PolarisId == polarisId)
		if len(t.sesh.Matches) > 0 && t.sesh.Matches[0].ReplayID == latestMatch.ReplayID {
			continue
		}
		if latestMatch.Victory {
			latestMatch.Wins = t.sesh.MatchesWon + 1
		} else {
			latestMatch.Losses = t.sesh.MatchesLost + 1
		}

		t.sesh.Matches = append([]*model.Match{&latestMatch}, t.sesh.Matches...)
		err = t.sqlDb.UpdateSession(ctx, t.sesh, latestMatch, t.sesh.Id)
		if err != nil {
			wails.EventsEmit(ctx, "stopped-tracking")
			t.stopPolling()
			return
		}

		trackingState := model.ConvMatchToTrackingState(latestMatch, t.user.Code, t.user.DisplayName)
		trackingState.Log()
		t.txtDb.SaveTrackingState(&trackingState)
		wails.EventsEmit(ctx, "cfn-data", trackingState)
	}
}

func (t *T8Tracker) Stop() {
	t.stopPolling()
}

func (t *T8Tracker) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	statChan <- tracker.AuthStatus{Progress: 100, Err: nil}
}

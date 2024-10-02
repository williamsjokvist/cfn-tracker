package t8

import (
	"context"
	"time"
	"log"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/utils"
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

func (t *T8Tracker) Authenticate(email string, password string, statChan chan tracker.AuthStatus) {
	statChan <- tracker.AuthStatus{Progress: 100, Err: nil}
}

func (t *T8Tracker) Start(ctx context.Context, userCode string, restore bool, pollRate time.Duration) error {
	beginPolling := func() {
		pollCtx, cancelFn := context.WithCancel(ctx)
		t.stopPolling = cancelFn
		go t.poll(pollCtx, 0, pollRate)
	}

	beginPolling()
	return nil
}

func (t *T8Tracker) poll(ctx context.Context, userCode uint64, pollRate time.Duration) {
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

		_, err := t.wavuClient.GetReplays(userCode)
		if err != nil {
			wails.EventsEmit(ctx, `stopped-tracking`)
			t.stopPolling()
			return
		}

	}
}

func (t *T8Tracker) Stop() {
	t.stopPolling()
}

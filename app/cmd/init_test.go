package cmd_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	cfgDb "github.com/williamsjokvist/cfn-tracker/pkg/storage/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

type eventRecord struct {
	Name string
	Data []interface{}
}

var testSuite = struct {
	trackingHandler *cmd.TrackingHandler
}{}

func TestMain(m *testing.M) {
	sqlDb, err := sql.NewStorage(true)
	if err != nil {
		log.Fatalf("init sql storage: %v", err)
	}
	cfgDb, err := cfgDb.NewStorage()
	if err != nil {
		log.Fatalf("init nosql storage: %v", err)
	}
	txtDb, err := txt.NewStorage()
	if err != nil {
		log.Fatalf("init txt storage: %v", err)
	}

	cfg := config.BuildConfig{
		AppVersion:        "4.0.0",
		Headless:          true,
		CapIDEmail:        "test",
		CapIDPassword:     "test",
		BrowserSourcePort: 4242,
	}

	testSuite.trackingHandler = cmd.NewTrackingHandler(
		// todo mock api clients
		wavu.NewClient(),
		cfn.NewClient(nil),
		sqlDb,
		cfgDb,
		txtDb,
		&cfg,
		nil,
	)
	testSuite.trackingHandler.SetEventEmitter(func(eventName string, optionalData ...interface{}) {
		if len(optionalData) > 0 {
			log.Println(fmt.Sprintf("[EVENT] %s", eventName), optionalData[0])
		} else {
			log.Println(fmt.Sprintf("[EVENT] %s", eventName))
		}
	})
	os.Exit(m.Run())
}

func setupTrackingTest(t *testing.T, mockResults ...tracker.MockPollResult) (*cmd.TrackingHandler, *sql.Storage, *model.Session, *[]eventRecord) {
	t.Helper()
	sqlDb, err := sql.NewStorage(true)
	if err != nil {
		t.Fatalf("create in-memory storage: %v", err)
	}
	t.Cleanup(func() { sqlDb.DB().Close() })

	ctx := context.Background()

	user := &model.User{DisplayName: "TestPlayer", Code: "test-user-123"}
	if err := sqlDb.SaveUser(ctx, *user); err != nil {
		t.Fatalf("save user: %v", err)
	}

	session, err := sqlDb.CreateSession(ctx, user.Code)
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	session.UserName = user.DisplayName
	session.LP = 1000
	session.MR = 0

	mock := tracker.NewMockGameTracker(user, mockResults...)

	events := &[]eventRecord{}

	cfg := &config.BuildConfig{}
	handler := cmd.NewTrackingHandler(nil, nil, sqlDb, nil, nil, cfg)
	handler.SetGameTracker(mock)
	handler.SetEventEmitter(func(eventName string, optionalData ...interface{}) {
		*events = append(*events, eventRecord{Name: eventName, Data: optionalData})
	})

	return handler, sqlDb, session, events
}

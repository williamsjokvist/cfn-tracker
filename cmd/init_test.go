package cmd_test

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	cfgDb "github.com/williamsjokvist/cfn-tracker/pkg/storage/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/sf6/cfn"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker/t8/wavu"
)

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
		log.Println(fmt.Sprintf("[EVENT] %s", eventName), optionalData[0])
	})
	os.Exit(m.Run())
}

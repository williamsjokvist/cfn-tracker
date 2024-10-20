package cmd_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/nosql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
)

var testSuite = struct {
	trackingHandler *cmd.TrackingHandler
}{}

func TestMain(t *testing.M) {
	sqlDb, err := sql.NewStorage(true)
	if err != nil {
		log.Fatal("failed to init sql store")
	}
	nosqlDb, err := nosql.NewStorage()
	if err != nil {
		log.Fatal("failed to init nosql store")
	}
	txtDb, err := txt.NewStorage()
	if err != nil {
		log.Fatal("failed to init txt store")
	}

	cfg := config.Config{
		AppVersion:        "4.0.0",
		Headless:          true,
		CapIDEmail:        "test",
		CapIDPassword:     "test",
		BrowserSourcePort: 4242,
	}

	testSuite.trackingHandler = cmd.NewTrackingHandler(nil, sqlDb, nosqlDb, txtDb, &cfg, nil)
	testSuite.trackingHandler.SetEventEmitter(func(eventName string, optionalData ...interface{}) {
		log.Println(fmt.Sprintf("[EVENT] %s", eventName), optionalData[0])
	})
	t.Run()
}

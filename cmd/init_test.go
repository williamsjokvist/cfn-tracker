package cmd_test

import (
	"fmt"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/nosql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
)

var tEnv = struct {
	trackingHandler *cmd.TrackingHandler
}{}

func TestMain(t *testing.T) {
	sqlDb, err := sql.NewStorage()
	if err != nil {
		t.Fatal("failed to init sql store")
	}
	nosqlDb, err := nosql.NewStorage()
	if err != nil {
		t.Fatal("failed to init nosql store")
	}
	txtDb, err := txt.NewStorage()
	if err != nil {
		t.Fatal("failed to init txt store")
	}

	tEnv.trackingHandler = cmd.NewTrackingHandler(nil, sqlDb, nosqlDb, txtDb, nil)
	tEnv.trackingHandler.SetEventEmitter(func(eventName string, optionalData ...interface{}) {
		t.Log(fmt.Sprintf("[EVENT] %s", eventName), optionalData[0])
	})
}

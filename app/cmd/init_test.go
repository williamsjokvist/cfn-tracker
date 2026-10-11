package cmd_test

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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
	sqlDb           *sql.Storage
	cfgDb           *cfgDb.Storage
	txtDb           *txt.Storage
}{}

func TestMain(m *testing.M) {
	sqlDb, err := sql.NewInMemoryStorage()
	if err != nil {
		log.Fatalf("init sql storage: %v", err)
	}
	// A temp dir, so tracking test users aren't remembered in the real config.
	cfgDir, err := os.MkdirTemp("", "cfn-tracker-test")
	if err != nil {
		log.Fatalf("create config dir: %v", err)
	}
	cfgDb, err := cfgDb.NewStorageIn(cfgDir)
	if err != nil {
		log.Fatalf("init nosql storage: %v", err)
	}
	txtDb, err := txt.NewStorage(filepath.Join(cfgDir, "results"))
	if err != nil {
		log.Fatalf("init txt storage: %v", err)
	}

	cfg := config.BuildConfig{
		AppVersion: "4.0.0",
	}

	testSuite.sqlDb, testSuite.cfgDb, testSuite.txtDb = sqlDb, cfgDb, txtDb
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
		log.Println(fmt.Sprintf("[EVENT] %s", eventName), optionalData)
	})
	code := m.Run()
	if err := os.RemoveAll(cfgDir); err != nil {
		log.Printf("remove config dir: %v", err)
	}
	os.Exit(code)
}

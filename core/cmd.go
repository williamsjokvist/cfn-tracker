package core

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/williamsjokvist/cfn-tracker/core/common"
	"github.com/williamsjokvist/cfn-tracker/core/sf6"
	"github.com/williamsjokvist/cfn-tracker/core/sfv"
)

var (
	SteamUsername   string
	SteamPassword   string
	CapIDEmail      string
	CapIDPassword   string
	AppVersion      *version.Version
	RefreshInterval time.Duration = 30 * time.Second
	RunHeadless     bool
)

// The CommandHandler is the interface between the GUI and the core
type CommandHandler struct {
	ctx     context.Context
	tracker common.GameTracker
	browser *common.Browser
}

func NewCommandHandler() *CommandHandler {
	return &CommandHandler{}
}

func (ch *CommandHandler) GetAppVersion() string {
	return AppVersion.Original()
}

func (ch *CommandHandler) StopTracking() {
	ch.tracker.Stop()
}

func (ch *CommandHandler) StartTracking(cfn string, restore bool) {
	err := ch.tracker.Start(cfn, restore, RefreshInterval)
	if err != nil {
		fmt.Println(err)
	}
}

func (ch *CommandHandler) GetMatchHistory() common.MatchHistory {
	mh := ch.tracker.GetMatchHistory()
	return *mh
}

func (ch *CommandHandler) OpenResultsDirectory() {
	var cmd *exec.Cmd

	if runtime.GOOS == `darwin` {
		cmd = exec.Command(`Open`, `./results`)
	} else {
		cmd = exec.Command(`explorer.exe`, `.\results`)
	}

	cmd.Run()
}

func (ch *CommandHandler) GetMatchLog(cfn string) []common.MatchHistory {
	mhLog, err := common.GetLog(cfn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return mhLog
}

func (ch *CommandHandler) GetAvailableLogs() []string {
	loggedCfns, err := common.GetLoggedCFNs()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return loggedCfns
}

func (ch *CommandHandler) DeleteMatchLog(cfn string) {
	err := common.DeleteLog(cfn)
	if err != nil {
		fmt.Println(err)
	}
}

func (ch *CommandHandler) ExportLogToCSV(cfn string) {
	err := common.ExportLog(cfn)
	if err != nil {
		fmt.Println(err)
	}
}

func (ch *CommandHandler) SelectGame(game string) {
	switch game {
	case `sf6`:
		ch.tracker, _ = sf6.MakeSF6Tracker(ch.ctx, ch.browser, CapIDEmail, CapIDPassword)
	case `sfv`:
		ch.tracker, _ = sfv.MakeSFVTracker(ch.ctx, ch.browser, SteamUsername, SteamPassword)
	}
}

func (ch *CommandHandler) ResultsJSONExist() bool {
	_, err := os.Stat(`results/results.json`)
	return !os.IsNotExist(err)
}

func (ch *CommandHandler) StartBrowser(ctx context.Context) {
	ch.ctx = ctx
	ch.browser = common.NewBrowser(ctx, RunHeadless)
	ch.browser.CheckForVersionUpdate(AppVersion)
}

func (ch *CommandHandler) CloseBrowser(ctx context.Context) {
	if ch.browser.Page != nil {
		ch.browser.Page.Browser().Close()
	}
}

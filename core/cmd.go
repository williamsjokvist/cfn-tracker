package core

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/hashicorp/go-version"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/shared"
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
	tracker GameTracker
	browser *shared.Browser
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
	log.Printf(`starting tracking for %s, restoring = %v`, cfn, restore)
	err := ch.tracker.Start(cfn, restore, RefreshInterval)
	if err != nil {
		wails.EventsEmit(ch.ctx, `error-cfn`, err.Error())
		fmt.Println(err)
	}
}

func (ch *CommandHandler) OpenResultsDirectory() {
	switch runtime.GOOS {
	case `darwin`:
		exec.Command(`Open`, `./results`).Run()
	case `windows`:
		exec.Command(`explorer.exe`, `.\results`).Run()
	}
}

func (ch *CommandHandler) GetMatchLog(cfn string) []data.MatchHistory {
	mhLog, err := data.GetLog(cfn)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return mhLog
}

func (ch *CommandHandler) GetAvailableLogs() []data.PlayerInfo {
	loggedCfns, err := data.GetLoggedCFNs()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return loggedCfns
}

func (ch *CommandHandler) GetThemeList() []string {
	themes, err := shared.GetThemeList()
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return themes
}

func (ch *CommandHandler) DeleteMatchLog(cfn string) {
	err := data.DeleteLog(cfn)
	if err != nil {
		fmt.Println(err)
	}
}

func (ch *CommandHandler) ExportLogToCSV(cfn string) {
	err := data.ExportLog(cfn)
	if err != nil {
		fmt.Println(err)
	}
}

func (ch *CommandHandler) SelectGame(game string) {
	switch game {
	case GameTypeSF6.String():
		ch.tracker, _ = MakeSF6Tracker(ch.ctx, ch.browser, CapIDEmail, CapIDPassword)
	case GameTypeSFV.String():
		ch.tracker, _ = MakeSFVTracker(ch.ctx, ch.browser, SteamUsername, SteamPassword)
	}
}

func (ch *CommandHandler) ResultsJSONExist() bool {
	_, err := os.Stat(`results/results.json`)
	return !os.IsNotExist(err)
}

func (ch *CommandHandler) StartBrowser(ctx context.Context) {
	ch.ctx = ctx
	ch.browser = shared.NewBrowser(ctx, RunHeadless)
	ch.browser.CheckForVersionUpdate(AppVersion)
	go shared.Serve(ctx)
}

func (ch *CommandHandler) CloseBrowser(ctx context.Context) {
	if ch.browser.Page != nil {
		ch.browser.Page.Browser().Close()
	}
}

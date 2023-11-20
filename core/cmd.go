package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/browser"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/model"
	"github.com/williamsjokvist/cfn-tracker/core/server"
	"github.com/williamsjokvist/cfn-tracker/core/tracker"
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
	tracker tracker.GameTracker
	browser *browser.Browser
	repo    *data.CFNTrackerRepository
}

func NewCommandHandler(trackerRepo *data.CFNTrackerRepository) *CommandHandler {
	return &CommandHandler{
		repo: trackerRepo,
	}
}

func (ch *CommandHandler) GetAppVersion() string {
	return AppVersion.Original()
}

func (ch *CommandHandler) StopTracking() {
	log.Println(`Stopped tracking`)
	ch.tracker.Stop()
}

func (ch *CommandHandler) StartTracking(cfn string, restore bool) {
	log.Printf(`starting tracking for %s, restoring = %v`, cfn, restore)
	err := ch.tracker.Start(ch.ctx, cfn, restore, RefreshInterval)
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

func (ch *CommandHandler) GetAllMatchesForUser(userId string) []*model.Match {
	matches, err := ch.repo.GetMatches(ch.ctx, 0, userId)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return matches
}

func (ch *CommandHandler) GetUsers() []*model.User {
	users, err := ch.repo.GetUsers(ch.ctx)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return users
}

func (ch *CommandHandler) GetThemeList() ([]string, error) {
	files, err := ioutil.ReadDir(`themes`)
	if err != nil {
		return nil, fmt.Errorf(`read themes directory: %w`, err)
	}
	themes := make([]string, 0, len(files))
	for _, file := range files {
		fileName := file.Name()
		if !strings.Contains(fileName, `.css`) {
			continue
		}
		theme := strings.Split(fileName, `.css`)[0]
		themes = append(themes, theme)
	}
	return themes, nil
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
	case tracker.GameTypeSF6.String():
		ch.tracker, _ = tracker.MakeSF6Tracker(ch.ctx, ch.browser, CapIDEmail, CapIDPassword, ch.repo)
	case tracker.GameTypeSFV.String():
		ch.tracker, _ = tracker.MakeSFVTracker(ch.ctx, ch.browser, SteamUsername, SteamPassword)
	}
}

func (ch *CommandHandler) ResultsJSONExist() bool {
	_, err := os.Stat(`results/results.json`)
	return !os.IsNotExist(err)
}

func (ch *CommandHandler) StartBrowser(ctx context.Context) {
	ch.ctx = ctx
	browser, err := browser.NewBrowser(ctx, RunHeadless)
	if err != nil {
		log.Fatal(err)
		return
	}
	ch.browser = browser

	ch.browser.CheckForVersionUpdate(AppVersion)
	go server.Start(ctx)
}

func (ch *CommandHandler) CloseBrowser(ctx context.Context) {
	if ch.browser.Page != nil {
		ch.browser.Page.Browser().Close()
	}
}

func (ch *CommandHandler) GetTrackingStateUnused() *data.TrackingState {
	return nil
}

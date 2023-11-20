package core

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/browser"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/model"
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

func NewCommandHandler(browser *browser.Browser, trackerRepo *data.CFNTrackerRepository) *CommandHandler {
	return &CommandHandler{
		repo:    trackerRepo,
		browser: browser,
	}
}

// The CommandHandler needs the wails runtime context in order to emit events
func (ch *CommandHandler) AssignRuntimeContext(ctx context.Context) {
	ch.ctx = ctx
}

func (ch *CommandHandler) GetAppVersion() string {
	return AppVersion.Original()
}

func (ch *CommandHandler) StopTracking() {
	log.Println(`Stopped tracking`)
	ch.tracker.Stop()
}

func (ch *CommandHandler) StartTracking(cfn string, restore bool) error {
	log.Printf(`Starting tracking for %s, restoring = %v`, cfn, restore)
	err := ch.tracker.Start(ch.ctx, cfn, restore, RefreshInterval)
	if err != nil {
		// TODO: remove and use a notification hook in the frontend instead
		wails.EventsEmit(ch.ctx, `error-cfn`, err.Error())
		log.Println(err)
	}
	return err
}

func (ch *CommandHandler) OpenResultsDirectory() {
	switch runtime.GOOS {
	case `darwin`:
		exec.Command(`Open`, `./results`).Run()
	case `windows`:
		exec.Command(`explorer.exe`, `.\results`).Run()
	}
}

func (ch *CommandHandler) GetAllMatchesForUser(userId string) ([]*model.Match, error) {
	matches, err := ch.repo.GetMatches(ch.ctx, 0, userId)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return matches, err
}

func (ch *CommandHandler) GetUsers() ([]*model.User, error) {
	users, err := ch.repo.GetUsers(ch.ctx)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return users, nil
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

func (ch *CommandHandler) DeleteMatchLog(cfn string) error {
	err := data.DeleteLog(cfn)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ch *CommandHandler) ExportLogToCSV(cfn string) error {
	err := data.ExportLog(cfn)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func (ch *CommandHandler) SelectGame(game string) {
	switch game {
	case tracker.GameTypeSF6.String():
		ch.tracker, _ = tracker.MakeSF6Tracker(ch.ctx, ch.browser, CapIDEmail, CapIDPassword, ch.repo)
	case tracker.GameTypeSFV.String():
		ch.tracker, _ = tracker.MakeSFVTracker(ch.ctx, ch.browser, SteamUsername, SteamPassword)
	}
}

func (ch *CommandHandler) GetTrackingStateUnused() *data.TrackingState {
	return nil
}

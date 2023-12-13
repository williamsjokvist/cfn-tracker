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

func (ch *CommandHandler) GetSessions(userId string) ([]*model.Session, error) {
	sessions, err := ch.repo.GetSessions(ch.ctx, userId, 0, 0)
	if err != nil {
		log.Println(err)
	}
	return sessions, err
}

func (ch *CommandHandler) GetMatches(sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	matches, err := ch.repo.GetMatches(ch.ctx, sessionId, userId, limit, offset)
	if err != nil {
		log.Println(err)
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

func (ch *CommandHandler) SelectGame(game string) error {
	var err error
	switch game {
	case tracker.GameTypeSF6.String():
		ch.tracker, err = tracker.MakeSF6Tracker(ch.ctx, ch.browser, CapIDEmail, CapIDPassword, ch.repo)
	case tracker.GameTypeSFV.String():
		ch.tracker, err = tracker.MakeSFVTracker(ch.ctx, ch.browser, SteamUsername, SteamPassword)
	}
	return err
}

func (ch *CommandHandler) GetTrackingStateUnused() *model.TrackingState {
	return nil
}

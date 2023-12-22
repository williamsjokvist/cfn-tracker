package core

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/williamsjokvist/cfn-tracker/core/browser"
	"github.com/williamsjokvist/cfn-tracker/core/config"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/errorsx"
	"github.com/williamsjokvist/cfn-tracker/core/i18n"
	"github.com/williamsjokvist/cfn-tracker/core/i18n/locales"
	"github.com/williamsjokvist/cfn-tracker/core/model"
	"github.com/williamsjokvist/cfn-tracker/core/tracker"
)

// The CommandHandler is the interface between the GUI and the core
type CommandHandler struct {
	ctx     context.Context
	tracker tracker.GameTracker
	browser *browser.Browser
	repo    *data.CFNTrackerRepository
	cfg     *config.Config
}

func NewCommandHandler(browser *browser.Browser, trackerRepo *data.CFNTrackerRepository, cfg *config.Config) *CommandHandler {
	return &CommandHandler{
		repo:    trackerRepo,
		browser: browser,
		cfg:     cfg,
	}
}

// The CommandHandler needs the wails runtime context in order to emit events
func (ch *CommandHandler) SetContext(ctx context.Context) {
	ch.ctx = ctx
}

func (ch *CommandHandler) GetAppVersion() string {
	return ch.cfg.AppVersion
}

func (ch *CommandHandler) GetTranslation(locale string) (*locales.Localization, error) {
	lng, err := i18n.GetTranslation(locale)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusNotFound, fmt.Errorf(`failed to get translation %w`, err))
		}
		return nil, err
	}
	return lng, nil
}

func (ch *CommandHandler) GetSupportedLanguages() []string {
	return i18n.GetSupportedLanguages()
}

func (ch *CommandHandler) CheckForUpdate() (bool, error) {
	currentVersion, err := version.NewVersion(ch.cfg.AppVersion)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf(`failed to parse current app version: %w`, err)
	}
	latestVersion, err := ch.browser.GetLatestAppVersion()
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf(`failed to check for update: %w`, err)
	}

	hasUpdate := currentVersion.LessThan(latestVersion)
	log.Println(`Has update: `, hasUpdate, `. Current: `, currentVersion.String(), ` Latest: `, latestVersion.String())
	return hasUpdate, nil
}

func (ch *CommandHandler) StopTracking() {
	log.Println(`Stopped tracking`)
	ch.tracker.Stop()
}

func (ch *CommandHandler) StartTracking(cfn string, restore bool) error {
	log.Printf(`Starting tracking for %s, restoring = %v`, cfn, restore)
	err := ch.tracker.Start(ch.ctx, cfn, restore, 30*time.Second)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusInternalServerError, fmt.Errorf(`failed to start tracking %w`, err))
		}
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
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusNotFound, fmt.Errorf(`failed to get sessions %w`, err))
		}
	}
	return sessions, err
}

func (ch *CommandHandler) GetMatches(sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	matches, err := ch.repo.GetMatches(ch.ctx, sessionId, userId, limit, offset)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusNotFound, fmt.Errorf(`failed to get matches %w`, err))
		}
	}
	return matches, err
}

func (ch *CommandHandler) GetUsers() ([]*model.User, error) {
	users, err := ch.repo.GetUsers(ch.ctx)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusNotFound, fmt.Errorf(`failed to get users %w`, err))
		}
	}
	return users, err
}

func (ch *CommandHandler) FetchPlayer(code string) (*model.Player, error) {
	player, err := ch.tracker.FetchPlayer(ch.ctx, code)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusNotFound, fmt.Errorf(`failed to fetch player %w`, err))
		}
	}
	return player, err
}

func (ch *CommandHandler) SearchPlayers(query string) ([]model.Player, error) {
	players, err := ch.tracker.SearchPlayers(ch.ctx, query)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusNotFound, fmt.Errorf(`failed to search players %w`, err))
		}
	}
	return players, err
}

func (ch *CommandHandler) GetThemeList() ([]string, error) {
	files, err := ioutil.ReadDir(`themes`)
	if err != nil {

		log.Println(err)
		return nil, errorsx.NewError(http.StatusInternalServerError, errors.New("failed to read themes directory"))
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
		ch.tracker, err = tracker.MakeSF6Tracker(ch.ctx, ch.cfg, ch.browser, ch.repo)
	case tracker.GameTypeSFV.String():
		ch.tracker, err = tracker.MakeSFVTracker(ch.ctx, ch.cfg, ch.browser)
	}

	if err != nil {
		log.Println(err)
		if !errorsx.ContainsTrackingError(err) {
			err = errorsx.NewError(http.StatusInternalServerError, fmt.Errorf(`failed to select game %w`, err))
		}
	}
	return err
}

func (ch *CommandHandler) SaveLocale(locale string) error {
	return ch.repo.SaveLocale(locale)
}

func (ch *CommandHandler) GetGuiConfig() (*model.GuiConfig, error) {
	return ch.repo.GetGuiConfig()
}

func (ch *CommandHandler) SaveSidebarMinimized(sidebarMinified bool) error {
	return ch.repo.SaveSidebarMinimized(sidebarMinified)
}

func (ch *CommandHandler) GetTrackingStateUnused() *model.TrackingState {
	return nil
}

func (ch *CommandHandler) GetFrontendErrorModelUnused() *errorsx.FrontEndError {
	return nil
}

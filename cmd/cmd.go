package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/hashicorp/go-version"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/data"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/i18n"
	"github.com/williamsjokvist/cfn-tracker/pkg/i18n/locales"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/server"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
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
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get translation %w`, err))
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
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to start tracking %w`, err))
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
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get sessions %w`, err))
		}
	}
	return sessions, err
}

func (ch *CommandHandler) GetMatches(sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	matches, err := ch.repo.GetMatches(ch.ctx, sessionId, userId, limit, offset)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get matches %w`, err))
		}
	}
	return matches, err
}

func (ch *CommandHandler) GetUsers() ([]*model.User, error) {
	users, err := ch.repo.GetUsers(ch.ctx)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get users %w`, err))
		}
	}
	return users, err
}

func (ch *CommandHandler) GetThemes() ([]model.Theme, error) {
	// get internal themes
	internalThemes := server.GetInternalThemes()

	// get custom themes
	files, err := os.ReadDir(`themes`)
	if err != nil {
		return internalThemes, nil
	}
	customThemes := make([]model.Theme, 0, len(files))
	for _, file := range files {
		fileName := file.Name()

		if !strings.Contains(fileName, `.css`) {
			continue
		}
		css, err := os.ReadFile(fmt.Sprintf(`themes/%s`, fileName))
		if err != nil {
			log.Println(err)
			return nil, errorsx.NewFormattedError(http.StatusInternalServerError, errors.New("failed to read theme css"))
		}
		name := strings.Split(fileName, `.css`)[0]

		customThemes = append(customThemes, model.Theme{
			Name: name,
			CSS:  string(css),
		})
	}

	combinedThemes := append(customThemes, internalThemes...)
	return combinedThemes, nil
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
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusInternalServerError, fmt.Errorf(`failed to select game %w`, err))
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

func (ch *CommandHandler) SaveTheme(theme model.ThemeName) error {
	return ch.repo.SaveTheme(theme)
}

func (ch *CommandHandler) GetTrackingStateUnused() *model.TrackingState {
	return nil
}

func (ch *CommandHandler) GetFormattedErrorModelUnused() *errorsx.FormattedError {
	return nil
}

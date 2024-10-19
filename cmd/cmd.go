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

	"github.com/hashicorp/go-version"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/i18n"
	"github.com/williamsjokvist/cfn-tracker/pkg/i18n/locales"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/server"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/nosql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/update/github"
)

type EventEmitFn func(eventName string, optionalData ...interface{})

type CmdHandler interface {
	SetEventEmitter(eventEmitter EventEmitFn)
}

// The CommandHandler is the interface between the GUI and the core
type CommandHandler struct {
	eventEmitter EventEmitFn

	browser      *browser.Browser
	githubClient github.GithubClient

	sqlDb   *sql.Storage
	nosqlDb *nosql.Storage
	txtDb   *txt.Storage

	cfg *config.Config
}

var _ CmdHandler = (*CommandHandler)(nil)

func NewCommandHandler(githubClient github.GithubClient, browser *browser.Browser, sqlDb *sql.Storage, nosqlDb *nosql.Storage, txtDb *txt.Storage, cfg *config.Config) *CommandHandler {
	return &CommandHandler{
		sqlDb:        sqlDb,
		nosqlDb:      nosqlDb,
		txtDb:        txtDb,
		browser:      browser,
		githubClient: githubClient,
		cfg:          cfg,
	}
}

func (ch *CommandHandler) SetEventEmitter(eventEmitter EventEmitFn) {
	ch.eventEmitter = eventEmitter
}

func (ch *CommandHandler) CheckForUpdate() (bool, error) {
	currentVersion, err := version.NewVersion(ch.cfg.AppVersion)
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf(`failed to parse current app version: %w`, err)
	}
	latestVersion, err := ch.githubClient.GetLatestAppVersion()
	if err != nil {
		log.Println(err)
		return false, fmt.Errorf(`failed to check for update: %w`, err)
	}

	hasUpdate := currentVersion.LessThan(latestVersion)
	log.Println(`Has update: `, hasUpdate, `. Current: `, currentVersion.String(), ` Latest: `, latestVersion.String())
	return hasUpdate, nil
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

func (ch *CommandHandler) GetAppVersion() string {
	return ch.cfg.AppVersion
}

func (ch *CommandHandler) OpenResultsDirectory() {
	switch runtime.GOOS {
	case `darwin`:
		if err := exec.Command(`Open`, `./results`).Run(); err != nil {
			log.Println(err)
		}
	case `windows`:
		if err := exec.Command(`explorer.exe`, `.\results`).Run(); err != nil {
			log.Println(err)
		}
	}
}

func (ch *CommandHandler) GetSessions(userId string) ([]*model.Session, error) {
	sessions, err := ch.sqlDb.GetSessions(context.Background(), userId, 0, 0)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get sessions %w`, err))
		}
	}
	return sessions, err
}

func (ch *CommandHandler) GetMatches(sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	matches, err := ch.sqlDb.GetMatches(context.Background(), sessionId, userId, limit, offset)
	if err != nil {
		log.Println(err)
		if !errorsx.ContainsFormattedError(err) {
			err = errorsx.NewFormattedError(http.StatusNotFound, fmt.Errorf(`failed to get matches %w`, err))
		}
	}
	return matches, err
}

func (ch *CommandHandler) GetUsers() ([]*model.User, error) {
	users, err := ch.sqlDb.GetUsers(context.Background())
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

func (ch *CommandHandler) GetSupportedLanguages() []string {
	return i18n.GetSupportedLanguages()
}

func (ch *CommandHandler) SaveLocale(locale string) error {
	return ch.nosqlDb.SaveLocale(locale)
}

func (ch *CommandHandler) GetGuiConfig() (*model.GuiConfig, error) {
	return ch.nosqlDb.GetGuiConfig()
}

func (ch *CommandHandler) SaveSidebarMinimized(sidebarMinified bool) error {
	return ch.nosqlDb.SaveSidebarMinimized(sidebarMinified)
}

func (ch *CommandHandler) SaveTheme(theme model.ThemeName) error {
	return ch.nosqlDb.SaveTheme(theme)
}

func (ch *CommandHandler) GetFormattedErrorModelUnused() *errorsx.FormattedError {
	return nil
}

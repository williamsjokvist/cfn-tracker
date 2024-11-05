package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hashicorp/go-version"

	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/i18n"
	"github.com/williamsjokvist/cfn-tracker/pkg/i18n/locales"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"github.com/williamsjokvist/cfn-tracker/pkg/server"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/nosql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
	"github.com/williamsjokvist/cfn-tracker/pkg/update/github"
)

// The CommandHandler is the interface between the GUI and the core
type CommandHandler struct {
	githubClient github.GithubClient

	sqlDb   *sql.Storage
	nosqlDb *nosql.Storage

	cfg *config.Config
}

func NewCommandHandler(githubClient github.GithubClient, sqlDb *sql.Storage, nosqlDb *nosql.Storage, txtDb *txt.Storage, cfg *config.Config) *CommandHandler {
	return &CommandHandler{
		sqlDb:        sqlDb,
		nosqlDb:      nosqlDb,
		githubClient: githubClient,
		cfg:          cfg,
	}
}

func (ch *CommandHandler) CheckForUpdate() (bool, error) {
	currentVersion, err := version.NewVersion(ch.cfg.AppVersion)
	if err != nil {
		return false, model.NewError(model.ErrCheckForUpdate, err)
	}
	latestVersion, err := ch.githubClient.GetLatestAppVersion()
	if err != nil {
		return false, model.NewError(model.ErrCheckForUpdate, err)
	}
	return currentVersion.LessThan(latestVersion), nil
}

func (ch *CommandHandler) GetTranslation(locale string) (*locales.Localization, error) {
	lng, err := i18n.GetTranslation(locale)
	if err != nil {
		return nil, model.NewError(model.ErrGetTranslations, err)
	}
	return lng, nil
}

func (ch *CommandHandler) GetAppVersion() string {
	return ch.cfg.AppVersion
}

func (ch *CommandHandler) OpenResultsDirectory() error {
	switch runtime.GOOS {
	case `darwin`:
		if err := exec.Command(`Open`, `./results`).Run(); err != nil {
			return model.NewError(model.ErrOpenResultsDirectory, err)
		}
	case `windows`:
		if err := exec.Command(`explorer.exe`, `.\results`).Run(); err != nil {
			return model.NewError(model.ErrOpenResultsDirectory, err)
		}
	}
	return nil
}

func (ch *CommandHandler) GetSessions(userId, date string, limit uint8, offset uint16) ([]*model.Session, error) {
	sessions, err := ch.sqlDb.GetSessions(context.Background(), userId, date, limit, offset)
	if err != nil {
		return nil, model.NewError(model.ErrGetSessions, err)
	}
	return sessions, nil
}

func (ch *CommandHandler) GetSessionsStatistics(userId string) (*model.SessionsStatistics, error) {
	sessionStatistics, err := ch.sqlDb.GetSessionsStatistics(context.Background(), userId)
	if err != nil {
		return nil, model.NewError(model.ErrGetSessionStatistics, err)
	}
	return sessionStatistics, nil
}

func (ch *CommandHandler) GetMatches(sessionId uint16, userId string, limit uint8, offset uint16) ([]*model.Match, error) {
	matches, err := ch.sqlDb.GetMatches(context.Background(), sessionId, userId, limit, offset)
	if err != nil {
		return nil, model.NewError(model.ErrGetMatches, err)
	}
	return matches, nil
}

func (ch *CommandHandler) GetUsers() ([]*model.User, error) {
	users, err := ch.sqlDb.GetUsers(context.Background())
	if err != nil {
		return nil, model.NewError(model.ErrGetUser, err)
	}
	return users, nil
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
			return nil, model.NewError(model.ErrReadThemeCSS, err)
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
	if err := ch.nosqlDb.SaveLocale(locale); err != nil {
		return model.NewError(model.ErrSaveLocale, err)
	}
	return nil
}

func (ch *CommandHandler) GetGuiConfig() (*model.GuiConfig, error) {
	guiCfg, err := ch.nosqlDb.GetGuiConfig()
	if err != nil {
		return nil, model.NewError(model.ErrGetGUIConfig, err)
	}
	return guiCfg, nil
}

func (ch *CommandHandler) SaveSidebarMinimized(sidebarMinified bool) error {
	if err := ch.nosqlDb.SaveSidebarMinimized(sidebarMinified); err != nil {
		return model.NewError(model.ErrSaveSidebarMinimized, err)
	}
	return nil
}

func (ch *CommandHandler) SaveTheme(theme model.ThemeName) error {
	if err := ch.nosqlDb.SaveTheme(theme); err != nil {
		return model.NewError(model.ErrSaveTheme, err)
	}
	return nil
}

func (ch *CommandHandler) GetFGCTrackerErrorModelUnused() *model.FGCTrackerError {
	return nil
}

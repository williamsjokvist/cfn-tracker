package model

import (
	"errors"
	"fmt"
)

type ErrorLocalizationKey string

const (
	tKeyErrSelectGame           ErrorLocalizationKey = "errSelectGame"
	tKeyErrAuth                 ErrorLocalizationKey = "errAuth"
	tKeyErrGetLatestSession     ErrorLocalizationKey = "errGetLatestSession"
	tKeyErrGetUser              ErrorLocalizationKey = "errGetUser"
	tKeyErrGetMatches           ErrorLocalizationKey = "errGetMatches"
	tKeyErrSaveLocale           ErrorLocalizationKey = "errSaveLocale"
	tKeyErrCheckForUpdate       ErrorLocalizationKey = "errCheckForUpdate"
	tKeyErrGetGuiConfig         ErrorLocalizationKey = "errGetGuiConfig"
	tKeyErrSaveTheme            ErrorLocalizationKey = "errSaveTheme"
	tKeyErrSaveUser             ErrorLocalizationKey = "errSaveUser"
	tKeyErrSaveSidebarMinimized ErrorLocalizationKey = "errSaveSidebarMinimized"
	tKeyErrGetSessions          ErrorLocalizationKey = "errGetSessions"
	tKeyErrGetTranslations      ErrorLocalizationKey = "errGetTranslations"
	tKeyErrGetSessionStatistics ErrorLocalizationKey = "errGetSessionStatistics"
	tKeyErrCreateSession        ErrorLocalizationKey = "errCreateSession"
	tKeyErrOpenResultsDirectory ErrorLocalizationKey = "errOpenResultsDirectory"
	tKeyErrReadThemeCSS         ErrorLocalizationKey = "errReadThemeCSS"
)

var AllErrors = []struct {
	Value  ErrorLocalizationKey
	TSName string
}{
	{tKeyErrSelectGame, "errSelectGame"},
	{tKeyErrAuth, "errAuth"},
	{tKeyErrGetLatestSession, "errGetLatestSession"},
	{tKeyErrGetUser, "errGetUser"},
	{tKeyErrGetMatches, "errGetMatches"},
	{tKeyErrSaveLocale, "errSaveLocale"},
	{tKeyErrCheckForUpdate, "errCheckForUpdate"},
	{tKeyErrGetGuiConfig, "errGetGuiConfig"},
	{tKeyErrSaveTheme, "errSaveTheme"},
	{tKeyErrSaveUser, "errSaveUser"},
	{tKeyErrSaveSidebarMinimized, "errSaveSidebarMinimized"},
	{tKeyErrGetSessions, "errGetSessions"},
	{tKeyErrGetTranslations, "errGetTranslations"},
	{tKeyErrGetSessionStatistics, "errGetSessionStatistics"},
	{tKeyErrCreateSession, "errCreateSession"},
	{tKeyErrOpenResultsDirectory, "errOpenResultsDirectory"},
	{tKeyErrReadThemeCSS, "errReadThemeCSS"},
}

var (
	ErrSelectGame           = newError(tKeyErrSelectGame, errors.New("select game"))
	ErrAuth                 = newError(tKeyErrAuth, errors.New("unauthenticated"))
	ErrGetLatestSession     = newError(tKeyErrGetLatestSession, errors.New("get latest session"))
	ErrGetUser              = newError(tKeyErrGetUser, errors.New("get user"))
	ErrGetMatches           = newError(tKeyErrGetMatches, errors.New("get matches"))
	ErrSaveLocale           = newError(tKeyErrSaveLocale, errors.New("save locale"))
	ErrCheckForUpdate       = newError(tKeyErrCheckForUpdate, errors.New("check for update"))
	ErrGetGUIConfig         = newError(tKeyErrGetGuiConfig, errors.New("get gui config"))
	ErrSaveTheme            = newError(tKeyErrSaveTheme, errors.New("save theme"))
	ErrSaveUser             = newError(tKeyErrSaveUser, errors.New("save user"))
	ErrSaveSidebarMinimized = newError(tKeyErrSaveSidebarMinimized, errors.New("save sidebar"))
	ErrGetSessions          = newError(tKeyErrGetSessions, errors.New("get sessions"))
	ErrGetTranslations      = newError(tKeyErrGetTranslations, errors.New("get translations"))
	ErrGetSessionStatistics = newError(tKeyErrGetSessionStatistics, errors.New("get session statistics"))
	ErrCreateSession        = newError(tKeyErrCreateSession, errors.New("create session"))
	ErrOpenResultsDirectory = newError(tKeyErrOpenResultsDirectory, errors.New("open results directory"))
	ErrReadThemeCSS         = newError(tKeyErrReadThemeCSS, errors.New("read theme css"))
)

type FGCTrackerError struct {
	LocalizationKey ErrorLocalizationKey `json:"localizationKey"`
	Message         string               `json:"message"`
	InnerError      error
}

func NewError(fmtErr *FGCTrackerError, err error) *FGCTrackerError {
	return &FGCTrackerError{
		LocalizationKey: fmtErr.LocalizationKey,
		Message:         fmtErr.Message,
		InnerError:      fmt.Errorf("%w: %w", fmtErr.InnerError, err),
	}
}

func newError(key ErrorLocalizationKey, err error) *FGCTrackerError {
	return &FGCTrackerError{
		LocalizationKey: key,
		InnerError:      err,
	}
}

func (e *FGCTrackerError) Error() string {
	return e.InnerError.Error()
}

func (e *FGCTrackerError) Unwrap() error {
	return e.InnerError
}

func ContainsFGCTrackerError(err error) bool {
	var trackingErr *FGCTrackerError
	return errors.As(err, &trackingErr)
}

func FormatError(err error) any {
	var formattedErr *FGCTrackerError
	var message string
	var localizationKey ErrorLocalizationKey
	if errors.As(err, &formattedErr) {
		localizationKey = formattedErr.LocalizationKey
		message = formattedErr.InnerError.Error()
	}

	return struct {
		LocalizationKey ErrorLocalizationKey `json:"localizationKey"`
		Message         string               `json:"message"`
	}{
		LocalizationKey: localizationKey,
		Message:         message,
	}
}

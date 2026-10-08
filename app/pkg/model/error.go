package model

import (
	"errors"
	"fmt"
)

type ErrorLocalizationKey string

const (
	tKeyErrUnknown               ErrorLocalizationKey = "errUnknown"
	tKeyErrSelectGame            ErrorLocalizationKey = "errSelectGame"
	tKeyErrAuth                  ErrorLocalizationKey = "errAuth"
	tKeyErrGetLatestSession      ErrorLocalizationKey = "errGetLatestSession"
	tKeyErrGetUser               ErrorLocalizationKey = "errGetUser"
	tKeyErrGetMatches            ErrorLocalizationKey = "errGetMatches"
	tKeyErrSaveLocale            ErrorLocalizationKey = "errSaveLocale"
	tKeyErrCheckForUpdate        ErrorLocalizationKey = "errCheckForUpdate"
	tKeyErrGetGuiConfig          ErrorLocalizationKey = "errGetGuiConfig"
	tKeyErrSaveTheme             ErrorLocalizationKey = "errSaveTheme"
	tKeyErrSaveUser              ErrorLocalizationKey = "errSaveUser"
	tKeyErrSaveSidebar           ErrorLocalizationKey = "errSaveSidebar"
	tKeyErrSaveLogFile           ErrorLocalizationKey = "errSaveLogFile"
	tKeyErrGetSessions           ErrorLocalizationKey = "errGetSessions"
	tKeyErrGetTranslations       ErrorLocalizationKey = "errGetTranslations"
	tKeyErrGetSessionStatistics  ErrorLocalizationKey = "errGetSessionStatistics"
	tKeyErrCreateSession         ErrorLocalizationKey = "errCreateSession"
	tKeyErrOpenResultsDirectory  ErrorLocalizationKey = "errOpenResultsDirectory"
	tKeyErrReadThemeCSS          ErrorLocalizationKey = "errReadThemeCSS"
	tKeyErrAuthManualLoginFailed ErrorLocalizationKey = "errAuthManualLoginFailed"
)

var AllErrorKeys = []struct {
	Value  ErrorLocalizationKey
	TSName string
}{
	{tKeyErrUnknown, string(tKeyErrUnknown)},
	{tKeyErrSelectGame, string(tKeyErrSelectGame)},
	{tKeyErrAuth, string(tKeyErrAuth)},
	{tKeyErrGetLatestSession, string(tKeyErrGetLatestSession)},
	{tKeyErrGetUser, string(tKeyErrGetUser)},
	{tKeyErrGetMatches, string(tKeyErrGetMatches)},
	{tKeyErrSaveLocale, string(tKeyErrSaveLocale)},
	{tKeyErrCheckForUpdate, string(tKeyErrCheckForUpdate)},
	{tKeyErrGetGuiConfig, string(tKeyErrGetGuiConfig)},
	{tKeyErrSaveTheme, string(tKeyErrSaveTheme)},
	{tKeyErrSaveUser, string(tKeyErrSaveUser)},
	{tKeyErrSaveSidebar, string(tKeyErrSaveSidebar)},
	{tKeyErrSaveLogFile, string(tKeyErrSaveLogFile)},
	{tKeyErrGetSessions, string(tKeyErrGetSessions)},
	{tKeyErrGetTranslations, string(tKeyErrGetTranslations)},
	{tKeyErrGetSessionStatistics, string(tKeyErrGetSessionStatistics)},
	{tKeyErrCreateSession, string(tKeyErrCreateSession)},
	{tKeyErrOpenResultsDirectory, string(tKeyErrOpenResultsDirectory)},
	{tKeyErrReadThemeCSS, string(tKeyErrReadThemeCSS)},
	{tKeyErrAuthManualLoginFailed, string(tKeyErrAuthManualLoginFailed)},
}

var (
	ErrUnknown              = newError(tKeyErrUnknown, errors.New("unknown error"))
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
	ErrSaveSidebar          = newError(tKeyErrSaveSidebar, errors.New("save sidebar"))
	ErrSaveLogFile          = newError(tKeyErrSaveLogFile, errors.New("save log file setting"))
	ErrGetSessions          = newError(tKeyErrGetSessions, errors.New("get sessions"))
	ErrGetTranslations      = newError(tKeyErrGetTranslations, errors.New("get translations"))
	ErrGetSessionStatistics = newError(tKeyErrGetSessionStatistics, errors.New("get session statistics"))
	ErrCreateSession        = newError(tKeyErrCreateSession, errors.New("create session"))
	ErrOpenResultsDirectory = newError(tKeyErrOpenResultsDirectory, errors.New("open results directory"))
	ErrReadThemeCSS         = newError(tKeyErrReadThemeCSS, errors.New("read theme css"))

	ErrAuthManualLoginFailed = newError(tKeyErrAuthManualLoginFailed, errors.New("manual login did not establish a buckler session"))
)

type FGCTrackerError struct {
	LocalizationKey ErrorLocalizationKey `json:"localizationKey"`
	Message         string               `json:"message"`
	InnerError      error                `json:"-"`
}

func WrapError(fmtErr *FGCTrackerError, err error) *FGCTrackerError {
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
	// FormatError leaves InnerError nil, and its result can be logged with %v,
	// so fall back instead of panicking.
	if e.InnerError != nil {
		return e.InnerError.Error()
	}
	if e.Message != "" {
		return e.Message
	}
	return string(e.LocalizationKey)
}

func (e *FGCTrackerError) Unwrap() error {
	return e.InnerError
}

func FormatError(err error) any {
	var formattedErr *FGCTrackerError
	if errors.As(err, &formattedErr) {
		return &FGCTrackerError{
			LocalizationKey: formattedErr.LocalizationKey,
			Message:         formattedErr.InnerError.Error(),
		}
	}
	return &FGCTrackerError{
		LocalizationKey: tKeyErrUnknown,
		Message:         err.Error(),
	}
}

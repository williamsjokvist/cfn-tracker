package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type GameTracker interface {
	GetUser(ctx context.Context, userId string) (*model.User, error)
	Poll(ctx context.Context, session *model.Session) (*model.Match, error)
	Authenticate(ctx context.Context, email string, password string, statusChan chan AuthStatus)
}

type AuthStatus struct {
	Progress int
	Err      error
	Action   *AuthAction
}

type AuthAction struct {
	LocalizationKey string `json:"localizationKey"`
	SecondsLeft     int    `json:"secondsLeft"`
}

func (s *AuthStatus) WithProgress(progress int) *AuthStatus {
	s.Progress = progress
	return s
}

func (s *AuthStatus) WithError(err error) *AuthStatus {
	s.Err = err
	return s
}

func (s *AuthStatus) WithAction(key string, secondsLeft int) *AuthStatus {
	s.Action = &AuthAction{LocalizationKey: key, SecondsLeft: secondsLeft}
	return s
}

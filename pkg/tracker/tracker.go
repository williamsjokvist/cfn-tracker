package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type GameTracker interface {
	GetUser(ctx context.Context, userId string) (*model.User, error)
	Poll(ctx context.Context, cancel context.CancelFunc, session *model.Session, onNewMatch func(model.Match))
	Authenticate(ctx context.Context, email string, password string, statusChan chan AuthStatus)
}

type AuthStatus struct {
	Progress int
	Err      error
}

func (s *AuthStatus) WithProgress(progress int) *AuthStatus {
	s.Progress = progress
	return s
}

func (s *AuthStatus) WithError(err error) *AuthStatus {
	s.Err = err
	return s
}

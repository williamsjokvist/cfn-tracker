package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type GameTracker interface {
	InitFn(ctx context.Context, polarisId string, restore bool) (*model.Session, error)
	PollFn(ctx context.Context, session *model.Session, matchChan chan model.Match, cancel context.CancelFunc)
	Authenticate(email string, password string, statusChan chan AuthStatus)
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

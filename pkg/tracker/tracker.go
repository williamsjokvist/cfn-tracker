package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type GameTracker interface {
	Init(ctx context.Context, polarisId string, restore bool) (*model.Session, error)
	Poll(ctx context.Context, cancel context.CancelFunc, session *model.Session, matchChan chan model.Match)
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

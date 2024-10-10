package tracker

import (
	"context"
	"time"
)

type GameTracker interface {
	Start(ctx context.Context, cfn string, restore bool, refreshInterval time.Duration) error
	Authenticate(email string, password string, statusChan chan AuthStatus)
	Stop()
	ForcePoll()
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

type GameType uint8

const (
	GameTypeUndefined GameType = iota
	GameTypeSFV
	GameTypeSF6
	GameTypeT8
)

func (s GameType) String() string {
	switch s {
	case GameTypeSFV:
		return `sfv`
	case GameTypeSF6:
		return `sf6`
	case GameTypeT8:
		return `t8`
	case GameTypeUndefined:
		return `undefined`
	}
	return `unknown`
}

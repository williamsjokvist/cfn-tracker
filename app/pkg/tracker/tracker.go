package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type GameTracker interface {
	GetUser(ctx context.Context, userId string) (*model.User, error)
	// Poll returns the outcome of a match not yet in the session, or nil if there is
	// none. Session stats such as wins and LP gain are left for Session to fill in.
	Poll(ctx context.Context, session *model.Session) (*model.Match, error)
	Authenticate(ctx context.Context, statusChan chan AuthStatus)
}

type AuthStatus struct {
	Done   bool
	Err    error
	Action *AuthAction
}

type AuthAction struct {
	LocalizationKey string `json:"localizationKey"`
	SecondsLeft     int    `json:"secondsLeft"`
}

func (s *AuthStatus) WithError(err error) *AuthStatus {
	s.Err = err
	return s
}

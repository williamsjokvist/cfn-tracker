package tracker

import (
	"context"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type GameTracker interface {
	GetUser(ctx context.Context, userId string) (*model.User, error)
	// Poll leaves session stats such as wins and LP gain to Session, since they depend
	// on every match in the session.
	Poll(ctx context.Context, session *model.Session) (*model.Match, error)
	Authenticate(ctx context.Context, statusChan chan AuthStatus)
}

// GameLister is implemented by trackers for platforms with many games, where a
// session tracks one of them.
type GameLister interface {
	GetPlayedGames(ctx context.Context, userId string) ([]model.PlayedGame, error)
}

type AuthStatus struct {
	Done   bool
	Err    error
	Action *AuthAction
	// Checking and restarting the browser take long enough to be worth showing.
	Step string
}

type AuthAction struct {
	LocalizationKey string `json:"localizationKey"`
	SecondsLeft     int    `json:"secondsLeft"`
}

func (s *AuthStatus) WithError(err error) *AuthStatus {
	s.Err = err
	return s
}

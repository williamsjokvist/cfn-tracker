package tracker

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/browser"
	"github.com/williamsjokvist/cfn-tracker/core/config"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/errorsx"
	"github.com/williamsjokvist/cfn-tracker/core/model"
	"github.com/williamsjokvist/cfn-tracker/core/tracker/sf6"
	"github.com/williamsjokvist/cfn-tracker/core/tracker/sfv"
)

type GameTracker interface {
	Start(ctx context.Context, cfn string, restore bool, refreshInterval time.Duration) error
	Stop()
	GetPlayer(ctx context.Context, code string) (*model.Player, error)
	Search(ctx context.Context, query string) ([]model.Player, error)
}

type GameType uint8

const (
	GameTypeUndefined GameType = iota
	GameTypeSFV
	GameTypeSF6
)

func (s GameType) String() string {
	switch s {
	case GameTypeSFV:
		return `sfv`
	case GameTypeSF6:
		return `sf6`
	case GameTypeUndefined:
		return `undefined`
	}
	return `unknown`
}

// Make a SF6Tracker and expose it as a GameTracker
func MakeSF6Tracker(ctx context.Context, cfg *config.Config, browser *browser.Browser, trackerRepository *data.CFNTrackerRepository) (GameTracker, error) {
	sf6Tracker := sf6.NewSF6Tracker(browser, trackerRepository)

	authChan := make(chan sf6.AuthStatus)
	go sf6Tracker.Authenticate(ctx, cfg.CapIDEmail, cfg.CapIDPassword, authChan)
	for status := range authChan {
		if status.Err != nil {
			return nil, errorsx.NewError(http.StatusUnauthorized, status.Err)
		}
		runtime.EventsEmit(ctx, "auth-progress", status.Progress)

		if status.Progress >= 100 {
			close(authChan)
			break
		}
	}

	var gt GameTracker = sf6Tracker
	return gt, nil
}

// Make a SFVTracker and expose it as a GameTracker
func MakeSFVTracker(ctx context.Context, cfg *config.Config, browser *browser.Browser) (GameTracker, error) {
	sfvTracker := sfv.NewSFVTracker(browser)
	err := sfvTracker.Authenticate(ctx, cfg.SteamUsername, cfg.SteamPassword, false)
	if err != nil {
		return nil, fmt.Errorf(`auth err: %v`, err)
	}
	var gt GameTracker = sfvTracker
	return gt, nil
}

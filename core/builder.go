package core

import (
	"context"
	"fmt"
	"time"

	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/impl/sf6"
	"github.com/williamsjokvist/cfn-tracker/core/impl/sfv"
	"github.com/williamsjokvist/cfn-tracker/core/shared"
)

type GameTracker interface {
	Start(cfn string, restore bool, refreshInterval time.Duration) error
	Stop()
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
func MakeSF6Tracker(ctx context.Context, browser *shared.Browser, username, password string, trackerRepository *data.CFNTrackerRepository) (GameTracker, error) {
	sf6Tracker := sf6.NewSF6Tracker(ctx, browser, trackerRepository)
	err := sf6Tracker.Authenticate(username, password, false)
	if err != nil {
		return nil, fmt.Errorf(`auth err: %v`, err)
	}
	var gt GameTracker = sf6Tracker
	return gt, nil
}

// Make a SFVTracker and expose it as a GameTracker
func MakeSFVTracker(ctx context.Context, browser *shared.Browser, username string, password string) (GameTracker, error) {
	sfvTracker := sfv.NewSFVTracker(ctx, browser)
	err := sfvTracker.Authenticate(username, password, false)
	if err != nil {
		return nil, fmt.Errorf(`auth err: %v`, err)
	}
	var gt GameTracker = sfvTracker
	return gt, nil
}

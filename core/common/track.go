package common

import "time"

type GameTracker interface {
	Start(cfn string, restore bool, refreshInterval time.Duration) error
	Stop()
	GetMatchHistory() *MatchHistory
}

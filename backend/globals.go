package backend

import (
	"time"

	"github.com/go-rod/rod"
	"github.com/hashicorp/go-version"
)

type MatchHistory struct {
	CFN               string `json:"cfn"`
	LP                int    `json:"lp"`
	LPGain            int    `json:"lpGain"`
	Wins              int    `json:"wins"`
	TotalWins         int    `json:"totalWins"`
	TotalLosses       int    `json:"totalLosses"`
	TotalMatches      int    `json:"totalMatches"`
	Losses            int    `json:"losses"`
	WinRate           int    `json:"winRate"`
	Opponent          string `json:"opponent"`
	OpponentCharacter string `json:"opponentCharacter"`
	OpponentLP        int    `json:"opponentLP"`
	OpponentLeague    string `json:"opponentLeague"`
	IsWin             bool   `json:"result"`
	TimeStamp         string `json:"timestamp"`
	Date              string `json:"date"`
	WinStreak         int    `json:"winStreak"`
}

var CurrentMatchHistory = MatchHistory{
	CFN:          ``,
	LP:           0,
	LPGain:       0,
	Wins:         0,
	Losses:       0,
	TotalWins:    0,
	TotalLosses:  0,
	TotalMatches: 0,
	WinRate:      0,
	WinStreak:    0,
	IsWin:        false,
	TimeStamp:    ``,
	Date:         ``,
}

var (
	FirstLPRecorded = 0
	IsTracking      = false
	IsInitialized   = false
	PageInstance    *rod.Page
	SteamUsername   string
	SteamPassword   string
	AppVersion      *version.Version
	RefreshInterval time.Duration = 30 * time.Second
)

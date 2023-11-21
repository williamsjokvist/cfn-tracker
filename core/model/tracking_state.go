package model

import (
	"fmt"
	"time"
)

type TrackingState struct {
	CFN               string `json:"cfn"`
	UserCode          string `json:"userCode"`
	LP                int    `json:"lp"`
	LPGain            int    `json:"lpGain"`
	MR                int    `json:"mr"`
	MRGain            int    `json:"mrGain"`
	Wins              int    `json:"wins"`
	TotalWins         int    `json:"totalWins"`
	TotalLosses       int    `json:"totalLosses"`
	TotalMatches      int    `json:"totalMatches"`
	Losses            int    `json:"losses"`
	WinRate           int    `json:"winRate"`
	Character         string `json:"character"`
	Opponent          string `json:"opponent"`
	OpponentCharacter string `json:"opponentCharacter"`
	OpponentLP        int    `json:"opponentLP"`
	OpponentLeague    string `json:"opponentLeague"`
	IsWin             bool   `json:"result"`
	TimeStamp         string `json:"timestamp"`
	Date              string `json:"date"`
	WinStreak         int    `json:"winStreak"`
}

func (mh *TrackingState) Log() {
	fmt.Printf(`[%s]\n`, time.Now().Format(`15:04`))
	fmt.Printf(`LP: %d\n`, mh.LP)
	fmt.Printf(`LPGain: %d\n`, mh.LPGain)
	fmt.Printf(`MR: %d\n`, mh.MR)
	fmt.Printf(`MRGain: %d\n`, mh.MRGain)
	fmt.Printf(`Wins: %d\n`, mh.Wins)
	fmt.Printf(`Losses: %d\n`, mh.Losses)
	fmt.Printf(`WinRate: %d\n`, mh.WinRate)
}

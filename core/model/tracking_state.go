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
	fmt.Println(`
		[`+time.Now().Format(`15:04`)+`]	
		LP:`, mh.LP, `/ 
		LP Gain:`, mh.LPGain, `/ 
		MR:`, mh.MR, `/ 
		MR Gain:`, mh.MRGain, `/ 
		Wins:`, mh.Wins, `/ 
		Losses:`, mh.Losses, `/ 
		Winrate:`, mh.WinRate, `%`,
	)
}

package model

type Session struct {
	Id          uint16   `db:"id" json:"id"`
	UserId      string   `db:"user_id" json:"userId"`
	UserName    string   `db:"user_name" json:"userName"`
	CreatedAt   string   `db:"created_at" json:"createdAt"`
	Matches     []*Match `json:"matches"`
	MatchesWon  int      `db:"matches_won" json:"matchesWon"`
	MatchesLost int      `db:"matches_lost" json:"matchesLost"`
	WinRate 		int      `db:"win_rate" json:"winRate"`
	EndingLP    int      `db:"ending_lp" json:"endingLp"`
	EndingMR    int      `db:"ending_mr" json:"endingMr"`
	StartingLP  int      `db:"starting_lp" json:"startingLp"`
	StartingMR  int      `db:"starting_mr" json:"startingMr"`
	LPGain      int      `db:"lp_gain" json:"lpGain"`
	MRGain      int      `db:"mr_gain" json:"mrGain"`
}

type SessionMonth struct {
	Date  string
	Count uint16
}

// future: extend with legendary stats
type SessionsStatistics struct {
	Months []SessionMonth
}

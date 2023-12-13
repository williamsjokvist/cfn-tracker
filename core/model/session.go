package model

type Session struct {
	Id          uint16   `db:"id" json:"id"`
	UserId      string   `db:"user_id" json:"userId"`
	UserName    string   `db:"user_name" json:"userName"`
	CreatedAt   string   `db:"created_at" json:"createdAt"`
	LP          int      `db:"lp" json:"lp"`
	MR          int      `db:"mr" json:"mr"`
	Matches     []*Match `json:"matches"`
	MatchesWon  int      `db:"matches_won" json:"matchesWon"`
	MatchesLost int      `db:"matches_lost" json:"matchesLost"`
	StartingLP  int      `db:"starting_lp" json:"startingLp"`
	EndingLP    int      `db:"ending_lp" json:"endingLp"`
	StartingMR  int      `db:"starting_mr" json:"startingMr"`
	EndingMR    int      `db:"ending_mr" json:"endingMr"`
	LPGain      int      `db:"lp_gain" json:"lpGain"`
	MRGain      int      `db:"mr_gain" json:"mrGain"`
}

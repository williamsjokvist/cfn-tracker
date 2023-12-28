package model

type Session struct {
	Id          uint16   `db:"id" json:"id"`
	UserId      string   `db:"user_id" json:"userId"`
	UserName    string   `db:"user_name" json:"userName"`
	CreatedAt   string   `db:"created_at" json:"createdAt"`
	LP          int      `db:"ending_lp" json:"lp"`
	MR          int      `db:"ending_mr" json:"mr"`
	Matches     []*Match `json:"matches"`
	MatchesWon  int      `db:"matches_won" json:"matchesWon"`
	MatchesLost int      `db:"matches_lost" json:"matchesLost"`
	StartingLP  int      `db:"starting_lp" json:"startingLp"`
	StartingMR  int      `db:"starting_mr" json:"startingMr"`
	LPGain      int      `db:"lp_gain" json:"lpGain"`
	MRGain      int      `db:"mr_gain" json:"mrGain"`
}

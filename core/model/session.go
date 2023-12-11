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
}

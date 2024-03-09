package model

type Match struct {
	UserId            string `db:"user_id" json:"userId"`
	SessionId         uint16 `db:"session_id" json:"sessionId"`
	Character         string `db:"character" json:"character"`
	LP                int    `db:"lp" json:"lp"`
	LPGain            int    `db:"lp_gain" json:"lpGain"`
	MR                int    `db:"mr" json:"mr"`
	MRGain            int    `db:"mr_gain" json:"mrGain"`
	Opponent          string `db:"opponent" json:"opponent"`
	OpponentCharacter string `db:"opponent_character" json:"opponentCharacter"`
	OpponentLP        int    `db:"opponent_lp" json:"opponentLp"`
	OpponentMR        int    `db:"opponent_mr" json:"opponentMr"`
	OpponentLeague    string `db:"opponent_league" json:"opponentLeague"`
	Victory           bool   `db:"victory" json:"victory"`
	Date              string `db:"date" json:"date"`
	Time              string `db:"time" json:"time"`
	WinStreak         int    `db:"win_streak" json:"winStreak"`
	Wins              int    `db:"wins" json:"wins"`
	Losses            int    `db:"losses" json:"losses"`
	WinRate           int    `db:"win_rate" json:"winRate"`
}

package model

type Match struct {
	UserId            string `db:"user_id" json:"userId"`
	UserName          string `db:"user_name" json:"userName"`
	SessionId         uint16 `db:"session_id" json:"sessionId"`
	ReplayID          string `db:"replay_id" json:"replayId"`
	Character         string `db:"character" json:"character"`
	LP                int    `db:"lp" json:"lp"`
	MR                int    `db:"mr" json:"mr"`
	Opponent          string `db:"opponent" json:"opponent"`
	OpponentCharacter string `db:"opponent_character" json:"opponentCharacter"`
	OpponentLP        int    `db:"opponent_lp" json:"opponentLp"`
	OpponentMR        int    `db:"opponent_mr" json:"opponentMr"`
	OpponentLeague    string `db:"opponent_league" json:"opponentLeague"`
	Victory           bool   `db:"victory" json:"victory"`
	Date              string `db:"date" json:"date"`
	Time              string `db:"time" json:"time"`
	WinStreak         int    `db:"win_streak" json:"winStreak"`
}

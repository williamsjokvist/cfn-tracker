package sql

type SessionStat struct {
	Id               int    `db:"id"`
	SessionId        int    `db:"sessionId"`
	Character        string `db:"character"`
	LP               int    `db:"lp"`
	MR               int    `db:"mr"`
	LPGain           int    `db:"lpGain"`
	MRGain           int    `db:"mrGain"`
	Wins             int    `db:"wins"`
	Losses           int    `db:"losses"`
	WinRate          int    `db:"winRate"`
	WinStreak        int    `db:"winStreak"`
	LongestWinStreak int    `db:"longestWinStreak"`
	MatchesPlayed    int    `db:"matchesPlayed"`
}

type SessionStatStorage interface {
	createSessionStatTable() error
	GetSessionStats(sessionId, character string) ([]*SessionStat, error)
	UpsertSessionStat(stat *SessionStat) ([]*SessionStat, error)
}

func (s *Storage) UpsertSessionStat(stat *SessionStat) ([]*SessionStat, error) {
	return nil, nil
}

func (s *Storage) GetSessionStats(sessionId, character string) ([]*SessionStat, error) {
	return nil, nil
}

func (s *Storage) createSessionStatTable() error {
	return nil
}

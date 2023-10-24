package sql

type Session struct {
	Id               int    `db:"id"`
	UserId           string `db:"userId"`
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
	StartedAt        string `db:"startedAt"`
}

type SessionStorage interface {
	createSessionsTable() error
	UpsertSession(userId, startedAt string, lp, mr, lpGain, mrGain, wins, losses, winRate, winStreak, longestWinStreak, matchesPlayed int)
	GetSessions(sessionId, userId string, directionOrder string, limit int) ([]*Session, error)
	GetLastSession(userId string) (Session, error)
	RemoveSession(sessionId string) error
}

func (s *Storage) GetSessions(sessionId, userId string, directionOrder string, limit int) ([]*Session, error) {
	return nil, nil
}

func (s *Storage) GetLastSession(userId string) (Session, error) {
	return Session{}, nil
}

func (s *Storage) createSessionsTable() error {
	return nil
}

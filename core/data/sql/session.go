package sql

type Session struct {
	Id      int            `db:"id"`
	UserId  string         `db:"userId"`
	Stats   []*SessionStat `db:"stats"`
	Matches []*Match       `db:"matches"`
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

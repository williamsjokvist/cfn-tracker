package sql

import "fmt"

type Session struct {
	Id      int    `db:"id"`
	UserId  string `db:"user_id"`
	Stats   []*CharacterSessionStats
	Matches []*Match
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
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY,
		user_id INTEGER,
		FOREIGN KEY(user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

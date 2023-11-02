package sql

import "fmt"

type Match struct {
	Id                int    `db:"id"`
	SessionId         string `db:"session_id"`
	Character         string `db:"character"`
	LP                int    `db:"lp"`
	MR                int    `db:"mr"`
	Opponent          string `db:"opponent"`
	OpponentCharacter string `db:"opponent_character"`
	OpponentLP        int    `db:"opponent_lp"`
	OpponentMR        int    `db:"opponent_mr"`
	Victory           bool   `db:"victory"`
	DateTime          string `db:"datetime"`
}

type MatchStorage interface {
	createMatchesTable() error
	GetMatches(sessionId string) ([]*Match, error)
	SaveMatch() error
	RemoveMatches(sessionId string) error
}

func (s *Storage) GetMatches(sessionId string) ([]*Match, error) {
	return nil, nil
}

func (s *Storage) SaveMatch() error {
	return nil
}

func (s *Storage) RemoveMatches(sessionId string) error {
	return nil
}

func (s *Storage) createMatchesTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS matches (
		id INTEGER PRIMARY KEY,
		session_id INTEGER,
		character TEXT NOT NULL,
		lp INTEGER,
		mr INTEGER,
		opponent TEXT,
		opponent_character TEXT,
		opponent_lp TEXT,
		opponent_mr INTEGER,
		victory BOOLEAN,
		datetime TEXT,
		FOREIGN KEY(session_id) REFERENCES sessions(id)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

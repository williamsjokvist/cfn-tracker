package sql

import (
	"context"
	"fmt"
)

type Match struct {
	UserId            uint8  `db:"user_id"`
	SessionCreatedAt  string `db:"session_created_at"`
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

func (s *Storage) SaveMatch(ctx context.Context, match Match) error {
	query := `
		INSERT OR IGNORE INTO matches (
			user_id,
			session_created_at,
			character,
			lp,
			mr,
			opponent,
			opponent_character,
			opponent_lp,
			opponent_mr,
			victory,
			datetime
		)
		VALUES (
			:user_id,
			:session_created_at,
			:character,
			:lp,
			:mr,
			:opponent,
			:opponent_character,
			:opponent_lp,
			:opponent_mr,
			:victory,
			:datetime
		)
	`
	_, err := s.db.NamedExecContext(ctx, query, match)
	if err != nil {
		return fmt.Errorf("create match: %w", err)
	}

	return nil
}

func (s *Storage) RemoveMatches(sessionId string) error {
	return nil
}

func (s *Storage) createMatchesTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS matches (
		user_id INTEGER,
		session_created_at TEXT,
		character TEXT NOT NULL,
		lp INTEGER,
		mr INTEGER,
		opponent TEXT,
		opponent_character TEXT,
		opponent_lp TEXT,
		opponent_mr INTEGER,
		victory BOOLEAN,
		datetime TEXT,
		PRIMARY KEY (user_id, datetime),
		FOREIGN KEY(user_id, session_created_at) REFERENCES sessions(user_id, created_at)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

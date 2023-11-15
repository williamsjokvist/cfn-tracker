package sql

import (
	"context"
	"fmt"
)

type Match struct {
	UserId            uint8  `db:"user_id"`
	SessionId         uint8  `db:"session_id"`
	Character         string `db:"character"`
	LP                int    `db:"lp"`
	LPGain            int    `db:"lp_gain"`
	MR                int    `db:"mr"`
	MRGain            int    `db:"mr_gain"`
	Opponent          string `db:"opponent"`
	OpponentCharacter string `db:"opponent_character"`
	OpponentLP        int    `db:"opponent_lp"`
	OpponentMR        int    `db:"opponent_mr"`
	OpponentLeague    string `db:"opponent_league"`
	Victory           bool   `db:"victory"`
	Date              string `db:"date"`
	Time              string `db:"time"`
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
			session_id,
			character,
			lp,
			lp_gain,
			mr,
			mr_gain,
			opponent,
			opponent_character,
			opponent_lp,
			opponent_mr,
			opponent_league,
			victory,
			date,
			time
		)
		VALUES (
			:user_id,
			:session_id,
			:character,
			:lp,
			:lp_gain,
			:mr,
			:mr_gain,
			:opponent,
			:opponent_character,
			:opponent_lp,
			:opponent_mr,
			:opponent_league,
			:victory,
			:date,
			:time
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
		session_id INTEGER,
		character TEXT NOT NULL,
		lp INTEGER,
		lp_gain INTEGER,
		mr INTEGER,
		mr_gain INTEGER,
		opponent TEXT,
		opponent_character TEXT,
		opponent_lp TEXT,
		opponent_mr INTEGER,
		opponent_league TEXT,
		victory BOOLEAN,
		date TEXT,
		time TEXT,
		PRIMARY KEY (user_id, date, time),
		FOREIGN KEY(session_id) REFERENCES sessions(id)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

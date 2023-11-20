package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
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
	WinStreak         int    `db:"win_streak"`
	Wins              int    `db:"wins"`
	Losses            int    `db:"losses"`
	WinRate           int    `db:"win_rate"`
}

type MatchStorage interface {
	GetMatchesFromSession(ctx context.Context, sessionId int) ([]*Match, error)
	SaveMatch(ctx context.Context, match Match) error
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
			wins,
			losses,
			win_rate,
			win_streak,
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
			:wins,
			:losses,
			:win_rate,
			:win_streak,
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

func (s *Storage) GetMatchesFromSession(ctx context.Context, sessionId int) ([]*Match, error) {
	query, args, err := sqlx.In(`
		SELECT * FROM matches 
		WHERE session_id = (?)
`, sessionId)
	if err != nil {
		return nil, fmt.Errorf("prepare matches by session query: %w", err)
	}
	var matches []*Match
	err = s.db.SelectContext(ctx, &matches, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute matches query: %w", err)
	}
	return matches, nil
}

func (s *Storage) createMatchesTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS matches (
		user_id TEXT,
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
		wins INTEGER,
		losses INTEGER,
		win_streak INTEGER,
		win_rate INTEGER,
		date TEXT,
		time TEXT,
		PRIMARY KEY (user_id, date, time),
		FOREIGN KEY (session_id) REFERENCES sessions(id)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

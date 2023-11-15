package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type Session struct {
	SessionId int64  `db:"id"`
	UserId    string `db:"user_id"`
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
	LP        int    `db:"lp"`
	LPGain    int    `db:"lp_gain"`
	MR        int    `db:"mr"`
	MRGain    int    `db:"mr_gain"`
	Wins      int    `db:"wins"`
	Losses    int    `db:"losses"`
	WinStreak int    `db:"win_streak"`
	WinRate   int    `db:"win_rate"`
}

type SessionStorage interface {
	CreateSession(ctx context.Context, userId string) error
	GetLatestSession(ctx context.Context, userId string) (*Session, error)
	UpdateLatestSession(ctx context.Context, userId string) error
}

func (s *Storage) CreateSession(ctx context.Context, userId string) (*Session, error) {
	sesh := Session{
		UserId:    userId,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	query := `
		INSERT OR IGNORE INTO sessions (user_id, created_at, updated_at)
		VALUES (:user_id, :created_at, :updated_at)
	`
	res, err := s.db.NamedExecContext(ctx, query, sesh)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	sesh.SessionId, err = res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return &sesh, nil
}

func (s *Storage) GetLatestSession(ctx context.Context, userId string) (*Session, error) {
	query, args, err := sqlx.In(`
		SELECT * FROM sessions 
		WHERE user_id = (?)
		LIMIT 1
`, userId)
	if err != nil {
		return nil, fmt.Errorf("prepare last session query: %w", err)
	}
	var sesh Session
	err = s.db.GetContext(ctx, &sesh, query, args...)
	if err != nil {
		return nil, fmt.Errorf("excute select last session query: %w", err)
	}
	return &sesh, nil
}

func (s *Storage) UpdateLatestSession(ctx context.Context, lp, mr, lpGain, mrGain, wins, losses, winStreak, winRate int, userId string) error {
	query, args, err := sqlx.In(`
		UPDATE sessions
		SET 
			lp = ?,
			mr = ?,
			lp_gain = ?,
			mr_gain = ?,
			wins = ?,
			losses = ?,
			win_streak = ?,
			win_rate = ?
		WHERE user_id = (?)
		ORDER BY id DESC
		LIMIT 1
`, lp, mr, lpGain, mrGain, wins, losses, winStreak, winRate, userId)
	if err != nil {
		return fmt.Errorf("prepare update session query: %w", err)
	}
	_, err = s.db.ExecContext(ctx, s.db.Rebind(query), args...)
	if err != nil {
		return fmt.Errorf("excute update last session query: %w", err)
	}
	return nil
}

func (s *Storage) createSessionsTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS sessions (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id TEXT,
		lp INTEGER,
		lp_gain INTEGER,
		mr INTEGER,
		mr_gain INTEGER,
		wins INTEGER,
		losses INTEGER,
		win_streak INTEGER,
		win_rate INTEGER,
		created_at TEXT,
		updated_at TEXT,
		FOREIGN KEY (user_id) REFERENCES users(code)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

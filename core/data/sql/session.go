package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/williamsjokvist/cfn-tracker/core/model"
)

type SessionStorage interface {
	CreateSession(ctx context.Context, userId string) error
	GetLatestSession(ctx context.Context, userId string) (*model.Session, error)
	UpdateLatestSession(ctx context.Context, userId string) error
}

func (s *Storage) CreateSession(ctx context.Context, userId string) (*model.Session, error) {
	sesh := model.Session{
		UserId:    userId,
		CreatedAt: time.Now().String(),
	}
	query := `
		INSERT OR IGNORE INTO sessions (user_id, created_at)
		VALUES (:user_id, :created_at)
	`
	res, err := s.db.NamedExecContext(ctx, query, sesh)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	sesh.Id = uint16(lastInsertId)
	return &sesh, nil
}

func (s *Storage) GetLatestSession(ctx context.Context, userId string) (*model.Session, error) {
	query, args, err := sqlx.In(`
		SELECT * FROM sessions 
		WHERE user_id = (?)
		LIMIT 1
`, userId)
	if err != nil {
		return nil, fmt.Errorf("prepare last session query: %w", err)
	}
	var sesh model.Session
	err = s.db.GetContext(ctx, &sesh, query, args...)
	if err != nil {
		return nil, fmt.Errorf("excute select last session query: %w", err)
	}
	return &sesh, nil
}

func (s *Storage) UpdateLatestSession(ctx context.Context, lp, mr int, sessionId uint16) error {
	query, args, err := sqlx.In(`
		UPDATE sessions
		SET 
			lp = ?,
			mr = ?
		WHERE id = (?)
`, lp, mr, sessionId)
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
		mr INTEGER,
		created_at TEXT,
		FOREIGN KEY (user_id) REFERENCES users(code)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

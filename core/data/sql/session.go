package sql

import (
	"context"
	"fmt"
	"time"
)

type Session struct {
	UserId    string `db:"user_id"`
	Stats     []*CharacterSessionStats
	Matches   []*Match
	CreatedAt string `db:"created_at"`
	UpdatedAt string `db:"updated_at"`
}

type SessionStorage interface {
	createSessionsTable() error
	CreateSession(ctx context.Context, userId string) error
	GetSessions(sessionId, userId string, directionOrder string, limit int) ([]*Session, error)
	GetLastSession(userId string) (Session, error)
	RemoveSession(sessionId string) error
}

func (s *Storage) CreateSession(ctx context.Context, userId string) error {
	sesh := Session{
		UserId:    userId,
		CreatedAt: time.Now().String(),
		UpdatedAt: time.Now().String(),
	}

	query := `
		INSERT OR IGNORE INTO sessions (user_id, created_at, updated_at)
		VALUES (:user_id, :created_at, :updated_at)
	`
	_, err := s.db.NamedExecContext(ctx, query, sesh)
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	return nil
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
		user_id INTEGER,
		created_at TEXT,
		updated_at TEXT,
		PRIMARY KEY (user_id, created_at),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

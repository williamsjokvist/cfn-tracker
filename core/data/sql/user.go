package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id          uint8  `db:"id"`
	DisplayName string `db:"displayName"`
	Code        string `db:"code"`
}

type UserStorage interface {
	createUsersTable() error
	GetUsers() ([]User, error)
	SaveUser(displayName, id string) error
	RemoveUser(id string) error
}

func (s *Storage) GetUsers(ctx context.Context) ([]*User, error) {
	var users []*User
	err := s.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("select sql users: %w", err)
	}
	return users, nil
}

func (s *Storage) SaveUser(ctx context.Context, displayName, code string) error {
	user := User{
		DisplayName: displayName,
		Code:        code,
	}
	query := `
		INSERT OR IGNORE INTO users (displayName, code)
		VALUES (:displayName, :code)
	`
	_, err := s.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Storage) RemoveUser(ctx context.Context, id string) error {
	query, args, err := sqlx.In(`
		DELETE * FROM users 
		WHERE id = (?)
	`, id)
	if err != nil {
		return fmt.Errorf("prepare delete user clause: %w", err)
	}
	_, err = s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (s *Storage) createUsersTable() error {
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY,
		code TEXT NOT NULL UNIQUE,
		displayName TEXT NOT NULL
	)`)
	if err != nil {
		return fmt.Errorf("create users table: %w", err)
	}
	return nil
}

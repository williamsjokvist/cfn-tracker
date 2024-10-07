package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type UserStorage interface {
	GetUserByCode(ctx context.Context, code string) (*model.User, error)
	GetUsers(ctx context.Context) ([]*model.User, error)
	SaveUser(ctx context.Context, user model.User) error
	RemoveUser(ctx context.Context, code string) error
}

func (s *Storage) GetUserByCode(ctx context.Context, code string) (*model.User, error) {
	query, args, err := sqlx.In(`
		SELECT * FROM users
		WHERE code = (?)
		LIMIT 1
	`, code)
	if err != nil {
		return nil, fmt.Errorf("prepare get user clause: %w", err)
	}
	var user model.User
	err = s.db.GetContext(ctx, &user, query, args...)
	if err != nil {
		return nil, fmt.Errorf("get user by code: %w", err)
	}
	return &user, nil
}

func (s *Storage) GetUsers(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	err := s.db.SelectContext(ctx, &users, "SELECT * FROM users")
	if err != nil {
		return nil, fmt.Errorf("select sql users: %w", err)
	}
	return users, nil
}

func (s *Storage) SaveUser(ctx context.Context, user model.User) error {
	query := `
		INSERT OR IGNORE INTO users (display_name, code)
		VALUES (:display_name, :code)
	`
	_, err := s.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (s *Storage) RemoveUser(ctx context.Context, code string) error {
	query, args, err := sqlx.In(`
		DELETE * FROM users
		WHERE code = (?)
	`, code)
	if err != nil {
		return fmt.Errorf("prepare delete user clause: %w", err)
	}
	_, err = s.db.NamedExecContext(ctx, query, args)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

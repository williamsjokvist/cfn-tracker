package storage

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage() (*Storage, error) {
	db, err := sqlx.Open("sqlite3", "cfn-tracker.db")
	if err != nil {
		return nil, fmt.Errorf("open sqlite3 connection: %w", err)
	}
	storage := &Storage{
		db: db,
	}

	err = storage.createUsersTable()
	if err != nil {
		return nil, err
	}

	err = storage.createSessionsTable()
	if err != nil {
		return nil, err
	}

	err = storage.createMatchesTable()
	if err != nil {
		return nil, err
	}

	return storage, nil
}

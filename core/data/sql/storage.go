package sql

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sqlx.DB
}

func NewStorage() (*Storage, error) {
	dataSource := getDataSource()
	db, err := sqlx.Open("sqlite3", dataSource)
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

func getDataSource() string {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	os.MkdirAll(dataDir, os.FileMode(0755))
	return filepath.Join(dataDir, "cfn-tracker.db")
}

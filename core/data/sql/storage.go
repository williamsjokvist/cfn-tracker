package sql

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"

	sqlitex "github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

//go:embed migrations
var migrationsFs embed.FS

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

	if err = migrateSchema(db); err != nil {
		return nil, fmt.Errorf("failed to perform migrations: %w", err)
	}

	return storage, nil
}

func getDataSource() string {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	os.MkdirAll(dataDir, os.FileMode(0755))
	return filepath.Join(dataDir, "cfn-tracker.db")
}

func migrateSchema(db *sqlx.DB) error {
	migrateDriver, err := sqlitex.WithInstance(db.DB, &sqlitex.Config{})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	return nil
}

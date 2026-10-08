package sql

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	migrate "github.com/golang-migrate/migrate/v4"
	sqlitex "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrationsFs embed.FS

type Storage struct {
	db *sqlx.DB
}

func (s *Storage) DB() *sqlx.DB {
	return s.db
}

func NewStorage(useInMemoryDb bool) (*Storage, error) {
	var db *sqlx.DB
	var err error

	if useInMemoryDb {
		db, err = sqlx.Open("sqlite", ":memory:")
		if err != nil {
			return nil, fmt.Errorf("open sqlite connection: %w", err)
		}
		if err := migrateSchema(db, nil); err != nil {
			db.Close()
			return nil, fmt.Errorf("perform sql migrations: %w", err)
		}
	} else {
		if err := migrateSchema(nil, nil); err != nil {
			return nil, fmt.Errorf("perform sql migrations: %w", err)
		}
		db, err = sqlx.Open("sqlite", getDataSource())
		if err != nil {
			return nil, fmt.Errorf("open sqlite connection: %w", err)
		}
	}

	return &Storage{
		db,
	}, nil
}

func getDataSource() string {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	if err := os.MkdirAll(dataDir, os.FileMode(0755)); err != nil {
		return "cfn-tracker.db"
	}

	return filepath.Join(dataDir, "cfn-tracker.db")
}

func migrateSchema(db *sqlx.DB, nSteps *int) error {
	slog.Debug("starting db migrations", slog.Any("steps", nSteps))

	shouldClose := false
	if db == nil {
		var err error
		db, err = sqlx.Open("sqlite", getDataSource())
		if err != nil {
			return fmt.Errorf("open sqlite connection: %w", err)
		}
		shouldClose = true
	}

	migrateDriver, err := sqlitex.WithInstance(db.DB, &sqlitex.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}
	srcDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("create migration source driver: %w", err)
	}
	preparedMigrations, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"",
		migrateDriver,
	)
	if err != nil {
		return fmt.Errorf("create migration tooling instance: %w", err)
	}
	defer func() {
		if shouldClose {
			preparedMigrations.Close()
			db.Close()
		}
	}()
	if nSteps != nil {
		fmt.Printf("stepping migrations %d...\n", *nSteps)
		err = preparedMigrations.Steps(*nSteps)
	} else {
		err = preparedMigrations.Up()
	}

	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("apply migrations: %w", err)
	}

	slog.Debug("applied db migrations")
	return nil
}

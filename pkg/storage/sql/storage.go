package sql

import (
	"embed"
	"errors"
	"fmt"
	"log"
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

func NewStorage(useInMemoryDb bool) (*Storage, error) {
	if err := migrateSchema(nil); err != nil {
		return nil, fmt.Errorf("perform sql migrations: %w", err)
	}
	dataSource := getDataSource()
	if useInMemoryDb {
		dataSource = ":memory:"
	}

	db, err := sqlx.Open("sqlite", dataSource)
	if err != nil {
		return nil, fmt.Errorf("open sqlite connection: %w", err)
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

func migrateSchema(nSteps *int) error {
	db, err := sqlx.Open("sqlite", getDataSource())
	if err != nil {
		return fmt.Errorf("open sqlite connection: %w", err)
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
		preparedMigrations.Close()
		db.Close()
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

	log.Println("Successfully applied db migrations")
	return nil
}

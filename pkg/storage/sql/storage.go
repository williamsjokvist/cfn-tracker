package sql

import (
	"context"
	"database/sql"
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
	"github.com/mattn/go-sqlite3"
)

//go:embed migrations
var migrationsFs embed.FS

type Storage struct {
	db *sqlx.DB
}

func NewStorage() (*Storage, error) {
	if err := migrateSchema(nil); err != nil {
		return nil, fmt.Errorf("failed to perform migrations: %w", err)
	}
	db, err := sqlx.Open("sqlite3", getDataSource())
	if err != nil {
		return nil, fmt.Errorf("open sqlite3 connection: %w", err)
	}
	return &Storage{
		db,
	}, nil
}

func getDataSource() string {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	os.MkdirAll(dataDir, os.FileMode(0755))
	return filepath.Join(dataDir, "cfn-tracker.db")
}

func migrateSchema(nSteps *int) error {
	db, err := sqlx.Open("sqlite3", getDataSource())
	if err != nil {
		return fmt.Errorf("open sqlite3 connection: %w", err)
	}

	migrateDriver, err := sqlitex.WithInstance(db.DB, &sqlitex.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return fmt.Errorf("failed to create migration driver: %w", err)
	}
	srcDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("failed to create migration source driver: %w", err)
	}
	preparedMigrations, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"",
		migrateDriver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migration tooling instance: %w", err)
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
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	log.Println("Successfully applied db migrations")
	return nil
}

func (s *Storage) CreateBackup(ctx context.Context) error {
	conns := make([]*sqlite3.SQLiteConn, 0, 1)
	sql.Register(
		"sqlite3-backup",
		&sqlite3.SQLiteDriver{
			ConnectHook: func(conn *sqlite3.SQLiteConn) error {
				conns = append(conns, conn)
				return nil
			},
		},
	)

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %w", err)
	}
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	os.MkdirAll(dataDir, os.FileMode(0755))
	backupFilepath := filepath.Join(dataDir, "cfn-tracker.backup.db")
	backupDb, err := sqlx.Open("sqlite3-backup", backupFilepath)
	if err != nil {
		return fmt.Errorf("open sqlite3 connection: %w", err)
	}
	defer func() {
		if closeErr := backupDb.Close(); closeErr != nil {
			log.Printf("failed to close backup: %v", closeErr)
		}
	}()
	// sql package offers no guarantee that a connection is established when Open() returns
	// sqlite driver lazily opens a connection when it is needed
	// so we need to ping the connection to ensure it is established
	backupDb.Ping()

	sourceDb, err := sqlx.Open("sqlite3-backup", getDataSource())
	if err != nil {
		return fmt.Errorf("open sqlite3 connection: %w", err)
	}
	defer func() {
		if closeErr := sourceDb.Close(); closeErr != nil {
			log.Printf("failed to close backup: %v", closeErr)
		}
	}()
	sourceDb.Ping()

	if len(conns) != 2 {
		return fmt.Errorf("expected 2 connections, got %d", len(conns))
	}

	backupConn := conns[0]
	sourceConn := conns[1]
	// TODO: Check error type
	// The database name is "main" as per sqlite3 documentation
	backup, err := backupConn.Backup("main", sourceConn, "main")
	if err != nil {
		return fmt.Errorf("backup db: %w", err)
	}
	// backup.Finish() just forwards to backup.Close()
	defer func() {
		if closeErr := backup.Close(); closeErr != nil {
			log.Printf("failed to close backup: %v", closeErr)
		}
	}()

	didComplete, err := backup.Step(-1)
	if err != nil {
		return fmt.Errorf("backup step: %w", err)
	}
	if !didComplete {
		return fmt.Errorf("backup did not complete") // Should never happen when using -1 step and error is nil
	}

	return nil
}

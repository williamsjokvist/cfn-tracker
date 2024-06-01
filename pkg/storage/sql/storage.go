package sql

import (
	"context"
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
	db       *sqlx.DB
	connChan chan *sqlite3.SQLiteConn
}

const backupDriverName string = "sqlite3-backup-db-driver"

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
		make(chan *sqlite3.SQLiteConn, 1),
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

/*
Documentation for sqlite3 backup API:
https://www.sqlite.org/c3ref/backup_finish.html#sqlite3backupinit
*/
func (s *Storage) CreateBackup(ctx context.Context) error {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %w", err)
	}
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	err = os.MkdirAll(dataDir, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}

	backupFilepath := filepath.Join(dataDir, "cfn-tracker.backup.db")
	backupDb, err := sqlx.Open(backupDriverName, backupFilepath)
	if err != nil {
		return fmt.Errorf("open sqlite3 connection: %w", err)
	}
	defer func() {
		if closeErr := backupDb.Close(); closeErr != nil {
			log.Printf("failed to close backup: %v", closeErr)
		}
	}()
	// sql.Open may not immediatly open a connection
	// call Ping to ensure the connection is established
	backupDb.Ping()
	var backupConn *sqlite3.SQLiteConn
	select {
	case backupConn = <-s.connChan:
	case <-ctx.Done():
		return fmt.Errorf("backup connection: %w", ctx.Err())
	}
	defer func() {
		if closeErr := backupConn.Close(); closeErr != nil {
			log.Printf("failed to close backup connection: %v", closeErr)
		}
	}()

	// Need new connection to perform backup
	// The source connection may still be used by the application during backup
	sourceDb, err := sqlx.Open(backupDriverName, getDataSource())
	if err != nil {
		return fmt.Errorf("open sqlite3 connection: %w", err)
	}
	defer func() {
		if closeErr := sourceDb.Close(); closeErr != nil {
			log.Printf("failed to close backup: %v", closeErr)
		}
	}()
	sourceDb.Ping()
	var sourceConn *sqlite3.SQLiteConn
	select {
	case sourceConn = <-s.connChan:
	case <-ctx.Done():
		return fmt.Errorf("source connection: %w", ctx.Err())
	}
	defer func() {
		if closeErr := sourceConn.Close(); closeErr != nil {
			log.Printf("failed to close source connection: %v", closeErr)
		}
	}()

	// The database name is "main"
	// https://www.sqlite.org/c3ref/backup_finish.html#sqlite3backupinit
	backup, err := backupConn.Backup("main", sourceConn, "main")
	if err != nil {
		return fmt.Errorf("backup db: %w", err)
	}
	defer func() {
		if closeErr := backup.Close(); closeErr != nil {
			log.Printf("failed to close backup: %v", closeErr)
		}
	}()

	// -1 means to copy all data
	// https://www.sqlite.org/c3ref/backup_finish.html#sqlite3backupstep
	didComplete, err := backup.Step(-1)
	if err != nil {
		return fmt.Errorf("backup step: %w", err)
	}
	if !didComplete {
		// Should never happen when using -1 step and error is nil
		return fmt.Errorf("backup did not complete")
	}

	return nil
}

package sql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-sqlite3"
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

func getDataSource() string {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	os.MkdirAll(dataDir, os.FileMode(0755))
	return filepath.Join(dataDir, "cfn-tracker.db")
}

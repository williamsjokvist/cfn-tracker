package sql

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNotABackup   = errors.New("not a cfn-tracker database")
	ErrBackupTooNew = errors.New("database is from a newer version of the app")
)

// importSuffix marks a backup waiting to replace the database on the next start.
// The database can't be swapped while it's open, as the handlers share it.
const importSuffix = ".import"

func (s *Storage) BackupTo(ctx context.Context, path string) error {
	// VACUUM INTO gives a consistent copy while the database is in use, but it won't
	// write over an existing file.
	tmp, err := tempPathIn(filepath.Dir(path))
	if err != nil {
		return err
	}
	defer removeIfExists(tmp)
	if _, err := s.db.ExecContext(ctx, "VACUUM INTO ?", tmp); err != nil {
		return fmt.Errorf("vacuum into %s: %w", tmp, err)
	}
	return replaceDatabase(tmp, path)
}

func CopyDatabase(src, dst string) error {
	tmp, err := tempPathIn(filepath.Dir(dst))
	if err != nil {
		return err
	}
	defer removeIfExists(tmp)
	if err := copyFile(src, tmp); err != nil {
		return err
	}
	return replaceDatabase(tmp, dst)
}

func CheckBackup(path string) error {
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("stat %s: %w", path, err)
	}
	db, err := sqlx.Open("sqlite", path)
	if err != nil {
		return fmt.Errorf("open sqlite connection: %w", err)
	}
	defer db.Close()

	var tables int
	// Fails on a file that isn't SQLite at all.
	if err := db.Get(&tables, "SELECT count(*) FROM sqlite_master WHERE type = 'table' AND name = 'migrations'"); err != nil {
		return errors.Join(ErrNotABackup, err)
	}
	if tables == 0 {
		return ErrNotABackup
	}
	var version int
	if err := db.Get(&version, "SELECT version FROM migrations LIMIT 1"); err != nil {
		return errors.Join(ErrNotABackup, err)
	}
	latest, err := latestMigration()
	if err != nil {
		return err
	}
	if version > latest {
		return fmt.Errorf("%w: schema version %d, app supports %d", ErrBackupTooNew, version, latest)
	}
	return nil
}

func StageImport(src, dbPath string) error {
	tmp, err := tempPathIn(filepath.Dir(dbPath))
	if err != nil {
		return err
	}
	defer removeIfExists(tmp)
	if err := copyFile(src, tmp); err != nil {
		return err
	}
	if err := os.Rename(tmp, StagedImportPath(dbPath)); err != nil {
		return fmt.Errorf("rename staged import: %w", err)
	}
	return nil
}

func StagedImportPath(dbPath string) string {
	return dbPath + importSuffix
}

func HasStagedImport(dbPath string) bool {
	_, err := os.Stat(StagedImportPath(dbPath))
	return err == nil
}

// ApplyStagedImport has to run before the database is opened.
func ApplyStagedImport(dbPath string) error {
	staged := StagedImportPath(dbPath)
	if !HasStagedImport(dbPath) {
		return nil
	}
	return replaceDatabase(staged, dbPath)
}

func replaceDatabase(src, dst string) error {
	// Journal files left by a crash belong to the old database, and SQLite would
	// apply them to the new one.
	for _, suffix := range []string{"-journal", "-wal", "-shm"} {
		if err := removeIfExists(dst + suffix); err != nil {
			return err
		}
	}
	if err := os.Rename(src, dst); err != nil {
		return fmt.Errorf("replace %s: %w", dst, err)
	}
	return nil
}

// tempPathIn returns an unused path in dir, so a rename to the real name is atomic.
func tempPathIn(dir string) (string, error) {
	f, err := os.CreateTemp(dir, ".cfn-tracker-*.tmp")
	if err != nil {
		return "", fmt.Errorf("create temp file: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close temp file: %w", err)
	}
	if err := os.Remove(f.Name()); err != nil {
		return "", fmt.Errorf("remove temp file: %w", err)
	}
	return f.Name(), nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open %s: %w", src, err)
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create %s: %w", dst, err)
	}
	if _, err := io.Copy(out, in); err != nil {
		return errors.Join(fmt.Errorf("copy %s: %w", src, err), out.Close())
	}
	if err := out.Close(); err != nil {
		return fmt.Errorf("close %s: %w", dst, err)
	}
	return nil
}

func removeIfExists(path string) error {
	if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("remove %s: %w", path, err)
	}
	return nil
}

func latestMigration() (int, error) {
	entries, err := fs.ReadDir(migrationsFs, "migrations")
	if err != nil {
		return 0, fmt.Errorf("read migrations: %w", err)
	}
	latest := 0
	for _, entry := range entries {
		version, err := strconv.Atoi(strings.SplitN(entry.Name(), "_", 2)[0])
		if err != nil {
			return 0, fmt.Errorf("parse migration version of %s: %w", entry.Name(), err)
		}
		latest = max(latest, version)
	}
	return latest, nil
}

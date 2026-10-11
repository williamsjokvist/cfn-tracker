package sql

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func newFileStorage(t *testing.T, path string) *Storage {
	t.Helper()
	storage, err := NewStorage(path)
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	t.Cleanup(func() { _ = storage.Close() })
	return storage
}

func saveUser(t *testing.T, storage *Storage, code string) {
	t.Helper()
	if err := storage.SaveUser(context.Background(), model.User{Code: code, DisplayName: code, Game: model.GameTypeSF6}); err != nil {
		t.Fatalf("SaveUser(%s) error = %v", code, err)
	}
}

func hasUser(t *testing.T, path, code string) bool {
	t.Helper()
	storage := newFileStorage(t, path)
	user, err := storage.GetUserByCode(context.Background(), code)
	return err == nil && user != nil
}

func TestBackupToWhileOpenAndOverExistingFile(t *testing.T) {
	dir := t.TempDir()
	storage := newFileStorage(t, filepath.Join(dir, DBFileName))
	saveUser(t, storage, "1111")

	backup := filepath.Join(dir, "backups", DBFileName+".bak")
	if err := os.MkdirAll(filepath.Dir(backup), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(backup, []byte("old backup"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := storage.BackupTo(context.Background(), backup); err != nil {
		t.Fatalf("BackupTo() error = %v", err)
	}

	if err := CheckBackup(backup); err != nil {
		t.Fatalf("CheckBackup() error = %v", err)
	}
	if !hasUser(t, backup, "1111") {
		t.Error("backup is missing the saved user")
	}
	entries, err := os.ReadDir(filepath.Dir(backup))
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 1 {
		t.Errorf("backup folder has %d files, want only the backup", len(entries))
	}
}

func TestCheckBackup(t *testing.T) {
	dir := t.TempDir()
	valid := filepath.Join(dir, DBFileName)
	newFileStorage(t, valid)

	notSQLite := filepath.Join(dir, "notes.bak")
	if err := os.WriteFile(notSQLite, []byte("not a database at all, just some text"), 0o644); err != nil {
		t.Fatal(err)
	}

	otherSQLite := filepath.Join(dir, "other.db")
	other := newFileStorage(t, otherSQLite)
	if _, err := other.db.Exec("DROP TABLE migrations"); err != nil {
		t.Fatal(err)
	}

	newer := filepath.Join(dir, "newer.db")
	newerStorage := newFileStorage(t, newer)
	if _, err := newerStorage.db.Exec("UPDATE migrations SET version = version + 1"); err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name string
		path string
		want error
	}{
		{"valid", valid, nil},
		{"not sqlite", notSQLite, ErrNotABackup},
		{"other sqlite database", otherSQLite, ErrNotABackup},
		{"newer schema", newer, ErrBackupTooNew},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := CheckBackup(tt.path)
			if tt.want == nil && err != nil {
				t.Fatalf("CheckBackup() error = %v, want nil", err)
			}
			if tt.want != nil && !errors.Is(err, tt.want) {
				t.Fatalf("CheckBackup() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestApplyStagedImport(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, DBFileName)

	backupPath := filepath.Join(dir, "backup.bak")
	backup := newFileStorage(t, backupPath)
	saveUser(t, backup, "2222")
	if err := backup.Close(); err != nil {
		t.Fatal(err)
	}

	current := newFileStorage(t, dbPath)
	saveUser(t, current, "1111")
	if err := current.Close(); err != nil {
		t.Fatal(err)
	}
	// A journal left by the old database would be applied to the imported one.
	if err := os.WriteFile(dbPath+"-journal", []byte("stale"), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := StageImport(backupPath, dbPath); err != nil {
		t.Fatalf("StageImport() error = %v", err)
	}
	if !HasStagedImport(dbPath) {
		t.Fatal("HasStagedImport() = false after StageImport()")
	}
	if err := ApplyStagedImport(dbPath); err != nil {
		t.Fatalf("ApplyStagedImport() error = %v", err)
	}
	if HasStagedImport(dbPath) {
		t.Error("staged import is still there after applying it")
	}
	if _, err := os.Stat(dbPath + "-journal"); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("stale journal wasn't removed: %v", err)
	}
	if !hasUser(t, dbPath, "2222") || hasUser(t, dbPath, "1111") {
		t.Error("database wasn't replaced by the import")
	}
}

func TestApplyStagedImportWithoutImport(t *testing.T) {
	dbPath := filepath.Join(t.TempDir(), DBFileName)
	if err := ApplyStagedImport(dbPath); err != nil {
		t.Fatalf("ApplyStagedImport() error = %v", err)
	}
	if _, err := os.Stat(dbPath); !errors.Is(err, os.ErrNotExist) {
		t.Errorf("ApplyStagedImport() created %s", dbPath)
	}
}

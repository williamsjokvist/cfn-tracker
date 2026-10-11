package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	cfgDb "github.com/williamsjokvist/cfn-tracker/pkg/storage/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
)

type dataTest struct {
	handler  *DataHandler
	tracking *TrackingHandler
	dialogs  *Dialogs
	sqlDb    *sql.Storage
	cfgDb    *cfgDb.Storage
	dir      string
}

// newDataTest points the user cache directory at a temp dir, so tests don't touch
// the real config or database.
func newDataTest(t *testing.T) *dataTest {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	t.Setenv("LocalAppData", filepath.Join(home, "AppData", "Local"))
	dir, err := cfgDb.DefaultDataDir()
	if err != nil {
		t.Fatal(err)
	}
	cfg, err := cfgDb.NewStorageIn(dir)
	if err != nil {
		t.Fatal(err)
	}
	sqlDb, err := sql.NewStorage(dbPathIn(dir))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDb.Close() })
	if err := sqlDb.SaveUser(context.Background(), model.User{Code: "1111", DisplayName: "Kakeru", Game: model.GameTypeSF6}); err != nil {
		t.Fatal(err)
	}
	tracking := &TrackingHandler{}
	dialogs := &Dialogs{}
	return &dataTest{
		handler:  NewDataHandler(sqlDb, cfg, tracking, dialogs, dir),
		tracking: tracking,
		dialogs:  dialogs,
		sqlDb:    sqlDb,
		cfgDb:    cfg,
		dir:      dir,
	}
}

func (dt *dataTest) startTracking() {
	dt.tracking.mu.Lock()
	dt.tracking.cancelPolling = func() {}
	dt.tracking.mu.Unlock()
}

func hasUser(t *testing.T, dbPath, code string) bool {
	t.Helper()
	storage, err := sql.NewStorage(dbPath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = storage.Close() }()
	user, err := storage.GetUserByCode(context.Background(), code)
	return err == nil && user != nil
}

func isFGCError(err error, want *model.FGCTrackerError) bool {
	var fgcErr *model.FGCTrackerError
	return errors.As(err, &fgcErr) && fgcErr.LocalizationKey == want.LocalizationKey
}

func TestSaveDataDirectory(t *testing.T) {
	dt := newDataTest(t)
	target := t.TempDir()

	if dt.handler.HasSaveData(target) {
		t.Fatal("HasSaveData() = true for an empty folder")
	}
	if err := dt.handler.SaveDataDirectory(target, false); err != nil {
		t.Fatalf("SaveDataDirectory() error = %v", err)
	}
	if !hasUser(t, dbPathIn(target), "1111") {
		t.Error("the moved database is missing the user's data")
	}
	if _, err := os.Stat(dbPathIn(dt.dir)); err != nil {
		t.Errorf("the old database wasn't left in place: %v", err)
	}
	loc, err := dt.handler.GetDataLocation()
	if err != nil {
		t.Fatal(err)
	}
	if loc.Dir != target || !loc.RestartRequired {
		t.Errorf("GetDataLocation() = %+v, want %s with a restart required", loc, target)
	}

	// Moving back to the open database cancels the move.
	if err := dt.handler.SaveDataDirectory(dt.dir, false); err != nil {
		t.Fatalf("SaveDataDirectory() back error = %v", err)
	}
	loc, err = dt.handler.GetDataLocation()
	if err != nil {
		t.Fatal(err)
	}
	if loc.Dir != dt.dir || loc.RestartRequired {
		t.Errorf("GetDataLocation() = %+v, want %s without a restart", loc, dt.dir)
	}
}

func TestSaveDataDirectoryWithExistingSaveData(t *testing.T) {
	dt := newDataTest(t)
	target := t.TempDir()
	existing, err := sql.NewStorage(dbPathIn(target))
	if err != nil {
		t.Fatal(err)
	}
	if err := existing.SaveUser(context.Background(), model.User{Code: "2222", DisplayName: "Punk", Game: model.GameTypeSF6}); err != nil {
		t.Fatal(err)
	}
	if err := existing.Close(); err != nil {
		t.Fatal(err)
	}

	if !dt.handler.HasSaveData(target) {
		t.Fatal("HasSaveData() = false for a folder with a database")
	}
	if err := dt.handler.SaveDataDirectory(target, true); err != nil {
		t.Fatalf("SaveDataDirectory(useExisting) error = %v", err)
	}
	if !hasUser(t, dbPathIn(target), "2222") || hasUser(t, dbPathIn(target), "1111") {
		t.Error("using the existing save data changed it")
	}

	if err := dt.handler.SaveDataDirectory(target, false); err != nil {
		t.Fatalf("SaveDataDirectory(replace) error = %v", err)
	}
	// The configured folder's database is what the app opens next, so replacing
	// it with itself keeps the existing data.
	if !hasUser(t, dbPathIn(target), "2222") {
		t.Error("replacing with the pending database lost it")
	}
}

func TestSaveDataDirectoryUseExistingRejectsInvalid(t *testing.T) {
	dt := newDataTest(t)
	target := t.TempDir()
	if err := os.WriteFile(dbPathIn(target), []byte("not a database at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := dt.handler.SaveDataDirectory(target, true); !isFGCError(err, model.ErrInvalidBackup) {
		t.Fatalf("SaveDataDirectory() error = %v, want ErrInvalidBackup", err)
	}
	if loc, err := dt.handler.GetDataLocation(); err != nil || loc.Dir != dt.dir {
		t.Errorf("GetDataLocation() = %+v, %v, want the setting unchanged", loc, err)
	}
}

func TestRefusesWhileTracking(t *testing.T) {
	dt := newDataTest(t)
	dt.startTracking()
	target := t.TempDir()

	if err := dt.handler.SaveDataDirectory(target, false); !isFGCError(err, model.ErrTrackingInProgress) {
		t.Errorf("SaveDataDirectory() error = %v, want ErrTrackingInProgress", err)
	}
	if _, err := os.Stat(dbPathIn(target)); !errors.Is(err, os.ErrNotExist) {
		t.Error("SaveDataDirectory() copied the database while tracking")
	}
	if err := dt.handler.ImportBackup(dbPathIn(dt.dir)); !isFGCError(err, model.ErrTrackingInProgress) {
		t.Errorf("ImportBackup() error = %v, want ErrTrackingInProgress", err)
	}
}

func TestBackupData(t *testing.T) {
	dt := newDataTest(t)
	backup := filepath.Join(t.TempDir(), "my-backup.bak")
	var defaultFilename string
	dt.dialogs.SaveBackup = func(name string) (string, error) {
		defaultFilename = name
		return backup, nil
	}
	var events []string
	dt.handler.SetEventEmitter(func(name string, data ...interface{}) {
		events = append(events, fmt.Sprint(name, data))
	})
	// Backing up doesn't write, so it's fine while tracking.
	dt.startTracking()

	if err := dt.handler.BackupData(); err != nil {
		t.Fatalf("BackupData() error = %v", err)
	}
	if len(events) != 1 || events[0] != "info[backupSaved]" {
		t.Errorf("events = %v, want one backupSaved info", events)
	}
	if defaultFilename != "cfn-tracker.db.bak" {
		t.Errorf("default filename = %q, want cfn-tracker.db.bak", defaultFilename)
	}
	if !hasUser(t, backup, "1111") {
		t.Error("backup is missing the user's data")
	}
}

func TestBackupDataCancelled(t *testing.T) {
	dt := newDataTest(t)
	dt.dialogs.SaveBackup = func(string) (string, error) { return "", nil }
	dt.handler.SetEventEmitter(func(name string, data ...interface{}) {
		t.Errorf("BackupData() emitted %s %v after a cancel", name, data)
	})
	if err := dt.handler.BackupData(); err != nil {
		t.Fatalf("BackupData() error = %v, want a no-op", err)
	}
}

func TestImportBackup(t *testing.T) {
	dt := newDataTest(t)
	backup := filepath.Join(t.TempDir(), "backup.bak")
	if err := dt.sqlDb.BackupTo(context.Background(), backup); err != nil {
		t.Fatal(err)
	}

	if err := dt.handler.ImportBackup(backup); err != nil {
		t.Fatalf("ImportBackup() error = %v", err)
	}
	if !sql.HasStagedImport(dbPathIn(dt.dir)) {
		t.Error("ImportBackup() didn't stage the backup next to the database")
	}
	loc, err := dt.handler.GetDataLocation()
	if err != nil {
		t.Fatal(err)
	}
	if !loc.RestartRequired {
		t.Error("GetDataLocation() doesn't require a restart after an import")
	}
}

func TestImportBackupFollowsPendingMove(t *testing.T) {
	dt := newDataTest(t)
	target := t.TempDir()
	if err := dt.handler.SaveDataDirectory(target, false); err != nil {
		t.Fatal(err)
	}
	backup := filepath.Join(t.TempDir(), "backup.bak")
	if err := dt.sqlDb.BackupTo(context.Background(), backup); err != nil {
		t.Fatal(err)
	}
	if err := dt.handler.ImportBackup(backup); err != nil {
		t.Fatalf("ImportBackup() error = %v", err)
	}
	if !sql.HasStagedImport(dbPathIn(target)) {
		t.Error("the import wasn't staged in the folder the app opens next")
	}
}

func TestImportBackupRejectsInvalid(t *testing.T) {
	dt := newDataTest(t)
	notABackup := filepath.Join(t.TempDir(), "notes.bak")
	if err := os.WriteFile(notABackup, []byte("not a database at all"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := dt.handler.ImportBackup(notABackup); !isFGCError(err, model.ErrInvalidBackup) {
		t.Fatalf("ImportBackup() error = %v, want ErrInvalidBackup", err)
	}
	if sql.HasStagedImport(dbPathIn(dt.dir)) {
		t.Error("an invalid backup was staged")
	}
}

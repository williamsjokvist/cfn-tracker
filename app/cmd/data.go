package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	cfgDb "github.com/williamsjokvist/cfn-tracker/pkg/storage/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
)

// OpenDialogFn returns "" if the user cancels.
type OpenDialogFn func() (string, error)

// Dialogs need the Wails runtime, which only exists once the app has started, so
// main sets them then.
type Dialogs struct {
	SaveBackup    SaveFileDialogFn
	OpenBackup    OpenDialogFn
	OpenDirectory OpenDialogFn
}

type DataHandler struct {
	sqlDb        *sql.Storage
	cfgDb        *cfgDb.Storage
	tracking     *TrackingHandler
	dialogs      *Dialogs
	eventEmitter EventEmitFn
	// The folder of the open database. The configured one can differ until a restart.
	activeDir string
}

func NewDataHandler(sqlDb *sql.Storage, cfgDb *cfgDb.Storage, tracking *TrackingHandler, dialogs *Dialogs, activeDir string) *DataHandler {
	return &DataHandler{sqlDb: sqlDb, cfgDb: cfgDb, tracking: tracking, dialogs: dialogs, activeDir: activeDir}
}

func (h *DataHandler) SetEventEmitter(eventEmitter EventEmitFn) { h.eventEmitter = eventEmitter }

func (h *DataHandler) GetDataLocation() (*model.DataLocation, error) {
	dir, err := h.cfgDb.DataDir()
	if err != nil {
		return nil, model.WrapError(model.ErrGetDataLocation, err)
	}
	defaultDir, err := cfgDb.DefaultDataDir()
	if err != nil {
		return nil, model.WrapError(model.ErrGetDataLocation, err)
	}
	return &model.DataLocation{
		Dir:             dir,
		DefaultDir:      defaultDir,
		RestartRequired: !sameDir(dir, h.activeDir) || sql.HasStagedImport(dbPathIn(dir)),
	}, nil
}

func (h *DataHandler) ChooseDataDirectory() (string, error) {
	if h.dialogs == nil || h.dialogs.OpenDirectory == nil {
		return "", model.WrapError(model.ErrSaveDataDirectory, errors.New("no directory dialog"))
	}
	dir, err := h.dialogs.OpenDirectory()
	if err != nil {
		return "", model.WrapError(model.ErrSaveDataDirectory, err)
	}
	return dir, nil
}

func (h *DataHandler) HasSaveData(dir string) bool {
	if sameDir(dir, h.activeDir) {
		return false
	}
	_, err := os.Stat(dbPathIn(dir))
	return err == nil
}

func (h *DataHandler) SaveDataDirectory(dir string, useExisting bool) error {
	if h.tracking.isTracking() {
		return model.ErrTrackingInProgress
	}
	if err := h.saveDataDirectory(dir, useExisting); err != nil {
		var fgcErr *model.FGCTrackerError
		if errors.As(err, &fgcErr) {
			return err
		}
		return model.WrapError(model.ErrSaveDataDirectory, err)
	}
	return nil
}

func (h *DataHandler) saveDataDirectory(dir string, useExisting bool) error {
	info, err := os.Stat(dir)
	if err != nil {
		return fmt.Errorf("stat %s: %w", dir, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s is not a folder", dir)
	}
	// Moving back to the open database only cancels the pending move.
	if sameDir(dir, h.activeDir) {
		return h.cfgDb.SaveDataDir(dir)
	}

	target := dbPathIn(dir)
	if useExisting {
		if err := sql.CheckBackup(target); err != nil {
			return model.WrapError(model.ErrInvalidBackup, err)
		}
		return h.cfgDb.SaveDataDir(dir)
	}
	if err := h.copyNextDatabase(target); err != nil {
		return err
	}
	return h.cfgDb.SaveDataDir(dir)
}

// copyNextDatabase copies the database the app would open on the next start, so
// a pending move or import isn't lost by moving again.
func (h *DataHandler) copyNextDatabase(target string) error {
	configuredDir, err := h.cfgDb.DataDir()
	if err != nil {
		return err
	}
	next := dbPathIn(configuredDir)
	switch {
	case sql.HasStagedImport(next):
		if err := sql.CopyDatabase(sql.StagedImportPath(next), target); err != nil {
			return err
		}
		return os.Remove(sql.StagedImportPath(next))
	case !sameDir(configuredDir, h.activeDir):
		return sql.CopyDatabase(next, target)
	default:
		return h.sqlDb.BackupTo(context.Background(), target)
	}
}

func (h *DataHandler) OpenDataDirectory() error {
	dir, err := h.cfgDb.DataDir()
	if err != nil {
		return model.WrapError(model.ErrOpenDataDirectory, err)
	}
	switch runtime.GOOS {
	case "darwin":
		err = exec.Command("open", dir).Run()
	case "windows":
		// explorer.exe exits with 1 even when it opens the folder.
		_ = exec.Command("explorer.exe", dir).Run()
	}
	if err != nil {
		return model.WrapError(model.ErrOpenDataDirectory, err)
	}
	return nil
}

// Cancelling the dialog is not an error.
func (h *DataHandler) BackupData() error {
	if h.dialogs == nil || h.dialogs.SaveBackup == nil {
		return model.WrapError(model.ErrBackupData, errors.New("no save file dialog"))
	}
	path, err := h.dialogs.SaveBackup(sql.DBFileName + ".bak")
	if err != nil {
		return model.WrapError(model.ErrBackupData, err)
	}
	if path == "" {
		return nil
	}
	if err := h.sqlDb.BackupTo(context.Background(), path); err != nil {
		return model.WrapError(model.ErrBackupData, err)
	}
	if h.eventEmitter != nil {
		h.eventEmitter("info", "backupSaved")
	}
	return nil
}

func (h *DataHandler) ChooseBackup() (string, error) {
	if h.dialogs == nil || h.dialogs.OpenBackup == nil {
		return "", model.WrapError(model.ErrImportBackup, errors.New("no open file dialog"))
	}
	path, err := h.dialogs.OpenBackup()
	if err != nil {
		return "", model.WrapError(model.ErrImportBackup, err)
	}
	return path, nil
}

func (h *DataHandler) ImportBackup(path string) error {
	if h.tracking.isTracking() {
		return model.ErrTrackingInProgress
	}
	if err := sql.CheckBackup(path); err != nil {
		if errors.Is(err, sql.ErrNotABackup) || errors.Is(err, sql.ErrBackupTooNew) {
			return model.WrapError(model.ErrInvalidBackup, err)
		}
		return model.WrapError(model.ErrImportBackup, err)
	}
	dir, err := h.cfgDb.DataDir()
	if err != nil {
		return model.WrapError(model.ErrImportBackup, err)
	}
	if err := sql.StageImport(path, dbPathIn(dir)); err != nil {
		return model.WrapError(model.ErrImportBackup, err)
	}
	return nil
}

func dbPathIn(dir string) string {
	return filepath.Join(dir, sql.DBFileName)
}

func sameDir(a, b string) bool {
	return filepath.Clean(a) == filepath.Clean(b)
}

package cmd

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	wailsRuntime "github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
)

type SettingHandler struct {
	ctx   context.Context
	sqlDb *sql.Storage
}

func NewSettingHandler(storage *sql.Storage) *SettingHandler {
	return &SettingHandler{
		sqlDb: storage,
	}
}

func (sh *SettingHandler) WithContext(ctx context.Context) {
	sh.ctx = ctx
}

func (sh *SettingHandler) CreateBackup() error {
	err := sh.sqlDb.CreateBackup(sh.ctx)
	if err != nil {
		return fmt.Errorf("create backup: %w", err)
	}

	return nil
}

func (sh *SettingHandler) RestoreBackup() error {
	pid := os.Getpid()
	currentExePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("get current executable path: %w", err)
	}

	cmd := exec.Command(currentExePath, "-restore", fmt.Sprintf("-previous-pid=%d", pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("start restore process: %w", err)
	}

	wailsRuntime.Quit(sh.ctx)
	return nil
}

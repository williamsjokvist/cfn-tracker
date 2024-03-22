package cmd

import (
	"context"
	"fmt"

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
	fmt.Println("Creating backup")
	//return sh.sqlDb.CreateBackup()
	return sh.RestoreBackup()
}

func (sh *SettingHandler) RestoreBackup() error {
	fmt.Println("Restoring backup")
	return sh.sqlDb.RestoreFromBackup(sh.ctx)
}

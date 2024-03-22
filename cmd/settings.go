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

func (sh *SettingHandler) WithContext(ctx context.Context) {
	sh.ctx = ctx
}

func (sh *SettingHandler) CreateBackup() error {
	fmt.Println("Creating backup")
	return sh.sqlDb.CreateBackup(sh.ctx)
}

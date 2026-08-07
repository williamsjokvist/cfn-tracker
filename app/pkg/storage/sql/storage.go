package sql

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	migrate "github.com/golang-migrate/migrate/v4"
	sqlitex "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrationsFs embed.FS

type Storage struct {
	db *sqlx.DB
}

// inMemoryDataSource は共有キャッシュのメモリDB。素の ":memory:" は接続ごとに
// 別々のDBになるため、マイグレーションを適用しても別接続からは見えない。
const inMemoryDataSource = "file:cfn-tracker-mem?mode=memory&cache=shared"

func NewStorage(useInMemoryDb bool) (*Storage, error) {
	dataSource := getDataSource()
	if useInMemoryDb {
		dataSource = inMemoryDataSource
	}

	db, err := sqlx.Open("sqlite", dataSource)
	if err != nil {
		return nil, fmt.Errorf("open sqlite connection: %w", err)
	}

	if useInMemoryDb {
		// 共有キャッシュのメモリDBは接続が1本も無くなった時点で消滅する。
		// マイグレーション用の接続が閉じても消えないよう、ここで1本張って保持する。
		db.SetMaxOpenConns(1)
		if err := db.Ping(); err != nil {
			db.Close()
			return nil, fmt.Errorf("open in-memory sqlite: %w", err)
		}
	}

	// マイグレーションは実際に使うデータソースへ適用する。
	// 以前はメモリDB指定時もディスク側だけを移行しており、返されるハンドルには
	// テーブルが1つも無かった（テストからDBを触れなかった原因）。
	if err := migrateSchemaAt(dataSource, nil); err != nil {
		db.Close()
		return nil, fmt.Errorf("perform sql migrations: %w", err)
	}

	return &Storage{
		db,
	}, nil
}

func getDataSource() string {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	if err := os.MkdirAll(dataDir, os.FileMode(0755)); err != nil {
		return "cfn-tracker.db"
	}

	return filepath.Join(dataDir, "cfn-tracker.db")
}

func migrateSchema(nSteps *int) error {
	return migrateSchemaAt(getDataSource(), nSteps)
}

func migrateSchemaAt(dataSource string, nSteps *int) error {
	slog.Debug("starting db migrations", slog.Any("steps", nSteps))

	db, err := sqlx.Open("sqlite", dataSource)
	if err != nil {
		return fmt.Errorf("open sqlite connection: %w", err)
	}

	migrateDriver, err := sqlitex.WithInstance(db.DB, &sqlitex.Config{
		MigrationsTable: "migrations",
	})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}
	srcDriver, err := iofs.New(migrationsFs, "migrations")
	if err != nil {
		return fmt.Errorf("create migration source driver: %w", err)
	}
	preparedMigrations, err := migrate.NewWithInstance(
		"iofs",
		srcDriver,
		"",
		migrateDriver,
	)
	if err != nil {
		return fmt.Errorf("create migration tooling instance: %w", err)
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
		return fmt.Errorf("apply migrations: %w", err)
	}

	slog.Debug("applied db migrations")
	return nil
}

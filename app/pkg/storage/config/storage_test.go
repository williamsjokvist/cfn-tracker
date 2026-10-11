package config

import (
	"maps"
	"path/filepath"
	"testing"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

// isolate points the user cache directory at a temp dir, so tests don't touch
// the real config.
func isolate(t *testing.T) *Storage {
	t.Helper()
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("XDG_CACHE_HOME", filepath.Join(home, "cache"))
	t.Setenv("LocalAppData", filepath.Join(home, "AppData", "Local"))
	dir, err := DefaultDataDir()
	if err != nil {
		t.Fatalf("DefaultDataDir() error = %v", err)
	}
	storage, err := NewStorageIn(dir)
	if err != nil {
		t.Fatalf("NewStorageIn() error = %v", err)
	}
	return storage
}

func TestDataDir(t *testing.T) {
	storage := isolate(t)
	defaultDir, err := DefaultDataDir()
	if err != nil {
		t.Fatal(err)
	}

	got, err := storage.DataDir()
	if err != nil || got != defaultDir {
		t.Fatalf("DataDir() = %q, %v, want the default %q", got, err, defaultDir)
	}

	custom := t.TempDir()
	if err := storage.SaveDataDir(custom); err != nil {
		t.Fatalf("SaveDataDir() error = %v", err)
	}
	if got, err := storage.DataDir(); err != nil || got != custom {
		t.Fatalf("DataDir() = %q, %v, want %q", got, err, custom)
	}

	// Saved as empty, so the setting follows the default.
	if err := storage.SaveDataDir(defaultDir); err != nil {
		t.Fatalf("SaveDataDir() error = %v", err)
	}
	cfg, err := storage.GetRuntimeConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Data.Dir != "" {
		t.Errorf("Data.Dir = %q after saving the default, want empty", cfg.Data.Dir)
	}
}

func TestDataDirMissing(t *testing.T) {
	storage := isolate(t)
	missing := filepath.Join(t.TempDir(), "unplugged-drive")
	if err := storage.SaveDataDir(missing); err != nil {
		t.Fatalf("SaveDataDir() error = %v", err)
	}
	if got, err := storage.DataDir(); err == nil {
		t.Fatalf("DataDir() = %q, want an error for a missing folder", got)
	}
}

func TestSaveRuntimeConfigKeepsDataDir(t *testing.T) {
	storage := isolate(t)
	custom := t.TempDir()
	if err := storage.SaveDataDir(custom); err != nil {
		t.Fatal(err)
	}
	cfg, err := storage.GetRuntimeConfig()
	if err != nil {
		t.Fatal(err)
	}
	cfg.GUI.Theme = "jade"
	if err := storage.SaveRuntimeConfig(cfg); err != nil {
		t.Fatal(err)
	}
	if got, err := storage.DataDir(); err != nil || got != custom {
		t.Errorf("DataDir() = %q, %v after saving the GUI config, want %q", got, err, custom)
	}
}

func TestLastGameAndUsers(t *testing.T) {
	storage := isolate(t)
	if err := storage.SaveLastGame(model.GameTypeSF6); err != nil {
		t.Fatalf("SaveLastGame() error = %v", err)
	}
	if err := storage.SaveLastUser(model.GameTypeSF6, "1234567890"); err != nil {
		t.Fatalf("SaveLastUser() error = %v", err)
	}
	if err := storage.SaveLastUser(model.GameTypeT8, "abcd1234"); err != nil {
		t.Fatalf("SaveLastUser() error = %v", err)
	}

	// Saving the rest of the config must keep them.
	cfg, err := storage.GetRuntimeConfig()
	if err != nil {
		t.Fatal(err)
	}
	if err := storage.SaveRuntimeConfig(cfg); err != nil {
		t.Fatal(err)
	}

	cfg, err = storage.GetRuntimeConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.GUI.LastGame != model.GameTypeSF6 {
		t.Errorf("LastGame = %q, want %q", cfg.GUI.LastGame, model.GameTypeSF6)
	}
	want := map[model.GameType]string{model.GameTypeSF6: "1234567890", model.GameTypeT8: "abcd1234"}
	if !maps.Equal(cfg.GUI.LastUsers, want) {
		t.Errorf("LastUsers = %v, want %v", cfg.GUI.LastUsers, want)
	}
}

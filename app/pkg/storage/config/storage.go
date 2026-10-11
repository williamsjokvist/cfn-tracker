package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"gopkg.in/ini.v1"
)

type Storage struct {
	// The ini stays here even when the database moves, as it stores where the database is.
	dir string
}

func NewStorage() (*Storage, error) {
	dir, err := DefaultDataDir()
	if err != nil {
		return nil, fmt.Errorf("get app data directory: %w", err)
	}
	return NewStorageIn(dir)
}

func NewStorageIn(dir string) (*Storage, error) {
	storage := &Storage{dir: dir}
	_, err := storage.GetRuntimeConfig()
	if err != nil {
		return nil, fmt.Errorf("get runtime config: %w", err)
	}
	return storage, nil
}

func DefaultDataDir() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("get user cache directory: %w", err)
	}
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	err = os.MkdirAll(dataDir, os.FileMode(0755))
	if err != nil {
		return "", fmt.Errorf("create directories: %w", err)
	}
	return dataDir, nil
}

func (s *Storage) iniPath() string {
	return filepath.Join(s.dir, "cfn-tracker.ini")
}

func (s *Storage) getIniFile() (*ini.File, error) {
	path := s.iniPath()
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if _, err = os.Create(path); err != nil {
			return nil, fmt.Errorf("create config file: %w", err)
		}
	}
	iniData, err := ini.Load(path)
	if err != nil {
		return nil, fmt.Errorf("load config data: %w", err)
	}
	return iniData, nil
}

func (s *Storage) GetRuntimeConfig() (*model.RuntimeConfig, error) {
	iniData, err := s.getIniFile()
	if err != nil {
		return nil, fmt.Errorf("get config file: %w", err)
	}
	// MapTo leaves fields without a key untouched, so this sets the defaults.
	cfg := model.RuntimeConfig{GUI: model.GUIConfig{LogFile: true, BrowserSourcePort: 4242}}
	err = iniData.MapTo(&cfg)
	if err != nil {
		return nil, fmt.Errorf("map config to struct: %w", err)
	}
	cfg.GUI.LastUsers = make(map[model.GameType]string)
	for _, key := range iniData.Section("last_user").Keys() {
		cfg.GUI.LastUsers[model.GameType(key.Name())] = key.Value()
	}
	return &cfg, nil
}

func (s *Storage) SaveRuntimeConfig(cfg *model.RuntimeConfig) error {
	iniData, err := s.getIniFile()
	if err != nil {
		return fmt.Errorf("get config file: %w", err)
	}
	iniData.Section("gui").Key("locale").SetValue(cfg.GUI.Locale)
	iniData.Section("gui").Key("theme").SetValue(string(cfg.GUI.Theme))
	iniData.Section("gui").Key("sidebar").SetValue(strconv.FormatBool(cfg.GUI.SideBar))
	iniData.Section("gui").Key("log_file").SetValue(strconv.FormatBool(cfg.GUI.LogFile))
	iniData.Section("gui").Key("browser_source_port").SetValue(strconv.FormatUint(uint64(cfg.GUI.BrowserSourcePort), 10))
	iniData.Section("gui").Key("last_game").SetValue(string(cfg.GUI.LastGame))
	for game, code := range cfg.GUI.LastUsers {
		iniData.Section("last_user").Key(string(game)).SetValue(code)
	}

	if err := iniData.SaveTo(s.iniPath()); err != nil {
		return fmt.Errorf("save ini config file: %w", err)
	}
	return nil
}

// DataDir is the folder the database is in. A configured folder that's missing
// (e.g. an unplugged drive) is an error rather than a fallback to the default,
// since a new empty database there would look like the data was lost.
func (s *Storage) DataDir() (string, error) {
	cfg, err := s.GetRuntimeConfig()
	if err != nil {
		return "", fmt.Errorf("get runtime config: %w", err)
	}
	if cfg.Data.Dir == "" {
		return DefaultDataDir()
	}
	info, err := os.Stat(cfg.Data.Dir)
	if err != nil {
		return "", fmt.Errorf("save data folder %s not found, reconnect it or remove the [data] dir setting in %s: %w", cfg.Data.Dir, s.iniPath(), err)
	}
	if !info.IsDir() {
		return "", fmt.Errorf("save data folder %s is not a folder", cfg.Data.Dir)
	}
	return cfg.Data.Dir, nil
}

// SaveDataDir stores the default folder as empty, so it follows the default if that changes.
func (s *Storage) SaveDataDir(dir string) error {
	iniData, err := s.getIniFile()
	if err != nil {
		return fmt.Errorf("get config file: %w", err)
	}
	defaultDir, err := DefaultDataDir()
	if err != nil {
		return fmt.Errorf("get app data directory: %w", err)
	}
	if filepath.Clean(dir) == filepath.Clean(defaultDir) {
		dir = ""
	}
	iniData.Section("data").Key("dir").SetValue(dir)
	if err := iniData.SaveTo(s.iniPath()); err != nil {
		return fmt.Errorf("save ini config file: %w", err)
	}
	return nil
}

func (s *Storage) SaveLastGame(game model.GameType) error {
	iniData, err := s.getIniFile()
	if err != nil {
		return fmt.Errorf("get config file: %w", err)
	}
	iniData.Section("gui").Key("last_game").SetValue(string(game))
	if err := iniData.SaveTo(s.iniPath()); err != nil {
		return fmt.Errorf("save ini config file: %w", err)
	}
	return nil
}

func (s *Storage) SaveLastUser(game model.GameType, code string) error {
	iniData, err := s.getIniFile()
	if err != nil {
		return fmt.Errorf("get config file: %w", err)
	}
	iniData.Section("last_user").Key(string(game)).SetValue(code)
	if err := iniData.SaveTo(s.iniPath()); err != nil {
		return fmt.Errorf("save ini config file: %w", err)
	}
	return nil
}

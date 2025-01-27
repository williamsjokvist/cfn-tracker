package nosql

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
	"gopkg.in/ini.v1"
)

type Storage struct{}

func NewStorage() (*Storage, error) {
	storage := &Storage{}
	_, err := storage.GetRuntimeConfig()
	if err != nil {
		return nil, fmt.Errorf("get runtime config: %w", err)
	}
	return storage, nil
}

func getAppDataDir() (string, error) {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	err := os.MkdirAll(dataDir, os.FileMode(0755))
	if err != nil {
		return "", fmt.Errorf("create directories: %w", err)
	}
	return dataDir, nil
}

func (s *Storage) getIniFile() (*ini.File, error) {
	appDataDir, err := getAppDataDir()
	if err != nil {
		return nil, fmt.Errorf("get app data directory: %w", err)
	}
	path := filepath.Join(appDataDir, "cfn-tracker.ini")
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		if _, err = os.Create(path); err != nil {
			return nil, fmt.Errorf("create config file: %w", err)
		}
	}
	iniData, err := ini.Load(filepath.Join(appDataDir, "cfn-tracker.ini"))
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
	var cfg model.RuntimeConfig
	err = iniData.MapTo(&cfg)
	if err != nil {
		return nil, fmt.Errorf("map config to struct: %w", err)
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

	appDataDir, err := getAppDataDir()
	if err != nil {
		return fmt.Errorf("get app data directory: %w", err)
	}
	iniDir := filepath.Join(appDataDir, "cfn-tracker.ini")
	err = iniData.SaveTo(iniDir)
	if err != nil {
		return fmt.Errorf("save ini config file: %w", err)
	}
	return nil
}

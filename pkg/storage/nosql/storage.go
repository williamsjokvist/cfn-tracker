package nosql

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

type Storage struct {
	directory string
}

func NewStorage() (*Storage, error) {
	appDataDir, err := getDocumentPath()
	if err != nil {
		return nil, fmt.Errorf(`get app data directory: %w`, err)
	}
	jsonDir := filepath.Join(appDataDir, "cfn-tracker.json")
	storage := &Storage{
		directory: jsonDir,
	}
	_, err = storage.GetGuiConfig()
	if err != nil {
		log.Println(err)
		_, err = os.Create(jsonDir)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf(`create file: %w`, err)
		}
		err = storage.writeConfig(model.NewGuiConfig())
		if err != nil {
			return nil, fmt.Errorf(`failed to write initial config: %w`, err)
		}
	}
	return storage, nil
}

func getDocumentPath() (string, error) {
	cacheDir, _ := os.UserCacheDir()
	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	err := os.MkdirAll(dataDir, os.FileMode(0755))
	if err != nil {
		return "", fmt.Errorf(`create directories: %w`, err)
	}
	return dataDir, nil
}

func (s *Storage) GetGuiConfig() (*model.GuiConfig, error) {
	data, err := os.ReadFile(s.directory)
	if err != nil {
		return nil, fmt.Errorf(`failed to read nosql db: %w`, err)
	}
	var cfg model.GuiConfig
	err = json.Unmarshal(data, &cfg)
	if err != nil {
		return nil, fmt.Errorf(`failed to parse nosql db: %w`, err)
	}
	return &cfg, nil
}

func (s *Storage) writeConfig(cfg *model.GuiConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return fmt.Errorf(`unmarshal config: %w`, err)
	}
	err = os.WriteFile(s.directory, data, os.FileMode(0755))
	if err != nil {
		return fmt.Errorf(`failed to read nosql db: %w`, err)
	}
	return nil
}

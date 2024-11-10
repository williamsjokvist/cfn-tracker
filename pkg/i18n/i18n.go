package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"strings"

	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

//go:embed locales
var localeFs embed.FS

func GetTranslation(locale string) (*model.Localization, error) {
	localeJson, err := localeFs.ReadFile(fmt.Sprintf("locales/%s.json", locale))
	if err != nil {
		return nil, fmt.Errorf("read locale json: %w", err)
	}
	var lng model.Localization
	if err := json.Unmarshal(localeJson, &lng); err != nil {
		return nil, fmt.Errorf("unmarshal locale json: %w", err)
	}
	return &lng, nil
}

func GetSupportedLanguages() ([]string, error) {
	dirEntries, err := fs.ReadDir(localeFs, "locales")
	if err != nil {
		return nil, fmt.Errorf("read locales directory: %w", err)
	}
	var languages = make([]string, 0, len(dirEntries))
	for _, d := range dirEntries {
		languages = append(languages, strings.Split(d.Name(), ".json")[0])
	}
	return languages, nil
}

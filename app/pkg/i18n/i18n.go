package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"slices"
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

// GetTranslationMap returns the locale's translations keyed by their JSON key.
func GetTranslationMap(locale string) (map[string]string, error) {
	lng, err := GetTranslation(locale)
	if err != nil {
		return nil, err
	}
	// Round-trip through the struct so only declared keys come back.
	b, err := json.Marshal(lng)
	if err != nil {
		return nil, fmt.Errorf("marshal translation: %w", err)
	}
	var m map[string]string
	if err := json.Unmarshal(b, &m); err != nil {
		return nil, fmt.Errorf("unmarshal translation: %w", err)
	}
	return m, nil
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
	if i := slices.Index(languages, "en-GB"); i > 0 {
		languages = slices.Insert(slices.Delete(languages, i, i+1), 0, "en-GB")
	}
	return languages, nil
}

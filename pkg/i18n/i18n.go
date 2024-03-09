package i18n

import (
	"errors"

	"github.com/williamsjokvist/cfn-tracker/pkg/i18n/locales"
)

var languages = map[string]locales.Localization{
	"en-GB": locales.EN_GB,
	"fr-FR": locales.FR_FR,
	"ja-JP": locales.JA_JP,
}

func GetTranslation(locale string) (*locales.Localization, error) {
	lng, ok := languages[locale]
	if !ok {
		return nil, errors.New(`locale does not exist`)
	}
	return &lng, nil
}

func GetSupportedLanguages() []string {
	keys := make([]string, 0, len(languages))
	for k := range languages {
		keys = append(keys, k)
	}
	return keys
}

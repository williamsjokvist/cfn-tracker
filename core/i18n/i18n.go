package i18n

import (
	"github.com/williamsjokvist/cfn-tracker/core/i18n/locales"
)

func GetTranslation(locale string) locales.Localization {
	switch locale {
	case "fr-FR":
		return locales.FR_FR
	case "ja-JP":
		return locales.JA_JP
	case "en-GB":
	default:
		return locales.EN_GB
	}
	return locales.EN_GB
}

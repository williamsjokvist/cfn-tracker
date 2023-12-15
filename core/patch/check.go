package patch

import (
	"fmt"
	"log"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func CheckForUpdate(appVersion string) (bool, string, error) {
	latest, found, err := selfupdate.DetectLatest(`williamsjokvist/cfn-tracker`)
	if err != nil {
		return false, "", fmt.Errorf(`get latest app version: %w`, err)
	}

	v, err := semver.Parse(appVersion)
	if err != nil {
		return false, "", fmt.Errorf(`parse app version: %w`, err)
	}

	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return false, "", nil
	}

	log.Println("New app version available: ", latest.Version.String())
	return true, latest.Version.String(), nil
}

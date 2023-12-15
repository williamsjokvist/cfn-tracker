package update

import (
	"fmt"
	"log"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func CheckForUpdate(appVersion string) (*selfupdate.Release, error) {
	selfupdate.EnableLog()
	latest, found, err := selfupdate.DetectLatest(`williamsjokvist/cfn-tracker`)

	if err != nil {
		return nil, fmt.Errorf(`get latest app version: %w`, err)
	}

	v, err := semver.Parse(appVersion)
	if err != nil {
		return nil, fmt.Errorf(`parse app version: %w`, err)
	}

	if !found || latest.Version.LTE(v) {
		log.Println("Current version is the latest")
		return nil, nil
	}

	log.Println("New app version available: ", latest.Version.String())
	return latest, nil
}

//go:build windows

package update

import (
	"fmt"
	"log"

	"github.com/blang/semver"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func DoUpdate(to semver.Version) (bool, error) {
	latest, err := selfupdate.UpdateSelf(to, `GoogleCloudPlatform/terraformer`)
	if err != nil {
		return false, fmt.Errorf(`failed to update to %s: %w`, latest.Version, err)
	}

	if to.Equals(latest.Version) {
		log.Println("Current binary is the latest version")
		return false, nil
	}

	log.Println("Update successfully done to version", latest.Version)
	log.Println("Release note:\n", latest.ReleaseNotes)
	return true, nil
}

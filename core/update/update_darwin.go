//go:build darwin

package update

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

func DoUpdate(to *selfupdate.Release) (bool, error) {
	log.Println("Started updating to", to.Version)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return false, fmt.Errorf(`get user home dir: %w`, err)
	}
	downloadPath := filepath.Join(homeDir, "Downloads", `CFNTracker.zip`)
	err = exec.Command("curl", "-L", to.AssetURL, "-o", downloadPath).Run()
	if err != nil {
		return false, fmt.Errorf("download latest app: %w", err)
	}

	var appPath string
	cmdPath, err := os.Executable()
	if err != nil {
		appPath = `/Applications/`
	} else {
		appPath = strings.TrimSuffix(cmdPath, filepath.Join(`CFN Tracker.app`, `Contents`, `MacOS`, `CFN Tracker`))
	}

	err = exec.Command("ditto", "-xk", downloadPath, appPath).Run()
	if err != nil {
		return false, fmt.Errorf("extract archive: %w", err)
	}
	err = exec.Command("rm", downloadPath).Run()
	if err != nil {
		return false, fmt.Errorf("remove archive: %w", err)
	}

	log.Println(`Update successfully done to version`, to.Version)
	log.Println(`Release note:\n`, to.ReleaseNotes)
	return true, nil
}

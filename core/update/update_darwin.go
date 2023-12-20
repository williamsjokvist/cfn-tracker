//go:build darwin

package update

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rhysd/go-github-selfupdate/selfupdate"
)

type ProgressStatus struct {
	Progress int
	Err      error
}

func (s *ProgressStatus) WithProgress(progress int) *ProgressStatus {
	s.Progress = progress
	return s
}

func (s *ProgressStatus) WithError(err error) *ProgressStatus {
	s.Err = err
	return s
}

func DoUpdate(to *selfupdate.Release, progChan chan ProgressStatus) {
	log.Println("Started updating to", to.Version)
	status := &ProgressStatus{Progress: 0, Err: nil}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		progChan <- *status.WithError(fmt.Errorf(`get user home dir: %w`, err))
		return
	}
	downloadPath := filepath.Join(homeDir, "Downloads", `CFN_Tracker.zip`)
	cmd := exec.Command("curl", "-#", "-L", to.AssetURL, "-o", downloadPath)
	fmt.Println(cmd.String())
	// curl outputs download progress to stderr
	stderr, _ := cmd.StderrPipe()
	defer stderr.Close()

	if err := cmd.Start(); err != nil {
		progChan <- *status.WithError(fmt.Errorf("failed to download new version: %w", err))
		return
	}

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		// only parse percentage numbers, not the bar
		if !strings.Contains(m, "#") {
			progress, err := strconv.Atoi(strings.Split(m, ".")[0])
			if err != nil {
				progChan <- *status.WithError(fmt.Errorf(`failed to parse progress: %w`, err))
			}
			fmt.Println(progress)
			progChan <- *status.WithProgress(progress)
		}
	}

	if err := cmd.Wait(); err != nil {
		progChan <- *status.WithError(fmt.Errorf(`failed to download new version: %w`, err))
		return
	}

	var appPath string
	// cmdPath, err := os.Executable()
	// if err != nil {
	// 	appPath = filepath.Join(homeDir, `Applications`)
	// } else {
	// 	appPath = strings.TrimSuffix(cmdPath, filepath.Join(`CFN Tracker.app`, `Contents`, `MacOS`, `CFN Tracker`))
	// }
	appPath = filepath.Join(homeDir, `Applications`)

	err = exec.Command("ditto", "-xk", downloadPath, appPath).Run()
	if err != nil {
		progChan <- *status.WithError(fmt.Errorf(`failed to extract archive: %w`, err))
		return
	}
	err = exec.Command("rm", downloadPath).Run()
	if err != nil {
		progChan <- *status.WithError(fmt.Errorf(`failed to remove archive: %w`, err))
		return
	}

	log.Println(`Update successfully done to version`, to.Version)
	log.Println(`Release note:\n`, to.ReleaseNotes)
}

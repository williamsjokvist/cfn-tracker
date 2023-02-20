package backend

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called at application startup
func (a *App) Startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	a.Initialize(SteamUsername, SteamPassword, AppVersion)
}

// domReady is called after front-end resources have been loaded
func (a App) DomReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) BeforeClose(ctx context.Context) (prevent bool) {
	PageInstance.Browser().Close()
	return false
}

// shutdown is called at application termination
func (a *App) Shutdown(ctx context.Context) {

}

func (a *App) GetAppVersion() string {
	if AppVersion != nil {
		return AppVersion.Original()
	}
	return ""
}

func (a *App) IsTracking() bool {
	return IsTracking
}

func (a *App) IsInitialized() bool {
	return IsInitialized
}

func (a *App) StopTracking() {
	IsTracking = false
}

func (a *App) GetMatchHistory() MatchHistory {
	return CurrentMatchHistory
}

func (a *App) OpenResultsDirectory() {
	var cmd *exec.Cmd

	if runtime.GOOS == `darwin` {
		cmd = exec.Command(`Open`, `./results`)
	} else {
		cmd = exec.Command(`explorer.exe`, `.\results`)
	}

	cmd.Run()
}

func (a *App) GetMatchLog(cfn string) []MatchHistory {
	var matchLog []MatchHistory

	pastMatches, _ := os.ReadFile(`results/` + cfn + `-log.json`)
	_ = json.Unmarshal(pastMatches, &matchLog)

	return matchLog
}

func (a *App) GetAvailableLogs() []string {
	files, err := ioutil.ReadDir("results")
	if err != nil {
		fmt.Println(`Failed to read result directory`, err)
	}

	cfns := []string{}

	for _, file := range files {
		fileName := file.Name()

		if !strings.Contains(fileName, `-log.json`) {
			continue
		}

		cfn := strings.Split(fileName, `-log.json`)[0]
		cfns = append(cfns, cfn)
	}

	return cfns
}

func (a *App) ResultsJSONExist() bool {
	_, err := os.Stat(`results/results.json`)
	return !os.IsNotExist(err)
}

func (a *App) DeleteMatchLog(cfn string) {
	err := os.Remove(`results/` + cfn + `-log.json`)
	if err != nil {
		fmt.Println(err)
	}
}

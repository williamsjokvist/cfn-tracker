package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"runtime"

	"os/exec"
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
func (a *App) startup(ctx context.Context) {
	// Perform your setup here
	a.ctx = ctx
	Initialize()
}

// domReady is called after front-end resources have been loaded
func (a App) domReady(ctx context.Context) {
	// Add your action here
}

// beforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue, false will continue shutdown as normal.
func (a *App) beforeClose(ctx context.Context) (prevent bool) {
	pageInstance.Browser().Close()
	return false
}

// shutdown is called at application termination
func (a *App) shutdown(ctx context.Context) {

}

func (a *App) Track(cfnName string) bool {
	go StartTracking(cfnName)
	return isInitialized
}

func (a *App) IsTracking() bool {
	return isTracking
}

func (a *App) IsInitialized() bool {
	return isInitialized
}

func (a *App) StopTracking() {
	isTracking = false
}

func (a *App) GetMatchHistory() MatchHistory {
	return matchHistory
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

func (a *App) GetMatchLog() []MatchHistory {
	var matchLog []MatchHistory

	pastMatches, _ := os.ReadFile(`results/match-history.json`)
	_ = json.Unmarshal(pastMatches, &matchLog)

	return matchLog
}

func (a *App) DeleteMatchLog() {
	err := os.Remove(`results/match-history.json`)
	if err != nil {
		fmt.Println(err)
	}
}

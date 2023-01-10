package main

import (
	"fmt"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/briandowns/spinner"
	"github.com/joho/godotenv"
)

type Config struct {
	CFN string
}

var profile string
var progressBar = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

func init() {
	err := godotenv.Load(".env")

	if err != nil {
		LogError(`Error loading environment variables. Are you missing a .env file?`)
	}
}

func main() {
	progressBar.Start()
	progressBar.HideCursor = true
	progressBar.Color(`yellow`)

	f := `config.toml`
	if _, err := os.Stat(f); err != nil {
		f = `config.toml`
	}

	var config Config

	_, err := toml.DecodeFile(f, &config)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
		return
	}

	if config.CFN == "" {
		LogError(`CFN profile not set`)
		return
	}

	profile = config.CFN

	StartTracking(profile)
}

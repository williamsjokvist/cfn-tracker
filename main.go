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

var progressBar = spinner.New(spinner.CharSets[9], 100*time.Millisecond)

func init() {
	if os.Getenv(`TESTING`) == `true` {
		return
	}

	err := godotenv.Load(".env")

	if err != nil {
		LogError(MissingEnvError)
	}
}

func main() {
	f := `cfn-scraper-config.toml`
	if _, err := os.Stat(f); err != nil {
		f = `cfn-scraper-config.toml`
	}

	var config Config
	var profile string

	_, err := toml.DecodeFile(f, &config)

	if err != nil {
		fmt.Println(`No CFN account configured, please input a valid CFN account to track. You can change it later in the config file.`)
		var inputText string

		_, err := fmt.Scanln(&inputText)

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
			return
		}

		profile = inputText

		SaveTextToFile(``, `cfn-scraper-config.toml`, `CFN = "`+profile+`"`)
	} else {
		profile = config.CFN
	}

	progressBar.Start()
	progressBar.HideCursor = true
	progressBar.Color(`yellow`)

	StartTracking(profile)
}

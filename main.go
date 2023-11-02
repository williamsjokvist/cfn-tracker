package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/williamsjokvist/cfn-tracker/core"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
)

var (
	steamUsername string = ``
	steamPassword string = ``
	capIDEmail    string = ``
	capIDPassword string = ``
	appVersion    string = ``
	runHeadless   string = ``
)

//go:embed all:gui/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func init() {
	cleanUpProcess()

	if capIDEmail == `` || capIDPassword == `` ||
		steamUsername == `` || steamPassword == `` || appVersion == `` {
		err := godotenv.Load(`.env`)
		if err != nil {
			log.Fatalf(`missing .env file: %v`, err)
		}

		steamUsername = os.Getenv(`STEAM_USERNAME`)
		steamPassword = os.Getenv(`STEAM_PASSWORD`)
		capIDEmail = os.Getenv(`CAP_ID_EMAIL`)
		capIDPassword = os.Getenv(`CAP_ID_PASSWORD`)
		appVersion = os.Getenv(`APP_VERSION`)
		runHeadless = os.Getenv(`RUN_HEADLESS`)
	}

	core.AppVersion, _ = version.NewVersion(appVersion)
	core.SteamUsername = steamUsername
	core.SteamPassword = steamPassword
	core.RunHeadless = runHeadless == `true`
	core.CapIDEmail = capIDEmail
	core.CapIDPassword = capIDPassword
}

func main() {
	sqlDb, err := sql.NewStorage()
	if err != nil {
		log.Fatalf("init db: %v", err)
	}
	trackerRepo := data.NewCFNTrackerRepository(sqlDb)
	cmdHandler := core.NewCommandHandler(trackerRepo)

	err = wails.Run(&options.App{
		Title:              `CFN Tracker v3`,
		Assets:             assets,
		Width:              920,
		Height:             450,
		MinWidth:           800,
		MinHeight:          450,
		DisableResize:      true,
		Fullscreen:         false,
		Frameless:          true,
		StartHidden:        false,
		HideWindowOnClose:  false,
		LogLevel:           logger.WARNING,
		LogLevelProduction: logger.ERROR,
		BackgroundColour:   options.NewRGBA(0, 0, 0, 1),
		CSSDragProperty:    `--draggable`,
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			Theme:                             windows.Theme(windows.Dark),
			DisableFramelessWindowDecorations: true,
		},
		Mac: &mac.Options{
			TitleBar: mac.TitleBarHiddenInset(),
			About: &mac.AboutInfo{
				Title:   fmt.Sprintf(`CFN Tracker v%s`, appVersion),
				Message: fmt.Sprintf(`CFN Tracker version %s © 2023 William Sjökvist <william.sjokvist@gmail.com>`, appVersion),
			},
		},
		OnStartup:  cmdHandler.StartBrowser,
		OnShutdown: cmdHandler.CloseBrowser,
		OnBeforeClose: func(ctx context.Context) (prevent bool) {
			cmdHandler.CloseBrowser(ctx)
			return false
		},
		Bind: []interface{}{
			cmdHandler,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

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
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core"
	"github.com/williamsjokvist/cfn-tracker/core/browser"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
	"github.com/williamsjokvist/cfn-tracker/core/server"
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
	appBrowser, err := browser.NewBrowser(runHeadless == `true`)
	if err != nil {
		log.Fatalf(`failed to launch browser: %v`, err)
	}
	closeWithError := func(err error) {
		appBrowser.Page.Browser().Close()

		// TODO: log error and/or show error message?
		log.Fatal(err)
	}

	sqlDb, err := sql.NewStorage()
	if err != nil {
		closeWithError(err)
	}
	trackerRepo := data.NewCFNTrackerRepository(sqlDb)
	cmdHandler := core.NewCommandHandler(appBrowser, trackerRepo)

	appVer, err := version.NewVersion(appVersion)
	if err != nil {
		closeWithError(err)
	}
	latestVersion, err := appBrowser.GetLatestAppVersion()
	if err != nil {
		closeWithError(err)
	}

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
		OnStartup: func(ctx context.Context) {
			cmdHandler.AssignRuntimeContext(ctx)

			if appVer.LessThan(latestVersion) {
				log.Println(`Has new version: `, latestVersion.String())
				runtime.EventsEmit(ctx, `version-update`, latestVersion.String())
			} else {
				log.Println(`No new version, running: `, appVer.String())
			}

			go server.Start(ctx)
		},
		OnShutdown: func(_ context.Context) {
			appBrowser.Page.Browser().Close()
		},
		OnBeforeClose: func(_ context.Context) (prevent bool) {
			appBrowser.Page.Browser().Close()
			return false
		},
		Bind: []interface{}{
			cmdHandler,
		},
	})
	if err != nil {
		closeWithError(err)
	}
}

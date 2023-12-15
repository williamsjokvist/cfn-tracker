package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/williamsjokvist/cfn-tracker/core"
	"github.com/williamsjokvist/cfn-tracker/core/browser"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
	"github.com/williamsjokvist/cfn-tracker/core/data/txt"
	"github.com/williamsjokvist/cfn-tracker/core/errorsx"
	"github.com/williamsjokvist/cfn-tracker/core/server"
	"github.com/williamsjokvist/cfn-tracker/core/update"
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

//go:embed gui/error/error.html
var errorTmpl []byte

//go:embed gui/update/update.html
var updateTmpl []byte

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

	core.AppVersion = appVersion
	core.SteamUsername = steamUsername
	core.SteamPassword = steamPassword
	core.RunHeadless = runHeadless == `true`
	core.CapIDEmail = capIDEmail
	core.CapIDPassword = capIDPassword
}

func main() {
	var appBrowser *browser.Browser
	closeWithError := func(err error) {
		if appBrowser != nil {
			appBrowser.Page.Browser().Close()
		}
		log.Println("close with error", err)
		wails.Run(&options.App{
			Title:                    `CFN Tracker - Error`,
			Width:                    400,
			Height:                   148,
			DisableResize:            true,
			Frameless:                true,
			EnableDefaultContextMenu: false,
			AssetServer: &assetserver.Options{
				Middleware: assetserver.ChainMiddleware(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						var b bytes.Buffer
						tmpl := template.Must(template.New("errorPage").Parse(string(errorTmpl)))
						params := struct {
							Error string
						}{
							Error: err.Error(),
						}
						tmpl.Execute(&b, params)
						w.Write(b.Bytes())
					})
				}),
			},
		})
	}

	hasNewVersion, version, err := update.CheckForUpdate(appVersion)
	if err != nil {
		closeWithError(fmt.Errorf(`failed to check for update: %w`, err))
		return
	}

	if hasNewVersion {
		wails.Run(&options.App{
			Title:                    `CFN Tracker - Update`,
			Width:                    400,
			Height:                   148,
			DisableResize:            true,
			Frameless:                true,
			EnableDefaultContextMenu: false,
			AssetServer: &assetserver.Options{
				Middleware: assetserver.ChainMiddleware(func(next http.Handler) http.Handler {
					return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
						var b bytes.Buffer
						tmpl := template.Must(template.New("updatePage").Parse(string(updateTmpl)))
						params := struct {
							Version string
						}{
							Version: version,
						}
						tmpl.Execute(&b, params)
						w.Write(b.Bytes())
					})
				}),
			},
		})
		return
	}

	appBrowser, err = browser.NewBrowser(runHeadless == `true`)
	if err != nil {
		closeWithError(fmt.Errorf(`failed to launch browser: %v`, err))
		return
	}
	sqlDb, err := sql.NewStorage()
	if err != nil {
		closeWithError(fmt.Errorf(`failed to initalize database: %w`, err))
		return
	}
	txtDb, err := txt.NewStorage()
	if err != nil {
		closeWithError(fmt.Errorf(`failed to initalize text store: %w`, err))
		return
	}
	trackerRepo := data.NewCFNTrackerRepository(sqlDb, txtDb)
	cmdHandler := core.NewCommandHandler(appBrowser, trackerRepo)

	err = wails.Run(&options.App{
		Title: fmt.Sprintf(`CFN Tracker v%s`, appVersion),
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
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
		ErrorFormatter:     errorsx.ConvertToFrontEndError,
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
		closeWithError(fmt.Errorf(`failed to launch app: %w`, err))
	}
}

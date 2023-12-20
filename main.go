package main

import (
	"bytes"
	"context"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"github.com/williamsjokvist/cfn-tracker/core"
	"github.com/williamsjokvist/cfn-tracker/core/browser"
	"github.com/williamsjokvist/cfn-tracker/core/config"
	"github.com/williamsjokvist/cfn-tracker/core/data"
	"github.com/williamsjokvist/cfn-tracker/core/data/sql"
	"github.com/williamsjokvist/cfn-tracker/core/data/txt"
	"github.com/williamsjokvist/cfn-tracker/core/errorsx"
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

//go:embed gui/error/error.html
var errorTmpl []byte

//go:embed build/appicon.png
var icon []byte

var cfg config.Config

func init() {
	cleanUpProcess()
	err := godotenv.Load(`.env`)
	if err != nil {
		log.Println(fmt.Errorf(`missing .env file: %w`, err))
		cfg = config.Config{
			AppVersion:        appVersion,
			Headless:          runHeadless == `true`,
			SteamUsername:     steamUsername,
			SteamPassword:     steamPassword,
			CapIDEmail:        capIDEmail,
			CapIDPassword:     capIDPassword,
			BrowserSourcePort: 4242,
		}
		return
	}
	err = envconfig.Process("", &cfg)
	if err != nil {
		log.Fatalf(`failed to process envconfig: %v`, err)
	}
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
	appBrowser, err := browser.NewBrowser(cfg.Headless)
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
	cmdHandler := core.NewCommandHandler(appBrowser, trackerRepo, &cfg)

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
			cmdHandler.SetContext(ctx)
			go server.Start(ctx, &cfg)
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

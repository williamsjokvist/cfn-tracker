package main

import (
	"cfnscraper/backend"
	"embed"
	"log"
	"os"

	"github.com/hashicorp/go-version"
	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var (
	steamUsername string = ``
	steamPassword string = ``
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var WailsApp *backend.App

func init() {
	if steamUsername == `` && steamPassword == `` {
		err := godotenv.Load(`.env`)
		if err != nil {
			log.Fatal(err)
		}
		steamUsername = os.Getenv(`STEAM_USERNAME`)
		steamPassword = os.Getenv(`STEAM_PASSWORD`)
	}
}

func main() {
	appVersion, _ := version.NewVersion(`1.9.0`)

	// Create an instance of the app structure
	WailsApp = backend.NewApp()
	backend.AppVersion = appVersion
	backend.SteamUsername = steamUsername
	backend.SteamPassword = steamPassword

	err := wails.Run(&options.App{
		Title:             `CFN Tracker v2`,
		Assets:            assets,
		Width:             800,
		Height:            450,
		MinWidth:          800,
		MinHeight:         450,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         true,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  options.NewRGBA(0, 0, 0, 0),
		CSSDragProperty:   `--wails-draggable`,
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			Theme:                             windows.Theme(windows.Dark),
			DisableFramelessWindowDecorations: true,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			Appearance:           mac.AppearanceType(mac.NSAppearanceNameAccessibilityHighContrastDarkAqua),
			About: &mac.AboutInfo{
				Title:   "CFN Tracker " + appVersion.Original(),
				Message: "Version " + appVersion.Original() + " © 2022 William Sjökvist <william.sjokvist@gmail.com>",
			},
		},
		OnStartup:     WailsApp.Startup,
		OnDomReady:    WailsApp.DomReady,
		OnShutdown:    WailsApp.Shutdown,
		OnBeforeClose: WailsApp.BeforeClose,
		Bind: []interface{}{
			WailsApp,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

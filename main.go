package main

import (
	"embed"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

var (
	steamUsername string = ``
	steamPassword string = ``
	profile       string
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

var WailsApp *App

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
	// Create an instance of the app structure
	WailsApp = NewApp()

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
				Title:   "CFN Tracker v2",
				Message: "Version 2.1.0 © 2022 William Sjökvist <william.sjokvist@gmail.com>",
			},
		},
		OnStartup:     WailsApp.startup,
		OnDomReady:    WailsApp.domReady,
		OnShutdown:    WailsApp.shutdown,
		OnBeforeClose: WailsApp.beforeClose,
		Bind: []interface{}{
			WailsApp,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

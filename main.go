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
	steamUsername string
	steamPassword string
	profile       string
)

//go:embed all:frontend/dist
var assets embed.FS

//go:embed build/appicon.png
var icon []byte

func init() {
	if os.Getenv(`EXECUTION_ENVIRONMENT`) != "test" && os.Getenv(`EXECUTION_ENVIRONMENT`) != "build" {
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
	app := NewApp()

	err := wails.Run(&options.App{
		Title:             `CFN Tracker v2`,
		Assets:            assets,
		Width:             800,
		Height:            450,
		MinWidth:          800,
		MinHeight:         450,
		DisableResize:     true,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		CSSDragProperty:   `--wails-draggable`,
		BackgroundColour:  &options.RGBA{R: 33, G: 37, B: 43, A: 255},
		Windows: &windows.Options{
			WebviewIsTransparent:              true,
			WindowIsTranslucent:               false,
			Theme:                             windows.Theme(windows.Dark),
			DisableFramelessWindowDecorations: true,
		},
		Mac: &mac.Options{
			TitleBar:             mac.TitleBarHiddenInset(),
			Appearance:           mac.NSAppearanceNameAccessibilityHighContrastVibrantDark,
			WebviewIsTransparent: true,
			WindowIsTranslucent:  true,
			About: &mac.AboutInfo{
				Title:   "CFN Tracker v2",
				Message: "Version 2.0.0 © 2022 William Sjökvist <william.sjokvist@gmail.com>",
			},
		},
		OnStartup:  app.startup,
		OnShutdown: app.shutdown,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		log.Fatal(err)
	}
}

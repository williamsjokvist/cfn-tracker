package main

import (
	"bytes"
	"context"
	"embed"
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/logger"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/cmd"
	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/errorsx"
	"github.com/williamsjokvist/cfn-tracker/pkg/server"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/nosql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/sql"
	"github.com/williamsjokvist/cfn-tracker/pkg/storage/txt"
)

var (
	steamUsername string = ``
	steamPassword string = ``
	capIDEmail    string = ``
	capIDPassword string = ``
	appVersion    string = ``
	isProduction  string = ``
)

type args struct {
	restoreBackup bool
	previousPID   int
}

//go:embed all:gui/dist
var assets embed.FS

//go:embed gui/error/error.html
var errorTmpl []byte

var cfg config.Config
var logFile *os.File

func logToFile() {
	file, err := os.OpenFile(`cfn-tracker.log`, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}

	logFile = file
	log.SetOutput(file)
	log.SetFlags(log.Ldate | log.LstdFlags | log.Lshortfile)
}

func init() {
	if isProduction == `true` {
		logToFile()
	}

	err := godotenv.Load(`.env`)
	if err != nil {
		log.Println(fmt.Errorf(`missing .env file: %w`, err))
		cfg = config.Config{
			AppVersion:        appVersion,
			Headless:          isProduction == `true`,
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
	defer func() {
		if logFile != nil {
			logFile.Close()
		}
	}()

	args := args{}
	flag.BoolVar(&args.restoreBackup, "restore", false, "Trigger backup restore on app start")
	flag.IntVar(&args.previousPID, "previous-pid", -1, "PID of previous instance")
	flag.Parse()

	if args.restoreBackup && args.previousPID > 0 {
		err := handleRestoreBackup(context.TODO(), args.previousPID)
		if err != nil {
			log.Printf("failed to restore backup: %s", err.Error())
		}
	}

	appBrowser, err := browser.NewBrowser(cfg.Headless)
	if err != nil {
		closeWithError(fmt.Errorf(`failed to launch browser: %v`, err))
		return
	}
	defer func() {
		if err := appBrowser.Page.Browser().Close(); err != nil {
			log.Printf("failed to close browser: %v", err)
		}
	}()

	sqlDb, err := sql.NewStorage()
	if err != nil {
		closeWithError(fmt.Errorf(`failed to initalize database: %w`, err))
		return
	}
	noSqlDb, err := nosql.NewStorage()
	if err != nil {
		closeWithError(fmt.Errorf(`failed to initalize app config: %w`, err))
		return
	}
	txtDb, err := txt.NewStorage()
	if err != nil {
		closeWithError(fmt.Errorf(`failed to initalize text store: %w`, err))
		return
	}
	cmdHandler := cmd.NewCommandHandler(appBrowser, sqlDb, noSqlDb, txtDb, &cfg)
	settingsHandler := cmd.NewSettingHandler(sqlDb)

	pprofServer := &http.Server{
		Addr: ":6060",
	}
	go func() {
		if err := pprofServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP pprof server error %v", err)
		}
	}()
	defer func() {
		if err := pprofServer.Close(); err != nil {
			log.Printf("error closing pprof server: %v", err)
		}
	}()

	err = wails.Run(newApp(
		cfg.AppVersion,
		assets,
		appBrowser,
		cmdHandler,
		settingsHandler,
	))
	if err != nil {
		closeWithError(fmt.Errorf(`failed to launch app: %w`, err))
	}
}

func newApp(
	appVersion string,
	assets embed.FS,
	appBrowser *browser.Browser,
	cmdHandler *cmd.CommandHandler,
	settingsHandler *cmd.SettingHandler,
) *options.App {
	var wailsCtx context.Context
	return &options.App{
		Title:              fmt.Sprintf(`CFN Tracker v%s`, appVersion),
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
		ErrorFormatter:     errorsx.FormatError,
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
			wailsCtx = ctx
			cmdHandler.SetContext(ctx)
			settingsHandler.WithContext(ctx)
			go server.Start(ctx, &cfg)
		},
		OnShutdown: func(_ context.Context) {
			appBrowser.Page.Browser().Close()
		},
		OnBeforeClose: func(_ context.Context) (prevent bool) {
			appBrowser.Page.Browser().Close()
			return false
		},
		SingleInstanceLock: &options.SingleInstanceLock{
			UniqueId: "d0ef6612-49f7-437a-9ffc-2076ec9e37db",
			OnSecondInstanceLaunch: func(secondInstanceData options.SecondInstanceData) {
				log.Println("user opened second instance", strings.Join(secondInstanceData.Args, ","))
				log.Println("user opened second from", secondInstanceData.WorkingDirectory)
				runtime.WindowUnminimise(wailsCtx)
				runtime.Show(wailsCtx)
				go runtime.EventsEmit(wailsCtx, "launchArgs", secondInstanceData.Args)
			},
		},
		Bind: []interface{}{
			cmdHandler,
			settingsHandler,
		},
	}
}

func closeWithError(err error) {
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

func handleRestoreBackup(ctx context.Context, prevPID int) error {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	err := waitForProcessClosure(ctx, prevPID)
	if err != nil {
		return fmt.Errorf("err waiting for process to close: %w", err)
	}

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("get user cache dir: %w", err)
	}

	dataDir := filepath.Join(cacheDir, "cfn-tracker")
	os.MkdirAll(dataDir, os.FileMode(0755))
	backupDbFilepath := filepath.Join(dataDir, "cfn-tracker.backup.db")
	currentDbFilepath := filepath.Join(dataDir, "cfn-tracker.db")

	err = os.Rename(currentDbFilepath, fmt.Sprintf("%s.restored.%d", currentDbFilepath, time.Now().Unix()))
	if err != nil {
		return fmt.Errorf("rename current db: %w", err)
	}
	err = os.Rename(backupDbFilepath, currentDbFilepath)
	if err != nil {
		return fmt.Errorf("rename backup db: %w", err)
	}

	return nil
}

func waitForProcessClosure(ctx context.Context, pid int) error {
	errChan := make(chan error)
	go func() {
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			p, err := os.FindProcess(pid)
			if err != nil {
				errChan <- err
				return
			}
			if p == nil {
				errChan <- nil
				return
			}

			err = p.Signal(syscall.Signal(0))
			if err != nil {
				errChan <- nil
				return
			}

			select {
			case <-ctx.Done():
				errChan <- ctx.Err()
				return
			case <-ticker.C:
				log.Println("ticking")
			}
		}
	}()

	return <-errChan
}

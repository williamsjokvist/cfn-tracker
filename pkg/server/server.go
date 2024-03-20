package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

//go:embed static
var staticFs embed.FS

func GetInternalThemes() []model.Theme {
	var themes = make([]model.Theme, 0, 10)

	err := fs.WalkDir(staticFs, "static/themes", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		b, err := fs.ReadFile(staticFs, path)
		if err != nil {
			return err
		}

		themes = append(themes, model.Theme{
			Name: strings.Split(d.Name(), ".css")[0],
			CSS:  string(b),
		})
		return nil
	})

	if err != nil {
		return []model.Theme{}
	}

	return themes
}

func Start(ctx context.Context, cfg *config.Config) error {
	log.Println(`Starting browser source server`)

	var mhJson *[]byte

	wails.EventsOn(ctx, `cfn-data`, func(incomingData ...interface{}) {
		mh, ok := incomingData[0].(*model.TrackingState)
		if !ok {
			return
		}
		js, _ := json.Marshal(mh)
		mhJson = &js
	})

	http.HandleFunc("GET /themes/{theme}", func(w http.ResponseWriter, req *http.Request) {
		fileName := req.PathValue("theme")
		css, err := staticFs.ReadFile(fmt.Sprintf("static/themes/%s", fileName))

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.Header().Set(`Content-Type`, `text/css`)
			w.WriteHeader(http.StatusOK)
			w.Write(css)
		}
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, _ *http.Request) {
		html, err := staticFs.ReadFile("static/index.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			w.Header().Set(`Content-Type`, `text/html`)
			w.WriteHeader(http.StatusOK)
			w.Write(html)
		}
	})

	http.HandleFunc(`GET /stream`, func(w http.ResponseWriter, _ *http.Request) {
		flusher, ok := w.(http.Flusher)
		var lastJson *[]byte = nil
		if !ok {
			http.Error(w, `SSE not supported`, http.StatusInternalServerError)
			return
		}

		w.Header().Set(`Content-Type`, `text/event-stream`)
		w.Header().Set(`Cache-Control`, `no-cache`)
		w.Header().Set(`Connection`, `keep-alive`)

		ticker := time.NewTicker(time.Second * 5)
		for range ticker.C {
			if mhJson == lastJson {
				continue
			}

			fmt.Fprint(w, "event: message\n\n")
			fmt.Fprintf(w, "data: %s\n\n", *mhJson)
			lastJson = mhJson
			flusher.Flush()
		}
	})

	// custom themes
	fs := http.FileServer(http.Dir("./themes"))
	http.Handle("/themes/", http.StripPrefix("/themes/", fs))

	if err := http.ListenAndServe(fmt.Sprintf(`:%d`, cfg.BrowserSourcePort), nil); err != nil {
		return fmt.Errorf(`failed to launch browser source server: %v`, err)
	}
	return nil
}

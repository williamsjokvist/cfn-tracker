package server

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/config"
	"github.com/williamsjokvist/cfn-tracker/core/model"
)

//go:embed static/index.html
var page []byte

//go:embed static/themes/default.css
var defaultCss []byte

//go:embed static/themes/blades.css
var bladesCss []byte

//go:embed static/themes/jaeger.css
var jaegerCss []byte

//go:embed static/themes/nord.css
var nordCss []byte

//go:embed static/themes/pills.css
var pillsCss []byte

func GetInternalThemes() []model.Theme {
	return []model.Theme{
		{
			Name: "default",
			CSS:  string(defaultCss),
		},
		{
			Name: "jaeger",
			CSS:  string(jaegerCss),
		},
		{
			Name: "nord",
			CSS:  string(nordCss),
		},
		{
			Name: "pills",
			CSS:  string(pillsCss),
		},
		{
			Name: "blades",
			CSS:  string(bladesCss),
		},
	}
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

	fs := http.FileServer(http.Dir("./themes"))
	http.Handle("/themes/", http.StripPrefix("/themes/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/html`)
		w.WriteHeader(http.StatusOK)
		w.Write(page)
	})

	http.HandleFunc("/default.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/css`)
		w.WriteHeader(http.StatusOK)
		w.Write(defaultCss)
	})

	http.HandleFunc("/blades.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/css`)
		w.WriteHeader(http.StatusOK)
		w.Write(bladesCss)
	})

	http.HandleFunc("/jaeger.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/css`)
		w.WriteHeader(http.StatusOK)
		w.Write(jaegerCss)
	})

	http.HandleFunc("/nord.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/css`)
		w.WriteHeader(http.StatusOK)
		w.Write(nordCss)
	})

	http.HandleFunc("/pills.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/css`)
		w.WriteHeader(http.StatusOK)
		w.Write(pillsCss)
	})

	http.HandleFunc(`/stream`, func(w http.ResponseWriter, _ *http.Request) {
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

	if err := http.ListenAndServe(fmt.Sprintf(`:%d`, cfg.BrowserSourcePort), nil); err != nil {
		return fmt.Errorf(`failed to launch browser source server: %v`, err)
	}
	return nil
}

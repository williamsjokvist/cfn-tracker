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
	"github.com/williamsjokvist/cfn-tracker/core/model"
)

//go:embed static/index.html
var page []byte

//go:embed static/index.js
var js []byte

//go:embed static/default.css
var css []byte

const PORT = 4242

func Start(ctx context.Context) error {
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

	http.HandleFunc("/index.js", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/javascript`)
		w.WriteHeader(http.StatusOK)
		w.Write(js)
	})

	http.HandleFunc("/default.css", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set(`Content-Type`, `text/css`)
		w.WriteHeader(http.StatusOK)
		w.Write(css)
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

	if err := http.ListenAndServe(fmt.Sprintf(`:%d`, PORT), nil); err != nil {
		return fmt.Errorf(`failed to launch browser source server: %v`, err)
	}
	return nil
}

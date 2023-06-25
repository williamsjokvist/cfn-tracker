package common

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func Serve(ctx context.Context) {
	fmt.Println("Starting server")

	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(`Content-Type`, `text/html`)
			http.ServeFile(w, r, "server/index.html")
		})

		http.HandleFunc("/index.js", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set(`Content-Type`, `text/javascript`)
			http.ServeFile(w, r, "server/index.js")
		})

		http.HandleFunc("/style.css", func(w http.ResponseWriter, r *http.Request) {
			fmt.Println(`style.css`)
			w.Header().Set(`Content-Type`, `text/css`)
			http.ServeFile(w, r, "server/style.css")
		})

		var mhJson *[]byte

		runtime.EventsOn(ctx, `cfn-data`, func(data ...interface{}) {
			mh, ok := data[0].(*MatchHistory)
			if !ok {
				return
			}

			js, _ := json.Marshal(mh)
			mhJson = &js
		})

		http.HandleFunc("/stream", func(w http.ResponseWriter, _ *http.Request) {
			flusher, ok := w.(http.Flusher)
			var lastJson *[]byte = nil
			if !ok {
				http.Error(w, "SSE not supported", http.StatusInternalServerError)
				return
			}

			w.Header().Set(`Content-Type`, `text/event-stream`)
			w.Header().Set("Cache-Control", "no-cache")
			w.Header().Set(`Connection`, `keep-alive`)

			for {
				if mhJson == lastJson {
					continue
				}

				fmt.Fprint(w, "event: message\n\n")
				fmt.Fprintf(w, "data: %s\n\n", *mhJson)
				lastJson = mhJson
				flusher.Flush()
				time.Sleep(1 * time.Second)
			}
		})

		if err := http.ListenAndServe(":8080", nil); err != nil {
			fmt.Println(err)
		}
	}()

	fmt.Println("Server listening on 8080")
}

package shared

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	wails "github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/williamsjokvist/cfn-tracker/core/data"
)

//go:embed static/index.html
var page []byte

//go:embed static/index.js
var js []byte

//go:embed static/default.css
var css []byte

const PORT = 4242

func Serve(ctx context.Context) {
	fmt.Println(`Starting Browser Source Server`)

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

	var mhJson *[]byte

	wails.EventsOn(ctx, `cfn-data`, func(incomingData ...interface{}) {
		mh, ok := incomingData[0].(*data.MatchHistory)
		if !ok {
			return
		}

		js, _ := json.Marshal(mh)
		mhJson = &js
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
		fmt.Println(err)
	}

	fmt.Println(`Browser Source Server listening on `, PORT)
}

func GetThemeList() ([]string, error) {
	files, err := ioutil.ReadDir(`themes`)
	if err != nil {
		return nil, fmt.Errorf(`read themes directory: %w`, err)
	}

	themes := []string{}

	for _, file := range files {
		fileName := file.Name()

		if !strings.Contains(fileName, `.css`) {
			continue
		}

		theme := strings.Split(fileName, `.css`)[0]
		themes = append(themes, theme)
	}

	return themes, nil
}

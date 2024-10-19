package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"strings"

	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

//go:embed static
var staticFs embed.FS

type BrowserSourceServer struct {
	matchChan chan model.Match
	lastMatch []byte
}

func NewBrowserSourceServer(matchChan chan model.Match) *BrowserSourceServer {
	return &BrowserSourceServer{
		matchChan: matchChan,
		lastMatch: nil,
	}
}

func (b *BrowserSourceServer) Start(ctx context.Context, cfg *config.Config) {
	log.Println("Starting browser source server")

	http.HandleFunc("/", b.handleRoot)
	http.HandleFunc(`GET /stream`, b.handleStream)
	http.HandleFunc("GET /themes/{theme}", b.handleTheme)

	// serve custom themes through "themes" directory in the same directory as the user's executable
	fs := http.FileServer(http.Dir("./themes"))
	http.Handle("/themes/", http.StripPrefix("/themes/", fs))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.BrowserSourcePort), nil); err != nil {
		log.Println("failed to launch browser source server", err)
	}
}

func (b *BrowserSourceServer) handleStream(w http.ResponseWriter, _ *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	if b.lastMatch != nil {
		fmt.Fprint(w, "event: message\n\n")
		fmt.Fprintf(w, "data: %s\n\n", b.lastMatch)
		flusher.Flush()
	}

	for match := range b.matchChan {
		log.Print("[w] straem match browser src", match)
		matchJson, err := json.Marshal(match)
		b.lastMatch = matchJson

		if err != nil {
			log.Println("browser source: failed to marshal match json", err)
			http.Error(w, "browser source: failed to marshal match data", http.StatusInternalServerError)
			break
		}
		fmt.Fprint(w, "event: message\n\n")
		fmt.Fprintf(w, "data: %s\n\n", matchJson)
		flusher.Flush()
	}
}

func (b *BrowserSourceServer) handleTheme(w http.ResponseWriter, req *http.Request) {
	fileName := req.PathValue("theme")
	css, err := staticFs.ReadFile(fmt.Sprintf("static/themes/%s", fileName))

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.Header().Set("Content-Type", "text/css")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(css)
		if err != nil {
			log.Println("failed to write browser source css")
		}
	}
}

func (b *BrowserSourceServer) handleRoot(w http.ResponseWriter, _ *http.Request) {
	html, err := staticFs.ReadFile("static/index.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write(html)
		if err != nil {
			log.Println("failed to write browser source html")
		}
	}
}

func GetInternalThemes() []model.Theme {
	var themes = make([]model.Theme, 0, 10)

	if err := fs.WalkDir(staticFs, "static/themes", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		b, _ := fs.ReadFile(staticFs, path)
		themes = append(themes, model.Theme{
			Name: strings.Split(d.Name(), ".css")[0],
			CSS:  string(b),
		})
		return nil
	}); err != nil {
		return []model.Theme{}
	}
	return themes
}

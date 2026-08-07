package server

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/config"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

//go:embed static
var staticFs embed.FS

type BrowserSourceServer struct {
	matchChan chan model.Match
	mu        sync.Mutex
	sseChans  map[chan []byte]struct{}
	lastMatch []byte
}

var heartbeatInterval = 15 * time.Second

func NewBrowserSourceServer(matchChan chan model.Match) *BrowserSourceServer {
	return &BrowserSourceServer{
		matchChan: matchChan,
		lastMatch: nil,
		sseChans:  make(map[chan []byte]struct{}),
	}
}

func (b *BrowserSourceServer) Start(ctx context.Context, cfg *config.BuildConfig) {
	go func() {
		for match := range b.matchChan {
			slog.Info("browser source: match received", slog.Any("replay_id", match.ReplayID))
			matchJson, err := json.Marshal(match)
			if err != nil {
				slog.Error("browser source: marshal match data", slog.Any("error", err))
			}
			b.mu.Lock()
			b.lastMatch = matchJson
			for sse := range b.sseChans {
				select {
				case sse <- matchJson:
				default:
				}
			}
			b.mu.Unlock()
		}
	}()

	slog.Debug("launching browser source server")

	http.HandleFunc("/", b.handleRoot)
	http.HandleFunc("GET /stream", b.handleStream)
	http.HandleFunc("GET /themes/{theme}", b.handleTheme)

	// serve custom themes through "themes" directory in the same directory as the user's executable
	fs := http.FileServer(http.Dir("./themes"))
	http.Handle("/themes/", http.StripPrefix("/themes/", fs))

	if err := http.ListenAndServe(fmt.Sprintf(":%d", cfg.BrowserSourcePort), nil); err != nil {
		slog.Error("browser source: launch server", slog.Any("error", err))
	}
}

func (b *BrowserSourceServer) handleStream(w http.ResponseWriter, req *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "SSE not supported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	b.mu.Lock()
	lastMatch := append([]byte(nil), b.lastMatch...)
	b.mu.Unlock()
	if lastMatch != nil {
		if _, err := fmt.Fprintf(w, "event: message\n\ndata: %s\n\n", lastMatch); err != nil {
			return
		}
		flusher.Flush()
	}

	sseChan := make(chan []byte, 1)
	b.mu.Lock()
	b.sseChans[sseChan] = struct{}{}
	b.mu.Unlock()
	defer func() {
		b.mu.Lock()
		delete(b.sseChans, sseChan)
		b.mu.Unlock()
	}()
	heartbeat := time.NewTicker(heartbeatInterval)
	defer heartbeat.Stop()
	for {
		select {
		case match := <-sseChan:
			if _, err := fmt.Fprintf(w, "event: message\n\ndata: %s\n\n", match); err != nil {
				return
			}
			flusher.Flush()
		case <-heartbeat.C:
			if _, err := fmt.Fprint(w, ": keep-alive\n\n"); err != nil {
				return
			}
			flusher.Flush()
		case <-req.Context().Done():
			return
		}
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
			slog.Error("browser source: send css", slog.Any("error", err))
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
			slog.Error("browser source: send html", slog.Any("error", err))
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

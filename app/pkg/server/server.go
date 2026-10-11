package server

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/i18n"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

//go:embed static
var staticFs embed.FS

type BrowserSourceServer struct {
	matchChan chan model.Match
	mux       *http.ServeMux
	mu        sync.Mutex
	sseChans  map[chan []byte]struct{}
	lastMatch []byte

	// Separate from mu: stopping a server waits for its handlers, which take mu.
	srvMu     sync.Mutex
	srv       *http.Server
	srvPort   uint16
	cancelSrv context.CancelFunc
}

var heartbeatInterval = 15 * time.Second

var shutdownTimeout = time.Second

func NewBrowserSourceServer(matchChan chan model.Match) *BrowserSourceServer {
	b := &BrowserSourceServer{
		matchChan: matchChan,
		mux:       http.NewServeMux(),
		lastMatch: nil,
		sseChans:  make(map[chan []byte]struct{}),
	}
	b.mux.HandleFunc("/", b.handleRoot)
	b.mux.HandleFunc("GET /stream", b.handleStream)
	b.mux.HandleFunc("GET /themes/{theme}", b.handleTheme)
	b.mux.HandleFunc("GET /translations/{locale}", b.handleTranslations)

	// serve custom themes through "themes" directory in the same directory as the user's executable
	fs := http.FileServer(http.Dir("./themes"))
	b.mux.Handle("/themes/", http.StripPrefix("/themes/", fs))
	return b
}

// Start forwards matches to connected clients and serves on port until ctx is done.
// It must only be called once.
func (b *BrowserSourceServer) Start(ctx context.Context, port uint16) error {
	go b.forwardMatches()
	go func() {
		<-ctx.Done()
		b.Close()
	}()
	slog.Debug("launching browser source server")
	return b.Listen(port)
}

func (b *BrowserSourceServer) forwardMatches() {
	for match := range b.matchChan {
		slog.Info("browser source: match received", slog.Any("replay_id", match.ReplayID))
		matchJson, err := json.Marshal(match)
		if err != nil {
			slog.Error("browser source: marshal match data", slog.Any("error", err))
		}
		b.mu.Lock()
		b.lastMatch = matchJson
		for sse := range b.sseChans {
			// Replace a match the client hasn't read yet, so it always gets the latest.
			select {
			case <-sse:
			default:
			}
			select {
			case sse <- matchJson:
			default:
			}
		}
		b.mu.Unlock()
	}
}

// Listen moves the server to port. The new port is bound before the old server
// stops, so on error the old server keeps running.
func (b *BrowserSourceServer) Listen(port uint16) error {
	b.srvMu.Lock()
	defer b.srvMu.Unlock()
	// Rebinding the port we're already on would fail with "address in use".
	if b.srv != nil && b.srvPort == port {
		return nil
	}

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("listen on port %d: %w", port, err)
	}
	b.stopLocked()

	// Requests inherit this context, so cancelling it ends SSE streams that
	// would otherwise keep Shutdown waiting.
	ctx, cancel := context.WithCancel(context.Background())
	srv := &http.Server{
		Handler:     b.mux,
		BaseContext: func(net.Listener) context.Context { return ctx },
	}
	b.srv, b.srvPort, b.cancelSrv = srv, port, cancel
	go func() {
		if err := srv.Serve(listener); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("browser source: serve", slog.Any("error", err))
		}
	}()
	return nil
}

func (b *BrowserSourceServer) Close() {
	b.srvMu.Lock()
	defer b.srvMu.Unlock()
	b.stopLocked()
}

func (b *BrowserSourceServer) stopLocked() {
	if b.srv == nil {
		return
	}
	b.cancelSrv()
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := b.srv.Shutdown(ctx); err != nil {
		slog.Warn("browser source: graceful shutdown", slog.Any("error", err))
		_ = b.srv.Close()
	}
	b.srv = nil
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

	// Register and snapshot under one lock, so a match arriving in between isn't missed.
	sseChan := make(chan []byte, 1)
	b.mu.Lock()
	b.sseChans[sseChan] = struct{}{}
	lastMatch := append([]byte(nil), b.lastMatch...)
	b.mu.Unlock()
	defer func() {
		b.mu.Lock()
		delete(b.sseChans, sseChan)
		b.mu.Unlock()
	}()
	if len(lastMatch) > 0 {
		if _, err := fmt.Fprintf(w, "event: message\n\ndata: %s\n\n", lastMatch); err != nil {
			return
		}
		flusher.Flush()
	}
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

func (b *BrowserSourceServer) handleTranslations(w http.ResponseWriter, req *http.Request) {
	translation, err := i18n.GetTranslation(req.PathValue("locale"))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(translation); err != nil {
		slog.Error("browser source: send translation", slog.Any("error", err))
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

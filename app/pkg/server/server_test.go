package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/i18n"
	"github.com/williamsjokvist/cfn-tracker/pkg/model"
)

func TestSSEClientDisconnectRemovesChannel(t *testing.T) {
	b := NewBrowserSourceServer(make(chan model.Match))
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/stream", nil).WithContext(ctx)
	recorder := httptest.NewRecorder()
	done := make(chan struct{})
	go func() { b.handleStream(recorder, req); close(done) }()
	deadline := time.After(time.Second)
	for {
		b.mu.Lock()
		count := len(b.sseChans)
		b.mu.Unlock()
		if count == 1 {
			break
		}
		select {
		case <-deadline:
			t.Fatal("SSE channel was not registered")
		default:
			time.Sleep(time.Millisecond)
		}
	}
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("handler did not return")
	}
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.sseChans) != 0 {
		t.Fatal("SSE channel was not removed")
	}
}

func TestSSEHeartbeat(t *testing.T) {
	old := heartbeatInterval
	heartbeatInterval = 5 * time.Millisecond
	defer func() { heartbeatInterval = old }()
	b := NewBrowserSourceServer(make(chan model.Match))
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/stream", nil).WithContext(ctx)
	recorder := httptest.NewRecorder()
	done := make(chan struct{})
	go func() { b.handleStream(recorder, req); close(done) }()
	time.Sleep(20 * time.Millisecond)
	cancel()
	<-done
	if !strings.Contains(recorder.Body.String(), ": keep-alive\n\n") {
		t.Fatal("heartbeat was not written")
	}
}

func freePort(t *testing.T) uint16 {
	t.Helper()
	l, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = l.Close() }()
	return uint16(l.Addr().(*net.TCPAddr).Port)
}

func serves(port uint16) bool {
	client := http.Client{Timeout: time.Second}
	res, err := client.Get(fmt.Sprintf("http://localhost:%d/", port))
	if err != nil {
		return false
	}
	_ = res.Body.Close()
	return res.StatusCode == http.StatusOK
}

func TestListenMovesServerToNewPort(t *testing.T) {
	b := NewBrowserSourceServer(make(chan model.Match))
	defer b.Close()
	oldPort, newPort := freePort(t), freePort(t)
	if err := b.Listen(oldPort); err != nil {
		t.Fatal(err)
	}
	if !serves(oldPort) {
		t.Fatal("old port is not serving")
	}
	if err := b.Listen(newPort); err != nil {
		t.Fatal(err)
	}
	if !serves(newPort) {
		t.Fatal("new port is not serving")
	}
	if serves(oldPort) {
		t.Fatal("old port is still serving")
	}
}

func TestListenOnSamePortKeepsServing(t *testing.T) {
	b := NewBrowserSourceServer(make(chan model.Match))
	defer b.Close()
	port := freePort(t)
	if err := b.Listen(port); err != nil {
		t.Fatal(err)
	}
	if err := b.Listen(port); err != nil {
		t.Fatal(err)
	}
	if !serves(port) {
		t.Fatal("port is not serving")
	}
}

func TestListenOnPortInUseKeepsOldServer(t *testing.T) {
	b := NewBrowserSourceServer(make(chan model.Match))
	defer b.Close()
	oldPort := freePort(t)
	if err := b.Listen(oldPort); err != nil {
		t.Fatal(err)
	}
	taken, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = taken.Close() }()
	if err := b.Listen(uint16(taken.Addr().(*net.TCPAddr).Port)); err == nil {
		t.Fatal("expected an error for a port in use")
	}
	if !serves(oldPort) {
		t.Fatal("old port stopped serving")
	}
}

func TestListenEndsOpenStreams(t *testing.T) {
	old := heartbeatInterval
	heartbeatInterval = 5 * time.Millisecond
	defer func() { heartbeatInterval = old }()
	b := NewBrowserSourceServer(make(chan model.Match))
	defer b.Close()
	oldPort := freePort(t)
	if err := b.Listen(oldPort); err != nil {
		t.Fatal(err)
	}
	res, err := http.Get(fmt.Sprintf("http://localhost:%d/stream", oldPort))
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = res.Body.Close() }()

	start := time.Now()
	if err := b.Listen(freePort(t)); err != nil {
		t.Fatal(err)
	}
	if elapsed := time.Since(start); elapsed >= shutdownTimeout {
		t.Fatalf("Listen waited %s for the open stream", elapsed)
	}
	closed := make(chan struct{})
	go func() { _, _ = io.Copy(io.Discard, res.Body); close(closed) }()
	select {
	case <-closed:
	case <-time.After(time.Second):
		t.Fatal("stream is still open")
	}
}

// A match field without a label, or a label key missing from a locale, would show an empty title on stream.
func TestStatLabelsCoverMatchInEveryLocale(t *testing.T) {
	matchType := reflect.TypeOf(model.Match{})
	for i := range matchType.NumField() {
		stat := matchType.Field(i).Tag.Get("json")
		if stat == "replayId" || stat == "sessionId" {
			continue
		}
		if _, ok := statLabelKeys[stat]; !ok {
			t.Errorf("model.Match field %q has no entry in statLabelKeys", stat)
		}
	}

	languages, err := i18n.GetSupportedLanguages()
	if err != nil {
		t.Fatalf("get supported languages: %v", err)
	}
	srv := httptest.NewServer(NewBrowserSourceServer(nil).mux)
	defer srv.Close()
	for _, lang := range languages {
		t.Run(lang, func(t *testing.T) {
			res, err := http.Get(fmt.Sprintf("%s/stat-labels/%s", srv.URL, lang))
			if err != nil {
				t.Fatalf("get stat labels: %v", err)
			}
			defer func() { _ = res.Body.Close() }()
			var labels map[string]string
			if err := json.NewDecoder(res.Body).Decode(&labels); err != nil {
				t.Fatalf("decode stat labels: %v", err)
			}
			for stat := range statLabelKeys {
				if labels[stat] == "" {
					t.Errorf("stat %q has no label", stat)
				}
			}
		})
	}
}

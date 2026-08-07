package server

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

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

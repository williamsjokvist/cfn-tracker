package sf6

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
)

// The most crucial test, to make sure authentication is always working.
func TestSF6Authentication(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	browser, err := browser.NewBrowser(true)

	if !assert.Nil(err) {
		t.Fatalf("failed to create browser: %v", err)
	}

	sf6Tracker := NewSF6Tracker(browser, nil)
	authChan := make(chan AuthStatus, 1)

	t.Cleanup(func() {
		browser.Page.Browser().Close()
		close(authChan)
	})

	go sf6Tracker.Authenticate(ctx, os.Getenv("CAP_ID_EMAIL"), os.Getenv("CAP_ID_PASSWORD"), authChan)
	for status := range authChan {
		if !assert.Nil(status.Err) {
			t.Fatalf("failed to authenticate: %v", status.Err)
		}
		if status.Progress >= 100 {
			break
		}
	}
}

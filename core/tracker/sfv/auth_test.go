package sfv

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/core/browser"
)

// The most crucial test, to make sure authentication is always working.
func TestSFVAuthentication(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	browser, err := browser.NewBrowser(false)
	if !assert.Nil(err) {
		t.Fatalf("failed to create browser: %v", err)
	}
	t.Cleanup(func() {
		browser.Page.Browser().Close()
	})

	sf5Tracker := NewSFVTracker(browser)
	err = sf5Tracker.Authenticate(ctx, os.Getenv(`STEAM_USERNAME`), os.Getenv(`STEAM_PASSWORD`), true)
	if !assert.Nil(err) {
		t.Fatalf("failed to authenticate: %v", err)
	}
}

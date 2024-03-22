package sfv

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
)

// The most crucial test, to make sure authentication is always working.
func TestSFVAuthentication(t *testing.T) {
	assert := assert.New(t)

	browser, err := browser.NewBrowser(true)

	if !assert.Nil(err) {
		t.Fatalf("failed to create browser: %v", err)
	}
	_, err = NewSFVTracker(browser, &config.Config{
		SteamUsername: os.Getenv(`STEAM_USERNAME`),
		SteamPassword: os.Getenv(`STEAM_PASSWORD`),
	}, nil)
	if !assert.Nil(err) {
		t.Fatalf("failed to authenticate: %v", err)
	}
}

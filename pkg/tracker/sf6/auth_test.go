package sf6

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

// The most crucial test, to make sure authentication is always working.
func TestSF6Authentication(t *testing.T) {
	assert := assert.New(t)

	browser, err := browser.NewBrowser(true)

	if !assert.Nil(err) {
		t.Fatalf("failed to create browser: %v", err)
	}

	sf6Tracker := NewSF6Tracker(browser, nil, nil)
	authChan := make(chan tracker.AuthStatus, 1)

	t.Cleanup(func() {
		browser.Page.Browser().Close()
		close(authChan)
	})

	go sf6Tracker.Authenticate(os.Getenv("CAP_ID_EMAIL"), os.Getenv("CAP_ID_PASSWORD"), authChan)
	for status := range authChan {
		if !assert.Nil(status.Err) {
			t.Fatalf("failed to authenticate: %v", status.Err)
		}
		if status.Progress >= 100 {
			break
		}
	}
}

package sf6

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/core/browser"
)

// The most crucial test, to make sure authentication is always working.
func TestSF6Authentication(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	browser, err := browser.NewBrowser(true)
	assert.Nil(err)

	sf6Tracker := NewSF6Tracker(browser, nil)
	progChan := make(chan int)
	sf6Tracker.Authenticate(ctx, os.Getenv(`CAP_ID_EMAIL`), os.Getenv(`CAP_ID_PASSWORD`), progChan)

	assert.Equal(nil, err)
}

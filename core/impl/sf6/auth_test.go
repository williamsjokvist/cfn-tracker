package sf6

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/core/shared"
)

// The most crucial test, to make sure authentication is always working.
func TestSF6Authentication(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	browser, _ := shared.NewBrowser(ctx, true)
	sf6Tracker := NewSF6Tracker(ctx, browser, nil)
	err := sf6Tracker.Authenticate(os.Getenv(`CAP_ID_EMAIL`), os.Getenv(`CAP_ID_PASSWORD`), true)

	assert.Equal(nil, err)
}

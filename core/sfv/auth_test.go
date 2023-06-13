package sfv

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/core/common"
)

// The most crucial test, to make sure authentication is always working.
func TestSFVAuthenticate(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	browser := common.NewBrowser(ctx, true)
	sf5Tracker := NewSFVTracker(ctx, browser)
	err := sf5Tracker.Authenticate(os.Getenv(`STEAM_USERNAME`), os.Getenv(`STEAM_PASSWORD`), true)

	assert.Equal(nil, err)
}

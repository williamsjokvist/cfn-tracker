package sf6

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/core/common"
)

// The most crucial test, to make sure authentication is always working.
func TestSF6Authenticate(t *testing.T) {
	assert := assert.New(t)

	ctx := context.Background()
	browser := common.NewBrowser(ctx, true)
	sf5Tracker := NewSF6Tracker(ctx, browser)
	err := sf5Tracker.Authenticate(os.Getenv(`CAP_ID_EMAIL`), os.Getenv(`CAP_ID_PASSWORD`), true)

	assert.Equal(nil, err)
}

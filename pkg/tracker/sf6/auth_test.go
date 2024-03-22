package sf6

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/williamsjokvist/cfn-tracker/pkg/browser"
	"github.com/williamsjokvist/cfn-tracker/pkg/config"
)

// The most crucial test, to make sure authentication is always working.
func TestSF6Authentication(t *testing.T) {
	assert := assert.New(t)

	browser, err := browser.NewBrowser(true)
	if !assert.Nil(err) {
		t.Fatalf("failed to create browser: %v", err)
	}

	_, err = NewSF6Tracker(browser, nil, nil, &config.Config{
		CapIDEmail:    os.Getenv("CAP_ID_EMAIL"),
		CapIDPassword: os.Getenv("CAP_ID_PASSWORD"),
	}, nil)
	if !assert.Nil(err) {
		t.Fatalf("failed to authenticate: %v", err)
	}
}

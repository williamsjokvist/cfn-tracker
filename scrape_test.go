package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	page, _ := SetupBrowser()

	loginStatus := Login(`GreenSoap`, page, os.Getenv(`STEAM_USERNAME`), os.Getenv(`STEAM_PASSWORD`))

	// Expect to run into CAPTCHA on github
	if os.Getenv(`GITHUB_TEST`) == `true` {
		assert.Equal(-3, <-loginStatus)
	} else {
		assert.Equal(1, <-loginStatus)
	}
}

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

	if !(<-loginStatus == 1 || <-loginStatus == -3) {
		assert.True(false, true, `Login attempt was not successful`)
	}
}

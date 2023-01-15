package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	assert := assert.New(t)
	page := SetupBrowser()

	loginStatus, _ := Login(`GreenSoap`, page, os.Getenv(`STEAM_USERNAME`), os.Getenv(`STEAM_PASSWORD`))
	assert.Equal(1, loginStatus)
}

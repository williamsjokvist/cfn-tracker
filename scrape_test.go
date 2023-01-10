package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogin(t *testing.T) {
	assert := assert.New(t)

	page, _ := SetupBrowser()

	// With false login information
	//isLoggedIn := Login(`GreenSoap`, page, `1234234`, `12341234`)
	//assert.Equal(<-isLoggedIn == LoginError.returnCode, true)

	// With nonexistant profile
	//isLoggedIn = Login(``, page, os.Getenv(`STEAM_USERNAME`), os.Getenv(`STEAM_PASSWORD`))
	//assert.Equal(<-isLoggedIn == ProfileError.returnCode, true)

	// With correct info
	isLoggedIn := Login(`GreenSoap`, page, os.Getenv(`STEAM_USERNAME`), os.Getenv(`STEAM_PASSWORD`))
	assert.Equal(true, (<-isLoggedIn == 1 || <-isLoggedIn == -3)) // Either log in or get CAPTCHA'd
}

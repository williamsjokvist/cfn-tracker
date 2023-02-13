package main

import (
	"fmt"
)

func LogError(err AppError) {
	fmt.Println(err.message)
}

type AppError struct {
	message    string
	returnCode int
}

var LoginError = AppError{
	message:    `Failed to Log in`,
	returnCode: -1,
}

var ProfileError = AppError{
	message:    `CFN profile does not exist`,
	returnCode: -2,
}

var CaptchaError = AppError{
	message:    `Encountered a captcha, please wait a while before opening the application again.`,
	returnCode: -3,
}

var ParseError = AppError{
	message:    `Could not locate data on CFN profile`,
	returnCode: -4,
}

var SaveError = AppError{
	message:    `Could not save data`,
	returnCode: -5,
}

var MissingEnvError = AppError{
	message:    `Error loading environment variables from .env file`,
	returnCode: -6,
}

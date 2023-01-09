package main

type AppError struct {
	message string
}

var LoginError = AppError{
	message: `Failed to Log in`,
}

var ProfileError = AppError{
	message: `CFN profile does not exist`,
}

var ParseError = AppError{
	message: `Could not locate data on CFN profile`,
}

var SaveErrror = AppError{
	message: `Could not save data`,
}

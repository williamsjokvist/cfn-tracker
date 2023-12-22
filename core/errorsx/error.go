package errorsx

import (
	"errors"
)

type AppError struct {
	ErrorCode  *int  `json:"code"`
	InnerError error `json:"message"`
}

func NewError(errorCode int, innerError error) *AppError {
	return &AppError{
		ErrorCode:  &errorCode,
		InnerError: innerError,
	}
}

func (e *AppError) Error() string {
	return e.InnerError.Error()
}

func (e *AppError) Unwrap() error {
	return e.InnerError
}

func ContainsTrackingError(err error) bool {
	var trackingErr *AppError
	return errors.As(err, &trackingErr)
}

func ConvertToAppError(err error) any {
	var trackingErr *AppError
	if errors.As(err, &trackingErr) {
		return trackingErr
	}

	return AppError{
		ErrorCode:  nil,
		InnerError: err,
	}
}

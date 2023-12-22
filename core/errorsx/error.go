package errorsx

import (
	"errors"
)

type FormattedError struct {
	ErrorCode  *int  `json:"code"`
	InnerError error `json:"message"`
}

func NewFormattedError(errorCode int, innerError error) *FormattedError {
	return &FormattedError{
		ErrorCode:  &errorCode,
		InnerError: innerError,
	}
}

func (e *FormattedError) Error() string {
	return e.InnerError.Error()
}

func (e *FormattedError) Unwrap() error {
	return e.InnerError
}

func ContainsFormattedError(err error) bool {
	var trackingErr *FormattedError
	return errors.As(err, &trackingErr)
}

func FormatError(err error) any {
	var trackingErr *FormattedError
	if errors.As(err, &trackingErr) {
		return trackingErr
	}

	return FormattedError{
		ErrorCode:  nil,
		InnerError: err,
	}
}

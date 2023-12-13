package errorsx

import (
	"errors"
	"strconv"
)

type TrackingError struct {
	errorCode  int
	innerError error
}

type FrontEndError struct {
	ErrorCode *int   `json:"code"`
	Message   string `json:"message"`
}

func NewError(errorCode int, innerError error) *TrackingError {
	return &TrackingError{
		errorCode:  errorCode,
		innerError: innerError,
	}
}

func (e *TrackingError) Error() string {
	return e.innerError.Error()
}

func (e *TrackingError) ErrorCode() string {
	return strconv.Itoa(e.errorCode)
}

func (e *TrackingError) Unwrap() error {
	return e.innerError
}

func FirstTrackingErrorOrDefault(err error) any {
	var trackingErr *TrackingError
	if errors.As(err, &trackingErr) {
		return FrontEndError{
			ErrorCode: &trackingErr.errorCode,
			Message:   trackingErr.innerError.Error(),
		}
	}

	return FrontEndError{
		ErrorCode: nil,
		Message:   err.Error(),
	}
}

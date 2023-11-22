package core

import "errors"

type TrackingError struct {
	errorCode  string
	innerError error
}

func (e *TrackingError) Error() string {
	return e.innerError.Error()
}

func (e *TrackingError) ErrorCode() string {
	return e.errorCode
}

func (e *TrackingError) Unwrap() error {
	return e.innerError
}

func UnwrapErrorWithCode(err error) error {
	var trackingErr *TrackingError
	if errors.As(err, &trackingErr) {
		return errors.New(trackingErr.ErrorCode())
	}
	return err
}

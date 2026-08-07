package model

import "fmt"

type HTTPStatusError struct {
	Op         string
	StatusCode int
}

func (e *HTTPStatusError) Error() string {
	return fmt.Sprintf("%s: http status %d", e.Op, e.StatusCode)
}

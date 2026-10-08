package model

import "fmt"

type ParseError struct {
	Op  string
	Err error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *ParseError) Unwrap() error { return e.Err }

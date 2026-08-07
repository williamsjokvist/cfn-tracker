package model

import (
	"context"
	"errors"
	"net"
	"syscall"
)

type ErrorClass int

const (
	ClassTransient ErrorClass = iota
	ClassAuth
	ClassParse
	ClassFatal
)

func ClassifyPollError(err error) ErrorClass {
	var statusErr *HTTPStatusError
	if errors.As(err, &statusErr) {
		switch {
		case statusErr.StatusCode == 401 || statusErr.StatusCode == 403:
			return ClassAuth
		case statusErr.StatusCode == 404:
			return ClassFatal
		case statusErr.StatusCode == 429 || statusErr.StatusCode >= 500:
			return ClassTransient
		case statusErr.StatusCode >= 400:
			return ClassFatal
		}
	}
	var parseErr *ParseError
	if errors.As(err, &parseErr) {
		return ClassParse
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return ClassTransient
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return ClassTransient
	}
	var opErr *net.OpError
	var dnsErr *net.DNSError
	if errors.As(err, &opErr) || errors.As(err, &dnsErr) ||
		errors.Is(err, syscall.ECONNREFUSED) || errors.Is(err, syscall.ECONNRESET) ||
		errors.Is(err, syscall.EHOSTUNREACH) {
		return ClassTransient
	}
	return ClassTransient
}

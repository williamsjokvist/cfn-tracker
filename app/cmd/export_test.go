package cmd

import "time"

// SetPollTimeoutForTest shortens the Poll timeout for external tests and returns a
// function that restores the original value. Only compiled into test builds.
func SetPollTimeoutForTest(timeout time.Duration) func() {
	original := pollTimeout
	pollTimeout = timeout
	return func() { pollTimeout = original }
}

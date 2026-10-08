package cmd

import (
	"time"

	"github.com/williamsjokvist/cfn-tracker/pkg/tracker"
)

// SetPollTimeoutForTest shortens the Poll timeout for external tests and returns a
// function that restores the original value. Only compiled into test builds.
func SetPollTimeoutForTest(timeout time.Duration) func() {
	original := pollTimeout
	pollTimeout = timeout
	return func() { pollTimeout = original }
}

// SetGameTracker replaces the GameTracker for external tests.
func (ch *TrackingHandler) SetGameTracker(gt tracker.GameTracker) { ch.gameTracker = gt }

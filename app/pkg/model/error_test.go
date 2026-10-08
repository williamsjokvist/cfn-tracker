package model

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// FormatError leaves InnerError nil, and main.go logs the result with `event=%s data=%v`.
// If Error() touched the nil InnerError, fmt would print %!v(PANIC=...) and lose the error.
func TestFormatErrorResultIsSafeToFormat(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"plain error", fmt.Errorf("poll: %w", errors.New("boom")), "boom"},
		{"wrapped in FGCTrackerError", WrapError(ErrGetMatches, errors.New("boom")), "boom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := fmt.Sprintf("%v", FormatError(tt.err))
			if strings.Contains(out, "PANIC") {
				t.Fatalf("formatting FormatError result panicked: %s", out)
			}
			if !strings.Contains(out, tt.want) {
				t.Errorf("original error was lost: got %q, want it to contain %q", out, tt.want)
			}
		})
	}
}

func TestErrorFallsBackWhenInnerErrorIsNil(t *testing.T) {
	tests := []struct {
		name string
		err  *FGCTrackerError
		want string
	}{
		{"InnerError takes precedence", &FGCTrackerError{LocalizationKey: tKeyErrAuth, Message: "detail", InnerError: errors.New("inner")}, "inner"},
		{"Message when InnerError is nil", &FGCTrackerError{LocalizationKey: tKeyErrUnknown, Message: "detail"}, "detail"},
		{"key when neither is set", &FGCTrackerError{LocalizationKey: tKeyErrAuth}, "errAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

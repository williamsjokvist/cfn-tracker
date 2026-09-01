package model

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

// FormatError はフロントへ送る値として InnerError を設定しない。その戻り値は
// main.go の `event=%s data=%v` でログへ出力されるため、Error() が nil の
// InnerError に触れると fmt が recover して %!v(PANIC=...) になり、
// エラーの正体がログから失われる。
func TestFormatErrorResultIsSafeToFormat(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"素のエラー", fmt.Errorf("poll: %w", errors.New("boom")), "boom"},
		{"FGCTrackerError で包んだエラー", WrapError(ErrGetMatches, errors.New("boom")), "boom"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			out := fmt.Sprintf("%v", FormatError(tt.err))
			if strings.Contains(out, "PANIC") {
				t.Fatalf("FormatError の戻り値を整形すると panic した: %s", out)
			}
			if !strings.Contains(out, tt.want) {
				t.Errorf("元のエラー内容が失われた: got %q, want it to contain %q", out, tt.want)
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
		{"InnerError が優先される", &FGCTrackerError{LocalizationKey: tKeyErrAuth, Message: "detail", InnerError: errors.New("inner")}, "inner"},
		{"InnerError が nil なら Message", &FGCTrackerError{LocalizationKey: tKeyErrUnknown, Message: "detail"}, "detail"},
		{"どちらも無ければキー", &FGCTrackerError{LocalizationKey: tKeyErrAuth}, "errAuth"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

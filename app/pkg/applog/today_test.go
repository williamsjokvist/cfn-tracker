package applog

import (
	"strings"
	"testing"
	"time"
)

func TestLinesFromDay(t *testing.T) {
	day := time.Date(2026, 10, 11, 12, 0, 0, 0, time.Local)
	log := strings.Join([]string{
		"2026/10/10 23:59:59 main.go:1: INFO yesterday",
		"yesterday continued",
		"2026/10/11 00:00:01 main.go:2: INFO today",
		"today continued",
		"2026/10/11 10:00:00 main.go:3: ERROR later today",
		"",
	}, "\n")

	got, err := linesFromDay(strings.NewReader(log), day)
	if err != nil {
		t.Fatal(err)
	}
	want := "2026/10/11 00:00:01 main.go:2: INFO today\ntoday continued\n2026/10/11 10:00:00 main.go:3: ERROR later today\n"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestLinesFromDayKeepsLatest(t *testing.T) {
	day := time.Date(2026, 10, 11, 12, 0, 0, 0, time.Local)
	line := "2026/10/11 10:00:00 main.go:1: INFO " + strings.Repeat("x", 100) + "\n"
	log := strings.Repeat(line, maxTodayBytes/len(line)*3) + "2026/10/11 11:00:00 main.go:2: INFO last\n"

	got, err := linesFromDay(strings.NewReader(log), day)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) > maxTodayBytes {
		t.Errorf("got %d bytes, want at most %d", len(got), maxTodayBytes)
	}
	if !strings.HasPrefix(got, "2026/10/11 ") {
		t.Errorf("starts mid-line: %q", got[:40])
	}
	if !strings.HasSuffix(got, "INFO last\n") {
		t.Errorf("dropped the latest line")
	}
}

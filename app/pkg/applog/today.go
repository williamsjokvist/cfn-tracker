package applog

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// GitHub caps issue bodies at 65536 characters, which the report text also has to fit in.
const maxTodayBytes = 48 * 1024

func Today() (string, error) {
	f, err := os.Open(fileName)
	if errors.Is(err, os.ErrNotExist) {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("open log file: %w", err)
	}
	defer func() { _ = f.Close() }()
	return linesFromDay(f, time.Now())
}

func linesFromDay(r io.Reader, day time.Time) (string, error) {
	prefix := []byte(day.Format("2006/01/02 "))
	var out []byte
	// Lines without a date prefix continue a multi-line entry, so they follow the entry they belong to.
	inDay := false
	sc := bufio.NewScanner(r)
	sc.Buffer(make([]byte, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if hasDatePrefix(line) {
			inDay = bytes.HasPrefix(line, prefix)
		}
		if !inDay {
			continue
		}
		out = append(out, line...)
		out = append(out, '\n')
		if len(out) > 2*maxTodayBytes {
			out = append([]byte(nil), out[len(out)-maxTodayBytes:]...)
		}
	}
	if err := sc.Err(); err != nil {
		return "", fmt.Errorf("read log file: %w", err)
	}
	if len(out) > maxTodayBytes {
		out = out[len(out)-maxTodayBytes:]
		if i := bytes.IndexByte(out, '\n'); i >= 0 {
			out = out[i+1:]
		}
	}
	return string(out), nil
}

func hasDatePrefix(line []byte) bool {
	if len(line) < 11 {
		return false
	}
	_, err := time.Parse("2006/01/02 ", string(line[:11]))
	return err == nil
}

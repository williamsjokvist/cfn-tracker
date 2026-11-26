//go:build dev

package applog

import "log/slog"

// wails dev builds with the dev tag, so debug logs show while developing but not
// in releases, where they'd fill the log file.
func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

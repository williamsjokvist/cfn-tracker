package applog

import (
	"fmt"
	"log"
	"os"
	"sync"
)

const fileName = "cfn-tracker.log"

var (
	mu   sync.Mutex
	file *os.File
)

// SetFileLogging switches log output between cfn-tracker.log and stderr.
func SetFileLogging(enabled bool) error {
	mu.Lock()
	defer mu.Unlock()

	if !enabled {
		log.SetOutput(os.Stderr)
		return closeFile()
	}
	if file != nil {
		return nil
	}
	f, err := os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return fmt.Errorf("open log file: %w", err)
	}
	file = f
	log.SetOutput(f)
	log.SetFlags(log.Ldate | log.LstdFlags | log.Lshortfile)
	return nil
}

// Close closes the log file, if one is open.
func Close() error {
	mu.Lock()
	defer mu.Unlock()
	return closeFile()
}

func closeFile() error {
	if file == nil {
		return nil
	}
	err := file.Close()
	file = nil
	return err
}

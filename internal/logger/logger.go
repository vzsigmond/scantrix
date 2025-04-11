// internal/logger/logger.go
package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logFile        *os.File
	logOnce        sync.Once
	loggingEnabled = false
)

// EnableDebug turns on debug logging to logs/debug.log
func EnableDebug() {
	loggingEnabled = true
}

func ensureLogFile() {
	logOnce.Do(func() {
		if !loggingEnabled {
			return
		}
		_ = os.MkdirAll("logs", os.ModePerm)
		absPath, _ := filepath.Abs("logs/debug.log")
		fmt.Println("🔎 Writing logs to:", absPath)

		var err error
		logFile, err = os.OpenFile(absPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
			os.Exit(1)
		}

		startLine := fmt.Sprintf("==== Scantrix Log Started at %s ====", time.Now().Format(time.RFC3339))
		logFile.WriteString(startLine + "\n")
		logFile.Sync()
	})
}

// Log writes a formatted debug message to the log file.
func Log(format string, args ...any) {
	if !loggingEnabled {
		return
	}
	ensureLogFile()
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintf(logFile, "%s %s\n", timestamp, fmt.Sprintf(format, args...))
	logFile.Sync()
}

// Close gracefully closes the log file.
func Close() {
	if logFile != nil {
		logFile.Close()
	}
}

package log

import (
	"fmt"
	"log"
	"os"
)

var (
	debugMode bool
	logger    *log.Logger
)

func init() {
	logger = log.New(os.Stderr, "", log.LstdFlags)
	debugMode = os.Getenv("CLAUDE_DEBUG") == "1"
}

// SetDebug enables or disables debug logging
func SetDebug(enabled bool) {
	debugMode = enabled
}

// Debug logs a debug message (only shown when debug mode is on)
func Debug(format string, args ...interface{}) {
	if debugMode {
		logger.Printf("[DEBUG] "+format, args...)
	}
}

// Info logs an informational message
func Info(format string, args ...interface{}) {
	logger.Printf("[INFO] "+format, args...)
}

// Warn logs a warning message
func Warn(format string, args ...interface{}) {
	logger.Printf("[WARN] "+format, args...)
}

// Error logs an error
func Error(err error) {
	if err != nil {
		logger.Printf("[ERROR] %v", err)
	}
}

// Errorf logs a formatted error
func Errorf(format string, args ...interface{}) {
	logger.Printf("[ERROR] "+format, args...)
}

// Fatal logs and exits
func Fatal(format string, args ...interface{}) {
	logger.Printf("[FATAL] "+format, args...)
	os.Exit(1)
}

// FormatError returns a user-friendly error string
func FormatError(err error) string {
	return fmt.Sprintf("Error: %v", err)
}

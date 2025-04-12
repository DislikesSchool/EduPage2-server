package util

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	logFile     *os.File
)

// InitLogging initializes the loggers
func InitLogging() {
	logDir := "./logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		fmt.Printf("Failed to create logs directory: %v\n", err)
		// Fall back to console logging
		InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
		ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		return
	}

	logPath := filepath.Join(logDir, fmt.Sprintf("edupage2-%s.log", time.Now().Format("2006-01-02")))
	var err error
	logFile, err = os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		// Fall back to console logging
		InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
		ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
		return
	}

	InfoLogger = log.New(io.MultiWriter(os.Stdout, logFile), "INFO: ", log.Ldate|log.Ltime)
	ErrorLogger = log.New(io.MultiWriter(os.Stderr, logFile), "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// CloseLogger closes the log file
func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

// Log user session lifecycle events
func LogUserSession(action, username, server string, err error) {
	if InfoLogger == nil {
		InitLogging()
	}

	if err != nil {
		ErrorLogger.Printf("User %s@%s %s failed: %v", username, server, action, err)
	} else {
		InfoLogger.Printf("User %s@%s %s successful", username, server, action)
	}
}

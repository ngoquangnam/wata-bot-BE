package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	errorLogFile *os.File
	errorLogMu   sync.Mutex
)

// InitErrorLog initializes error log file
func InitErrorLog(logPath string) error {
	errorLogMu.Lock()
	defer errorLogMu.Unlock()

	errorLogDir := filepath.Join(logPath, "error")
	
	// Create error log directory if it doesn't exist
	if err := os.MkdirAll(errorLogDir, 0755); err != nil {
		return fmt.Errorf("failed to create error log directory: %v", err)
	}

	// Create error log file with timestamp
	timestamp := time.Now().Format("2006-01-02")
	filePath := filepath.Join(errorLogDir, fmt.Sprintf("error-%s.log", timestamp))
	
	var err error
	errorLogFile, err = os.OpenFile(
		filePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0644,
	)
	if err != nil {
		return fmt.Errorf("failed to open error log file: %v", err)
	}

	return nil
}

// WriteErrorLog writes error message to error log file
func WriteErrorLog(message string, err error) {
	errorLogMu.Lock()
	defer errorLogMu.Unlock()

	if errorLogFile == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s", timestamp, message)
	if err != nil {
		logMessage += fmt.Sprintf(" - Error: %v", err)
	}
	logMessage += "\n"

	errorLogFile.WriteString(logMessage)
	errorLogFile.Sync()
}

// WriteErrorLogWithContext writes error message with context to error log file
func WriteErrorLogWithContext(message string, err error, context map[string]interface{}) {
	errorLogMu.Lock()
	defer errorLogMu.Unlock()

	if errorLogFile == nil {
		return
	}

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s", timestamp, message)
	if err != nil {
		logMessage += fmt.Sprintf(" - Error: %v", err)
	}
	if len(context) > 0 {
		logMessage += fmt.Sprintf(" - Context: %+v", context)
	}
	logMessage += "\n"

	errorLogFile.WriteString(logMessage)
	errorLogFile.Sync()
}

// CloseErrorLog closes error log file
func CloseErrorLog() error {
	errorLogMu.Lock()
	defer errorLogMu.Unlock()

	if errorLogFile != nil {
		err := errorLogFile.Close()
		errorLogFile = nil
		return err
	}
	return nil
}


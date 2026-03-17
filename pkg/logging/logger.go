package logging

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Logger struct {
	path string
	mu   sync.RWMutex
}

func NewLogger(fileName string) *Logger {
	// absolute path to the project root or executable dir
	execPath, _ := os.Executable()
	dir := filepath.Dir(execPath)
	
	return &Logger{
		path: filepath.Join(dir, fileName),
	}
}

func (l *Logger) Write(message string) {
	l.mu.Lock() // reduce writing crashes
    defer l.mu.Unlock()

	f, err := os.OpenFile(l.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("LOGGING ERROR: %v\n", err)
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
}

// ReadLogs safely retrieves the logs for the UI
func (l *Logger) ReadLogs(limit int) (string, error) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	data, err := os.ReadFile(l.path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
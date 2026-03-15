package logging

import (
	"time"
	"os"
	"fmt"
)

type Logger struct {
	logFile string
}

func NewLogger(logFile string) *Logger {
	return &Logger{
		logFile: logFile,
	}
}

// Write appends a timestamped message to cura.log
func (m *Logger) Write(message string) {
	f, err := os.OpenFile("cura.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	f.WriteString(fmt.Sprintf("[%s] %s\n", timestamp, message))
}
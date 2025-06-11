package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

// emojiMap — эмодзи для разных уровней логирования.
var emojiMap = map[string]string{
	"DEBUG": "🐛",
	"INFO":  "ℹ️",
	"WARN":  "⚠️",
	"ERROR": "❌",
}

// colorMap — цветовые escape-коды для разных уровней.
var colorMap = map[string]string{
	"DEBUG": "\033[36m", // cyan
	"INFO":  "\033[32m", // green
	"WARN":  "\033[33m", // yellow
	"ERROR": "\033[31m", // red
}

const colorReset = "\033[0m"

type Logger struct {
	logger *log.Logger
}

// NewLogger — конструктор логгера.
func NewLogger() *Logger {
	return &Logger{
		logger: log.New(os.Stdout, "", 0),
	}
}

// logf — внутренний метод для форматирования и вывода лога.
func (l *Logger) logf(level string, format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	shortFile := file
	if idx := len(file) - 1; idx >= 0 {
		for i := len(file) - 1; i >= 0; i-- {
			if file[i] == '/' {
				shortFile = file[i+1:]
				break
			}
		}
	}
	emoji := emojiMap[level]
	color := colorMap[level]
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, v...)
	l.logger.Printf("%s%s [%s] %s %s:%d | %s%s", color, emoji, level, timestamp, shortFile, line, msg, colorReset)
}

// Debug — логирование на уровне debug.
func (l *Logger) Debug(format string, v ...interface{}) {
	l.logf("DEBUG", format, v...)
}

// Info — логирование на уровне info.
func (l *Logger) Info(format string, v ...interface{}) {
	l.logf("INFO", format, v...)
}

// Warn — логирование на уровне warn.
func (l *Logger) Warn(format string, v ...interface{}) {
	l.logf("WARN", format, v...)
}

// Error — логирование на уровне error.
func (l *Logger) Error(format string, v ...interface{}) {
	l.logf("ERROR", format, v...)
}

// Sync — для совместимости с zap, ничего не делает.
func (l *Logger) Sync() error { return nil }

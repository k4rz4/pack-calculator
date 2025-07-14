package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
)

// Level represents log levels
type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var levelNames = map[Level]string{
	DEBUG: "DEBUG",
	INFO:  "INFO",
	WARN:  "WARN",
	ERROR: "ERROR",
}

// Logger provides structured logging
type Logger struct {
	level  Level
	format string // "json" or "text"
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp string                 `json:"timestamp"`
	Level     string                 `json:"level"`
	Message   string                 `json:"message"`
	Fields    map[string]interface{} `json:"fields,omitempty"`
	RequestID string                 `json:"request_id,omitempty"`
	Method    string                 `json:"method,omitempty"`
	Path      string                 `json:"path,omitempty"`
	Duration  string                 `json:"duration,omitempty"`
	Status    int                    `json:"status,omitempty"`
}

// New creates a new logger
func New(level, format string) *Logger {
	logLevel := parseLevel(level)
	if format != "json" && format != "text" {
		format = "json"
	}

	return &Logger{
		level:  logLevel,
		format: format,
	}
}

// parseLevel converts string to Level
func parseLevel(level string) Level {
	switch strings.ToUpper(level) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO
	}
}

// Debug logs a debug message
func (l *Logger) Debug(message string, fields ...map[string]interface{}) {
	l.log(DEBUG, message, fields...)
}

// Info logs an info message
func (l *Logger) Info(message string, fields ...map[string]interface{}) {
	l.log(INFO, message, fields...)
}

// Warn logs a warning message
func (l *Logger) Warn(message string, fields ...map[string]interface{}) {
	l.log(WARN, message, fields...)
}

// Error logs an error message
func (l *Logger) Error(message string, fields ...map[string]interface{}) {
	l.log(ERROR, message, fields...)
}

// HTTP logs an HTTP request
func (l *Logger) HTTP(
	method, path, requestID string,
	status int,
	duration time.Duration,
	fields ...map[string]interface{},
) {
	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     levelNames[INFO],
		Message:   "HTTP Request",
		RequestID: requestID,
		Method:    method,
		Path:      path,
		Status:    status,
		Duration:  duration.String(),
	}

	if len(fields) > 0 {
		entry.Fields = fields[0]
	}

	l.output(entry)
}

// log is the internal logging method
func (l *Logger) log(level Level, message string, fields ...map[string]interface{}) {
	if level < l.level {
		return
	}

	entry := LogEntry{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Level:     levelNames[level],
		Message:   message,
	}

	if len(fields) > 0 {
		entry.Fields = fields[0]
	}

	l.output(entry)
}

// output writes the log entry
func (l *Logger) output(entry LogEntry) {
	if l.format == "json" {
		data, err := json.Marshal(entry)
		if err != nil {
			log.Printf("Failed to marshal log entry: %v", err)
			return
		}
		fmt.Println(string(data))
	} else {
		// Text format
		output := fmt.Sprintf("%s [%s] %s", entry.Timestamp, entry.Level, entry.Message)

		if entry.RequestID != "" {
			output += fmt.Sprintf(" request_id=%s", entry.RequestID)
		}

		if entry.Method != "" && entry.Path != "" {
			output += fmt.Sprintf(" %s %s", entry.Method, entry.Path)
		}

		if entry.Status > 0 {
			output += fmt.Sprintf(" status=%d", entry.Status)
		}

		if entry.Duration != "" {
			output += fmt.Sprintf(" duration=%s", entry.Duration)
		}

		if entry.Fields != nil {
			for k, v := range entry.Fields {
				output += fmt.Sprintf(" %s=%v", k, v)
			}
		}

		fmt.Println(output)
	}
}

// Global logger instance
var defaultLogger *Logger

// Initialize sets up the global logger
func Initialize(level, format string) {
	defaultLogger = New(level, format)
}

// Global logging functions
func Debug(message string, fields ...map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(message, fields...)
	}
}

func Info(message string, fields ...map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(message, fields...)
	}
}

func Warn(message string, fields ...map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(message, fields...)
	}
}

func Error(message string, fields ...map[string]interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(message, fields...)
	}
}

func HTTP(
	method, path, requestID string,
	status int,
	duration time.Duration,
	fields ...map[string]interface{},
) {
	if defaultLogger != nil {
		defaultLogger.HTTP(method, path, requestID, status, duration, fields...)
	}
}

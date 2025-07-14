package logger

import (
	"bytes"
	"encoding/json"
	"io"
	"os"
	"strings"
	"testing"
	"time"
)

// captureOutput captures stdout for testing
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestNew(t *testing.T) {
	tests := []struct {
		name           string
		level          string
		format         string
		expectedLevel  Level
		expectedFormat string
	}{
		{"Debug level JSON", "debug", "json", DEBUG, "json"},
		{"Info level text", "info", "text", INFO, "text"},
		{"Invalid level defaults to info", "invalid", "json", INFO, "json"},
		{"Invalid format defaults to json", "info", "invalid", INFO, "json"},
		{"Case insensitive", "INFO", "JSON", INFO, "json"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := New(tt.level, tt.format)

			if logger.level != tt.expectedLevel {
				t.Errorf("Expected level %d, got %d", tt.expectedLevel, logger.level)
			}

			if logger.format != tt.expectedFormat {
				t.Errorf("Expected format %s, got %s", tt.expectedFormat, logger.format)
			}
		})
	}
}

func TestParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected Level
	}{
		{"debug", DEBUG},
		{"DEBUG", DEBUG},
		{"info", INFO},
		{"INFO", INFO},
		{"warn", WARN},
		{"WARN", WARN},
		{"error", ERROR},
		{"ERROR", ERROR},
		{"invalid", INFO}, // defaults to INFO
		{"", INFO},        // defaults to INFO
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseLevel(%s): expected %d, got %d", tt.input, tt.expected, result)
			}
		})
	}
}

func TestLogger_LogLevels(t *testing.T) {
	logger := New("info", "json")

	tests := []struct {
		name        string
		logFunc     func()
		shouldLog   bool
		expectedMsg string
	}{
		{
			name:      "Debug should not log when level is INFO",
			logFunc:   func() { logger.Debug("debug message") },
			shouldLog: false,
		},
		{
			name:        "Info should log when level is INFO",
			logFunc:     func() { logger.Info("info message") },
			shouldLog:   true,
			expectedMsg: "info message",
		},
		{
			name:        "Warn should log when level is INFO",
			logFunc:     func() { logger.Warn("warn message") },
			shouldLog:   true,
			expectedMsg: "warn message",
		},
		{
			name:        "Error should log when level is INFO",
			logFunc:     func() { logger.Error("error message") },
			shouldLog:   true,
			expectedMsg: "error message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := captureOutput(tt.logFunc)

			if tt.shouldLog {
				if output == "" {
					t.Errorf("Expected log output but got none")
					return
				}

				var entry LogEntry
				if err := json.Unmarshal([]byte(output), &entry); err != nil {
					t.Errorf("Failed to parse JSON output: %v", err)
					return
				}

				if entry.Message != tt.expectedMsg {
					t.Errorf("Expected message %s, got %s", tt.expectedMsg, entry.Message)
				}
			} else {
				if output != "" {
					t.Errorf("Expected no log output but got: %s", output)
				}
			}
		})
	}
}

func TestLogger_JSONFormat(t *testing.T) {
	logger := New("info", "json")

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	output := captureOutput(func() {
		logger.Info("test message", fields)
	})

	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry.Level != "INFO" {
		t.Errorf("Expected level INFO, got %s", entry.Level)
	}

	if entry.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", entry.Message)
	}

	if entry.Fields["key1"] != "value1" {
		t.Errorf("Expected field key1 to be 'value1', got %v", entry.Fields["key1"])
	}

	if entry.Fields["key2"] != float64(42) { // JSON unmarshals numbers as float64
		t.Errorf("Expected field key2 to be 42, got %v", entry.Fields["key2"])
	}

	if entry.Timestamp == "" {
		t.Errorf("Timestamp should not be empty")
	}
}

func TestLogger_TextFormat(t *testing.T) {
	logger := New("info", "text")

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	output := captureOutput(func() {
		logger.Info("test message", fields)
	})

	if !strings.Contains(output, "[INFO]") {
		t.Errorf("Expected output to contain [INFO], got: %s", output)
	}

	if !strings.Contains(output, "test message") {
		t.Errorf("Expected output to contain 'test message', got: %s", output)
	}

	if !strings.Contains(output, "key1=value1") {
		t.Errorf("Expected output to contain 'key1=value1', got: %s", output)
	}

	if !strings.Contains(output, "key2=42") {
		t.Errorf("Expected output to contain 'key2=42', got: %s", output)
	}
}

func TestLogger_HTTP(t *testing.T) {
	logger := New("info", "json")

	duration := 100 * time.Millisecond
	fields := map[string]interface{}{
		"user_agent": "test-agent",
	}

	output := captureOutput(func() {
		logger.HTTP("POST", "/api/test", "req123", 200, duration, fields)
	})

	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Fatalf("Failed to parse JSON output: %v", err)
	}

	if entry.Message != "HTTP Request" {
		t.Errorf("Expected message 'HTTP Request', got %s", entry.Message)
	}

	if entry.Method != "POST" {
		t.Errorf("Expected method POST, got %s", entry.Method)
	}

	if entry.Path != "/api/test" {
		t.Errorf("Expected path /api/test, got %s", entry.Path)
	}

	if entry.RequestID != "req123" {
		t.Errorf("Expected request_id req123, got %s", entry.RequestID)
	}

	if entry.Status != 200 {
		t.Errorf("Expected status 200, got %d", entry.Status)
	}

	if entry.Duration != duration.String() {
		t.Errorf("Expected duration %s, got %s", duration.String(), entry.Duration)
	}

	if entry.Fields["user_agent"] != "test-agent" {
		t.Errorf("Expected user_agent 'test-agent', got %v", entry.Fields["user_agent"])
	}
}

func TestGlobalLogger(t *testing.T) {
	// Save original state
	originalLogger := defaultLogger
	defer func() {
		defaultLogger = originalLogger
	}()

	// Test without initialization (should not panic)
	output := captureOutput(func() {
		Info("test message")
	})

	if output != "" {
		t.Errorf("Expected no output when logger not initialized, got: %s", output)
	}

	// Test with initialization
	Initialize("info", "json")

	output = captureOutput(func() {
		Info("test message")
	})

	if output == "" {
		t.Errorf("Expected output after initialization")
		return
	}

	var entry LogEntry
	if err := json.Unmarshal([]byte(output), &entry); err != nil {
		t.Errorf("Failed to parse JSON output: %v", err)
		return
	}

	if entry.Message != "test message" {
		t.Errorf("Expected message 'test message', got %s", entry.Message)
	}
}

func TestGlobalLoggerMethods(t *testing.T) {
	// Save original state
	originalLogger := defaultLogger
	defer func() {
		defaultLogger = originalLogger
	}()

	Initialize("debug", "json")

	methods := []struct {
		name    string
		logFunc func()
		level   string
	}{
		{"Debug", func() { Debug("debug msg") }, "DEBUG"},
		{"Info", func() { Info("info msg") }, "INFO"},
		{"Warn", func() { Warn("warn msg") }, "WARN"},
		{"Error", func() { Error("error msg") }, "ERROR"},
	}

	for _, method := range methods {
		t.Run(method.name, func(t *testing.T) {
			output := captureOutput(method.logFunc)

			var entry LogEntry
			if err := json.Unmarshal([]byte(output), &entry); err != nil {
				t.Errorf("Failed to parse JSON output: %v", err)
				return
			}

			if entry.Level != method.level {
				t.Errorf("Expected level %s, got %s", method.level, entry.Level)
			}
		})
	}
}

func BenchmarkLogger_Info(b *testing.B) {
	logger := New("info", "json")

	// Redirect to discard output for benchmarking
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Info("benchmark message", fields)
	}
}

func BenchmarkLogger_Debug_Filtered(b *testing.B) {
	logger := New("info", "json") // Debug logs will be filtered out

	fields := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		logger.Debug("debug message", fields) // Should be filtered out
	}
}

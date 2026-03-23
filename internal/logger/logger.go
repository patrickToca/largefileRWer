package logger

import (
	"io"
	"log/slog"
	"os"
	"strings"
)

// Setup configures the global slog logger based on configuration
func Setup(level string, format string, outputPath string) error {
	var logLevel slog.Level
	switch strings.ToLower(level) {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	var output io.Writer = os.Stdout
	if outputPath != "" {
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		output = file
	}

	var handler slog.Handler
	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: false,
		})
	default:
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{
			Level:     logLevel,
			AddSource: false,
		})
	}

	slog.SetDefault(slog.New(handler))
	return nil
}

// WithContext adds context fields to the logger
func WithContext(key string, value interface{}) *slog.Logger {
	return slog.Default().With(key, value)
}

// Debug logs a debug message
func Debug(msg string, args ...interface{}) {
	slog.Debug(msg, args...)
}

// Info logs an info message
func Info(msg string, args ...interface{}) {
	slog.Info(msg, args...)
}

// Warn logs a warning message
func Warn(msg string, args ...interface{}) {
	slog.Warn(msg, args...)
}

// Error logs an error message
func Error(msg string, args ...interface{}) {
	slog.Error(msg, args...)
}

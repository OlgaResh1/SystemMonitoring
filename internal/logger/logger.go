package logger

import (
	"log/slog"
	"os"
)

type Logger struct {
	logger *slog.Logger
}

func New(level string, format string, isAddSource bool) *Logger {
	l := slog.LevelDebug
	switch level {
	case "debug":
		l = slog.LevelDebug
	case "info":
		l = slog.LevelInfo
	case "warn":
		l = slog.LevelWarn
	case "error":
		l = slog.LevelError
	}

	logConfig := &slog.HandlerOptions{
		AddSource:   isAddSource,
		Level:       l,
		ReplaceAttr: nil,
	}
	var logger *slog.Logger
	if format == "json" {
		logger = slog.New(slog.NewJSONHandler(os.Stderr, logConfig))
	} else if format == "text" {
		logger = slog.New(slog.NewTextHandler(os.Stderr, logConfig))
	}

	slog.SetDefault(logger)
	return &Logger{logger: logger}
}

func (l Logger) Info(msg string, args ...any) {
	slog.Info(msg, args...)
}

func (l Logger) Error(msg string, args ...any) {
	slog.Error(msg, args...)
}

func (l Logger) Warn(msg string, args ...any) {
	slog.Warn(msg, args...)
}

func (l Logger) Debug(msg string, args ...any) {
	slog.Debug(msg, args...)
}

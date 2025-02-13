package logger

import (
	"context"
	"log/slog"
	"os"
)

var logger *slog.Logger

const (
	LevelTrace = slog.Level(-8)
	LevelFatal = slog.Level(12)
)

var LevelNames = map[slog.Leveler]string{
	LevelTrace: "TRACE",
	LevelFatal: "FATAL",
}

func init() {
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == slog.LevelKey {
				level := a.Value.Any().(slog.Level)
				levelLabel, exists := LevelNames[level]
				if !exists {
					levelLabel = level.String()
				}
				a.Value = slog.StringValue(levelLabel)
			}

			return a
		},
	}

	if envLevel := os.Getenv("LOG_LEVEL"); envLevel != "" {
		var level slog.Level
		if envLevel == "trace" {
			opts.Level = LevelTrace
		} else if err := level.UnmarshalText([]byte(envLevel)); err == nil {
			opts.Level = level
		}
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	logger = slog.New(handler)

	slog.SetDefault(logger)
}

func Trace(msg string, args ...any) {
	logger.Log(context.Background(), LevelTrace, msg, args...)
}

func Debug(msg string, args ...any) {
	logger.Debug(msg, args...)
}

func Info(msg string, args ...any) {
	logger.Info(msg, args...)
}

func Warn(msg string, args ...any) {
	logger.Warn(msg, args...)
}

func Error(msg string, args ...any) {
	logger.Error(msg, args...)
}

func Panic(msg string, args ...any) {
	logger.Log(context.Background(), LevelFatal, msg, args...)
	panic(msg)
}

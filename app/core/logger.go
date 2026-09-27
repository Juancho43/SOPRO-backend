package core

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger(config *Config) *slog.Logger {

	var logLevel slog.Level
	if config.Environment == "production" {
		logLevel = slog.LevelInfo
	} else {
		logLevel = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level: logLevel,
	}

	logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	var logOutput io.Writer = os.Stdout
	if err == nil {
		logOutput = io.MultiWriter(os.Stdout, logFile)
	}

	return slog.New(slog.NewJSONHandler(logOutput, opts))
}

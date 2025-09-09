package server

import (
	"context"
	"log/slog"

	"github.com/creachadair/jrpc2"
)

func SlogToLogBridge(logger *slog.Logger) jrpc2.Logger {
	return func(text string) {
		if err := logger.Handler().Handle(context.Background(), slog.Record{
			Level:   slog.LevelInfo,
			Message: text,
		}); err != nil {
			logger.Error("error logging message", "error", err)
		}
	}
}

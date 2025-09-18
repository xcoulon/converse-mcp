package server

import (
	"context"
	"log/slog"
	"time"

	"github.com/creachadair/jrpc2"
)

func SlogToLogBridge(logger *slog.Logger) jrpc2.Logger {
	return func(text string) {
		if err := logger.Handler().Handle(context.Background(), slog.Record{
			Time:    time.Now(),
			Level:   slog.LevelDebug,
			Message: text,
		}); err != nil {
			logger.Error("error logging message", "error", err)
		}
	}
}

func SlogToRPCLogBridge(logger *slog.Logger) jrpc2.RPCLogger {
	return &rpcLogBridge{
		logger: logger,
	}
}

type rpcLogBridge struct {
	logger *slog.Logger
}

func (b *rpcLogBridge) LogRequest(ctx context.Context, req *jrpc2.Request) {
	if b.logger.Enabled(ctx, slog.LevelDebug) {
		b.logger.Debug("request", "method", req.Method(), "params", req.ParamString())
	}
}

func (b *rpcLogBridge) LogResponse(ctx context.Context, res *jrpc2.Response) {
	if b.logger.Enabled(ctx, slog.LevelDebug) {
		b.logger.Debug("response", "result", res.ResultString())
	}
}

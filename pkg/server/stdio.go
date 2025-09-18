package server

import (
	"log/slog"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
)

type StdioServer struct {
	*jrpc2.Server
}

func NewStdioServer(logger *slog.Logger, router Router) *StdioServer {
	srv := jrpc2.NewServer(handler.Map(router), &jrpc2.ServerOptions{
		Logger: SlogToLogBridge(logger),
	})
	return &StdioServer{
		Server: srv,
	}
}

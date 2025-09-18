package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/jhttp"
)

const DefaultHTTPPort = 8080

type StreamableHTTPServer struct {
	srv    *http.Server
	logger *slog.Logger
}

// Start starts an HTTP server in a separate go routine and returns a Server interface that can be used to stop the server.
// Use `srv.Wait()` to wait for the server to receive a shutdown signal (`syscall.SIGINT` or `syscall.SIGTERM`)
func NewStreamableHTTPServer(logger *slog.Logger, router Router, port int) *StreamableHTTPServer {
	mux := http.NewServeMux()
	mux.Handle("/_health", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("Health check request", "method", r.Method, "uri", r.RequestURI)
		w.WriteHeader(http.StatusOK)
	}))
	mux.Handle("/mcp", NewHTTPHandler(router, logger))
	srv := &http.Server{
		Addr:        fmt.Sprintf("127.0.0.1:%d", port),
		Handler:     mux,
		ReadTimeout: 10 * time.Second,
	}
	return &StreamableHTTPServer{
		srv:    srv,
		logger: logger,
	}
}

// Start starts the server
func (s *StreamableHTTPServer) Start() {
	// see https://dev.to/mokiat/proper-http-shutdown-in-go-3fji
	go func() {
		if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			s.logger.Error("Streamable HTTP server error", "error", err.Error())
		}
		s.logger.Info("Stopped serving new connections.")
	}()
	s.logger.Info("Streamable HTTP server started")
}

// Wait waits for the server to receive a shutdown signal (`syscall.SIGINT` or `syscall.SIGTERM`)
func (s *StreamableHTTPServer) Wait() error {
	// Shutdown gracefully shuts down the server without interrupting any active connections
	s.logger.Info("Streamable HTTP server waiting for shutdown signal...")
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	s.logger.Info("Streamable HTTP server received shutdown signal", "signal", sig)
	return nil
}

// Stop stops the server gracefully
func (s *StreamableHTTPServer) Stop() error {
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownRelease()
	if err := s.srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("HTTP shutdown error: %v", err.Error())
	}
	s.logger.Info("Streamable HTTP server gracefully stopped")
	return nil
}

func (s *StreamableHTTPServer) Addr() string {
	return s.srv.Addr
}

func NewHTTPHandler(router Router, logger *slog.Logger) http.Handler {
	return jhttp.NewBridge(handler.Map(router), &jhttp.BridgeOptions{
		Client: &jrpc2.ClientOptions{
			Logger: SlogToLogBridge(logger),
		},
		Server: &jrpc2.ServerOptions{
			Logger: SlogToLogBridge(logger),
			RPCLog: SlogToRPCLogBridge(logger),
		},
	})
}

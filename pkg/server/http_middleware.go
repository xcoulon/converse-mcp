package server

import (
	"log/slog"
	"net/http"
)

func LoggingMiddleware(logger *slog.Logger, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Debug("HTTP request", "method", r.Method, "uri", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}

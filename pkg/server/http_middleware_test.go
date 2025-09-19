package server_test

import (
	"bytes"
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/xcoulon/converse-mcp/pkg/server"

	"github.com/stretchr/testify/require"
)

func TestLoggingMiddleware(t *testing.T) {
	// given
	client := http.Client{}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	testCases := []struct {
		name         string
		loggingLevel slog.Level
		expectedLog  bool
	}{
		{
			name:         "logging debug",
			loggingLevel: slog.LevelDebug,
			expectedLog:  true,
		},
		{
			name:         "logging info",
			loggingLevel: slog.LevelInfo,
			expectedLog:  false,
		},
		{
			name:         "logging warn",
			loggingLevel: slog.LevelWarn,
			expectedLog:  false,
		},
		{
			name:         "logging error",
			loggingLevel: slog.LevelError,
			expectedLog:  false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// given
			out := bytes.NewBuffer(nil)
			logger := slog.New(slog.NewTextHandler(out, &slog.HandlerOptions{Level: testCase.loggingLevel}))
			next := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(http.StatusOK)
			})
			httpSrv := httptest.NewServer(server.LoggingMiddleware(logger, next))
			defer httpSrv.Close()
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, httpSrv.URL, nil)
			require.NoError(t, err)

			// when
			resp, err := client.Do(req)

			// then
			require.NoError(t, err)
			defer resp.Body.Close()
			require.Equal(t, http.StatusOK, resp.StatusCode)
			if testCase.expectedLog {
				require.Contains(t, out.String(), "HTTP request")
			} else {
				require.NotContains(t, out.String(), "HTTP request")
			}
		})
	}
}

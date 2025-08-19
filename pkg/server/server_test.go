package server_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/jhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	api "github.com/xcoulon/converse-mcp/pkg/api"
	"github.com/xcoulon/converse-mcp/pkg/server"
)

var EmptyPromptHandle server.PromptHandleFunc = func(_ context.Context, _ *slog.Logger, _ api.GetPromptRequestParams) (any, error) {
	return nil, nil
}

var EmptyResourceHandle server.ResourceHandleFunc = func(_ context.Context, _ *slog.Logger, _ api.ReadResourceRequestParams) (any, error) {
	return nil, nil
}

var EmptyToolHandle server.ToolHandleFunc = func(_ context.Context, _ *slog.Logger, _ api.CallToolRequestParams) (any, error) {
	return nil, nil
}

func TestServer(t *testing.T) {

	// given
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	mux := server.NewMux("converse-mcp", "0.1", logger).
		WithPrompt(api.NewPrompt("my-first-prompt"), EmptyPromptHandle).
		WithPrompt(api.NewPrompt("my-second-prompt"), EmptyPromptHandle).
		WithResource(api.NewResource("my-first-resource", "https://example.com/my-first-resource"), EmptyResourceHandle).
		WithResource(api.NewResource("my-second-resource", "https://example.com/my-second-resource"), EmptyResourceHandle).
		WithTool(api.NewTool("my-first-tool"), EmptyToolHandle).
		WithTool(api.NewTool("my-second-tool"), EmptyToolHandle).
		Build()
	// stdio server
	c2s, s2c := channel.Direct()
	directCl := jrpc2.NewClient(c2s, &jrpc2.ClientOptions{})
	stdioSrv := server.NewStdioServer(mux, logger)
	stdioSrv.Start(s2c)
	defer directCl.Close()
	defer stdioSrv.Stop()

	// http server
	hsrv := httptest.NewServer(server.NewHTTPHandler(mux, logger))
	httpCl := jrpc2.NewClient(jhttp.NewChannel(hsrv.URL, nil), nil)
	defer httpCl.Close()
	defer hsrv.Close()

	for name, cl := range map[string]*jrpc2.Client{
		"stdio": directCl,
		"http":  httpCl,
	} {
		t.Run(name, func(t *testing.T) {
			t.Run("initialize", func(t *testing.T) {
				// when
				resp, err := cl.Call(context.Background(), "initialize", api.InitializeRequestParams{})

				// then
				require.NoError(t, err)
				expected := api.InitializeResult{
					ProtocolVersion: "2025-06-18",
					ServerInfo: api.Implementation{
						Name:    "converse-mcp",
						Version: "0.1",
					},
					Capabilities: api.DefaultCapabilities,
				}
				expectedJSON, err := json.Marshal(expected)
				require.NoError(t, err)
				assert.JSONEq(t, string(expectedJSON), resp.ResultString())
			})

			t.Run("list prompts", func(t *testing.T) {
				// when
				resp, err := cl.Call(context.Background(), "prompts/list", api.ListResourcesRequestParams{})

				// then
				require.NoError(t, err)
				expected := api.ListPromptsResult{
					Prompts: []api.Prompt{
						{
							Name: "my-first-prompt",
						},
						{
							Name: "my-second-prompt",
						},
					},
				}
				expectedJSON, _ := json.Marshal(expected)
				assert.JSONEq(t, string(expectedJSON), resp.ResultString())
			})

			t.Run("list resources", func(t *testing.T) {
				// when
				resp, err := cl.Call(context.Background(), "resources/list", api.ListResourcesRequestParams{})

				// then
				require.NoError(t, err)
				expected := api.ListResourcesResult{
					Resources: []api.Resource{
						{
							Name: "my-first-resource",
							Uri:  "https://example.com/my-first-resource",
						},
						{
							Name: "my-second-resource",
							Uri:  "https://example.com/my-second-resource",
						},
					},
				}
				expectedJSON, _ := json.Marshal(expected)
				assert.JSONEq(t, string(expectedJSON), resp.ResultString())
			})

			t.Run("list tools", func(t *testing.T) {
				// when
				resp, err := cl.Call(context.Background(), "tools/list", api.ListResourcesRequestParams{})

				// then
				require.NoError(t, err)
				expected := api.ListToolsResult{
					Tools: []api.Tool{
						{
							Name: "my-first-tool",
							InputSchema: api.ToolInputSchema{
								Type: "object",
							},
						},
						{
							Name: "my-second-tool",
							InputSchema: api.ToolInputSchema{
								Type: "object",
							},
						},
					},
				}
				expectedJSON, _ := json.Marshal(expected)
				assert.JSONEq(t, string(expectedJSON), resp.ResultString())
			})
		})
	}
}

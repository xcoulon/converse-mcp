package server_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/xcoulon/converse-mcp/pkg/api"
	"github.com/xcoulon/converse-mcp/pkg/server"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/jhttp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var EmptyPromptHandle server.PromptHandleFunc = func(_ context.Context, _ api.GetPromptRequestParams) (api.GetPromptResult, error) {
	return api.GetPromptResult{}, nil
}

var EmptyResourceHandle server.ResourceHandleFunc = func(_ context.Context, _ api.ReadResourceRequestParams) (api.ReadResourceResult, error) {
	return api.ReadResourceResult{}, nil
}

var EmptyToolHandle server.ToolHandleFunc = func(_ context.Context, _ api.CallToolRequestParams) (api.CallToolResult, error) {
	return api.CallToolResult{}, nil
}

func TestServer(t *testing.T) {

	// given
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	router := server.NewRouterBuilder("converse-mcp", "0.1", logger).
		WithPrompt(api.NewPrompt("my-first-prompt"), EmptyPromptHandle).
		WithPrompt(api.NewPrompt("my-second-prompt"), EmptyPromptHandle).
		WithResource(api.NewResource("my-first-resource", "https://example.com/my-first-resource"), EmptyResourceHandle).
		WithResource(api.NewResource("my-second-resource", "https://example.com/my-second-resource"), EmptyResourceHandle).
		WithTool(api.NewTool("my-first-tool"), EmptyToolHandle).
		WithTool(api.NewTool("my-second-tool"), EmptyToolHandle).
		Build()
	// stdio server
	c2s, s2c := channel.Direct()
	stdioCl := jrpc2.NewClient(c2s, &jrpc2.ClientOptions{})
	stdioSrv := server.NewStdioServer(logger, router)
	stdioSrv.Start(s2c)
	defer func() {
		require.NoError(t, stdioCl.Close())
		stdioSrv.Stop()
	}()

	// http server
	httpSrv := httptest.NewServer(server.NewHTTPHandler(router, logger))
	httpCl := jrpc2.NewClient(jhttp.NewChannel(httpSrv.URL, nil), nil)
	defer func() {
		require.NoError(t, httpCl.Close())
	}()
	defer httpSrv.Close()

	for name, cl := range map[string]*jrpc2.Client{
		"stdio": stdioCl,
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
					Capabilities: api.ServerCapabilities{
						Prompts: &api.ServerCapabilitiesPrompts{
							ListChanged: api.BoolPtr(true),
						},
						Resources: &api.ServerCapabilitiesResources{
							ListChanged: api.BoolPtr(true),
						},
						Tools: &api.ServerCapabilitiesTools{
							ListChanged: api.BoolPtr(true),
						},
					},
				}
				expectedJSON, _ := json.Marshal(expected)
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
							Name:        "my-first-tool",
							Annotations: &api.ToolAnnotations{},
							InputSchema: api.ToolInputSchema{
								Type: "object",
							},
							OutputSchema: &api.ToolOutputSchema{
								Type: "object",
							},
						},
						{
							Name:        "my-second-tool",
							Annotations: &api.ToolAnnotations{},
							InputSchema: api.ToolInputSchema{
								Type: "object",
							},
							OutputSchema: &api.ToolOutputSchema{
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

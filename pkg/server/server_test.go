package server_test

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"testing"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xcoulon/converse/pkg/server"
	"github.com/xcoulon/converse/pkg/types"
)

func TestServer(t *testing.T) {
	// given
	logger := slog.New(slog.NewTextHandler(os.Stderr, nil))
	c2s, s2c := channel.Direct()
	cl := jrpc2.NewClient(c2s, &jrpc2.ClientOptions{})
	srv := server.New("converse", "0.1").
		Prompt(types.Prompt{Name: "my-first-prompt"}, server.EmptyPromptHandle).
		Prompt(types.Prompt{Name: "my-second-prompt"}, server.EmptyPromptHandle).
		Resource(types.Resource{Name: "my-first-resource"}, server.EmptyResourceHandle).
		Resource(types.Resource{Name: "my-second-resource"}, server.EmptyResourceHandle).
		Tool(types.Tool{Name: "my-first-tool"}, server.EmptyToolHandle).
		Tool(types.Tool{Name: "my-second-tool"}, server.EmptyToolHandle).
		Start(logger, s2c)
	defer func(cl *jrpc2.Client, srv *jrpc2.Server) {
		// close the streams
		cl.Close()
		srv.Stop()
	}(cl, srv)

	t.Run("initialize", func(t *testing.T) {
		// when
		resp, err := cl.Call(context.Background(), "initialize", types.InitializeRequestParams{})

		// then
		require.NoError(t, err)
		expected := types.InitializeResult{
			ProtocolVersion: "2025-03-26",
			ServerInfo: types.Implementation{
				Name:    "converse",
				Version: "0.1",
			},
			Capabilities: types.DefaultCapabilities,
		}
		expectedJSON, err := json.Marshal(expected)
		require.NoError(t, err)
		assert.JSONEq(t, string(expectedJSON), resp.ResultString())
	})

	t.Run("list prompts", func(t *testing.T) {
		// when
		resp, err := cl.Call(context.Background(), "prompts/list", types.ListResourcesRequestParams{})

		// then
		require.NoError(t, err)
		expected := types.ListPromptsResult{
			Prompts: []types.Prompt{
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
		resp, err := cl.Call(context.Background(), "resources/list", types.ListResourcesRequestParams{})

		// then
		require.NoError(t, err)
		expected := types.ListResourcesResult{
			Resources: []types.Resource{
				{
					Name: "my-first-resource",
				},
				{
					Name: "my-second-resource",
				},
			},
		}
		expectedJSON, _ := json.Marshal(expected)
		assert.JSONEq(t, string(expectedJSON), resp.ResultString())
	})

	t.Run("list tools", func(t *testing.T) {
		// when
		resp, err := cl.Call(context.Background(), "tools/list", types.ListResourcesRequestParams{})

		// then
		require.NoError(t, err)
		expected := types.ListToolsResult{
			Tools: []types.Tool{
				{
					Name: "my-first-tool",
				},
				{
					Name: "my-second-tool",
				},
			},
		}
		expectedJSON, _ := json.Marshal(expected)
		assert.JSONEq(t, string(expectedJSON), resp.ResultString())
	})
}

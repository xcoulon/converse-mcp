package server

import (
	"context"
	"log/slog"

	api "github.com/xcoulon/converse-mcp/pkg/api"
)

type PromptHandleFunc func(ctx context.Context, logger *slog.Logger, params api.GetPromptRequestParams) (any, error)

var EmptyPromptHandle PromptHandleFunc = func(_ context.Context, _ *slog.Logger, _ api.GetPromptRequestParams) (any, error) {
	return nil, nil
}

type PromptHandler struct {
	api.Prompt
	Handle PromptHandleFunc
}

type ResourceHandleFunc func(ctx context.Context, logger *slog.Logger, params api.ReadResourceRequestParams) (any, error)

var EmptyResourceHandle ResourceHandleFunc = func(_ context.Context, _ *slog.Logger, _ api.ReadResourceRequestParams) (any, error) {
	return nil, nil
}

type ResourceHandler struct {
	api.Resource
	Handle ResourceHandleFunc
}

type ToolHandleFunc func(ctx context.Context, logger *slog.Logger, params api.CallToolRequestParams) (any, error)

var EmptyToolHandle ToolHandleFunc = func(_ context.Context, _ *slog.Logger, _ api.CallToolRequestParams) (any, error) {
	return nil, nil
}

type ToolHandler struct {
	api.Tool
	Handle ToolHandleFunc
}

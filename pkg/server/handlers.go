package server

import (
	"context"
	"log/slog"

	"github.com/xcoulon/converse/pkg/types"
)

type PromptHandleFunc func(ctx context.Context, logger *slog.Logger, params types.GetPromptRequestParams) (any, error)

var EmptyPromptHandle PromptHandleFunc = func(_ context.Context, _ *slog.Logger, _ types.GetPromptRequestParams) (any, error) {
	return nil, nil
}

type PromptHandler struct {
	types.Prompt
	Handle PromptHandleFunc
}

type ResourceHandleFunc func(ctx context.Context, logger *slog.Logger, params types.ReadResourceRequestParams) (any, error)

var EmptyResourceHandle ResourceHandleFunc = func(_ context.Context, _ *slog.Logger, _ types.ReadResourceRequestParams) (any, error) {
	return nil, nil
}

type ResourceHandler struct {
	types.Resource
	Handle ResourceHandleFunc
}

type ToolHandleFunc func(ctx context.Context, logger *slog.Logger, params types.CallToolRequestParams) (any, error)

var EmptyToolHandle ToolHandleFunc = func(_ context.Context, _ *slog.Logger, _ types.CallToolRequestParams) (any, error) {
	return nil, nil
}

type ToolHandler struct {
	types.Tool
	Handle ToolHandleFunc
}

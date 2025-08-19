package server

import (
	"context"
	"log/slog"

	api "github.com/xcoulon/converse-mcp/pkg/api"
)

// TODO: remove logger
type PromptHandleFunc func(ctx context.Context, logger *slog.Logger, params api.GetPromptRequestParams) (any, error)

type PromptHandler struct {
	Prompt api.Prompt
	Handle PromptHandleFunc
}

// TODO: remove logger
type ResourceHandleFunc func(ctx context.Context, logger *slog.Logger, params api.ReadResourceRequestParams) (any, error)

type ResourceHandler struct {
	Resource api.Resource
	Handle   ResourceHandleFunc
}

// TODO: remove logger
type ToolHandleFunc func(ctx context.Context, logger *slog.Logger, params api.CallToolRequestParams) (any, error)

type ToolHandler struct {
	Tool   api.Tool
	Handle ToolHandleFunc
}

package server

import (
	"context"

	api "github.com/xcoulon/converse-mcp/pkg/api"
)

type PromptHandleFunc func(ctx context.Context, params api.GetPromptRequestParams) (any, error)

type PromptHandler struct {
	Prompt api.Prompt
	Handle PromptHandleFunc
}

type ResourceHandleFunc func(ctx context.Context, params api.ReadResourceRequestParams) (any, error)

type ResourceHandler struct {
	Resource api.Resource
	Handle   ResourceHandleFunc
}

type ToolHandleFunc func(ctx context.Context, params api.CallToolRequestParams) (any, error)

type ToolHandler struct {
	Tool   api.Tool
	Handle ToolHandleFunc
}

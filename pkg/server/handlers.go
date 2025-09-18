package server

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/handler"
	api "github.com/xcoulon/converse-mcp/pkg/api"
)

const protocolVersion = "2025-06-18"

type PromptHandleFunc func(ctx context.Context, params api.GetPromptRequestParams) (api.GetPromptResult, error)

type PromptHandler struct {
	Prompt api.Prompt
	Handle PromptHandleFunc
}

type ResourceHandleFunc func(ctx context.Context, params api.ReadResourceRequestParams) (api.ReadResourceResult, error)

type ResourceHandler struct {
	Resource api.Resource
	Handle   ResourceHandleFunc
}

type ToolHandleFunc func(ctx context.Context, params api.CallToolRequestParams) (api.CallToolResult, error)

type ToolHandler struct {
	Tool   api.Tool
	Handle ToolHandleFunc
}

type Router handler.Map
type RouterBuilder struct {
	capabilities api.ServerCapabilities
	serverInfo   api.Implementation
	prompts      []PromptHandler
	resources    []ResourceHandler
	tools        []ToolHandler
	logger       *slog.Logger
}

func NewRouterBuilder(name, version string, logger *slog.Logger) *RouterBuilder {
	return &RouterBuilder{
		capabilities: api.ServerCapabilities{
			Prompts: &api.ServerCapabilitiesPrompts{
				ListChanged: api.BoolPtr(false), // default to false, until a prompt is added
			},
			Resources: &api.ServerCapabilitiesResources{
				ListChanged: api.BoolPtr(false), // default to false, until a resource is added
			},
			Tools: &api.ServerCapabilitiesTools{
				ListChanged: api.BoolPtr(false), // default to false, until a tool is added
			},
			// Logging:
		},
		serverInfo: api.Implementation{
			Name:    name,
			Version: version,
		},
		prompts:   []PromptHandler{},
		resources: []ResourceHandler{},
		tools:     []ToolHandler{},
		logger:    logger,
	}
}

func (b *RouterBuilder) WithPrompt(prompt api.Prompt, handle PromptHandleFunc) *RouterBuilder {
	b.logger.Debug("with prompt", "prompt", prompt.Name)
	b.prompts = append(b.prompts, PromptHandler{
		Prompt: prompt,
		Handle: handle,
	})
	// Servers that support prompts MUST declare the prompts capability
	b.capabilities.Prompts.ListChanged = api.BoolPtr(true)
	return b
}

func (b *RouterBuilder) WithResource(resource api.Resource, handle ResourceHandleFunc) *RouterBuilder {
	b.logger.Debug("with resource", "resource", resource.Name)
	b.resources = append(b.resources, ResourceHandler{
		Resource: resource,
		Handle:   handle,
	})
	// Servers that support resources MUST declare the resources capability
	b.capabilities.Resources.ListChanged = api.BoolPtr(true)
	return b
}

func (b *RouterBuilder) WithTool(tool api.Tool, handle ToolHandleFunc) *RouterBuilder {
	b.logger.Debug("with tool", "tool", tool.Name)
	b.tools = append(b.tools, ToolHandler{
		Tool:   tool,
		Handle: handle,
	})
	// Servers that support tools MUST declare the tools capability
	b.capabilities.Tools.ListChanged = api.BoolPtr(true)
	return b
}

func (b *RouterBuilder) Build() Router {
	return Router(handler.Map{
		"initialize":     initialize(b.capabilities, b.serverInfo, b.logger),
		"prompts/list":   listPrompts(b.prompts, b.logger),
		"prompts/get":    getPrompt(b.prompts, b.logger),
		"resources/list": listResources(b.resources, b.logger),
		"resources/read": readResource(b.resources, b.logger),
		"tools/list":     listTools(b.tools, b.logger),
		"tools/call":     callTool(b.tools, b.logger),
	})
}

func initialize(capabilities api.ServerCapabilities, serverInfo api.Implementation, logger *slog.Logger) jrpc2.Handler {
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		logger.Debug("initialize")
		return &api.InitializeResult{
			ProtocolVersion: protocolVersion,
			ServerInfo:      serverInfo,
			Capabilities:    capabilities,
		}, nil
	}
}

func listPrompts(handlers []PromptHandler, logger *slog.Logger) jrpc2.Handler {
	prompts := make([]api.Prompt, 0, len(handlers))
	for _, h := range handlers {
		prompts = append(prompts, h.Prompt)
	}
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		logger.Debug("list prompts")
		return &api.ListPromptsResult{
			Prompts: prompts,
		}, nil
	}
}

func getPrompt(handlers []PromptHandler, logger *slog.Logger) jrpc2.Handler {
	prompts := make(map[string]PromptHandler, len(handlers))
	for _, h := range handlers {
		prompts[h.Prompt.Name] = h
	}
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		params := api.GetPromptRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			return nil, fmt.Errorf("error while unmarshalling '%s' request parameters: %w", req.Method(), err)
		}
		logger.Debug("get prompt", "name", params.Name)
		if h, ok := prompts[params.Name]; ok {
			return h.Handle(ctx, params)
		}
		return nil, fmt.Errorf("prompt '%s' does not exist", params.Name)
	}
}

func listResources(handlers []ResourceHandler, logger *slog.Logger) jrpc2.Handler {
	resources := make([]api.Resource, 0, len(handlers))
	for _, h := range handlers {
		resources = append(resources, h.Resource)
	}
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		logger.Debug("list resources")
		return &api.ListResourcesResult{
			Resources: resources,
		}, nil
	}
}

func readResource(handlers []ResourceHandler, logger *slog.Logger) jrpc2.Handler {
	resources := make(map[string]ResourceHandler, len(handlers))
	for _, h := range handlers {
		resources[h.Resource.Name] = h
	}
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		params := api.ReadResourceRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			return nil, fmt.Errorf("error while unmarshalling '%s' request parameters: %w", req.Method(), err)
		}
		logger.Debug("read resource", "uri", params.Uri)
		if h, ok := resources[params.Uri]; ok {
			return h.Handle(ctx, params)
		}
		return nil, fmt.Errorf("resource '%s' does not exist", params.Uri)
	}
}

func listTools(handlers []ToolHandler, logger *slog.Logger) jrpc2.Handler {
	tools := make([]api.Tool, 0, len(handlers))
	for _, h := range handlers {
		tools = append(tools, h.Tool)
	}
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		logger.Debug("list tools")
		return &api.ListToolsResult{
			Tools: tools,
		}, nil
	}
}

func callTool(handlers []ToolHandler, logger *slog.Logger) jrpc2.Handler {
	tools := make(map[string]ToolHandler, len(handlers))
	for _, h := range handlers {
		tools[h.Tool.Name] = h
	}
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		params := api.CallToolRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			return nil, fmt.Errorf("error while unmarshalling '%s' request parameters: %w", req.Method(), err)
		}
		logger.Debug("call tool", "name", params.Name)
		if h, ok := tools[params.Name]; ok {
			return h.Handle(ctx, params)
		}
		return nil, fmt.Errorf("tool '%s' does not exist", params.Name)
	}
}

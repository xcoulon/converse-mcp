package server

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/creachadair/jrpc2/jhttp"
	"github.com/xcoulon/converse-mcp/pkg/api"
)

const protocolVersion = "2025-06-18"

var StdioChannel = channel.Line(os.Stdin, os.Stdout)

func NewStdioServer(mux handler.Map, logger *slog.Logger) *jrpc2.Server {
	return jrpc2.NewServer(mux, &jrpc2.ServerOptions{
		Logger: SlogToLogBridge(logger),
	})
}

func NewHTTPHandler(mux handler.Map, logger *slog.Logger) http.Handler {
	return jhttp.NewBridge(mux, &jhttp.BridgeOptions{
		Server: &jrpc2.ServerOptions{
			Logger: SlogToLogBridge(logger),
		},
	})
}

type MuxBuilder struct {
	capabilities api.ServerCapabilities
	serverInfo   api.Implementation
	prompts      []PromptHandler
	resources    []ResourceHandler
	tools        []ToolHandler
	logger       *slog.Logger
}

func NewMux(name, version string, logger *slog.Logger) *MuxBuilder {
	return &MuxBuilder{
		capabilities: api.ServerCapabilities{
			Prompts: &api.ServerCapabilitiesPrompts{
				ListChanged: api.BoolPtr(false),
			},
			Resources: &api.ServerCapabilitiesResources{
				ListChanged: api.BoolPtr(false),
			},
			Tools: &api.ServerCapabilitiesTools{
				ListChanged: api.BoolPtr(false),
			},
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

func (b *MuxBuilder) WithPrompt(prompt api.Prompt, handle PromptHandleFunc) *MuxBuilder {
	b.prompts = append(b.prompts, PromptHandler{
		Prompt: prompt,
		Handle: handle,
	})
	// Servers that support prompts MUST declare the prompts capability
	b.capabilities.Prompts.ListChanged = api.BoolPtr(true)
	return b
}

func (b *MuxBuilder) WithResource(resource api.Resource, handle ResourceHandleFunc) *MuxBuilder {
	b.resources = append(b.resources, ResourceHandler{
		Resource: resource,
		Handle:   handle,
	})
	// Servers that support resources MUST declare the resources capability
	b.capabilities.Resources.ListChanged = api.BoolPtr(true)
	return b
}

func (b *MuxBuilder) WithTool(tool api.Tool, handle ToolHandleFunc) *MuxBuilder {
	b.tools = append(b.tools, ToolHandler{
		Tool:   tool,
		Handle: handle,
	})
	// Servers that support tools MUST declare the tools capability
	b.capabilities.Tools.ListChanged = api.BoolPtr(true)
	return b
}

func (b *MuxBuilder) Build() handler.Map {
	return handler.Map{
		"initialize":     initialize(b.capabilities, b.serverInfo),
		"prompts/list":   listPrompts(b.prompts),
		"prompts/get":    getPrompt(b.logger, b.prompts),
		"resources/list": listResources(b.resources),
		"resources/read": readResource(b.logger, b.resources),
		"tools/list":     listTools(b.tools),
		"tools/call":     callTool(b.logger, b.tools),
	}
}

func SlogToLogBridge(logger *slog.Logger) jrpc2.Logger {
	return func(text string) {
		if err := logger.Handler().Handle(context.Background(), slog.Record{
			Level:   slog.LevelInfo,
			Message: text,
		}); err != nil {
			logger.Error("error logging message", "error", err)
		}
	}
}

func initialize(capabilities api.ServerCapabilities, serverInfo api.Implementation) jrpc2.Handler {
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		return &api.InitializeResult{
			ProtocolVersion: protocolVersion,
			ServerInfo:      serverInfo,
			Capabilities:    capabilities,
		}, nil
	}
}

func listPrompts(handlers []PromptHandler) jrpc2.Handler {
	prompts := make([]api.Prompt, 0, len(handlers))
	for _, h := range handlers {
		prompts = append(prompts, h.Prompt)
	}
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		return &api.ListPromptsResult{
			Prompts: prompts,
		}, nil
	}
}

func getPrompt(logger *slog.Logger, handlers []PromptHandler) jrpc2.Handler {
	prompts := make(map[string]PromptHandler, len(handlers))
	for _, h := range handlers {
		prompts[h.Prompt.Name] = h
	}
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		params := api.GetPromptRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			return nil, fmt.Errorf("error while unmarshalling '%s' request parameters: %w", req.Method(), err)
		}
		if h, ok := prompts[params.Name]; ok {
			return h.Handle(ctx, logger, params)
		}
		return nil, fmt.Errorf("prompt '%s' does not exist", params.Name)
	}
}

func listResources(handlers []ResourceHandler) jrpc2.Handler {
	resources := make([]api.Resource, 0, len(handlers))
	for _, h := range handlers {
		resources = append(resources, h.Resource)
	}
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		return &api.ListResourcesResult{
			Resources: resources,
		}, nil
	}
}

func readResource(logger *slog.Logger, handlers []ResourceHandler) jrpc2.Handler {
	resources := make(map[string]ResourceHandler, len(handlers))
	for _, h := range handlers {
		resources[h.Resource.Name] = h
	}
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		params := api.ReadResourceRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			return nil, fmt.Errorf("error while unmarshalling '%s' request parameters: %w", req.Method(), err)
		}
		if h, ok := resources[params.Uri]; ok {
			return h.Handle(ctx, logger, params)
		}
		return nil, fmt.Errorf("resource '%s' does not exist", params.Uri)
	}
}

func listTools(handlers []ToolHandler) jrpc2.Handler {
	tools := make([]api.Tool, 0, len(handlers))
	for _, h := range handlers {
		tools = append(tools, h.Tool)
	}
	return func(_ context.Context, _ *jrpc2.Request) (any, error) {
		return &api.ListToolsResult{
			Tools: tools,
		}, nil
	}
}

func callTool(logger *slog.Logger, handlers []ToolHandler) jrpc2.Handler {
	tools := make(map[string]ToolHandler, len(handlers))
	for _, h := range handlers {
		tools[h.Tool.Name] = h
	}
	return func(ctx context.Context, req *jrpc2.Request) (any, error) {
		params := api.CallToolRequestParams{}
		if err := req.UnmarshalParams(&params); err != nil {
			return nil, fmt.Errorf("error while unmarshalling '%s' request parameters: %w", req.Method(), err)
		}
		if h, ok := tools[params.Name]; ok {
			return h.Handle(ctx, logger, params)
		}
		return nil, fmt.Errorf("tool '%s' does not exist", params.Name)
	}
}

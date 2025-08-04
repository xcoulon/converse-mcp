package server

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/creachadair/jrpc2"
	"github.com/creachadair/jrpc2/channel"
	"github.com/creachadair/jrpc2/handler"
	"github.com/xcoulon/converse-mcp/pkg/api"
)

var StdioChannel = channel.Line(os.Stdin, os.Stdout)

type Builder struct {
	capabilities api.ServerCapabilities
	serverInfo   api.Implementation
	prompts      []PromptHandler
	resources    []ResourceHandler
	tools        []ToolHandler
}

func New(name, version string, capabilities ...api.ServerCapability) *Builder {
	sc := api.DefaultCapabilities
	for _, apply := range capabilities {
		apply(&sc)
	}
	return &Builder{
		capabilities: sc,
		serverInfo: api.Implementation{
			Name:    name,
			Version: version,
		},
		prompts:   []PromptHandler{},
		resources: []ResourceHandler{},
		tools:     []ToolHandler{},
	}
}

func (b *Builder) Prompt(prompt api.Prompt, handle PromptHandleFunc) *Builder {
	b.prompts = append(b.prompts, PromptHandler{
		Prompt: prompt,
		Handle: handle,
	})
	return b
}

func (b *Builder) Resource(resource api.Resource, handle ResourceHandleFunc) *Builder {
	b.resources = append(b.resources, ResourceHandler{
		Resource: resource,
		Handle:   handle,
	})
	return b
}

func (b *Builder) Tools(tools ...ToolHandler) *Builder {
	b.tools = tools
	return b
}

func (b *Builder) Tool(tool api.Tool, handle ToolHandleFunc) *Builder {
	b.tools = append(b.tools, ToolHandler{
		Tool:   tool,
		Handle: handle,
	})
	return b
}

func (b *Builder) Start(logger *slog.Logger, c channel.Channel) *jrpc2.Server {
	mux := handler.Map{
		"initialize":     initialize(b.capabilities, b.serverInfo),
		"prompts/list":   listPrompts(b.prompts),
		"prompts/get":    getPrompt(logger, b.prompts),
		"resources/list": listResources(b.resources),
		"resources/read": readResource(logger, b.resources),
		"tools/list":     listTools(b.tools),
		"tools/call":     callTool(logger, b.tools),
	}
	opts := &jrpc2.ServerOptions{
		// Logger: jrpc2.StdLogger(logger),
	}
	s := jrpc2.NewServer(mux, opts)
	return s.Start(c)
}

var protocolVersion = "2025-03-26"

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
		prompts[h.Name] = h
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
		resources[h.Name] = h
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
		tools[h.Name] = h
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

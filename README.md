# Converse

Converse empowers developers to rapidly build Model Context Protocol Servers in Go with unprecedented ease.

## Building and Installing

Requires Go 1.24 and [Task](https://taskfile.dev/)

```
task test lint
```

### Generating the MCP types

1. Download the latest stable version of the JSON schema from the [modelcontextprotocol/specification](https://github.com/modelcontextprotocol/specification/blob/main/schema/) repository

2. Install the `go-jsonschema` generator from https://github.com/omissis/go-jsonschema then run:
   ```
   go install github.com/atombender/go-jsonschema@latest
   go-jsonschema -p internal/types resources/schema.json > pkg/types/types.go
   ```

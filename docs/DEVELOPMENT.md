# Development Guide

## Prerequisites

- Go 1.24.2 or higher

## Building

```bash
# Build MCP server with stdio transport
go build -o bin/mcp-server-stdio ./cmd/mcp-server-stdio

# Build MCP server with HTTP transport
go build -o bin/mcp-server-http ./cmd/mcp-server-http

# Build Mock API server
go build -o bin/mock-api ./cmd/mock-api
```

## Running Tests

### Manual Testing with Mock API

```bash
# Start Mock API server
go run ./cmd/mock-api &

# Test MCP server (stdio)
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}' | go run ./cmd/mcp-server-stdio

# Test MCP server (HTTP)
go run ./cmd/mcp-server-http -port 8080 &
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "ping"}'
```

### Mock API Configuration

The Mock API server can be configured with different API sets:

```bash
# Start with default configuration (users API)
./bin/mock-api

# Start with products configuration
MOCK_CONFIG=configs/mock/products.json ./bin/mock-api

# Or run directly
go run ./cmd/mock-api
```

For detailed Mock API documentation, see **[Mock API Documentation](MOCK-API.md)**.

## Command Line Options

### MCP Server (Stdio)
```bash
./bin/mcp-server-stdio -config ./config.json
./bin/mcp-server-stdio -api-url http://localhost:8080
./bin/mcp-server-stdio -verbose
```

### MCP Server (HTTP)
```bash
./bin/mcp-server-http -port 8080 -host localhost -cors
./bin/mcp-server-http -config ./example-config.json -port 8080
./bin/mcp-server-http -verbose
```

## Adding New Features

### Adding Custom Endpoints

1. Update your configuration file's `endpoints` section
2. No code changes required - the bridge automatically handles new endpoints

### Adding New Transport Types

1. Implement the transport interface in `internal/transport/`
2. Add command-line handling in the appropriate `cmd/` directory
3. Update configuration parsing if needed

### Adding Authentication Methods

1. Extend the auth configuration in `internal/config/`
2. Implement the auth handler in `internal/bridge/`
3. Update the configuration documentation

## Code Structure

- `cmd/`: Entry points for different executables
- `internal/mcp/`: MCP protocol implementation
- `internal/bridge/`: REST API to MCP conversion logic
- `internal/transport/`: Transport layer implementations
- `internal/config/`: Configuration management
- `pkg/types/`: Shared type definitions

## Debugging

Enable verbose logging with the `-verbose` flag:

```bash
go run ./cmd/mcp-server-stdio -verbose
go run ./cmd/mcp-server-http -verbose
```

This will output detailed information about:
- MCP protocol messages
- REST API calls
- Configuration parsing
- Error handling
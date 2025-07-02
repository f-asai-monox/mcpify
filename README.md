# MCP Bridge

A proxy server that enables REST APIs to be used as MCP (Model Context Protocol) servers.

## Overview

MCP Bridge is a proxy server that allows existing REST APIs to be used as Model Context Protocol (MCP) servers. This enables REST APIs to be directly utilized by MCP clients such as Claude Code.

## Features

- **REST API to MCP Conversion**: Automatically converts REST API endpoints to MCP tools
- **JSON-RPC 2.0 Compliant**: Fully compliant with the MCP protocol
- **Configurable**: Flexible customization through configuration files
- **Mock API Server**: Built-in simple REST API server for testing

## Project Structure

```
mcp-bridge/
├── cmd/
│   ├── mcp-server/     # MCP server executable
│   └── mock-api/       # REST API mock server
├── internal/
│   ├── mcp/           # MCP implementation
│   ├── bridge/        # REST API conversion logic
│   └── config/        # Configuration management
├── pkg/
│   └── types/         # Common type definitions
├── go.mod
├── go.sum
├── README.md
└── README-ja.md       # Japanese version
```

## Installation & Build

### Dependencies
- Go 1.21 or higher

### Build

```bash
# Build MCP server
go build -o bin/mcp-server ./cmd/mcp-server

# Build Mock API server
go build -o bin/mock-api ./cmd/mock-api
```

## Usage

### 1. Start Mock API Server

Start the test REST API server:

```bash
./bin/mock-api

# Or run directly
go run ./cmd/mock-api
```

The API server starts at `http://localhost:8080` with the following endpoints:

- `GET /health` - Health check
- `GET /users` - Get all users
- `POST /users` - Create user
- `GET /users/{id}` - Get specific user
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user
- `GET /products` - Get all products
- `GET /products?category={category}` - Get products by category
- `POST /products` - Create product

### 2. Start MCP Server

Start the MCP bridge server:

```bash
./bin/mcp-server

# Or specify config file
./bin/mcp-server -config ./config.json

# Or specify API base URL directly
./bin/mcp-server -api-url http://localhost:8080

# Enable verbose logging
./bin/mcp-server -verbose
```

### 3. Configuration File

Example configuration file (`config.json`):

```json
{
  "api": {
    "baseUrl": "http://localhost:8080",
    "timeout": 30
  },
  "server": {
    "name": "mcp-bridge",
    "version": "1.0.0",
    "description": "REST API to MCP Bridge Server"
  },
  "headers": {
    "Content-Type": "application/json",
    "Authorization": "Bearer your-token-here"
  },
  "endpoints": [
    {
      "name": "custom_endpoint",
      "description": "Custom endpoint",
      "method": "GET",
      "path": "/api/custom",
      "parameters": [
        {
          "name": "param1",
          "type": "string",
          "required": true,
          "description": "Parameter 1",
          "in": "query"
        }
      ]
    }
  ]
}
```

### 4. Usage with Claude Code

Configuration example for Claude Code:

```json
{
  "mcpServers": {
    "rest-api-bridge": {
      "command": "/path/to/mcp-server",
      "args": ["-api-url", "http://localhost:8080"]
    }
  }
}
```

## Available Tools

Tools provided by the MCP bridge server:

### Default Tools (when using Mock API server)

- `get_users` - Get all users
- `create_user` - Create user
- `get_user` - Get specific user
- `update_user` - Update user
- `delete_user` - Delete user
- `get_products` - Get products
- `create_product` - Create product
- `health_check` - Health check

### Usage Examples

```javascript
// Get user list
await callTool("get_users", {});

// Create new user
await callTool("create_user", {
  name: "John Doe",
  email: "john@example.com"
});

// Get specific user
await callTool("get_user", {
  id: 1
});

// Filter products by category
await callTool("get_products", {
  category: "Electronics"
});
```

## Resources

The MCP server provides the following resources:

- `rest-api://docs` - REST API specification (JSON format)

## Development

### Running Tests

```bash
# Start Mock API server
go run ./cmd/mock-api &

# Test MCP server
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}' | go run ./cmd/mcp-server
```

### Adding Custom Endpoints

You can use custom API endpoints by adding them to the `endpoints` section in the configuration file.

## License

MIT License

## Contributing

Pull requests and issue reports are welcome.
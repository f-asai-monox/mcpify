# mcpify

A proxy server that enables REST APIs to be used as MCP (Model Context Protocol) servers.

## Overview

MCP Bridge is a proxy server that allows existing REST APIs to be used as Model Context Protocol (MCP) servers. This enables REST APIs to be directly utilized by MCP clients such as Claude Code.

## Project Significance

### Technical Value
- **Protocol Unification**: Converts existing REST APIs to MCP protocol, standardizing AI tool integration
- **Adapter Layer**: Provides unified interface for handling different API formats through bridge functionality
- **Type Safety**: JSON-RPC 2.0 compliant with schema-based type checking

### Practical Value
- **Leverage Existing Assets**: Directly utilize existing REST APIs with MCP clients (Claude Code, etc.) without creating new APIs
- **Development Efficiency**: Add new APIs through configuration files only, no individual implementation required
- **Authentication & Security**: Unified header management and error handling

### Ecosystem Contribution
- **MCP Adoption**: Contributes to MCP ecosystem expansion by making REST APIs MCP-compatible
- **Interoperability**: Provides standardized method for promoting integration between different services

## Features

- **REST API to MCP Conversion**: Automatically converts REST API endpoints to MCP tools
- **JSON-RPC 2.0 Compliant**: Fully compliant with the MCP protocol
- **Multiple Transport Support**: Supports both stdio and HTTP communication
- **Configurable**: Flexible customization through configuration files
- **Mock API Server**: Built-in simple REST API server for testing

## Project Structure

```
mcp-bridge/
├── cmd/
│   ├── mcp-server-stdio/  # MCP server with stdio transport
│   ├── mcp-server-http/   # MCP server with HTTP transport
│   └── mock-api/          # Configurable mock API server
├── internal/
│   ├── mcp/              # MCP implementation
│   ├── bridge/           # REST API conversion logic
│   ├── transport/        # Transport layer (stdio/HTTP)
│   └── config/           # Configuration management
├── pkg/
│   └── types/            # Common type definitions
├── go.mod
├── go.sum
├── README.md
└── README-ja.md          # Japanese version
```

## Installation & Build

### Dependencies
- Go 1.21 or higher

### Build

```bash
# Build MCP server with stdio transport
go build -o bin/mcp-server-stdio ./cmd/mcp-server-stdio

# Build MCP server with HTTP transport
go build -o bin/mcp-server-http ./cmd/mcp-server-http

# Build Mock API server
go build -o bin/mock-api ./cmd/mock-api
```

## Usage

### 1. Start Mock API Server

Start the configurable mock REST API server:

```bash
# Start with default configuration (users API)
./bin/mock-api

# Start with products configuration
MOCK_CONFIG=configs/mock/products.json ./bin/mock-api

# Or run directly
go run ./cmd/mock-api
```

For detailed Mock API documentation, configuration options, and usage examples, see **[Mock API Documentation](docs/MOCK-API.md)**.

### 2. Start MCP Server

Start the MCP bridge server:

#### Stdio Transport (Traditional)

```bash
./bin/mcp-server-stdio

# Or specify config file
./bin/mcp-server-stdio -config ./config.json

# Or specify API base URL directly
./bin/mcp-server-stdio -api-url http://localhost:8080

# Enable verbose logging
./bin/mcp-server-stdio -verbose
```

#### HTTP Transport (New)

```bash
./bin/mcp-server-http

# Specify HTTP settings
./bin/mcp-server-http -port 8080 -host localhost -cors

# Specify config file with HTTP settings
./bin/mcp-server-http -config ./example-config.json -port 8080

# Enable verbose logging
./bin/mcp-server-http -verbose
```

The HTTP server will accept MCP JSON-RPC requests at `http://localhost:8080/mcp`.

### 3. Configuration File

Example configuration file (`example-config.json`):

```json
{
  "apis": [
    {
      "name": "users-api",
      "baseUrl": "http://localhost:8081",
      "timeout": 30,
      "headers": {
        "X-API-Key": "your-api-key-here",
        "X-Custom-Header": "custom-value"
      },
      "auth": {
        "type": "basic",
        "basic": {
          "username": "admin",
          "password": "password"
        }
      },
      "endpoints": [
        {
          "name": "health",
          "description": "Health check endpoint",
          "method": "GET",
          "path": "/health",
          "parameters": []
        },
        {
          "name": "get_users",
          "description": "Get all users",
          "method": "GET",
          "path": "/users",
          "parameters": []
        },
        {
          "name": "create_user",
          "description": "Create a new user",
          "method": "POST",
          "path": "/users",
          "parameters": [
            {
              "name": "email",
              "type": "string",
              "required": true,
              "description": "User email address",
              "in": "body"
            },
            {
              "name": "name",
              "type": "string",
              "required": true,
              "description": "User name",
              "in": "body"
            }
          ]
        },
        {
          "name": "get_user",
          "description": "Get a specific user by ID",
          "method": "GET",
          "path": "/users/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "User ID",
              "in": "path"
            }
          ]
        },
        {
          "name": "update_user",
          "description": "Update a specific user by ID",
          "method": "PUT",
          "path": "/users/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "User ID",
              "in": "path"
            },
            {
              "name": "email",
              "type": "string",
              "required": true,
              "description": "User email address",
              "in": "body"
            },
            {
              "name": "name",
              "type": "string",
              "required": true,
              "description": "User name",
              "in": "body"
            }
          ]
        },
        {
          "name": "delete_user",
          "description": "Delete a specific user by ID",
          "method": "DELETE",
          "path": "/users/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "User ID",
              "in": "path"
            }
          ]
        }
      ]
    },
    {
      "name": "products-api",
      "baseUrl": "http://localhost:8082",
      "timeout": 30,
      "headers": {
        "Authorization": "Bearer your-token-here"
      },
      "endpoints": [
        {
          "name": "get_products",
          "description": "Get all products",
          "method": "GET",
          "path": "/products",
          "parameters": []
        },
        {
          "name": "get_product",
          "description": "Get a specific product by ID",
          "method": "GET",
          "path": "/products/{id}",
          "parameters": [
            {
              "name": "id",
              "type": "integer",
              "required": true,
              "description": "Product ID",
              "in": "path"
            }
          ]
        }
      ]
    }
  ],
  "server": {
    "name": "mcp-bridge",
    "version": "1.0.0",
    "description": "REST API to MCP Bridge Server"
  },
  "headers": {
    "Content-Type": "application/json"
  }
}
```

**Note**: The configuration file is the same for both stdio and HTTP transports. The transport type is determined by which server executable you use (`mcp-server-stdio` or `mcp-server-http`).

For HTTP transport, you can specify the server settings via command line arguments:

```bash
# HTTP transport with custom settings
go run ./cmd/mcp-server-http --config ./example-config.json --port 8080 --host localhost --cors
```

### 4. Usage with Claude Code

Configuration example for Claude Code:

#### Stdio Transport
```json
{
  "mcpServers": {
    "mcp-bridge": {
      "command": "go",
      "args": ["run", "./cmd/mcp-server-stdio", "--config", "./example-config.json"]
    }
  }
}
```

#### HTTP Transport
```json
{
  "mcpServers": {
    "mcp-bridge-http": {
      "transport": {
        "type": "http",
        "url": "http://localhost:8080/mcp"
      }
    }
  }
}
```

Note: You need to start the HTTP MCP server separately:
```bash
go run ./cmd/mcp-server-http --config ./example-config.json --port 8080
```

## Available Tools

Tools provided by the MCP bridge server:

### Default Tools (when using example configuration)

- `health` - Health check endpoint
- `get_users` - Get all users
- `create_user` - Create user
- `get_user` - Get specific user
- `update_user` - Update user
- `delete_user` - Delete user
- `get_products` - Get products
- `get_product` - Get specific product

### Usage Examples

```javascript
// Health check
await callTool("health", {});

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

// Get all products
await callTool("get_products", {});

// Get specific product
await callTool("get_product", {
  id: 1
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

# Test MCP server (stdio)
echo '{"jsonrpc": "2.0", "id": 1, "method": "initialize", "params": {"protocolVersion": "2024-11-05", "capabilities": {}, "clientInfo": {"name": "test", "version": "1.0.0"}}}' | go run ./cmd/mcp-server-stdio

# Test MCP server (HTTP)
go run ./cmd/mcp-server-http -port 8080 &
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc": "2.0", "id": 1, "method": "ping"}'
```

### Adding Custom Endpoints

You can use custom API endpoints by adding them to the `endpoints` section in the configuration file.

## License

MIT License

## Contributing

Pull requests and issue reports are welcome.

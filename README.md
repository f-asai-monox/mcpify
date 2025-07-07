# mcpify

A proxy server that enables REST APIs to be used as MCP (Model Context Protocol) servers.

## Features

- **REST API to MCP Conversion**: Automatically converts REST API endpoints to MCP tools
- **Multiple Transport Support**: Supports both stdio and HTTP communication
- **JSON-RPC 2.0 Compliant**: Fully compliant with the MCP protocol
- **Configurable**: Flexible customization through configuration files
- **Mock API Server**: Built-in simple REST API server for testing

## Quick Start

### 1. Install Dependencies
```bash
# Requires Go 1.24.2+
go version
```

### 2. Build the Server
```bash
# Build MCP server
go build -o bin/mcp-server-stdio ./cmd/mcp-server-stdio

# Build Mock API for testing
go build -o bin/mock-api ./cmd/mock-api
```

### 3. Start Mock API (for testing)
```bash
./bin/mock-api
```

### 4. Start MCP Server
```bash
# Basic usage
./bin/mcp-server-stdio

# With configuration file
./bin/mcp-server-stdio -config ./example-config.json

# With API URL
./bin/mcp-server-stdio -api-url http://localhost:8080
```

## Basic Usage

### Configuration Example
Create a `config.json` file:

```json
{
  "apis": [
    {
      "name": "users-api",
      "baseUrl": "http://localhost:8081",
      "endpoints": [
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
              "name": "name",
              "type": "string",
              "required": true,
              "description": "User name",
              "in": "body"
            },
            {
              "name": "email",
              "type": "string",
              "required": true,
              "description": "User email",
              "in": "body"
            }
          ]
        }
      ]
    }
  ]
}
```

### Usage with Claude Code
```json
{
  "mcpServers": {
    "mcp-bridge": {
      "command": "go",
      "args": ["run", "./cmd/mcp-server-stdio", "--config", "./config.json"]
    }
  }
}
```

## Available Tools

With the example configuration, you get these tools:
- `get_users` - Get all users
- `create_user` - Create a new user
- `get_user` - Get specific user by ID
- `update_user` - Update user information
- `delete_user` - Delete user

## HTTP Transport

For HTTP transport instead of stdio:

```bash
# Start HTTP server
go build -o bin/mcp-server-http ./cmd/mcp-server-http
./bin/mcp-server-http -port 8080

# Configure Claude Code
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

## Documentation

- **[Architecture](docs/ARCHITECTURE.md)** - Project structure and technical details
- **[Configuration](docs/CONFIGURATION.md)** - Complete configuration guide
- **[Development](docs/DEVELOPMENT.md)** - Development and testing guide
- **[API Reference](docs/API-REFERENCE.md)** - Available tools and usage examples

## License

MIT License

## Contributing

Pull requests and issue reports are welcome.
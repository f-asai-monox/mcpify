# mcpify

A proxy server that enables REST APIs to be used as MCP (Model Context Protocol) servers.

## Features

- **REST API to MCP Conversion**: Automatically converts REST API endpoints to MCP tools
- **Multiple Transport Support**: Supports both stdio and HTTP communication
- **JSON-RPC 2.0 Compliant**: Fully compliant with the MCP protocol
- **Configurable**: Flexible customization through configuration files
- **Mock API Server**: Built-in simple REST API server for testing

## Installation

### Option 1: Download Pre-built Binary (Recommended)

Download the latest release for your platform from the [Releases page](https://github.com/f-asai-monox/mcpify/releases).

#### Linux/macOS

Using curl:
```bash
curl -sSL https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.sh | bash
```

Using wget:
```bash
wget -qO- https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.sh | bash
```

#### Windows

Using PowerShell (recommended):
```powershell
iwr -useb https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.ps1 | iex
```

Using Command Prompt:
```cmd
curl -sSL https://raw.githubusercontent.com/f-asai-monox/mcpify/main/install.bat -o install.bat && install.bat
```

#### Manual Download
1. Go to [Releases](https://github.com/f-asai-monox/mcpify/releases)
2. Download the appropriate archive for your OS and architecture
3. Extract and move the binary to your PATH

### Option 2: Install with go install
```bash
go install github.com/f-asai-monox/mcpify/cmd/mcp-server-stdio@latest
```

### Option 3: Build from Source
```bash
# Clone the repository
git clone https://github.com/f-asai-monox/mcpify.git
cd mcpify

# Build MCP server
go build -o bin/mcp-server-stdio ./cmd/mcp-server-stdio

# Build Mock API for testing
go build -o bin/mock-api ./cmd/mock-api
```

## Quick Start

### 1. Start Mock API (for testing)
```bash
# If built from source
./bin/mock-api
```

### 2. Start MCP Server

#### Using stdio transport (for Claude Code):
```bash
# If installed via binary
mcp-server-stdio

# With configuration file
mcp-server-stdio -config ./example-config.json

# With API URL
mcp-server-stdio -api-url http://localhost:8080

# If built from source
./bin/mcp-server-stdio -config ./example-config.json
```

#### Using HTTP transport:
```bash
# If installed via binary
mcp-server-http -port 8080

# With configuration file
mcp-server-http -config ./example-config.json -port 8080

# If built from source
./bin/mcp-server-http -port 8080
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
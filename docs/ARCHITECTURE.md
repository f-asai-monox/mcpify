# Project Architecture

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

## Architecture Overview

The MCP Bridge acts as a proxy server that translates REST API calls into MCP protocol messages and vice versa. It supports both stdio and HTTP transports, making it compatible with various MCP clients.

### Key Components

1. **MCP Layer**: Handles MCP protocol compliance and JSON-RPC 2.0 communication
2. **Bridge Layer**: Converts REST API specifications into MCP tool definitions
3. **Transport Layer**: Manages communication via stdio or HTTP
4. **Configuration Layer**: Handles API endpoint definitions and authentication
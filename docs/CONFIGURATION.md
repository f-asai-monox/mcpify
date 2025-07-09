# Configuration Guide

## Configuration File Format

The MCP Bridge uses JSON configuration files to define API endpoints, authentication, and server settings. The same configuration file works for both stdio and HTTP transports.

## Complete Configuration Example

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
      "auth": [
        {
          "type": "basic",
          "basic": {
            "username": "admin",
            "password": "password"
          }
        }
      ],
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

## Configuration Options

### API Definition
- `name`: Unique identifier for the API
- `baseUrl`: Base URL for the REST API
- `timeout`: Request timeout in seconds
- `headers`: Default headers for all requests
- `auth`: Authentication configuration

### Authentication Types
- **Basic Auth**: Username/password authentication
- **Bearer Token**: Authorization header with token
- **API Key**: Custom header with API key

### Endpoint Parameters
- `in`: Parameter location (`path`, `query`, `body`, `header`)
- `type`: Parameter type (`string`, `integer`, `boolean`, `number`)
- `required`: Whether parameter is required
- `description`: Human-readable description

## Usage with Claude Code

### Stdio Transport
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

### HTTP Transport
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

## Adding Custom Endpoints

To add new endpoints, extend the `endpoints` array in your configuration:

```json
{
  "name": "custom_endpoint",
  "description": "Description of the endpoint",
  "method": "POST",
  "path": "/custom/{id}",
  "parameters": [
    {
      "name": "id",
      "type": "string",
      "required": true,
      "description": "Custom ID",
      "in": "path"
    },
    {
      "name": "data",
      "type": "string",
      "required": true,
      "description": "Custom data",
      "in": "body"
    }
  ]
}
```
# Mock API Server

A configurable mock REST API server for testing the MCP Bridge functionality.

## Overview

The mock API server provides a flexible testing environment with configurable endpoints, authentication, and data. It supports multiple configuration files for different API scenarios.

## Quick Start

```bash
# Build the mock API server
go build -o bin/mock-api ./cmd/mock-api

# Start with default configuration (users API)
./bin/mock-api

# Start with products configuration
MOCK_CONFIG=configs/mock/products.json ./bin/mock-api

# Start on specific port
PORT=8081 ./bin/mock-api

# Or run directly
go run ./cmd/mock-api
```

## Configuration

### Available Configurations

The mock API server uses configuration files in the `configs/mock/` directory:

- `configs/mock/users.json` - Users API (default, port 8080)
- `configs/mock/products.json` - Products API (port 8081)

### Configuration Structure

```json
{
  "server": {
    "port": "8080",
    "name": "Mock API Server"
  },
  "auth": {
    "enabled": false,
    "username": "admin",
    "password": "password"
  },
  "resources": [
    {
      "name": "users",
      "path": "/users",
      "enabled": true,
      "data": [...user objects...],
      "methods": ["GET", "POST", "PUT", "DELETE"],
      "supportsId": true
    }
  ],
  "endpoints": [
    {
      "path": "/health",
      "method": "GET",
      "enabled": true,
      "response": {"status": "healthy"}
    }
  ]
}
```

### Environment Variables

- `MOCK_CONFIG` - Path to configuration file (default: `configs/mock/users.json`)
- `PORT` - Server port (overrides config file)
- `AUTH_ENABLED` - Enable Basic Authentication (`true`/`false`)
- `AUTH_USERNAME` - Username for authentication (default: `admin`)
- `AUTH_PASSWORD` - Password for authentication (default: `password`)

## Basic Authentication

The mock API server supports Basic Authentication through environment variables:

```bash
# Start with Basic Authentication enabled
AUTH_ENABLED=true AUTH_USERNAME=admin AUTH_PASSWORD=secret PORT=8081 go run ./cmd/mock-api

# Start with custom credentials
AUTH_ENABLED=true AUTH_USERNAME=myuser AUTH_PASSWORD=mypass PORT=8081 go run ./cmd/mock-api

# Test authenticated endpoints
curl -u admin:secret http://localhost:8081/users

# Or with Authorization header
curl -H "Authorization: Basic YWRtaW46c2VjcmV0" http://localhost:8081/users
```

When Basic Authentication is enabled, all endpoints require valid credentials. Without authentication, requests will return `401 Unauthorized`.

## Available Endpoints

### With Users Configuration (default)

- `GET /health` - Health check
- `GET /users` - Get all users  
- `POST /users` - Create user
- `GET /users/{id}` - Get specific user
- `PUT /users/{id}` - Update user
- `DELETE /users/{id}` - Delete user

### With Products Configuration

- `GET /health` - Health check
- `GET /products` - Get all products
- `GET /products/{id}` - Get specific product

## Usage Examples

### Starting Different Services

```bash
# Users API (default)
./bin/mock-api

# Products API
MOCK_CONFIG=configs/mock/products.json ./bin/mock-api

# Users API with authentication
AUTH_ENABLED=true AUTH_USERNAME=admin AUTH_PASSWORD=secret ./bin/mock-api

# Products API on custom port
MOCK_CONFIG=configs/mock/products.json PORT=9000 ./bin/mock-api
```

### Testing Endpoints

```bash
# Health check
curl http://localhost:8080/health

# Get all users
curl http://localhost:8080/users

# Create a new user
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name": "John Doe", "email": "john@example.com"}'

# Get specific user
curl http://localhost:8080/users/1

# Update user
curl -X PUT http://localhost:8080/users/1 \
  -H "Content-Type: application/json" \
  -d '{"name": "John Smith", "email": "johnsmith@example.com"}'

# Delete user
curl -X DELETE http://localhost:8080/users/1

# Test with authentication
curl -u admin:secret http://localhost:8081/users
```

## Creating Custom Configurations

You can create custom configuration files in the `configs/mock/` directory:

```json
{
  "server": {
    "port": "8082",
    "name": "Custom Mock API"
  },
  "auth": {
    "enabled": true,
    "username": "custom",
    "password": "secret"
  },
  "resources": [
    {
      "name": "orders",
      "path": "/orders",
      "enabled": true,
      "data": [
        {
          "id": 1,
          "customerId": 1,
          "total": 99.99,
          "status": "pending"
        }
      ],
      "methods": ["GET", "POST"],
      "supportsId": true
    }
  ],
  "endpoints": [
    {
      "path": "/status",
      "method": "GET",
      "enabled": true,
      "response": {
        "service": "orders",
        "status": "running"
      }
    }
  ]
}
```

Then start with:
```bash
MOCK_CONFIG=configs/mock/custom.json ./bin/mock-api
```

## CORS Support

The mock API server automatically includes CORS headers for cross-origin requests:

- `Access-Control-Allow-Origin: *`
- `Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS`
- `Access-Control-Allow-Headers: Content-Type, Authorization`

## Features

- **Dynamic Configuration**: Load different API configurations via environment variables
- **Resource Management**: CRUD operations for configured resources
- **Custom Endpoints**: Define static response endpoints
- **Authentication**: Optional Basic Authentication
- **CORS Support**: Built-in cross-origin request support
- **Timestamp Templating**: Use `{{timestamp}}` in responses for dynamic timestamps
- **Flexible Data Types**: Support for various JSON data structures
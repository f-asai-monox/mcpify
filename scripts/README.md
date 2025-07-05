# Scripts

This directory contains utility scripts for the MCP Bridge project.

## start-mock-servers.sh

A script to start multiple mock API servers for testing the MCP Bridge functionality.

### Usage

```bash
# Make sure the script is executable
chmod +x scripts/start-mock-servers.sh

# Run the script from project root
./scripts/start-mock-servers.sh
```

### What it does

1. **Builds binaries** - Automatically builds `mock-api` and `mcp-server` if they don't exist
2. **Starts multiple servers**:
   - Users API on port 8081 (no authentication)
   - Products API on port 8082 (with Basic authentication: admin/password)
3. **Provides helpful information**:
   - Available endpoints for each server
   - Quick test commands
   - Instructions for starting MCP Bridge with example-config.json
4. **Graceful shutdown** - Press Ctrl+C to stop all servers

### Features

- **Color-coded output** for better readability
- **Automatic health checks** to verify servers started correctly
- **PID tracking** for proper cleanup
- **Signal handling** for graceful shutdown
- **Comprehensive usage instructions**

### Example Output

```
=== Mock API Servers Startup Script ===

Starting mock API servers...

Starting Users Mock API Server on port 8081...
✓ Users Mock API Server started successfully (PID: 12345)

Starting Products Mock API Server on port 8082...
✓ Products Mock API Server started successfully (PID: 12346)

=== All mock servers started successfully ===

Available endpoints:
Users API (port 8081):
  GET    http://localhost:8081/health
  GET    http://localhost:8081/users
  POST   http://localhost:8081/users
  GET    http://localhost:8081/users/{id}
  PUT    http://localhost:8081/users/{id}
  DELETE http://localhost:8081/users/{id}

Products API (port 8082):
  GET    http://localhost:8082/health
  GET    http://localhost:8082/products
  GET    http://localhost:8082/products/{id}
  Note: Products API requires authentication
  Authentication: admin:password

=== Quick Test Commands ===
Test Users API:
  curl http://localhost:8081/users
  curl http://localhost:8081/health

Test Products API (with auth):
  curl -u admin:password http://localhost:8082/products
  curl -u admin:password http://localhost:8082/health

=== MCP Bridge Usage ===
To start MCP Bridge server with example configuration:
  ./bin/mcp-server -config example-config.json

Or run directly:
  go run ./cmd/mcp-server -config example-config.json

The example-config.json includes:
  - Users API on http://localhost:8081 (with basic auth)
  - Products API on http://localhost:8082 (no auth in example config)

Press Ctrl+C to stop all servers
```
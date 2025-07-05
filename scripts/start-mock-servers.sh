#!/bin/bash

# Script to start multiple mock API servers
# This script starts both users and products mock servers

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}=== Mock API Servers Startup Script ===${NC}"
echo

# Check if mock-api binary exists
if [ ! -f "$PROJECT_DIR/bin/mock-api" ]; then
    echo -e "${YELLOW}Building mock-api binary...${NC}"
    cd "$PROJECT_DIR"
    go build -o bin/mock-api ./cmd/mock-api
    echo -e "${GREEN}✓ mock-api binary built${NC}"
    echo
fi

# Check if mcp-server binary exists
if [ ! -f "$PROJECT_DIR/bin/mcp-server" ]; then
    echo -e "${YELLOW}Building mcp-server binary...${NC}"
    cd "$PROJECT_DIR"
    go build -o bin/mcp-server ./cmd/mcp-server
    echo -e "${GREEN}✓ mcp-server binary built${NC}"
    echo
fi

# Function to start a mock server
start_mock_server() {
    local config_file="$1"
    local server_name="$2"
    local port="$3"
    
    echo -e "${YELLOW}Starting $server_name on port $port...${NC}"
    
    # Start the server in the background
    cd "$PROJECT_DIR"
    MOCK_CONFIG="$config_file" ./bin/mock-api &
    local pid=$!
    
    # Wait a moment for the server to start
    sleep 2
    
    # Check if the server is running
    if curl -s "http://localhost:$port/health" > /dev/null 2>&1; then
        echo -e "${GREEN}✓ $server_name started successfully (PID: $pid)${NC}"
        echo "$pid" >> /tmp/mock-servers.pid
        return 0
    else
        echo -e "${RED}✗ Failed to start $server_name${NC}"
        return 1
    fi
}

# Function to cleanup on exit
cleanup() {
    echo
    echo -e "${YELLOW}Stopping mock servers...${NC}"
    if [ -f /tmp/mock-servers.pid ]; then
        while read -r pid; do
            if kill -0 "$pid" 2>/dev/null; then
                kill "$pid" 2>/dev/null
                echo -e "${GREEN}✓ Stopped server (PID: $pid)${NC}"
            fi
        done < /tmp/mock-servers.pid
        rm -f /tmp/mock-servers.pid
    fi
    echo -e "${BLUE}=== Mock servers stopped ===${NC}"
}

# Set up signal handlers
trap cleanup EXIT INT TERM

# Remove existing PID file
rm -f /tmp/mock-servers.pid

# Start mock servers
echo -e "${BLUE}Starting mock API servers...${NC}"
echo

# Start Users API server (port 8081)
start_mock_server "configs/mock/users.json" "Users Mock API Server" "8081"

# Start Products API server (port 8082)  
start_mock_server "configs/mock/products.json" "Products Mock API Server" "8082"

echo
echo -e "${GREEN}=== All mock servers started successfully ===${NC}"
echo
echo -e "${BLUE}Available endpoints:${NC}"
echo -e "${YELLOW}Users API (port 8081):${NC}"
echo "  GET    http://localhost:8081/health"
echo "  GET    http://localhost:8081/users"
echo "  POST   http://localhost:8081/users"
echo "  GET    http://localhost:8081/users/{id}"
echo "  PUT    http://localhost:8081/users/{id}"
echo "  DELETE http://localhost:8081/users/{id}"
echo
echo -e "${YELLOW}Products API (port 8082):${NC}"
echo "  GET    http://localhost:8082/health"
echo "  GET    http://localhost:8082/products"
echo "  GET    http://localhost:8082/products/{id}"
echo -e "${RED}  Note: Products API requires authentication${NC}"
echo -e "${YELLOW}  Authentication: admin:password${NC}"
echo

echo -e "${BLUE}=== Quick Test Commands ===${NC}"
echo -e "${YELLOW}Test Users API:${NC}"
echo "  curl http://localhost:8081/users"
echo "  curl http://localhost:8081/health"
echo
echo -e "${YELLOW}Test Products API (with auth):${NC}"
echo "  curl -u admin:password http://localhost:8082/products"
echo "  curl -u admin:password http://localhost:8082/health"
echo

echo -e "${BLUE}=== MCP Bridge Usage ===${NC}"
echo -e "${YELLOW}To start MCP Bridge server with example configuration:${NC}"
echo -e "${GREEN}  ./bin/mcp-server -config example-config.json${NC}"
echo
echo -e "${YELLOW}Or run directly:${NC}"
echo -e "${GREEN}  go run ./cmd/mcp-server -config example-config.json${NC}"
echo
echo -e "${YELLOW}The example-config.json includes:${NC}"
echo "  - Users API on http://localhost:8081 (with basic auth)"
echo "  - Products API on http://localhost:8082 (no auth in example config)"
echo

echo -e "${BLUE}Press Ctrl+C to stop all servers${NC}"
echo

# Wait for interrupt signal
while true; do
    sleep 1
done
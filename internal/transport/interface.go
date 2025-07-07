package transport

import (
	"mcp-bridge/pkg/types"
)

// Transport interface defines the communication layer for MCP messages
type Transport interface {
	// Start begins listening for messages and returns when the transport is closed
	Start() error
	
	// ReadMessage reads a single JSON-RPC message from the transport
	ReadMessage() (*types.JSONRPCMessage, error)
	
	// WriteMessage writes a JSON-RPC message to the transport
	WriteMessage(msg *types.JSONRPCMessage) error
	
	// Close closes the transport
	Close() error
}

// Config represents transport-specific configuration
type Config interface {
	GetType() string
}

// StdioConfig represents configuration for stdio transport
type StdioConfig struct{}

func (c *StdioConfig) GetType() string {
	return "stdio"
}

// HTTPConfig represents configuration for HTTP transport
type HTTPConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	CORS bool   `json:"cors"`
}

func (c *HTTPConfig) GetType() string {
	return "http"
}
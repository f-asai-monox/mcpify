package transport_test

import (
	"testing"

	"mcp-bridge/internal/transport"
	"mcp-bridge/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestNewStdioTransport(t *testing.T) {
	stdioTransport := transport.NewStdioTransport()
	assert.NotNil(t, stdioTransport)
}

func TestNewHTTPTransport(t *testing.T) {
	config := &transport.HTTPConfig{
		Host: "localhost",
		Port: 8080,
		CORS: true,
	}
	httpTransport := transport.NewHTTPTransport(config)
	assert.NotNil(t, httpTransport)
}

func TestHTTPConfig_GetType(t *testing.T) {
	config := &transport.HTTPConfig{
		Host: "localhost",
		Port: 8080,
		CORS: true,
	}
	assert.Equal(t, "http", config.GetType())
}

func TestStdioConfig_GetType(t *testing.T) {
	config := &transport.StdioConfig{}
	assert.Equal(t, "stdio", config.GetType())
}

func TestStdioTransport_WriteMessage(t *testing.T) {
	stdioTransport := transport.NewStdioTransport()
	
	msg := &types.JSONRPCMessage{
		JSONRpc: "2.0",
		ID:      1,
		Method:  "test",
	}
	
	// Test that WriteMessage doesn't panic
	err := stdioTransport.WriteMessage(msg)
	assert.NoError(t, err)
}

func TestStdioTransport_Close(t *testing.T) {
	stdioTransport := transport.NewStdioTransport()
	
	err := stdioTransport.Close()
	assert.NoError(t, err)
}

func TestHTTPTransport_Start(t *testing.T) {
	config := &transport.HTTPConfig{
		Host: "localhost",
		Port: 0, // Use port 0 to let OS choose available port
		CORS: true,
	}
	httpTransport := transport.NewHTTPTransport(config)
	
	err := httpTransport.Start()
	assert.NoError(t, err)
	
	// Clean up
	httpTransport.Close()
}

func TestHTTPTransport_Close(t *testing.T) {
	config := &transport.HTTPConfig{
		Host: "localhost",
		Port: 0,
		CORS: true,
	}
	httpTransport := transport.NewHTTPTransport(config)
	
	// Start and then close
	httpTransport.Start()
	err := httpTransport.Close()
	assert.NoError(t, err)
}
package bridge_test

import (
	"testing"

	"mcp-bridge/internal/bridge"
	"mcp-bridge/internal/transport"

	"github.com/stretchr/testify/assert"
)

func TestNewMCPBridge(t *testing.T) {
	mockTransport := transport.NewStdioTransport()
	mcpBridge := bridge.NewMCPBridge(mockTransport)
	assert.NotNil(t, mcpBridge)
}

func TestMCPBridge_AddCustomEndpoint(t *testing.T) {
	mockTransport := transport.NewStdioTransport()
	mcpBridge := bridge.NewMCPBridge(mockTransport)
	endpoint := bridge.APIEndpoint{
		Name:        "test-endpoint",
		Description: "Test endpoint",
		Method:      "GET",
		Path:        "/test",
		Parameters:  []bridge.APIParameter{},
	}

	mcpBridge.AddCustomEndpoint(endpoint)
}

func TestMCPBridge_SetAPIHeader(t *testing.T) {
	mockTransport := transport.NewStdioTransport()
	mcpBridge := bridge.NewMCPBridge(mockTransport)
	mcpBridge.SetAPIHeader("Authorization", "Bearer token123")
}

func TestMCPBridge_FullWorkflow(t *testing.T) {
	mockTransport := transport.NewStdioTransport()
	mcpBridge := bridge.NewMCPBridge(mockTransport)

	endpoint := bridge.APIEndpoint{
		Name:        "test-workflow",
		Description: "Test workflow endpoint",
		Method:      "GET",
		Path:        "/test",
		APIName:     "test-api",
		Parameters: []bridge.APIParameter{
			{
				Name:        "param1",
				Type:        "string",
				Required:    true,
				Description: "Test parameter",
				In:          "query",
			},
		},
	}

	mcpBridge.AddCustomEndpoint(endpoint)
}

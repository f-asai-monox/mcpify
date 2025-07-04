package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"mcp-bridge/pkg/types"
)

func main() {
	log.SetOutput(os.Stderr)
	log.Println("MVP MCP Server starting...")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		log.Printf("Received: %s", line)
		if line == "" {
			continue
		}

		var msg types.JSONRPCMessage
		if err := json.Unmarshal([]byte(line), &msg); err != nil {
			log.Printf("Parse error: %v", err)
			sendError(nil, -32700, "Parse error", err)
			continue
		}

		log.Printf("Handling method: %s", msg.Method)
		switch msg.Method {
		case "initialize":
			handleInitialize(&msg)
		case "initialized":
			// Nothing to do for initialized notification
		case "tools/list":
			handleToolsList(&msg)
		case "resources/list":
			handleResourcesList(&msg)
		case "ping":
			handlePing(&msg)
		default:
			if msg.ID != nil {
				sendError(msg.ID, -32601, "Method not found", nil)
			}
		}
	}
}

func handleInitialize(msg *types.JSONRPCMessage) {
	log.Println("Handling initialize request")
	result := types.InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: types.ServerCapabilities{
			Tools: &types.ToolsCapability{
				ListChanged: true,
			},
			Resources: &types.ResourcesCapability{
				Subscribe:   false,
				ListChanged: true,
			},
			Logging: &types.LoggingCapability{},
		},
		ServerInfo: types.ServerInfo{
			Name:    "mvp-mcp-server",
			Version: "1.0.0",
		},
	}
	sendResult(msg.ID, result)
}

func handleToolsList(msg *types.JSONRPCMessage) {
	log.Println("Handling tools/list request")
	result := map[string]interface{}{
		"tools": []map[string]interface{}{
			{
				"name":        "echo",
				"description": "Echo back the input",
				"inputSchema": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"message": map[string]interface{}{
							"type":        "string",
							"description": "Message to echo back",
						},
					},
					"required": []string{"message"},
				},
			},
		},
	}
	sendResult(msg.ID, result)
}

func handleResourcesList(msg *types.JSONRPCMessage) {
	log.Println("Handling resources/list request")
	result := map[string]interface{}{
		"resources": []map[string]interface{}{},
	}
	sendResult(msg.ID, result)
}

func handlePing(msg *types.JSONRPCMessage) {
	log.Println("Handling ping request")
	sendResult(msg.ID, map[string]interface{}{})
}

func sendResult(id interface{}, result interface{}) {
	response := types.JSONRPCMessage{
		JSONRpc: "2.0",
		ID:      id,
		Result:  result,
	}
	sendMessage(response)
}

func sendError(id interface{}, code int, message string, data interface{}) {
	response := types.JSONRPCMessage{
		JSONRpc: "2.0",
		ID:      id,
		Error: &types.JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	sendMessage(response)
}

func sendMessage(msg types.JSONRPCMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	log.Printf("Sending: %s", string(data))
	fmt.Printf("%s\n", string(data))
}

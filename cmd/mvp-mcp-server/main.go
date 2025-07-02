package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

// JSON-RPC 2.0 message structure
type JSONRPCMessage struct {
	JSONRpc string      `json:"jsonrpc"`
	ID      interface{} `json:"id,omitempty"`
	Method  string      `json:"method,omitempty"`
	Params  interface{} `json:"params,omitempty"`
	Result  interface{} `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// MCP Initialize structures
type InitializeParams struct {
	ProtocolVersion string                 `json:"protocolVersion"`
	Capabilities    map[string]interface{} `json:"capabilities"`
	ClientInfo      ClientInfo             `json:"clientInfo"`
}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string           `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo       `json:"serverInfo"`
}

type ServerCapabilities struct {
	Tools     *ToolsCapability     `json:"tools,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Logging   *LoggingCapability   `json:"logging,omitempty"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe"`
	ListChanged bool `json:"listChanged"`
}

type LoggingCapability struct{}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

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

		var msg JSONRPCMessage
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

func handleInitialize(msg *JSONRPCMessage) {
	log.Println("Handling initialize request")
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: ServerCapabilities{
			Tools: &ToolsCapability{
				ListChanged: true,
			},
			Resources: &ResourcesCapability{
				Subscribe:   false,
				ListChanged: true,
			},
			Logging: &LoggingCapability{},
		},
		ServerInfo: ServerInfo{
			Name:    "mvp-mcp-server",
			Version: "1.0.0",
		},
	}
	sendResult(msg.ID, result)
}

func handleToolsList(msg *JSONRPCMessage) {
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

func handleResourcesList(msg *JSONRPCMessage) {
	log.Println("Handling resources/list request")
	result := map[string]interface{}{
		"resources": []map[string]interface{}{},
	}
	sendResult(msg.ID, result)
}

func handlePing(msg *JSONRPCMessage) {
	log.Println("Handling ping request")
	sendResult(msg.ID, map[string]interface{}{})
}

func sendResult(id interface{}, result interface{}) {
	response := JSONRPCMessage{
		JSONRpc: "2.0",
		ID:      id,
		Result:  result,
	}
	sendMessage(response)
}

func sendError(id interface{}, code int, message string, data interface{}) {
	response := JSONRPCMessage{
		JSONRpc: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	sendMessage(response)
}

func sendMessage(msg JSONRPCMessage) {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}
	log.Printf("Sending: %s", string(data))
	fmt.Printf("%s\n", string(data))
}
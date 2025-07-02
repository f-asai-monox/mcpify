package bridge

import (
	"encoding/json"
	"fmt"
	"strconv"

	"mcp-bridge/internal/mcp"
	"mcp-bridge/pkg/types"
)

type MCPBridge struct {
	server     *mcp.Server
	restClient *RestClient
	endpoints  []APIEndpoint
}

func NewMCPBridge(apiBaseURL string) *MCPBridge {
	restClient := NewRestClient(apiBaseURL)
	restClient.SetHeader("Content-Type", "application/json")
	
	bridge := &MCPBridge{
		server:     mcp.NewServer(),
		restClient: restClient,
		endpoints:  []APIEndpoint{}, // Initialize empty, will be populated via AddCustomEndpoint
	}
	
	bridge.setupMCPServer()
	return bridge
}

func (b *MCPBridge) setupMCPServer() {
	for _, endpoint := range b.endpoints {
		tool := b.createToolFromEndpoint(endpoint)
		b.server.AddTool(tool)
	}
	
	b.server.SetToolHandler(b.handleToolCall)
	b.server.SetResourceHandler(b.handleResourceRead)
	
	apiDocsResource := types.Resource{
		URI:         "rest-api://docs",
		Name:        "REST API Documentation",
		Description: "Documentation for available REST API endpoints",
		MimeType:    "application/json",
	}
	b.server.AddResource(apiDocsResource)
}

func (b *MCPBridge) createToolFromEndpoint(endpoint APIEndpoint) types.Tool {
	schema := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
		"required":   []string{},
	}
	
	properties := schema["properties"].(map[string]interface{})
	required := []string{}
	
	for _, param := range endpoint.Parameters {
		paramSchema := map[string]interface{}{
			"type":        b.convertParamType(param.Type),
			"description": param.Description,
		}
		
		if param.Default != nil {
			paramSchema["default"] = param.Default
		}
		
		properties[param.Name] = paramSchema
		
		if param.Required {
			required = append(required, param.Name)
		}
	}
	
	schema["required"] = required
	
	return types.Tool{
		Name:        endpoint.Name,
		Description: fmt.Sprintf("%s (%s %s)", endpoint.Description, endpoint.Method, endpoint.Path),
		InputSchema: schema,
	}
}

func (b *MCPBridge) convertParamType(paramType string) string {
	switch paramType {
	case "integer", "int":
		return "number"
	case "float", "double":
		return "number"
	case "bool", "boolean":
		return "boolean"
	default:
		return "string"
	}
}

func (b *MCPBridge) handleToolCall(name string, args map[string]interface{}) (*types.CallToolResult, error) {
	var endpoint *APIEndpoint
	for _, ep := range b.endpoints {
		if ep.Name == name {
			endpoint = &ep
			break
		}
	}
	
	if endpoint == nil {
		return &types.CallToolResult{
			Content: []types.ToolResult{
				{
					Type: "text",
					Text: fmt.Sprintf("Unknown tool: %s", name),
				},
			},
			IsError: true,
		}, nil
	}
	
	processedArgs := b.processArguments(args, endpoint.Parameters)
	
	response, err := b.restClient.MakeRequest(*endpoint, processedArgs)
	if err != nil {
		return &types.CallToolResult{
			Content: []types.ToolResult{
				{
					Type: "text",
					Text: fmt.Sprintf("Error calling API: %v", err),
				},
			},
			IsError: true,
		}, nil
	}
	
	return b.formatAPIResponse(response), nil
}

func (b *MCPBridge) processArguments(args map[string]interface{}, params []APIParameter) map[string]interface{} {
	processed := make(map[string]interface{})
	
	for key, value := range args {
		var paramType string
		for _, param := range params {
			if param.Name == key {
				paramType = param.Type
				break
			}
		}
		
		switch paramType {
		case "integer", "int":
			if str, ok := value.(string); ok {
				if intVal, err := strconv.Atoi(str); err == nil {
					processed[key] = intVal
				} else {
					processed[key] = value
				}
			} else {
				processed[key] = value
			}
		case "float", "double":
			if str, ok := value.(string); ok {
				if floatVal, err := strconv.ParseFloat(str, 64); err == nil {
					processed[key] = floatVal
				} else {
					processed[key] = value
				}
			} else {
				processed[key] = value
			}
		case "bool", "boolean":
			if str, ok := value.(string); ok {
				if boolVal, err := strconv.ParseBool(str); err == nil {
					processed[key] = boolVal
				} else {
					processed[key] = value
				}
			} else {
				processed[key] = value
			}
		default:
			processed[key] = value
		}
	}
	
	return processed
}

func (b *MCPBridge) formatAPIResponse(response *APIResponse) *types.CallToolResult {
	if response.Error != "" {
		return &types.CallToolResult{
			Content: []types.ToolResult{
				{
					Type: "text",
					Text: fmt.Sprintf("API Error: %s", response.Error),
				},
			},
			IsError: true,
		}
	}
	
	var resultText string
	if response.Data != nil {
		if jsonData, err := json.MarshalIndent(response.Data, "", "  "); err == nil {
			resultText = fmt.Sprintf("Status: %d\n\nResponse:\n%s", response.StatusCode, string(jsonData))
		} else {
			resultText = fmt.Sprintf("Status: %d\n\nResponse:\n%s", response.StatusCode, response.Body)
		}
	} else {
		resultText = fmt.Sprintf("Status: %d\n\nResponse:\n%s", response.StatusCode, response.Body)
	}
	
	return &types.CallToolResult{
		Content: []types.ToolResult{
			{
				Type: "text",
				Text: resultText,
			},
		},
		IsError: false,
	}
}

func (b *MCPBridge) handleResourceRead(uri string) (*types.ReadResourceResult, error) {
	switch uri {
	case "rest-api://docs":
		docsData := map[string]interface{}{
			"endpoints": b.endpoints,
			"baseURL":   b.restClient.baseURL,
		}
		
		jsonData, err := json.MarshalIndent(docsData, "", "  ")
		if err != nil {
			return nil, fmt.Errorf("error marshaling docs: %w", err)
		}
		
		return &types.ReadResourceResult{
			Contents: []types.ResourceContent{
				{
					URI:      uri,
					MimeType: "application/json",
					Text:     string(jsonData),
				},
			},
		}, nil
	default:
		return &types.ReadResourceResult{
			Contents: []types.ResourceContent{
				{
					URI:      uri,
					MimeType: "text/plain",
					Text:     fmt.Sprintf("Resource not found: %s", uri),
				},
			},
		}, nil
	}
}

func (b *MCPBridge) Start() error {
	return b.server.Start()
}

func (b *MCPBridge) SetAPIBaseURL(baseURL string) {
	b.restClient = NewRestClient(baseURL)
	b.restClient.SetHeader("Content-Type", "application/json")
}

func (b *MCPBridge) SetAPIHeader(key, value string) {
	b.restClient.SetHeader(key, value)
}

func (b *MCPBridge) AddCustomEndpoint(endpoint APIEndpoint) {
	b.endpoints = append(b.endpoints, endpoint)
	tool := b.createToolFromEndpoint(endpoint)
	b.server.AddTool(tool)
}
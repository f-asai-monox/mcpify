package mcp

import (
	"encoding/json"
	"fmt"
	"io"
	"log"

	"mcp-bridge/internal/transport"
	"mcp-bridge/pkg/types"
)

type Server struct {
	capabilities    types.ServerCapabilities
	tools           []types.Tool
	resources       []types.Resource
	prompts         []types.Prompt
	transport       transport.Transport
	initialized     bool
	toolHandler     func(string, map[string]interface{}) (*types.CallToolResult, error)
	resourceHandler func(string) (*types.ReadResourceResult, error)
	promptHandler   func(string, map[string]interface{}) (*types.GetPromptResult, error)
}

func NewServer(t transport.Transport) *Server {
	return &Server{
		capabilities: types.ServerCapabilities{
			Tools: &types.ToolsCapability{
				ListChanged: true,
			},
			Resources: &types.ResourcesCapability{
				Subscribe:   false,
				ListChanged: true,
			},
			Logging: &types.LoggingCapability{},
		},
		tools:     []types.Tool{},
		resources: []types.Resource{},
		prompts:   []types.Prompt{},
		transport: t,
	}
}

func (s *Server) Start() error {
	if err := s.transport.Start(); err != nil {
		return fmt.Errorf("error starting transport: %w", err)
	}

	for {
		msg, err := s.transport.ReadMessage()
		if err != nil {
			if err == io.EOF {
				break
			}
			s.sendError(nil, -32700, "Parse error", err)
			continue
		}

		if msg == nil {
			continue // No message available (for HTTP transport polling)
		}

		if err := s.handleMessage(msg); err != nil {
			log.Printf("Error handling message: %v", err)
		}
	}

	return nil
}

func (s *Server) handleMessage(msg *types.JSONRPCMessage) error {
	switch msg.Method {
	case "initialize":
		return s.handleInitialize(msg)
	case "initialized":
		return s.handleInitialized(msg)
	case "tools/list":
		return s.handleToolsList(msg)
	case "tools/call":
		return s.handleToolsCall(msg)
	case "resources/list":
		return s.handleResourcesList(msg)
	case "resources/read":
		return s.handleResourcesRead(msg)
	case "prompts/list":
		return s.handlePromptsList(msg)
	case "prompts/get":
		return s.handlePromptsGet(msg)
	case "ping":
		return s.handlePing(msg)
	default:
		if msg.ID != nil {
			s.sendError(msg.ID, -32601, "Method not found", nil)
		}
		return nil
	}
}

func (s *Server) handleInitialize(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("initialize request must have an ID")
	}

	var params types.InitializeParams
	if msg.Params != nil {
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
	}

	result := types.InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities:    s.capabilities,
		ServerInfo: types.ServerInfo{
			Name:    "mcp-bridge",
			Version: "1.0.0",
		},
	}

	return s.sendResult(msg.ID, result)
}

func (s *Server) handleInitialized(_ *types.JSONRPCMessage) error {
	s.initialized = true
	return nil
}

func (s *Server) handleToolsList(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("tools/list request must have an ID")
	}

	result := types.ToolsListResult{
		Tools: s.tools,
	}

	return s.sendResult(msg.ID, result)
}

func (s *Server) handleToolsCall(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("tools/call request must have an ID")
	}

	var params types.CallToolParams
	if msg.Params != nil {
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
	}

	result, err := s.callTool(params.Name, params.Arguments)
	if err != nil {
		s.sendError(msg.ID, -32603, "Internal error", err)
		return nil
	}

	return s.sendResult(msg.ID, result)
}

func (s *Server) handleResourcesList(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("resources/list request must have an ID")
	}

	result := types.ResourcesListResult{
		Resources: s.resources,
	}

	return s.sendResult(msg.ID, result)
}

func (s *Server) handleResourcesRead(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("resources/read request must have an ID")
	}

	var params types.ReadResourceParams
	if msg.Params != nil {
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
	}

	result, err := s.readResource(params.URI)
	if err != nil {
		s.sendError(msg.ID, -32603, "Internal error", err)
		return nil
	}

	return s.sendResult(msg.ID, result)
}

func (s *Server) handlePromptsList(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("prompts/list request must have an ID")
	}

	result := types.PromptsListResult{
		Prompts: s.prompts,
	}

	return s.sendResult(msg.ID, result)
}

func (s *Server) handlePromptsGet(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("prompts/get request must have an ID")
	}

	var params types.GetPromptParams
	if msg.Params != nil {
		paramsBytes, err := json.Marshal(msg.Params)
		if err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
		if err := json.Unmarshal(paramsBytes, &params); err != nil {
			s.sendError(msg.ID, -32602, "Invalid params", err)
			return nil
		}
	}

	result, err := s.getPrompt(params.Name, params.Arguments)
	if err != nil {
		s.sendError(msg.ID, -32603, "Internal error", err)
		return nil
	}

	return s.sendResult(msg.ID, result)
}

func (s *Server) handlePing(msg *types.JSONRPCMessage) error {
	if msg.ID == nil {
		return fmt.Errorf("ping request must have an ID")
	}

	return s.sendResult(msg.ID, map[string]interface{}{})
}

func (s *Server) sendResult(id interface{}, result interface{}) error {
	response := types.JSONRPCMessage{
		JSONRpc: "2.0",
		ID:      id,
		Result:  result,
	}
	return s.sendMessage(response)
}

func (s *Server) sendError(id interface{}, code int, message string, data interface{}) error {
	response := types.JSONRPCMessage{
		JSONRpc: "2.0",
		ID:      id,
		Error: &types.JSONRPCError{
			Code:    code,
			Message: message,
			Data:    data,
		},
	}
	return s.sendMessage(response)
}

func (s *Server) sendMessage(msg types.JSONRPCMessage) error {
	return s.transport.WriteMessage(&msg)
}

func (s *Server) AddTool(tool types.Tool) {
	s.tools = append(s.tools, tool)
}

func (s *Server) AddResource(resource types.Resource) {
	s.resources = append(s.resources, resource)
}

func (s *Server) AddPrompt(prompt types.Prompt) {
	s.prompts = append(s.prompts, prompt)
}

func (s *Server) SetToolHandler(handler func(name string, args map[string]interface{}) (*types.CallToolResult, error)) {
	s.toolHandler = handler
}

func (s *Server) SetResourceHandler(handler func(uri string) (*types.ReadResourceResult, error)) {
	s.resourceHandler = handler
}

func (s *Server) SetPromptHandler(handler func(name string, args map[string]interface{}) (*types.GetPromptResult, error)) {
	s.promptHandler = handler
}

func (s *Server) callTool(name string, args map[string]interface{}) (*types.CallToolResult, error) {
	if s.toolHandler != nil {
		return s.toolHandler(name, args)
	}

	return &types.CallToolResult{
		Content: []types.ToolResult{
			{
				Type: "text",
				Text: fmt.Sprintf("Tool '%s' not implemented", name),
			},
		},
		IsError: true,
	}, nil
}

func (s *Server) readResource(uri string) (*types.ReadResourceResult, error) {
	if s.resourceHandler != nil {
		return s.resourceHandler(uri)
	}

	return &types.ReadResourceResult{
		Contents: []types.ResourceContent{
			{
				URI:      uri,
				MimeType: "text/plain",
				Text:     fmt.Sprintf("Resource '%s' not found", uri),
			},
		},
	}, nil
}

func (s *Server) getPrompt(name string, args map[string]interface{}) (*types.GetPromptResult, error) {
	if s.promptHandler != nil {
		return s.promptHandler(name, args)
	}

	return &types.GetPromptResult{
		Description: fmt.Sprintf("Prompt '%s' not found", name),
		Messages: []types.PromptMessage{
			{
				Role:    "user",
				Content: fmt.Sprintf("Prompt '%s' not implemented", name),
			},
		},
	}, nil
}

package types

type JSONRPCMessage struct {
	JSONRpc string        `json:"jsonrpc"`
	ID      interface{}   `json:"id,omitempty"`
	Method  string        `json:"method,omitempty"`
	Params  interface{}   `json:"params,omitempty"`
	Result  interface{}   `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

type JSONRPCError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type InitializeParams struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ClientCapabilities `json:"capabilities"`
	ClientInfo      ClientInfo         `json:"clientInfo"`
}

type ClientCapabilities struct {
	Roots    *RootsCapability    `json:"roots,omitempty"`
	Sampling *SamplingCapability `json:"sampling,omitempty"`
}

type RootsCapability struct {
	ListChanged bool `json:"listChanged"`
}

type SamplingCapability struct{}

type ClientInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type InitializeResult struct {
	ProtocolVersion string             `json:"protocolVersion"`
	Capabilities    ServerCapabilities `json:"capabilities"`
	ServerInfo      ServerInfo         `json:"serverInfo"`
}

type ServerCapabilities struct {
	Logging   *LoggingCapability   `json:"logging,omitempty"`
	Prompts   *PromptsCapability   `json:"prompts,omitempty"`
	Resources *ResourcesCapability `json:"resources,omitempty"`
	Tools     *ToolsCapability     `json:"tools,omitempty"`
}

type LoggingCapability struct{}

type PromptsCapability struct {
	ListChanged bool `json:"listChanged"`
}

type ResourcesCapability struct {
	Subscribe   bool `json:"subscribe"`
	ListChanged bool `json:"listChanged"`
}

type ToolsCapability struct {
	ListChanged bool `json:"listChanged"`
}

type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Tool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

type ToolsListResult struct {
	Tools []Tool `json:"tools"`
}

type CallToolParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments"`
}

type CallToolResult struct {
	Content []ToolResult `json:"content"`
	IsError bool         `json:"isError,omitempty"`
}

type ToolResult struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type Resource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

type ResourcesListResult struct {
	Resources []Resource `json:"resources"`
}

type ReadResourceParams struct {
	URI string `json:"uri"`
}

type ReadResourceResult struct {
	Contents []ResourceContent `json:"contents"`
}

type ResourceContent struct {
	URI      string `json:"uri"`
	MimeType string `json:"mimeType"`
	Text     string `json:"text,omitempty"`
	Blob     string `json:"blob,omitempty"`
}

type Prompt struct {
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Arguments   []Argument `json:"arguments,omitempty"`
}

type Argument struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required,omitempty"`
}

type PromptsListResult struct {
	Prompts []Prompt `json:"prompts"`
}

type GetPromptParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

type GetPromptResult struct {
	Description string          `json:"description,omitempty"`
	Messages    []PromptMessage `json:"messages"`
}

type PromptMessage struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type LoggingLevel string

const (
	LoggingLevelDebug     LoggingLevel = "debug"
	LoggingLevelInfo      LoggingLevel = "info"
	LoggingLevelNotice    LoggingLevel = "notice"
	LoggingLevelWarning   LoggingLevel = "warning"
	LoggingLevelError     LoggingLevel = "error"
	LoggingLevelCritical  LoggingLevel = "critical"
	LoggingLevelAlert     LoggingLevel = "alert"
	LoggingLevelEmergency LoggingLevel = "emergency"
)

type LoggingParams struct {
	Level  LoggingLevel `json:"level"`
	Data   interface{}  `json:"data,omitempty"`
	Logger string       `json:"logger,omitempty"`
}

type NotificationParams struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
}

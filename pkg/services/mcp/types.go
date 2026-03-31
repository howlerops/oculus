package mcp

// ServerConfig holds MCP server configuration from settings.json
type ServerConfig struct {
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	Transport string            `json:"transport,omitempty"` // "stdio" or "http"
	URL       string            `json:"url,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

// ServerConnection represents a connected MCP server
type ServerConnection struct {
	Name      string
	Config    ServerConfig
	Status    ConnectionStatus
	Tools     []ServerTool
	Resources []ServerResource
}

// ConnectionStatus tracks server lifecycle
type ConnectionStatus string

const (
	StatusDisconnected ConnectionStatus = "disconnected"
	StatusConnecting   ConnectionStatus = "connecting"
	StatusConnected    ConnectionStatus = "connected"
	StatusError        ConnectionStatus = "error"
)

// ServerTool is a tool exposed by an MCP server
type ServerTool struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"inputSchema"`
}

// ServerResource is a resource exposed by an MCP server
type ServerResource struct {
	URI         string `json:"uri"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	MimeType    string `json:"mimeType,omitempty"`
}

// ToolCallRequest is sent to an MCP server to invoke a tool
type ToolCallRequest struct {
	Method string         `json:"method"` // "tools/call"
	Params ToolCallParams `json:"params"`
}

// ToolCallParams are the parameters for a tool call
type ToolCallParams struct {
	Name      string                 `json:"name"`
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}

// ToolCallResponse is received from an MCP server after a tool call
type ToolCallResponse struct {
	Content []ToolCallContent `json:"content"`
	IsError bool              `json:"isError,omitempty"`
}

// ToolCallContent is a piece of content in a tool response
type ToolCallContent struct {
	Type string `json:"type"` // "text", "image", "resource"
	Text string `json:"text,omitempty"`
}

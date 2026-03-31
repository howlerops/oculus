package bridge

// BridgeMessageType identifies protocol message types
type BridgeMessageType string

const (
	MsgTypeInit          BridgeMessageType = "init"
	MsgTypeQuery         BridgeMessageType = "query"
	MsgTypeResponse      BridgeMessageType = "response"
	MsgTypeToolUse       BridgeMessageType = "tool_use"
	MsgTypeToolResult    BridgeMessageType = "tool_result"
	MsgTypePermission    BridgeMessageType = "permission_request"
	MsgTypePermissionRes BridgeMessageType = "permission_response"
	MsgTypeStatus        BridgeMessageType = "status"
	MsgTypePing          BridgeMessageType = "ping"
	MsgTypePong          BridgeMessageType = "pong"
	MsgTypeError         BridgeMessageType = "error"
	MsgTypeShutdown      BridgeMessageType = "shutdown"
	MsgTypeAttachment    BridgeMessageType = "attachment"
	MsgTypeStream        BridgeMessageType = "stream"
)

// BridgeConfig holds connection settings
type BridgeConfig struct {
	Enabled      bool   `json:"enabled"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	SessionID    string `json:"session_id"`
	AuthToken    string `json:"auth_token"`
	PollInterval int    `json:"poll_interval_ms"`
	Transport    string `json:"transport"` // "websocket", "sse", "polling"
}

// ProtocolMessage is the wire format for bridge messages
type ProtocolMessage struct {
	ID        string            `json:"id"`
	Type      BridgeMessageType `json:"type"`
	SessionID string            `json:"session_id"`
	Timestamp int64             `json:"timestamp"`
	Payload   interface{}       `json:"payload"`
}

// InitPayload for session initialization
type InitPayload struct {
	Version   string `json:"version"`
	CWD       string `json:"cwd"`
	Model     string `json:"model"`
	GitBranch string `json:"git_branch,omitempty"`
}

// QueryPayload for incoming queries
type QueryPayload struct {
	Text        string   `json:"text"`
	Attachments []string `json:"attachments,omitempty"`
}

// StreamPayload for streaming text chunks
type StreamPayload struct {
	Text       string `json:"text"`
	IsComplete bool   `json:"is_complete"`
}

// PermissionPayload for tool permission requests
type PermissionPayload struct {
	ToolName    string                 `json:"tool_name"`
	Input       map[string]interface{} `json:"input"`
	Description string                 `json:"description"`
	IsReadOnly  bool                   `json:"is_read_only"`
}

// StatusPayload for status updates
type StatusPayload struct {
	State   string `json:"state"` // "idle", "thinking", "tool_use", "streaming"
	Message string `json:"message,omitempty"`
}

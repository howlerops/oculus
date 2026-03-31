package entrypoints

// SDKStatus tracks the SDK session state
type SDKStatus struct {
	Connected    bool   `json:"connected"`
	SessionID    string `json:"session_id,omitempty"`
	Streaming    bool   `json:"streaming"`
	LastActivity int64  `json:"last_activity,omitempty"`
}

// EntrypointMode identifies how the CLI was invoked
type EntrypointMode string

const (
	EntrypointREPL     EntrypointMode = "repl"
	EntrypointPrint    EntrypointMode = "print"
	EntrypointSDK      EntrypointMode = "sdk"
	EntrypointHeadless EntrypointMode = "headless"
	EntrypointMCP      EntrypointMode = "mcp"
)

// InitOptions configures CLI initialization
type InitOptions struct {
	Mode               EntrypointMode
	Model              string
	PermissionMode     string
	Verbose            bool
	Debug              bool
	CustomSystemPrompt string
	AppendSystemPrompt string
	AllowedTools       []string
	DisallowedTools    []string
	MaxBudgetUSD       float64
}

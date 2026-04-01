package bridge

import "context"

type Message struct {
	Role    string      `json:"role"`
	Content interface{} `json:"content"`
}

type StreamChunk struct {
	Type     string `json:"type"`
	Text     string `json:"text,omitempty"`
	ToolID   string `json:"tool_id,omitempty"`
	ToolName string `json:"tool_name,omitempty"`
	Input    string `json:"input,omitempty"`
	Error    string `json:"error,omitempty"`
}

type BridgeConfig struct {
	Provider     string            `json:"provider"`
	Model        string            `json:"model"`
	APIKey       string            `json:"api_key,omitempty"`
	BaseURL      string            `json:"base_url,omitempty"`
	ExtraHeaders map[string]string `json:"extra_headers,omitempty"`
}

type Bridge interface {
	Execute(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef) (*Response, error)
	Stream(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef, handler func(StreamChunk)) error
	Name() string
	IsAvailable() bool
}

type ToolDef struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	InputSchema interface{} `json:"input_schema"`
}

type Response struct {
	Content    string     `json:"content"`
	StopReason string     `json:"stop_reason"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	Usage      Usage      `json:"usage"`
	Model      string     `json:"model"`
}

type ToolCall struct {
	ID    string                 `json:"id"`
	Name  string                 `json:"name"`
	Input map[string]interface{} `json:"input"`
}

type Usage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

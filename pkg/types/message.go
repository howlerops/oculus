package types

import "time"

// ContentBlockType identifies the kind of content block
type ContentBlockType string

const (
	ContentBlockText             ContentBlockType = "text"
	ContentBlockToolUse          ContentBlockType = "tool_use"
	ContentBlockToolResult       ContentBlockType = "tool_result"
	ContentBlockThinking         ContentBlockType = "thinking"
	ContentBlockRedactedThinking ContentBlockType = "redacted_thinking"
	ContentBlockImage            ContentBlockType = "image"
)

// ContentBlock represents a single block within a message
type ContentBlock struct {
	Type      ContentBlockType       `json:"type"`
	Text      string                 `json:"text,omitempty"`
	ID        string                 `json:"id,omitempty"`           // for tool_use
	Name      string                 `json:"name,omitempty"`         // for tool_use
	Input     map[string]interface{} `json:"input,omitempty"`        // for tool_use
	ToolUseID string                 `json:"tool_use_id,omitempty"`  // for tool_result
	Content   interface{}            `json:"content,omitempty"`      // for tool_result (string or []ContentBlock)
	IsError   bool                   `json:"is_error,omitempty"`     // for tool_result
	Thinking  string                 `json:"thinking,omitempty"`     // for thinking blocks
	Data      string                 `json:"data,omitempty"`         // for redacted_thinking
}

// MessageRole identifies the sender
type MessageRole string

const (
	RoleUser      MessageRole = "user"
	RoleAssistant MessageRole = "assistant"
	RoleSystem    MessageRole = "system"
)

// StopReason from the API
type StopReason string

const (
	StopReasonEndTurn      StopReason = "end_turn"
	StopReasonToolUse      StopReason = "tool_use"
	StopReasonMaxTokens    StopReason = "max_tokens"
	StopReasonStopSequence StopReason = "stop_sequence"
)

// Usage tracks token consumption
type Usage struct {
	InputTokens              int `json:"input_tokens"`
	OutputTokens             int `json:"output_tokens"`
	CacheReadInputTokens     int `json:"cache_read_input_tokens,omitempty"`
	CacheCreationInputTokens int `json:"cache_creation_input_tokens,omitempty"`
}

// SystemMessageLevel for system message severity
type SystemMessageLevel string

const (
	SystemMessageLevelInfo    SystemMessageLevel = "info"
	SystemMessageLevelWarning SystemMessageLevel = "warning"
	SystemMessageLevelError   SystemMessageLevel = "error"
)

// SystemMessageType discriminates system message subtypes
type SystemMessageType string

const (
	SystemMsgInformational        SystemMessageType = "informational"
	SystemMsgAPIError             SystemMessageType = "api_error"
	SystemMsgCompactBoundary      SystemMessageType = "compact_boundary"
	SystemMsgMicrocompactBoundary SystemMessageType = "microcompact_boundary"
	SystemMsgLocalCommand         SystemMessageType = "local_command"
	SystemMsgMemorySaved          SystemMessageType = "memory_saved"
	SystemMsgBridgeStatus         SystemMessageType = "bridge_status"
	SystemMsgTurnDuration         SystemMessageType = "turn_duration"
	SystemMsgAgentsKilled         SystemMessageType = "agents_killed"
	SystemMsgPermissionRetry      SystemMessageType = "permission_retry"
	SystemMsgStopHookSummary      SystemMessageType = "stop_hook_summary"
	SystemMsgAwaySummary          SystemMessageType = "away_summary"
	SystemMsgAPIMetrics           SystemMessageType = "api_metrics"
	SystemMsgScheduledTaskFire    SystemMessageType = "scheduled_task_fire"
)

// UserMessage from the user
type UserMessage struct {
	Role             MessageRole    `json:"role"`
	Content          []ContentBlock `json:"content"`
	Timestamp        time.Time      `json:"timestamp"`
	IsInputPipelined bool           `json:"isInputPipelined,omitempty"`
}

// AssistantMessage from the model
type AssistantMessage struct {
	Role          MessageRole            `json:"role"`
	Content       []ContentBlock         `json:"content"`
	Model         string                 `json:"model,omitempty"`
	StopReason    StopReason             `json:"stop_reason,omitempty"`
	Usage         *Usage                 `json:"usage,omitempty"`
	Timestamp     time.Time              `json:"timestamp"`
	ToolUseResult map[string]interface{} `json:"toolUseResult,omitempty"`
}

// SystemMessage for system-level events
type SystemMessage struct {
	Role      MessageRole            `json:"role"`
	Type      SystemMessageType      `json:"type"`
	Level     SystemMessageLevel     `json:"level,omitempty"`
	Text      string                 `json:"text,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// AttachmentMessage for CLAUDE.md and memory injections
type AttachmentMessage struct {
	Role      MessageRole `json:"role"`
	Content   string      `json:"content"`
	FilePath  string      `json:"filePath,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

// ToolUseSummaryMessage for compacted tool use display
type ToolUseSummaryMessage struct {
	Role      MessageRole `json:"role"`
	Summary   string      `json:"summary"`
	ToolName  string      `json:"toolName"`
	ToolUseID string      `json:"toolUseId"`
	Timestamp time.Time   `json:"timestamp"`
}

// TombstoneMessage placeholder for removed/compacted messages
type TombstoneMessage struct {
	Role          MessageRole `json:"role"`
	OriginalRole  MessageRole `json:"originalRole"`
	RemovedTokens int         `json:"removedTokens,omitempty"`
	Timestamp     time.Time   `json:"timestamp"`
}

// ProgressMessage for real-time tool progress updates
type ProgressMessage struct {
	ToolUseID string      `json:"toolUseId"`
	Data      interface{} `json:"data"` // ToolProgressData or HookProgress
}

// Message is the union type for all message kinds
type Message struct {
	// Discriminator
	Kind string `json:"kind"` // "user", "assistant", "system", "attachment", "tool_use_summary", "tombstone"

	User           *UserMessage           `json:"user,omitempty"`
	Assistant      *AssistantMessage      `json:"assistant,omitempty"`
	System         *SystemMessage         `json:"system,omitempty"`
	Attachment     *AttachmentMessage     `json:"attachment,omitempty"`
	ToolUseSummary *ToolUseSummaryMessage `json:"toolUseSummary,omitempty"`
	Tombstone      *TombstoneMessage      `json:"tombstone,omitempty"`
}

// MessageOrigin tracks where a message came from
type MessageOrigin struct {
	Source  string `json:"source"` // "user", "api", "system", "tool"
	AgentID string `json:"agentId,omitempty"`
}

// StreamEventType for SSE streaming
type StreamEventType string

const (
	StreamEventMessageStart      StreamEventType = "message_start"
	StreamEventContentBlockStart StreamEventType = "content_block_start"
	StreamEventContentBlockDelta StreamEventType = "content_block_delta"
	StreamEventContentBlockStop  StreamEventType = "content_block_stop"
	StreamEventMessageDelta      StreamEventType = "message_delta"
	StreamEventMessageStop       StreamEventType = "message_stop"
	StreamEventPing              StreamEventType = "ping"
	StreamEventError             StreamEventType = "error"
)

// StreamEvent represents a single SSE event from the API
type StreamEvent struct {
	Type         StreamEventType        `json:"type"`
	Index        int                    `json:"index,omitempty"`
	Delta        map[string]interface{} `json:"delta,omitempty"`
	Message      *AssistantMessage      `json:"message,omitempty"`      // for message_start
	ContentBlock *ContentBlock          `json:"content_block,omitempty"` // for content_block_start
	StopReason   StopReason             `json:"stop_reason,omitempty"`  // for message_delta
	Usage        *Usage                 `json:"usage,omitempty"`        // for message_delta
	Error        *APIError              `json:"error,omitempty"`        // for error
}

// APIError from streaming
type APIError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

// RequestStartEvent marks the beginning of an API request
type RequestStartEvent struct {
	RequestID string    `json:"requestId"`
	Model     string    `json:"model"`
	Timestamp time.Time `json:"timestamp"`
}

// StopHookInfo captures post-sampling hook results
type StopHookInfo struct {
	HookName string `json:"hookName"`
	Outcome  string `json:"outcome"` // "success", "blocking", "cancelled"
	Message  string `json:"message,omitempty"`
}

// PartialCompactDirection for compact boundary
type PartialCompactDirection string

const (
	CompactDirectionPre  PartialCompactDirection = "pre"
	CompactDirectionPost PartialCompactDirection = "post"
)

// NormalizedUserMessage is UserMessage prepared for API submission
type NormalizedUserMessage struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// NormalizedAssistantMessage is AssistantMessage prepared for API submission
type NormalizedAssistantMessage struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// NormalizedMessage is the union for API-ready messages
type NormalizedMessage struct {
	Role    string         `json:"role"`
	Content []ContentBlock `json:"content"`
}

// Helper constructors

func NewUserMessage(text string) Message {
	return Message{
		Kind: "user",
		User: &UserMessage{
			Role:      RoleUser,
			Content:   []ContentBlock{{Type: ContentBlockText, Text: text}},
			Timestamp: time.Now(),
		},
	}
}

func NewAssistantMessage(content []ContentBlock) Message {
	return Message{
		Kind: "assistant",
		Assistant: &AssistantMessage{
			Role:      RoleAssistant,
			Content:   content,
			Timestamp: time.Now(),
		},
	}
}

func NewSystemMessage(msgType SystemMessageType, text string) Message {
	return Message{
		Kind: "system",
		System: &SystemMessage{
			Role:      RoleSystem,
			Type:      msgType,
			Level:     SystemMessageLevelInfo,
			Text:      text,
			Timestamp: time.Now(),
		},
	}
}

func NewAttachmentMessage(content, filePath string) Message {
	return Message{
		Kind: "attachment",
		Attachment: &AttachmentMessage{
			Role:      RoleUser,
			Content:   content,
			FilePath:  filePath,
			Timestamp: time.Now(),
		},
	}
}

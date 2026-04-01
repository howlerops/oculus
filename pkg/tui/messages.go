package tui

import (
	"github.com/howlerops/oculus/pkg/types"
)

// ResponseMsg is sent when the query engine finishes a full turn
type ResponseMsg struct {
	Messages []types.Message
	Err      error
}

// StreamTextMsg delivers streamed text chunks
type StreamTextMsg struct {
	Text string
}

// StreamThinkingMsg delivers thinking block text
type StreamThinkingMsg struct {
	Text string
}

// ToolStartMsg signals a tool call has begun
type ToolStartMsg struct {
	ToolID   string
	ToolName string
}

// ToolResultMsg signals a tool call completed
type ToolResultMsg struct {
	ToolID  string
	Result  string
	IsError bool
}

// PermissionRequestMsg asks the user to approve a tool
type PermissionRequestMsg struct {
	Request PermissionRequest
}

// PermissionResponseMsg carries the user's decision
type PermissionResponseMsg struct {
	Response PermissionResponse
}

// TickMsg for progress spinner animation
type TickMsg struct{}

// TaskUpdateMsg signals a task list change
type TaskUpdateMsg struct {
	Tasks []TaskItem
}

// ErrorMsg wraps an error for display
type ErrorMsg struct {
	Err error
}

package state

import (
	"github.com/jbeck018/claude-go/pkg/task"
	"github.com/jbeck018/claude-go/pkg/types"
)

// AppState holds all application state
// This is a 1:1 port of the TypeScript AppState from old-src/state/AppStateStore.ts
type AppState struct {
	// Settings
	Verbose        bool   `json:"verbose"`
	MainLoopModel  string `json:"mainLoopModel"`
	StatusLineText string `json:"statusLineText,omitempty"`

	// Permission state
	ToolPermissionContext types.ToolPermissionContext `json:"toolPermissionContext"`

	// Messages and conversation
	Messages []types.Message `json:"messages"`

	// Task management
	Tasks map[string]task.TaskState `json:"tasks"`

	// UI state
	ExpandedView string `json:"expandedView"` // "none", "tasks", "teammates"
	IsBriefOnly  bool   `json:"isBriefOnly"`

	// Response tracking
	ResponseLength int  `json:"responseLength"`
	IsStreaming    bool `json:"isStreaming"`

	// Tool state
	InProgressToolUseIDs map[string]bool `json:"inProgressToolUseIds"`

	// Context
	ConversationID string `json:"conversationId,omitempty"`

	// Thinking mode
	ThinkingEnabled bool `json:"thinkingEnabled"`

	// Cost tracking
	TotalInputTokens  int     `json:"totalInputTokens"`
	TotalOutputTokens int     `json:"totalOutputTokens"`
	TotalCostUSD      float64 `json:"totalCostUsd"`
}

// NewAppState creates a default AppState
func NewAppState(model string) AppState {
	return AppState{
		MainLoopModel:        model,
		ExpandedView:         "none",
		ToolPermissionContext: types.NewToolPermissionContext(),
		Tasks:                make(map[string]task.TaskState),
		InProgressToolUseIDs: make(map[string]bool),
	}
}

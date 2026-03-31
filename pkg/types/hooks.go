package types

// HookEvent identifies when a hook fires
type HookEvent string

const (
	HookEventPreToolUse       HookEvent = "PreToolUse"
	HookEventPostToolUse      HookEvent = "PostToolUse"
	HookEventNotification     HookEvent = "Notification"
	HookEventPreCompact       HookEvent = "PreCompact"
	HookEventPostCompact      HookEvent = "PostCompact"
	HookEventSessionStart     HookEvent = "SessionStart"
	HookEventStop             HookEvent = "Stop"
	HookEventSubagentStop     HookEvent = "SubagentStop"
	HookEventUserPromptSubmit HookEvent = "UserPromptSubmit"
)

// HookProgress reports hook execution progress
type HookProgress struct {
	Type          string    `json:"type"` // always "hook_progress"
	HookEvent     HookEvent `json:"hookEvent"`
	HookName      string    `json:"hookName"`
	Command       string    `json:"command"`
	PromptText    string    `json:"promptText,omitempty"`
	StatusMessage string    `json:"statusMessage,omitempty"`
}

// PromptRequestOption is a single choice in a prompt
type PromptRequestOption struct {
	Key         string `json:"key"`
	Label       string `json:"label"`
	Description string `json:"description,omitempty"`
}

// PromptRequest asks the user to choose from options
type PromptRequest struct {
	Prompt  string                `json:"prompt"`
	Message string                `json:"message"`
	Options []PromptRequestOption `json:"options"`
}

// PromptResponse is the user's choice
type PromptResponse struct {
	PromptResponse string `json:"prompt_response"`
	Selected       string `json:"selected"`
}

// HookResult is the outcome of running a hook
type HookResult struct {
	Message                      *Message               `json:"message,omitempty"`
	SystemMessage                *Message               `json:"systemMessage,omitempty"`
	BlockingError                *HookBlockingError     `json:"blockingError,omitempty"`
	Outcome                      string                 `json:"outcome"` // "success", "blocking", "non_blocking_error", "cancelled"
	PreventContinuation          bool                   `json:"preventContinuation,omitempty"`
	StopReason                   string                 `json:"stopReason,omitempty"`
	PermissionBehavior           string                 `json:"permissionBehavior,omitempty"`
	HookPermissionDecisionReason string                 `json:"hookPermissionDecisionReason,omitempty"`
	AdditionalContext            string                 `json:"additionalContext,omitempty"`
	InitialUserMessage           string                 `json:"initialUserMessage,omitempty"`
	UpdatedInput                 map[string]interface{} `json:"updatedInput,omitempty"`
	Retry                        bool                   `json:"retry,omitempty"`
}

// HookBlockingError when a hook blocks execution
type HookBlockingError struct {
	BlockingError string `json:"blockingError"`
	Command       string `json:"command"`
}

// AggregatedHookResult combines results from multiple hooks
type AggregatedHookResult struct {
	Message                      *Message               `json:"message,omitempty"`
	BlockingErrors               []HookBlockingError    `json:"blockingErrors,omitempty"`
	PreventContinuation          bool                   `json:"preventContinuation,omitempty"`
	StopReason                   string                 `json:"stopReason,omitempty"`
	HookPermissionDecisionReason string                 `json:"hookPermissionDecisionReason,omitempty"`
	PermissionBehavior           string                 `json:"permissionBehavior,omitempty"`
	AdditionalContexts           []string               `json:"additionalContexts,omitempty"`
	InitialUserMessage           string                 `json:"initialUserMessage,omitempty"`
	UpdatedInput                 map[string]interface{} `json:"updatedInput,omitempty"`
	Retry                        bool                   `json:"retry,omitempty"`
}

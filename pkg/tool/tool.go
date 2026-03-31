package tool

import (
	"context"

	"github.com/jbeck018/claude-go/pkg/types"
)

// InputSchema defines the JSON Schema for a tool's input
type InputSchema struct {
	Type       string                 `json:"type"` // always "object"
	Properties map[string]interface{} `json:"properties,omitempty"`
	Required   []string               `json:"required,omitempty"`
}

// Result is the return value from a tool call
type Result struct {
	Data        interface{}     `json:"data"`
	NewMessages []types.Message `json:"newMessages,omitempty"`
	// MCPMeta holds MCP protocol metadata
	MCPMeta *MCPMeta `json:"mcpMeta,omitempty"`
}

// MCPMeta holds MCP-specific metadata
type MCPMeta struct {
	Meta              map[string]interface{} `json:"_meta,omitempty"`
	StructuredContent map[string]interface{} `json:"structuredContent,omitempty"`
}

// ValidationResult from input validation
type ValidationResult struct {
	Valid     bool   `json:"valid"`
	Message   string `json:"message,omitempty"`
	ErrorCode int    `json:"errorCode,omitempty"`
}

// SearchOrReadInfo describes if a tool operation is search/read
type SearchOrReadInfo struct {
	IsSearch bool `json:"isSearch"`
	IsRead   bool `json:"isRead"`
	IsList   bool `json:"isList,omitempty"`
}

// InterruptBehavior defines what happens on user interrupt
type InterruptBehavior string

const (
	InterruptCancel InterruptBehavior = "cancel"
	InterruptBlock  InterruptBehavior = "block"
)

// Tool is the interface all tools must implement
// This is a 1:1 port of the TypeScript Tool type from old-src/Tool.ts
type Tool interface {
	// Name returns the tool's primary identifier
	Name() string

	// Aliases returns alternative names for backwards compatibility
	Aliases() []string

	// SearchHint returns a keyword phrase for ToolSearch matching
	SearchHint() string

	// Description returns a dynamic description based on input
	Description(ctx context.Context, input map[string]interface{}) (string, error)

	// InputSchema returns the JSON Schema for the tool's input
	GetInputSchema() InputSchema

	// Call executes the tool with the given input
	Call(ctx context.Context, input map[string]interface{}, onProgress func(types.ToolProgressData)) (*Result, error)

	// IsEnabled returns whether the tool is currently available
	IsEnabled() bool

	// IsConcurrencySafe returns true if the tool can run in parallel
	IsConcurrencySafe(input map[string]interface{}) bool

	// IsReadOnly returns true if the tool only reads (no writes)
	IsReadOnly(input map[string]interface{}) bool

	// IsDestructive returns true for irreversible operations
	IsDestructive(input map[string]interface{}) bool

	// CheckPermissions determines if permission is needed
	CheckPermissions(ctx context.Context, input map[string]interface{}) (*types.PermissionResult, error)

	// Prompt returns the tool's system prompt contribution
	Prompt(ctx context.Context) (string, error)

	// UserFacingName returns the display name for UI
	UserFacingName(input map[string]interface{}) string

	// MaxResultSizeChars returns max result size before disk persistence
	MaxResultSizeChars() int

	// InterruptBehavior defines behavior on user interrupt
	GetInterruptBehavior() InterruptBehavior

	// IsSearchOrReadCommand categorizes the operation type
	IsSearchOrReadCommand(input map[string]interface{}) *SearchOrReadInfo

	// ValidateInput checks if the input is valid
	ValidateInput(ctx context.Context, input map[string]interface{}) *ValidationResult

	// ShouldDefer returns true if this tool should be deferred (needs ToolSearch)
	ShouldDefer() bool

	// IsMCP returns true for MCP-provided tools
	IsMCP() bool

	// IsLSP returns true for LSP-provided tools
	IsLSP() bool
}

// Tools is a collection of tools
type Tools []Tool

// FindByName looks up a tool by name or alias
func (t Tools) FindByName(name string) Tool {
	for _, tool := range t {
		if tool.Name() == name {
			return tool
		}
		for _, alias := range tool.Aliases() {
			if alias == name {
				return tool
			}
		}
	}
	return nil
}

// FilterEnabled returns only enabled tools
func (t Tools) FilterEnabled() Tools {
	var result Tools
	for _, tool := range t {
		if tool.IsEnabled() {
			result = append(result, tool)
		}
	}
	return result
}

// MatchesName checks if a tool matches by name or alias
func MatchesName(tool Tool, name string) bool {
	if tool.Name() == name {
		return true
	}
	for _, alias := range tool.Aliases() {
		if alias == name {
			return true
		}
	}
	return false
}

// BaseTool provides default implementations for the Tool interface
// Embed this in concrete tool implementations to get sensible defaults
type BaseTool struct {
	ToolName          string
	ToolAliases       []string
	ToolSearchHint    string
	ToolMaxResultSize int
	ToolShouldDefer   bool
}

func (b *BaseTool) Name() string    { return b.ToolName }
func (b *BaseTool) Aliases() []string { return b.ToolAliases }
func (b *BaseTool) SearchHint() string { return b.ToolSearchHint }
func (b *BaseTool) IsEnabled() bool  { return true }
func (b *BaseTool) IsConcurrencySafe(_ map[string]interface{}) bool { return false }
func (b *BaseTool) IsReadOnly(_ map[string]interface{}) bool        { return false }
func (b *BaseTool) IsDestructive(_ map[string]interface{}) bool     { return false }
func (b *BaseTool) ShouldDefer() bool                               { return b.ToolShouldDefer }
func (b *BaseTool) IsMCP() bool                                     { return false }
func (b *BaseTool) IsLSP() bool                                     { return false }
func (b *BaseTool) GetInterruptBehavior() InterruptBehavior         { return InterruptBlock }
func (b *BaseTool) IsSearchOrReadCommand(_ map[string]interface{}) *SearchOrReadInfo { return nil }
func (b *BaseTool) ValidateInput(_ context.Context, _ map[string]interface{}) *ValidationResult {
	return nil
}
func (b *BaseTool) UserFacingName(_ map[string]interface{}) string { return b.ToolName }
func (b *BaseTool) MaxResultSizeChars() int {
	if b.ToolMaxResultSize > 0 {
		return b.ToolMaxResultSize
	}
	return 100000 // 100k default
}
func (b *BaseTool) CheckPermissions(_ context.Context, input map[string]interface{}) (*types.PermissionResult, error) {
	return &types.PermissionResult{
		Behavior:     types.PermissionAllow,
		UpdatedInput: input,
	}, nil
}

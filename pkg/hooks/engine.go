package hooks

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/jbeck018/claude-go/pkg/types"
)

// HookConfig defines a single hook from settings.json
type HookConfig struct {
	Matcher  string `json:"matcher,omitempty"`
	Command  string `json:"command"`
	Timeout  int    `json:"timeout,omitempty"` // ms, default 60000
	Internal bool   `json:"internal,omitempty"`
}

// HookInput is the JSON passed to hook commands via stdin
type HookInput struct {
	SessionID  string                 `json:"session_id"`
	Event      types.HookEvent        `json:"hook_event"`
	ToolName   string                 `json:"tool_name,omitempty"`
	ToolInput  map[string]interface{} `json:"tool_input,omitempty"`
	ToolOutput interface{}            `json:"tool_output,omitempty"`
}

// HookOutput is the JSON returned from hook commands
type HookOutput struct {
	Decision          string                 `json:"decision,omitempty"` // "approve", "deny", "block"
	Reason            string                 `json:"reason,omitempty"`
	Message           string                 `json:"message,omitempty"`
	UpdatedInput      map[string]interface{} `json:"updated_input,omitempty"`
	AdditionalContext string                 `json:"additional_context,omitempty"`
}

// Registry holds all configured hooks
type Registry struct {
	hooks map[types.HookEvent][]HookConfig
}

// NewRegistry creates an empty hook registry
func NewRegistry() *Registry {
	return &Registry{
		hooks: make(map[types.HookEvent][]HookConfig),
	}
}

// Register adds a hook for the given event
func (r *Registry) Register(event types.HookEvent, cfg HookConfig) {
	r.hooks[event] = append(r.hooks[event], cfg)
}

// RegisterFromSettings loads hooks from the settings.json hooks map
func (r *Registry) RegisterFromSettings(hooksMap map[string][]HookConfig) {
	for eventStr, configs := range hooksMap {
		event := types.HookEvent(eventStr)
		for _, cfg := range configs {
			r.Register(event, cfg)
		}
	}
}

// GetHooks returns all hooks for an event
func (r *Registry) GetHooks(event types.HookEvent) []HookConfig {
	return r.hooks[event]
}

// Execute runs all hooks for an event and returns aggregated results
func (r *Registry) Execute(ctx context.Context, event types.HookEvent, input HookInput) (*types.AggregatedHookResult, error) {
	hooks := r.GetHooks(event)
	if len(hooks) == 0 {
		return &types.AggregatedHookResult{}, nil
	}

	result := &types.AggregatedHookResult{}

	for _, hook := range hooks {
		// Check matcher
		if hook.Matcher != "" && input.ToolName != "" {
			if !matchesHookPattern(hook.Matcher, input.ToolName, input.ToolInput) {
				continue
			}
		}

		output, err := executeHook(ctx, hook, input)
		if err != nil {
			// Non-blocking error - continue
			continue
		}

		// Aggregate results
		if output.AdditionalContext != "" {
			result.AdditionalContexts = append(result.AdditionalContexts, output.AdditionalContext)
		}
		if output.UpdatedInput != nil {
			result.UpdatedInput = output.UpdatedInput
		}
		if output.Decision == "block" || output.Decision == "deny" {
			result.PreventContinuation = true
			result.StopReason = output.Reason
			if output.Decision == "deny" {
				result.PermissionBehavior = "deny"
			}
			break
		}
		if output.Decision == "approve" {
			result.PermissionBehavior = "allow"
		}
	}

	return result, nil
}

func executeHook(ctx context.Context, hook HookConfig, input HookInput) (*HookOutput, error) {
	timeout := 60 * time.Second
	if hook.Timeout > 0 {
		timeout = time.Duration(hook.Timeout) * time.Millisecond
	}

	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// Parse command
	parts := strings.Fields(hook.Command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty hook command")
	}

	cmd := exec.CommandContext(cmdCtx, parts[0], parts[1:]...)

	// Pass input via stdin
	inputJSON, _ := json.Marshal(input)
	cmd.Stdin = bytes.NewReader(inputJSON)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("hook %q failed: %w (stderr: %s)", hook.Command, err, stderr.String())
	}

	// Parse output
	var output HookOutput
	if stdout.Len() > 0 {
		if err := json.Unmarshal(stdout.Bytes(), &output); err != nil {
			// Non-JSON output - treat as additional context
			output.AdditionalContext = stdout.String()
		}
	}

	return &output, nil
}

func matchesHookPattern(matcher, toolName string, input map[string]interface{}) bool {
	// Simple matching: "ToolName" or "ToolName(pattern)"
	if matcher == toolName {
		return true
	}
	if strings.HasPrefix(matcher, toolName+"(") && strings.HasSuffix(matcher, ")") {
		return true // Content matching handled by the hook itself
	}
	// Wildcard
	if matcher == "*" {
		return true
	}
	return false
}

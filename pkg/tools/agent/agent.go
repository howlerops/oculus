package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/config"
	"github.com/howlerops/oculus/pkg/query"
	"github.com/howlerops/oculus/pkg/state"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type AgentTool struct {
	tool.BaseTool
	ParentClient *api.Client
	ParentTools  tool.Tools
}

func NewAgentTool(client *api.Client, tools tool.Tools) *AgentTool {
	return &AgentTool{
		BaseTool: tool.BaseTool{
			ToolName:          "Agent",
			ToolSearchHint:    "spawn subagent delegate task",
			ToolMaxResultSize: 100000,
		},
		ParentClient: client,
		ParentTools:  tools,
	}
}

func (t *AgentTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"prompt":            map[string]interface{}{"type": "string", "description": "The task for the agent to perform"},
			"description":       map[string]interface{}{"type": "string", "description": "Short 3-5 word description"},
			"subagent_type":     map[string]interface{}{"type": "string", "description": "Type of specialized agent"},
			"model":             map[string]interface{}{"type": "string", "description": "Model override (sonnet/opus/haiku)"},
			"run_in_background": map[string]interface{}{"type": "boolean", "description": "Run asynchronously"},
			"name":              map[string]interface{}{"type": "string", "description": "Agent name for addressing"},
		},
		Required: []string{"prompt", "description"},
	}
}

func (t *AgentTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Launch a new agent to handle complex, multi-step tasks autonomously.", nil
}

func (t *AgentTool) Prompt(_ context.Context) (string, error) {
	return "Launch a new agent to handle complex, multi-step tasks autonomously.\n\nThe Agent tool launches specialized agents that autonomously handle complex tasks.\n\nUsage:\n- Always include a short description (3-5 words)\n- Launch multiple agents concurrently when possible\n- The agent's result is not visible to the user - summarize it\n- Provide clear, detailed prompts for autonomous work", nil
}

func (t *AgentTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }

func (t *AgentTool) Call(ctx context.Context, input map[string]interface{}, onProgress func(types.ToolProgressData)) (*tool.Result, error) {
	prompt, _ := input["prompt"].(string)
	if prompt == "" {
		return &tool.Result{Data: "Error: prompt is required"}, nil
	}

	description, _ := input["description"].(string)
	model, _ := input["model"].(string)
	runInBg, _ := input["run_in_background"].(bool)

	// Resolve model
	if model == "" {
		model = config.GetModel()
	} else {
		// Map short names to full model IDs
		switch strings.ToLower(model) {
		case "opus":
			model = "claude-opus-4-20250514"
		case "sonnet":
			model = "claude-sonnet-4-20250514"
		case "haiku":
			model = "claude-haiku-4-20250506"
		}
	}

	agentID := types.NewAgentId(description)

	if onProgress != nil {
		onProgress(types.ToolProgressData{Type: types.ProgressTypeAgent})
	}

	// Create isolated state for the subagent
	subStore := state.NewStore(state.NewAppState(model))

	// Create query engine for subagent
	engine := query.NewEngine(t.ParentClient, t.ParentTools, subStore, model)

	messages := []types.Message{
		types.NewUserMessage(prompt),
	}

	systemPrompt := "You are a specialized agent. Complete the assigned task thoroughly and return the result."

	if runInBg {
		// Run in background goroutine
		var resultOnce sync.Once
		var bgResult string

		go func() {
			handler := &agentStreamHandler{}
			resultMsgs, err := engine.RunQuery(ctx, messages, systemPrompt, handler)
			resultOnce.Do(func() {
				if err != nil {
					bgResult = fmt.Sprintf("Agent error: %v", err)
				} else {
					bgResult = extractAssistantText(resultMsgs)
				}
			})
			_ = bgResult
		}()

		return &tool.Result{
			Data: fmt.Sprintf("Agent %q launched in background (model: %s, id: %s)", description, model, agentID),
		}, nil
	}

	// Run synchronously
	handler := &agentStreamHandler{onProgress: onProgress, agentID: string(agentID)}
	resultMsgs, err := engine.RunQuery(ctx, messages, systemPrompt, handler)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Agent error: %v", err)}, nil
	}

	result := extractAssistantText(resultMsgs)

	if onProgress != nil {
		onProgress(types.ToolProgressData{Type: types.ProgressTypeAgent})
	}

	return &tool.Result{Data: result}, nil
}

func extractAssistantText(messages []types.Message) string {
	// Find the last assistant message and extract text
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Kind == "assistant" && msg.Assistant != nil {
			var parts []string
			for _, block := range msg.Assistant.Content {
				if block.Type == types.ContentBlockText && block.Text != "" {
					parts = append(parts, block.Text)
				}
			}
			if len(parts) > 0 {
				return strings.Join(parts, "\n")
			}
		}
	}
	return "(no response from agent)"
}

type agentStreamHandler struct {
	onProgress func(types.ToolProgressData)
	agentID    string
}

func (h *agentStreamHandler) OnText(_ string)                          {}
func (h *agentStreamHandler) OnToolUseStart(_, _ string)               {}
func (h *agentStreamHandler) OnToolUseResult(_ string, _ interface{})  {}
func (h *agentStreamHandler) OnThinking(_ string)                      {}
func (h *agentStreamHandler) OnComplete(_ types.StopReason, _ *types.Usage) {}
func (h *agentStreamHandler) OnError(_ error)                          {}

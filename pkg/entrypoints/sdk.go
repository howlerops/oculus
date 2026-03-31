package entrypoints

import (
	"context"

	"github.com/jbeck018/claude-go/pkg/api"
	"github.com/jbeck018/claude-go/pkg/query"
	"github.com/jbeck018/claude-go/pkg/state"
	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

// SDKRunner provides a programmatic interface (no TUI)
type SDKRunner struct {
	Engine       *query.Engine
	SystemPrompt interface{}
}

func NewSDKRunner(client *api.Client, tools tool.Tools, store *state.Store, model string, systemPrompt interface{}) *SDKRunner {
	return &SDKRunner{
		Engine:       query.NewEngine(client, tools, store, model),
		SystemPrompt: systemPrompt,
	}
}

// RunOnce sends a single message and returns the response
func (r *SDKRunner) RunOnce(ctx context.Context, prompt string) (string, error) {
	messages := []types.Message{types.NewUserMessage(prompt)}
	handler := &sdkHandler{}
	resultMsgs, err := r.Engine.RunQuery(ctx, messages, r.SystemPrompt, handler)
	if err != nil {
		return "", err
	}
	// Extract last assistant text
	for i := len(resultMsgs) - 1; i >= 0; i-- {
		msg := resultMsgs[i]
		if msg.Kind == "assistant" && msg.Assistant != nil {
			for _, block := range msg.Assistant.Content {
				if block.Type == types.ContentBlockText {
					return block.Text, nil
				}
			}
		}
	}
	return "", nil
}

// RunStream sends a message and streams the response via callback
func (r *SDKRunner) RunStream(ctx context.Context, prompt string, onText func(string)) error {
	messages := []types.Message{types.NewUserMessage(prompt)}
	handler := &sdkHandler{onText: onText}
	_, err := r.Engine.RunQuery(ctx, messages, r.SystemPrompt, handler)
	return err
}

type sdkHandler struct {
	onText func(string)
}

func (h *sdkHandler) OnText(text string) {
	if h.onText != nil {
		h.onText(text)
	}
}
func (h *sdkHandler) OnToolUseStart(id, name string)              {}
func (h *sdkHandler) OnToolUseResult(id string, result interface{}) {}
func (h *sdkHandler) OnThinking(text string)                      {}
func (h *sdkHandler) OnComplete(stopReason types.StopReason, usage *types.Usage) {}
func (h *sdkHandler) OnError(err error)                           {}

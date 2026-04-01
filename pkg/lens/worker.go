package lens

import (
	"context"
	"fmt"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/bridge"
	"github.com/howlerops/oculus/pkg/query"
	"github.com/howlerops/oculus/pkg/state"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

// LensWorker wraps a query.Engine with lens-specific configuration
type LensWorker struct {
	Lens   LensConfig
	Engine *query.Engine
	Store  *state.Store
	Active bool
	Bridge bridge.Bridge // optional - if set, used instead of raw Engine
}

// NewLensWorker creates a worker for a specific lens
func NewLensWorker(cfg LensConfig, client *api.Client, tools tool.Tools, store *state.Store) *LensWorker {
	engine := query.NewEngine(client, tools, store, cfg.Model)
	return &LensWorker{
		Lens:   cfg,
		Engine: engine,
		Store:  store,
		Active: cfg.Enabled,
	}
}

// SetBridge sets an optional bridge to use instead of the raw Engine
func (w *LensWorker) SetBridge(b bridge.Bridge) {
	w.Bridge = b
}

// RunQuery executes a query through this lens with its persona injected.
// If a Bridge is set, it routes through the bridge for streaming.
// Falls back to Engine for tool loops (bridges don't handle tool dispatch).
func (w *LensWorker) RunQuery(ctx context.Context, messages []types.Message, baseSystemPrompt interface{}, handler query.StreamHandler) ([]types.Message, error) {
	systemPrompt := w.buildSystemPrompt(baseSystemPrompt)

	// If bridge is set and available, use it for the streaming call
	if w.Bridge != nil && w.Bridge.IsAvailable() {
		return w.runViaBridge(ctx, messages, systemPrompt, handler)
	}

	// Default: use the query engine directly (handles tool loops natively)
	return w.Engine.RunQuery(ctx, messages, systemPrompt, handler)
}

// runViaBridge routes through the bridge for streaming, then falls back to
// Engine for tool loop handling since bridges don't dispatch tools.
func (w *LensWorker) runViaBridge(ctx context.Context, messages []types.Message, systemPrompt interface{}, handler query.StreamHandler) ([]types.Message, error) {
	// Convert types.Message -> bridge.Message
	var bridgeMsgs []bridge.Message
	for _, msg := range messages {
		switch msg.Kind {
		case "user":
			if msg.User != nil {
				var text string
				for _, b := range msg.User.Content {
					if b.Type == types.ContentBlockText {
						text += b.Text
					}
				}
				bridgeMsgs = append(bridgeMsgs, bridge.Message{Role: "user", Content: text})
			}
		case "assistant":
			if msg.Assistant != nil {
				var text string
				for _, b := range msg.Assistant.Content {
					if b.Type == types.ContentBlockText {
						text += b.Text
					}
				}
				bridgeMsgs = append(bridgeMsgs, bridge.Message{Role: "assistant", Content: text})
			}
		}
	}

	// Convert system prompt to string
	promptStr := ""
	switch p := systemPrompt.(type) {
	case string:
		promptStr = p
	default:
		promptStr = fmt.Sprintf("%v", p)
	}

	// Stream through bridge, forwarding chunks to the handler
	err := w.Bridge.Stream(ctx, bridgeMsgs, promptStr, nil, func(chunk bridge.StreamChunk) {
		switch chunk.Type {
		case "text":
			handler.OnText(chunk.Text)
		case "done":
			handler.OnComplete(types.StopReasonEndTurn, nil)
		case "error":
			handler.OnError(fmt.Errorf("%s", chunk.Error))
		}
	})

	if err != nil {
		return messages, err
	}

	// Note: Bridge streaming doesn't handle tool loops.
	// If the model requests tool use, the response will contain the text only.
	// For full tool-loop support, use Engine directly (the default path).
	return messages, nil
}

// buildSystemPrompt prepends the lens persona to the base system prompt
func (w *LensWorker) buildSystemPrompt(base interface{}) interface{} {
	if w.Lens.Persona == "" {
		return base
	}

	switch b := base.(type) {
	case string:
		return w.Lens.Persona + "\n\n" + b
	default:
		// For non-string (e.g., SystemBlock slices), return as-is with persona prepended
		return w.Lens.Persona + "\n\n" + fmt.Sprintf("%v", b)
	}
}

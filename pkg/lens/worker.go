package lens

import (
	"context"
	"fmt"

	"github.com/howlerops/oculus/pkg/api"
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

// RunQuery executes a query through this lens with its persona injected
func (w *LensWorker) RunQuery(ctx context.Context, messages []types.Message, baseSystemPrompt interface{}, handler query.StreamHandler) ([]types.Message, error) {
	// Prepend lens persona to system prompt
	systemPrompt := w.buildSystemPrompt(baseSystemPrompt)
	return w.Engine.RunQuery(ctx, messages, systemPrompt, handler)
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

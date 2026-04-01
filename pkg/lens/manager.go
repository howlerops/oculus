package lens

import (
	"context"
	"fmt"
	"sync"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/bridge"
	"github.com/howlerops/oculus/pkg/query"
	"github.com/howlerops/oculus/pkg/state"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

// Manager orchestrates the three lenses (Focus, Scan, Craft)
type Manager struct {
	mu      sync.RWMutex
	workers map[LensType]*LensWorker
	router  Router
	config  OculusConfig
}

// NewManager creates a LensManager with the given configuration
func NewManager(cfg OculusConfig, client *api.Client, tools tool.Tools, store *state.Store) *Manager {
	m := &Manager{
		workers: make(map[LensType]*LensWorker),
		router:  NewStaticRouter(),
		config:  cfg,
	}

	// Create workers for each enabled lens
	if cfg.Focus.Enabled {
		m.workers[LensFocus] = NewLensWorker(cfg.Focus, client, tools, store)
	}
	if cfg.Scan.Enabled {
		m.workers[LensScan] = NewLensWorker(cfg.Scan, client, tools, store)
	}
	if cfg.Craft.Enabled {
		m.workers[LensCraft] = NewLensWorker(cfg.Craft, client, tools, store)
	}

	// Wire bridges for workers
	// If a provider is non-anthropic, always use a bridge.
	// If provider is anthropic but the api.Client has no key, fall back to claude-code bridge.
	apiKeyAvailable := client.GetAPIKey() != ""
	for lensType, worker := range m.workers {
		lensCfg := getLensConfig(cfg, lensType)
		provider := lensCfg.Provider
		if provider == "" {
			provider = "anthropic"
		}

		needsBridge := false
		if provider != "anthropic" {
			needsBridge = true
		} else if !apiKeyAvailable {
			// No API key - try claude-code CLI as fallback
			provider = "claude-code"
			needsBridge = true
		}

		if needsBridge {
			bridgeCfg := bridge.BridgeConfig{
				Provider: provider,
				Model:    lensCfg.Model,
			}
			b, err := bridge.CreateBridge(bridgeCfg)
			if err == nil && b.IsAvailable() {
				worker.SetBridge(b)
			}
		}
	}

	return m
}

// getLensConfig returns the LensConfig for a given LensType
func getLensConfig(cfg OculusConfig, lensType LensType) LensConfig {
	switch lensType {
	case LensFocus:
		return cfg.Focus
	case LensScan:
		return cfg.Scan
	case LensCraft:
		return cfg.Craft
	default:
		return cfg.Focus
	}
}

// GetWorker returns the worker for a lens type, falling back to Focus
func (m *Manager) GetWorker(lensType LensType) *LensWorker {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if w, ok := m.workers[lensType]; ok && w.Active {
		return w
	}
	// Fallback to Focus lens
	if w, ok := m.workers[LensFocus]; ok {
		return w
	}
	return nil
}

// GetFocusWorker returns the main orchestration lens
func (m *Manager) GetFocusWorker() *LensWorker {
	return m.GetWorker(LensFocus)
}

// RouteToolCall determines which lens should handle a tool call
func (m *Manager) RouteToolCall(toolName string) *LensWorker {
	lensType := m.router.RouteToolCall(toolName)
	return m.GetWorker(lensType)
}

// RouteMessage determines which lens should handle a user message
func (m *Manager) RouteMessage(text string) *LensWorker {
	lensType := m.router.RouteMessage(text)
	return m.GetWorker(lensType)
}

// RunQuery routes and executes a query through the appropriate lens
func (m *Manager) RunQuery(ctx context.Context, messages []types.Message, systemPrompt interface{}, handler query.StreamHandler) ([]types.Message, error) {
	// Determine which lens based on the last user message
	var userText string
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Kind == "user" && messages[i].User != nil {
			for _, b := range messages[i].User.Content {
				if b.Type == types.ContentBlockText {
					userText = b.Text
					break
				}
			}
			break
		}
	}

	worker := m.RouteMessage(userText)
	if worker == nil {
		return nil, fmt.Errorf("no active lens available")
	}

	return worker.RunQuery(ctx, messages, systemPrompt, handler)
}

// Handoff creates an episode summary from one lens and passes it to another
func (m *Manager) Handoff(from LensType, to LensType, summary EpisodeSummary) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	toWorker, ok := m.workers[to]
	if !ok || !toWorker.Active {
		return fmt.Errorf("target lens %s not available", to)
	}

	// Inject episode summary as context for the target lens
	// In v0.3.0+, this will use RLM's episode store
	_ = summary
	return nil
}

// GetActiveCount returns how many lenses are active
func (m *Manager) GetActiveCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	count := 0
	for _, w := range m.workers {
		if w.Active {
			count++
		}
	}
	return count
}

// GetConfig returns the current lens configuration
func (m *Manager) GetConfig() OculusConfig {
	return m.config
}

// IsMultiLens returns true if more than one lens is configured with different models
func (m *Manager) IsMultiLens() bool {
	models := make(map[string]bool)
	for _, w := range m.workers {
		if w.Active {
			models[w.Lens.Model] = true
		}
	}
	return len(models) > 1
}

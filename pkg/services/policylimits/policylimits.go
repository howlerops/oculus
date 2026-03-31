package policylimits

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/howlerops/oculus/pkg/config"
)

type PolicyLimits struct {
	MaxBudgetUSD       *float64 `json:"max_budget_usd,omitempty"`
	MaxTurnsPerSession *int     `json:"max_turns_per_session,omitempty"`
	AllowedTools       []string `json:"allowed_tools,omitempty"`
	BlockedTools       []string `json:"blocked_tools,omitempty"`
	AllowedModels      []string `json:"allowed_models,omitempty"`
	DisableWebSearch   bool     `json:"disable_web_search,omitempty"`
	DisableWebFetch    bool     `json:"disable_web_fetch,omitempty"`
	DisableMCP         bool     `json:"disable_mcp,omitempty"`
	RequirePermissions bool     `json:"require_permissions,omitempty"`
	EnforceMode        string   `json:"enforce_mode,omitempty"`
}

func LoadPolicyLimits() (*PolicyLimits, error) {
	paths := []string{
		filepath.Join(config.GetOculusDir(), "policy.json"),
		"/etc/claude-code/policy.json",
	}
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		var limits PolicyLimits
		if err := json.Unmarshal(data, &limits); err != nil {
			continue
		}
		return &limits, nil
	}
	return nil, nil
}

func (p *PolicyLimits) IsToolAllowed(toolName string) bool {
	if p == nil {
		return true
	}
	if len(p.BlockedTools) > 0 {
		for _, t := range p.BlockedTools {
			if t == toolName {
				return false
			}
		}
	}
	if len(p.AllowedTools) > 0 {
		for _, t := range p.AllowedTools {
			if t == toolName {
				return true
			}
		}
		return false
	}
	return true
}

func (p *PolicyLimits) IsModelAllowed(model string) bool {
	if p == nil || len(p.AllowedModels) == 0 {
		return true
	}
	for _, m := range p.AllowedModels {
		if m == model {
			return true
		}
	}
	return false
}

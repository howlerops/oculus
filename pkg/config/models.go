package config

import "strings"

// ModelProvider identifies the API provider
type ModelProvider string

const (
	ProviderFirstParty ModelProvider = "firstParty" // api.anthropic.com
	ProviderBedrock    ModelProvider = "bedrock"     // AWS Bedrock
	ProviderVertex     ModelProvider = "vertex"      // GCP Vertex AI
)

// ModelInfo holds metadata about a model
type ModelInfo struct {
	ID               string
	DisplayName      string
	Provider         ModelProvider
	ContextWindow    int
	MaxOutput        int
	SupportsImages   bool
	SupportsThinking bool
	SupportsPDFs     bool
	SupportsTools    bool
	CostInput        float64 // per million tokens
	CostOutput       float64 // per million tokens
}

// ModelRegistry holds all known models
var ModelRegistry = map[string]ModelInfo{
	"claude-opus-4-20250514": {
		ID: "claude-opus-4-20250514", DisplayName: "Claude Opus 4",
		Provider: ProviderFirstParty, ContextWindow: 200000, MaxOutput: 16384,
		SupportsImages: true, SupportsThinking: true, SupportsPDFs: true, SupportsTools: true,
		CostInput: 15.0, CostOutput: 75.0,
	},
	"claude-sonnet-4-20250514": {
		ID: "claude-sonnet-4-20250514", DisplayName: "Claude Sonnet 4",
		Provider: ProviderFirstParty, ContextWindow: 200000, MaxOutput: 16384,
		SupportsImages: true, SupportsThinking: true, SupportsPDFs: true, SupportsTools: true,
		CostInput: 3.0, CostOutput: 15.0,
	},
	"claude-haiku-4-20250506": {
		ID: "claude-haiku-4-20250506", DisplayName: "Claude Haiku 4",
		Provider: ProviderFirstParty, ContextWindow: 200000, MaxOutput: 8192,
		SupportsImages: true, SupportsThinking: false, SupportsPDFs: true, SupportsTools: true,
		CostInput: 0.80, CostOutput: 4.0,
	},
}

// ModelAliasMap maps short names to full model IDs
var ModelAliasMap = map[string]string{
	"opus":   "claude-opus-4-20250514",
	"sonnet": "claude-sonnet-4-20250514",
	"haiku":  "claude-haiku-4-20250506",
	"o":      "claude-opus-4-20250514",
	"s":      "claude-sonnet-4-20250514",
	"h":      "claude-haiku-4-20250506",
}

// ResolveModel resolves aliases and validates model IDs
func ResolveModel(name string) (ModelInfo, bool) {
	name = strings.ToLower(strings.TrimSpace(name))
	if alias, ok := ModelAliasMap[name]; ok {
		name = alias
	}
	info, ok := ModelRegistry[name]
	return info, ok
}

// GetModelInfo returns info for a model ID, with fallback
func GetModelInfo(modelID string) ModelInfo {
	if info, ok := ModelRegistry[modelID]; ok {
		return info
	}
	// Fallback to sonnet defaults
	return ModelInfo{
		ID: modelID, DisplayName: modelID,
		Provider: ProviderFirstParty, ContextWindow: 200000, MaxOutput: 16384,
		SupportsImages: true, SupportsTools: true,
		CostInput: 3.0, CostOutput: 15.0,
	}
}

// ListModels returns all available model IDs
func ListModels() []string {
	var models []string
	for id := range ModelRegistry {
		models = append(models, id)
	}
	return models
}

// GetProvider determines the API provider for a model
func GetProvider(modelID string) ModelProvider {
	info := GetModelInfo(modelID)
	return info.Provider
}

// EstimateCost calculates cost for token usage
func EstimateCost(modelID string, inputTokens, outputTokens int) float64 {
	info := GetModelInfo(modelID)
	return (float64(inputTokens)/1_000_000)*info.CostInput +
		(float64(outputTokens)/1_000_000)*info.CostOutput
}

// SupportsFeature checks if a model supports a specific feature
func SupportsFeature(modelID, feature string) bool {
	info := GetModelInfo(modelID)
	switch feature {
	case "images":
		return info.SupportsImages
	case "thinking":
		return info.SupportsThinking
	case "pdfs":
		return info.SupportsPDFs
	case "tools":
		return info.SupportsTools
	default:
		return false
	}
}

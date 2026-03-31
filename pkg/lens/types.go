package lens

// LensType identifies the three Oculus lenses
type LensType string

const (
	LensFocus LensType = "focus" // Main reasoning, planning, orchestration
	LensScan  LensType = "scan"  // Exploration, file search, codebase analysis
	LensCraft LensType = "craft" // Code writing, editing, command execution
)

// LensConfig configures a single lens
type LensConfig struct {
	Type     LensType `json:"type"`
	Model    string   `json:"model"`             // Model ID for this lens
	Provider string   `json:"provider,omitempty"` // Provider (anthropic, openai, ollama)
	Enabled  bool     `json:"enabled"`
	Persona  string   `json:"persona,omitempty"` // System prompt persona
}

// OculusConfig holds the three-lens configuration
type OculusConfig struct {
	Focus LensConfig `json:"focus"` // Cortex equivalent
	Scan  LensConfig `json:"scan"`  // Scout equivalent
	Craft LensConfig `json:"craft"` // Forge equivalent
}

// DefaultConfig returns the default lens configuration (all Anthropic)
func DefaultConfig() OculusConfig {
	return OculusConfig{
		Focus: LensConfig{
			Type:    LensFocus,
			Model:   "claude-sonnet-4-20250514",
			Enabled: true,
			Persona: "You are the Focus lens - primary orchestrator. Analyze tasks, plan approaches, and coordinate work across lenses.",
		},
		Scan: LensConfig{
			Type:    LensScan,
			Model:   "claude-sonnet-4-20250514",
			Enabled: true,
			Persona: "You are the Scan lens - specialized in exploration and research. Search files, analyze code, gather information.",
		},
		Craft: LensConfig{
			Type:    LensCraft,
			Model:   "claude-sonnet-4-20250514",
			Enabled: true,
			Persona: "You are the Craft lens - specialized in execution. Write code, run commands, make changes.",
		},
	}
}

// EpisodeSummary is a compressed handoff between lenses
// (Matches Iron Rain's episode summary concept)
type EpisodeSummary struct {
	LensType     LensType          `json:"lens_type"`
	Summary      string            `json:"summary"`
	Keywords     []string          `json:"keywords"`
	FilesTouched []string          `json:"files_touched,omitempty"`
	TokensUsed   int               `json:"tokens_used"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

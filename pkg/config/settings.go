package config

// SettingsJson represents ~/.claude/settings.json
type SettingsJson struct {
	// Permission settings
	DefaultMode string                `json:"defaultMode,omitempty"`
	CustomModes map[string]ModeConfig `json:"customModes,omitempty"`

	// Tool permissions
	AllowedTools    []string `json:"allowedTools,omitempty"`
	DisallowedTools []string `json:"disallowedTools,omitempty"`

	// Hooks
	Hooks map[string][]HookConfig `json:"hooks,omitempty"`

	// MCP servers
	MCPServers map[string]MCPServerConfig `json:"mcpServers,omitempty"`

	// Environment variables
	Env map[string]string `json:"env,omitempty"`

	// UI settings
	Theme      string            `json:"theme,omitempty"`
	StatusLine *StatusLineConfig `json:"statusLine,omitempty"`
	Verbose    bool              `json:"verbose,omitempty"`

	// Model settings
	Model          string `json:"model,omitempty"`
	SmallFastModel string `json:"smallFastModel,omitempty"`

	// Plugin settings
	EnabledPlugins         map[string]bool        `json:"enabledPlugins,omitempty"`
	ExtraKnownMarketplaces map[string]interface{} `json:"extraKnownMarketplaces,omitempty"`

	// Behavior settings
	SkipDangerousModePrompt bool   `json:"skipDangerousModePermissionPrompt,omitempty"`
	TeammateMode            string `json:"teammateMode,omitempty"`

	// Additional working directories
	AdditionalDirectories []string `json:"additionalDirectories,omitempty"`

	// Lens settings
	Lenses *LensSettings `json:"lenses,omitempty"`

	// Permission rules
	Permissions *PermissionSettings `json:"permissions,omitempty"`
}

type LensSettings struct {
	Focus *LensModelConfig `json:"focus,omitempty"`
	Scan  *LensModelConfig `json:"scan,omitempty"`
	Craft *LensModelConfig `json:"craft,omitempty"`
}

type LensModelConfig struct {
	Model    string `json:"model,omitempty"`
	Provider string `json:"provider,omitempty"`
	Enabled  *bool  `json:"enabled,omitempty"`
}

type ModeConfig struct {
	Tools       []string `json:"tools,omitempty"`
	Description string   `json:"description,omitempty"`
}

type HookConfig struct {
	Matcher  string `json:"matcher,omitempty"`
	Command  string `json:"command,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
	Internal bool   `json:"internal,omitempty"`
}

type MCPServerConfig struct {
	Command   string            `json:"command,omitempty"`
	Args      []string          `json:"args,omitempty"`
	Env       map[string]string `json:"env,omitempty"`
	Transport string            `json:"transport,omitempty"`
	URL       string            `json:"url,omitempty"`
	Headers   map[string]string `json:"headers,omitempty"`
}

type StatusLineConfig struct {
	Type    string `json:"type,omitempty"`
	Command string `json:"command,omitempty"`
}

type PermissionSettings struct {
	Allow []PermissionRuleSetting `json:"allow,omitempty"`
	Deny  []PermissionRuleSetting `json:"deny,omitempty"`
	Ask   []PermissionRuleSetting `json:"ask,omitempty"`
}

type PermissionRuleSetting struct {
	Tool    string `json:"tool"`
	Content string `json:"content,omitempty"`
}

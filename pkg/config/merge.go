package config

// SettingsSource identifies where settings came from
type SettingsSource string

const (
	SourceUser    SettingsSource = "user"    // ~/.claude/settings.json
	SourceProject SettingsSource = "project" // .claude/settings.json
	SourceLocal   SettingsSource = "local"   // .claude/settings.local.json
	SourceFlag    SettingsSource = "flag"    // CLI flags
	SourcePolicy  SettingsSource = "policy"  // Enterprise policy
)

// MergedSettings holds the final computed settings with source tracking
type MergedSettings struct {
	Settings SettingsJson
	Sources  map[string]SettingsSource // tracks which source each key came from
}

// LoadMergedSettings loads and merges settings from all sources in priority order:
// policy > flag > local > project > user (highest priority first)
func LoadMergedSettings() (*MergedSettings, error) {
	result := &MergedSettings{Sources: make(map[string]SettingsSource)}

	// 1. User settings (lowest priority)
	user, _ := LoadSettings()
	if user != nil {
		result.Settings = *user
		trackUserSources(user, result.Sources)
	}

	// 2. Project settings
	project, _ := LoadProjectSettings()
	if project != nil {
		mergeInto(&result.Settings, project, SourceProject, result.Sources)
	}

	// 3. Local settings (.claude/settings.local.json)
	local, _ := loadSettingsFrom(GetLocalSettingsPath())
	if local != nil {
		mergeInto(&result.Settings, local, SourceLocal, result.Sources)
	}

	return result, nil
}

// GetLocalSettingsPath returns the path to the local settings override file
func GetLocalSettingsPath() string {
	return ".oculus/settings.local.json"
}

// trackUserSources marks all non-zero fields in s as coming from SourceUser
func trackUserSources(s *SettingsJson, sources map[string]SettingsSource) {
	if s.DefaultMode != "" {
		sources["defaultMode"] = SourceUser
	}
	if s.Model != "" {
		sources["model"] = SourceUser
	}
	if s.Theme != "" {
		sources["theme"] = SourceUser
	}
	if len(s.AllowedTools) > 0 {
		sources["allowedTools"] = SourceUser
	}
	if len(s.DisallowedTools) > 0 {
		sources["disallowedTools"] = SourceUser
	}
	if s.Hooks != nil {
		sources["hooks"] = SourceUser
	}
	if s.MCPServers != nil {
		sources["mcpServers"] = SourceUser
	}
	if s.Env != nil {
		sources["env"] = SourceUser
	}
}

func mergeInto(base *SettingsJson, override *SettingsJson, source SettingsSource, sources map[string]SettingsSource) {
	if override.DefaultMode != "" {
		base.DefaultMode = override.DefaultMode
		sources["defaultMode"] = source
	}
	if override.Model != "" {
		base.Model = override.Model
		sources["model"] = source
	}
	if override.Theme != "" {
		base.Theme = override.Theme
		sources["theme"] = source
	}
	if len(override.AllowedTools) > 0 {
		base.AllowedTools = append(base.AllowedTools, override.AllowedTools...)
		sources["allowedTools"] = source
	}
	if len(override.DisallowedTools) > 0 {
		base.DisallowedTools = append(base.DisallowedTools, override.DisallowedTools...)
		sources["disallowedTools"] = source
	}
	if override.Hooks != nil {
		if base.Hooks == nil {
			base.Hooks = make(map[string][]HookConfig)
		}
		for k, v := range override.Hooks {
			base.Hooks[k] = append(base.Hooks[k], v...)
		}
		sources["hooks"] = source
	}
	if override.MCPServers != nil {
		if base.MCPServers == nil {
			base.MCPServers = make(map[string]MCPServerConfig)
		}
		for k, v := range override.MCPServers {
			base.MCPServers[k] = v
		}
		sources["mcpServers"] = source
	}
	if override.Env != nil {
		if base.Env == nil {
			base.Env = make(map[string]string)
		}
		for k, v := range override.Env {
			base.Env[k] = v
		}
		sources["env"] = source
	}
}

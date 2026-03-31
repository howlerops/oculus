package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// GetClaudeConfigDir returns the Claude config directory, respecting CLAUDE_CONFIG_DIR
func GetClaudeConfigDir() string {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".claude")
}

// GetSettingsPath returns the path to settings.json
func GetSettingsPath() string {
	return filepath.Join(GetClaudeConfigDir(), "settings.json")
}

// GetProjectSettingsPath returns the path to project-level settings
func GetProjectSettingsPath() string {
	return filepath.Join(".claude", "settings.json")
}

// LoadSettings reads and parses settings.json
func LoadSettings() (*SettingsJson, error) {
	return loadSettingsFrom(GetSettingsPath())
}

// LoadProjectSettings reads project-level settings
func LoadProjectSettings() (*SettingsJson, error) {
	return loadSettingsFrom(GetProjectSettingsPath())
}

func loadSettingsFrom(path string) (*SettingsJson, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &SettingsJson{}, nil
		}
		return nil, err
	}

	var settings SettingsJson
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, err
	}
	return &settings, nil
}

// MergeSettings combines user and project settings (project overrides user)
func MergeSettings(user, project *SettingsJson) *SettingsJson {
	merged := *user

	if project.DefaultMode != "" {
		merged.DefaultMode = project.DefaultMode
	}
	if project.Theme != "" {
		merged.Theme = project.Theme
	}
	if project.Model != "" {
		merged.Model = project.Model
	}
	if len(project.AllowedTools) > 0 {
		merged.AllowedTools = append(merged.AllowedTools, project.AllowedTools...)
	}
	if len(project.DisallowedTools) > 0 {
		merged.DisallowedTools = append(merged.DisallowedTools, project.DisallowedTools...)
	}
	if project.Hooks != nil {
		if merged.Hooks == nil {
			merged.Hooks = make(map[string][]HookConfig)
		}
		for k, v := range project.Hooks {
			merged.Hooks[k] = append(merged.Hooks[k], v...)
		}
	}
	if project.MCPServers != nil {
		if merged.MCPServers == nil {
			merged.MCPServers = make(map[string]MCPServerConfig)
		}
		for k, v := range project.MCPServers {
			merged.MCPServers[k] = v
		}
	}
	if project.Env != nil {
		if merged.Env == nil {
			merged.Env = make(map[string]string)
		}
		for k, v := range project.Env {
			merged.Env[k] = v
		}
	}

	return &merged
}

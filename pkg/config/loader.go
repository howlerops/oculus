package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetOculusDir returns the Claude config directory, respecting CLAUDE_CONFIG_DIR
func GetOculusDir() string {
	if dir := os.Getenv("OCULUS_CONFIG_DIR"); dir != "" {
		return dir
	}
	// Legacy fallback
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		return dir
	}
	home, _ := os.UserHomeDir()
	oculusDir := filepath.Join(home, ".oculus")
	// Auto-migrate: if ~/.oculus/ doesn't exist but ~/.claude/ does, copy
	if _, err := os.Stat(oculusDir); os.IsNotExist(err) {
		claudeDir := filepath.Join(home, ".claude")
		if _, err := os.Stat(claudeDir); err == nil {
			MigrateConfigDir(claudeDir, oculusDir)
		}
	}
	return oculusDir
}

// MigrateConfigDir copies contents from old config dir to new one
func MigrateConfigDir(src, dst string) {
	fmt.Fprintf(os.Stderr, "Migrating config from %s to %s...\n", src, dst)
	os.MkdirAll(dst, 0o755)
	entries, err := os.ReadDir(src)
	if err != nil {
		return
	}
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())
		if entry.IsDir() {
			// Recursive copy
			cpDir(srcPath, dstPath)
		} else {
			data, err := os.ReadFile(srcPath)
			if err == nil {
				os.WriteFile(dstPath, data, 0o644)
			}
		}
	}
	fmt.Fprintf(os.Stderr, "Migration complete. Your config is now at %s\n", dst)
}

func cpDir(src, dst string) {
	os.MkdirAll(dst, 0o755)
	entries, _ := os.ReadDir(src)
	for _, e := range entries {
		s := filepath.Join(src, e.Name())
		d := filepath.Join(dst, e.Name())
		if e.IsDir() {
			cpDir(s, d)
		} else {
			data, _ := os.ReadFile(s)
			os.WriteFile(d, data, 0o644)
		}
	}
}

// GetSettingsPath returns the path to settings.json
func GetSettingsPath() string {
	return filepath.Join(GetOculusDir(), "settings.json")
}

// GetProjectSettingsPath returns the path to project-level settings
func GetProjectSettingsPath() string {
	return filepath.Join(".oculus", "settings.json")
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

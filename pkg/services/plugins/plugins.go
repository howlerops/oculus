package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PluginManifest struct {
	Name        string                     `json:"name"`
	Version     string                     `json:"version"`
	Description string                     `json:"description"`
	Author      string                     `json:"author,omitempty"`
	Tools       []PluginToolDef            `json:"tools,omitempty"`
	Commands    []PluginCommandDef         `json:"commands,omitempty"`
	Hooks       map[string][]PluginHookDef `json:"hooks,omitempty"`
	Agents      []PluginAgentDef           `json:"agents,omitempty"`
}

type PluginToolDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"`
}

type PluginCommandDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type PluginHookDef struct {
	Matcher string `json:"matcher,omitempty"`
	Command string `json:"command"`
}

type PluginAgentDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Model       string `json:"model,omitempty"`
}

type LoadedPlugin struct {
	Manifest PluginManifest
	Path     string
	Enabled  bool
}

type Manager struct {
	plugins map[string]*LoadedPlugin
}

func NewManager() *Manager {
	return &Manager{plugins: make(map[string]*LoadedPlugin)}
}

// LoadFromDirectory scans a directory for plugin manifests
func (m *Manager) LoadFromDirectory(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		manifestPath := filepath.Join(dir, entry.Name(), "manifest.json")
		data, err := os.ReadFile(manifestPath)
		if err != nil {
			// Try plugin.json
			manifestPath = filepath.Join(dir, entry.Name(), "plugin.json")
			data, err = os.ReadFile(manifestPath)
			if err != nil {
				continue
			}
		}

		var manifest PluginManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			continue
		}

		m.plugins[manifest.Name] = &LoadedPlugin{
			Manifest: manifest,
			Path:     filepath.Join(dir, entry.Name()),
			Enabled:  true,
		}
	}

	return nil
}

// ApplySettings loads plugin enable/disable state from settings
func (m *Manager) ApplySettings(enabledPlugins map[string]bool) {
	for name, plugin := range m.plugins {
		if enabled, ok := enabledPlugins[name]; ok {
			plugin.Enabled = enabled
		}
	}
}

// GetEnabled returns all enabled plugins
func (m *Manager) GetEnabled() []*LoadedPlugin {
	var result []*LoadedPlugin
	for _, p := range m.plugins {
		if p.Enabled {
			result = append(result, p)
		}
	}
	return result
}

// Get returns a plugin by name
func (m *Manager) Get(name string) *LoadedPlugin {
	return m.plugins[name]
}

// List returns all plugin names
func (m *Manager) List() []string {
	var names []string
	for name := range m.plugins {
		names = append(names, name)
	}
	return names
}

// FormatList returns a human-readable plugin list
func (m *Manager) FormatList() string {
	var lines []string
	for name, p := range m.plugins {
		status := "enabled"
		if !p.Enabled {
			status = "disabled"
		}
		lines = append(lines, fmt.Sprintf("  %s v%s [%s] - %s", name, p.Manifest.Version, status, p.Manifest.Description))
	}
	if len(lines) == 0 {
		return "No plugins installed."
	}
	return fmt.Sprintf("Plugins (%d):\n%s", len(lines), strings.Join(lines, "\n"))
}

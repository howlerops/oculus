package plugins

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/howlerops/oculus/pkg/config"
)

// PluginManifest describes a plugin
type PluginManifest struct {
	Name        string               `json:"name"`
	Version     string               `json:"version"`
	Description string               `json:"description"`
	Author      string               `json:"author,omitempty"`
	Repository  string               `json:"repository,omitempty"`
	Skills      []SkillDef           `json:"skills,omitempty"`
	Tools       []ToolDef            `json:"tools,omitempty"`
	Commands    []CommandDef         `json:"commands,omitempty"`
	Hooks       map[string][]HookDef `json:"hooks,omitempty"`
	Agents      []AgentDef           `json:"agents,omitempty"`
}

// SkillDef defines a skill provided by a plugin
type SkillDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	File        string `json:"file"` // relative path to .md file
}

// ToolDef defines a tool provided by a plugin
type ToolDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command,omitempty"` // shell command for tool
}

// CommandDef defines a command provided by a plugin
type CommandDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	File        string `json:"file,omitempty"` // skill file to invoke
}

// HookDef defines a hook provided by a plugin
type HookDef struct {
	Matcher string `json:"matcher,omitempty"`
	Command string `json:"command"`
	Timeout int    `json:"timeout,omitempty"`
}

// AgentDef defines an agent provided by a plugin
type AgentDef struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Model       string `json:"model,omitempty"`
	File        string `json:"file,omitempty"` // .md system prompt
}

// InstalledPlugin is a plugin loaded from disk
type InstalledPlugin struct {
	Manifest PluginManifest
	Path     string
	Enabled  bool
	Source   string // "local", "marketplace", "git"
}

// Manager handles plugin lifecycle
type Manager struct {
	mu      sync.RWMutex
	plugins map[string]*InstalledPlugin
	dir     string
}

// NewManager creates a plugin manager
func NewManager() *Manager {
	dir := filepath.Join(config.GetOculusDir(), "plugins")
	os.MkdirAll(dir, 0o755)
	return &Manager{
		plugins: make(map[string]*InstalledPlugin),
		dir:     dir,
	}
}

// LoadAll discovers and loads all installed plugins
func (m *Manager) LoadAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load from ~/.oculus/plugins/
	m.loadFromDir(m.dir, "local")

	// Load from ~/.oculus/plugins/cache/ (marketplace installs)
	cacheDir := filepath.Join(m.dir, "cache")
	m.loadFromDir(cacheDir, "marketplace")

	// Load from ~/.claude/plugins/ (backward compat)
	home, _ := os.UserHomeDir()
	claudePlugins := filepath.Join(home, ".claude", "plugins")
	m.loadFromDir(claudePlugins, "local")

	return nil
}

func (m *Manager) loadFromDir(dir, source string) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		pluginDir := filepath.Join(dir, entry.Name())

		// Try manifest.json, then plugin.json
		manifest, err := loadManifest(pluginDir)
		if err != nil {
			continue
		}

		m.plugins[manifest.Name] = &InstalledPlugin{
			Manifest: *manifest,
			Path:     pluginDir,
			Enabled:  true,
			Source:   source,
		}
	}
}

func loadManifest(dir string) (*PluginManifest, error) {
	for _, name := range []string{"manifest.json", "plugin.json", "package.json"} {
		data, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			continue
		}
		var manifest PluginManifest
		if err := json.Unmarshal(data, &manifest); err != nil {
			continue
		}
		if manifest.Name != "" {
			return &manifest, nil
		}
	}
	return nil, fmt.Errorf("no manifest found in %s", dir)
}

// Get returns a plugin by name
func (m *Manager) Get(name string) *InstalledPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.plugins[name]
}

// List returns all installed plugins
func (m *Manager) List() []*InstalledPlugin {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var result []*InstalledPlugin
	for _, p := range m.plugins {
		result = append(result, p)
	}
	return result
}

// Enable enables a plugin by name
func (m *Manager) Enable(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.plugins[name]; ok {
		p.Enabled = true
		return true
	}
	return false
}

// Disable disables a plugin by name
func (m *Manager) Disable(name string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if p, ok := m.plugins[name]; ok {
		p.Enabled = false
		return true
	}
	return false
}

// Remove uninstalls a plugin
func (m *Manager) Remove(name string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	p, ok := m.plugins[name]
	if !ok {
		return fmt.Errorf("plugin %q not found", name)
	}
	os.RemoveAll(p.Path)
	delete(m.plugins, name)
	return nil
}

// FormatList returns a human-readable plugin list
func (m *Manager) FormatList() string {
	plugins := m.List()
	if len(plugins) == 0 {
		return "No plugins installed."
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Installed plugins (%d):\n", len(plugins)))
	for _, p := range plugins {
		status := "✓"
		if !p.Enabled {
			status = "✗"
		}
		sb.WriteString(fmt.Sprintf("  %s %s v%s - %s [%s]\n", status, p.Manifest.Name, p.Manifest.Version, p.Manifest.Description, p.Source))
		if len(p.Manifest.Skills) > 0 {
			sb.WriteString(fmt.Sprintf("    Skills: %d\n", len(p.Manifest.Skills)))
		}
		if len(p.Manifest.Agents) > 0 {
			sb.WriteString(fmt.Sprintf("    Agents: %d\n", len(p.Manifest.Agents)))
		}
	}
	return sb.String()
}

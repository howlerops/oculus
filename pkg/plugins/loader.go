package plugins

import (
	"os"
	"path/filepath"
	"strings"
)

// PluginSkill is a loaded skill from a plugin
type PluginSkill struct {
	Name        string
	Description string
	Content     string
	PluginName  string
}

// PluginAgent is a loaded agent from a plugin
type PluginAgent struct {
	Name         string
	Description  string
	Model        string
	SystemPrompt string
	PluginName   string
}

// PluginHook is a loaded hook from a plugin
type PluginHook struct {
	Matcher    string
	Command    string
	Timeout    int
	PluginName string
}

// PluginCommand is a loaded command from a plugin
type PluginCommand struct {
	Name        string
	Description string
	File        string
	PluginName  string
}

// LoadPluginSkills returns all skills from enabled plugins
func (m *Manager) LoadPluginSkills() []PluginSkill {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var skills []PluginSkill
	for _, p := range m.plugins {
		if !p.Enabled {
			continue
		}
		for _, sd := range p.Manifest.Skills {
			content, err := os.ReadFile(filepath.Join(p.Path, sd.File))
			if err != nil {
				continue
			}
			skills = append(skills, PluginSkill{
				Name:        p.Manifest.Name + ":" + sd.Name,
				Description: sd.Description,
				Content:     string(content),
				PluginName:  p.Manifest.Name,
			})
		}
	}
	return skills
}

// LoadPluginAgents returns all agents from enabled plugins
func (m *Manager) LoadPluginAgents() []PluginAgent {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var agents []PluginAgent
	for _, p := range m.plugins {
		if !p.Enabled {
			continue
		}
		for _, ad := range p.Manifest.Agents {
			systemPrompt := ""
			if ad.File != "" {
				data, err := os.ReadFile(filepath.Join(p.Path, ad.File))
				if err == nil {
					systemPrompt = string(data)
				}
			}
			agents = append(agents, PluginAgent{
				Name:         p.Manifest.Name + ":" + ad.Name,
				Description:  ad.Description,
				Model:        ad.Model,
				SystemPrompt: systemPrompt,
				PluginName:   p.Manifest.Name,
			})
		}
	}
	return agents
}

// LoadPluginHooks returns all hooks from enabled plugins
func (m *Manager) LoadPluginHooks() map[string][]PluginHook {
	m.mu.RLock()
	defer m.mu.RUnlock()

	hooks := make(map[string][]PluginHook)
	for _, p := range m.plugins {
		if !p.Enabled {
			continue
		}
		for event, defs := range p.Manifest.Hooks {
			for _, hd := range defs {
				hooks[event] = append(hooks[event], PluginHook{
					Matcher:    hd.Matcher,
					Command:    hd.Command,
					Timeout:    hd.Timeout,
					PluginName: p.Manifest.Name,
				})
			}
		}
	}
	return hooks
}

// LoadPluginCommands returns all commands from enabled plugins
func (m *Manager) LoadPluginCommands() []PluginCommand {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var commands []PluginCommand
	for _, p := range m.plugins {
		if !p.Enabled {
			continue
		}
		for _, cd := range p.Manifest.Commands {
			commands = append(commands, PluginCommand{
				Name:        cd.Name,
				Description: cd.Description,
				File:        filepath.Join(p.Path, cd.File),
				PluginName:  p.Manifest.Name,
			})
		}
	}
	return commands
}

// FindSkill searches all plugins for a skill by name
func (m *Manager) FindSkill(name string) *PluginSkill {
	skills := m.LoadPluginSkills()
	for _, s := range skills {
		if s.Name == name || strings.HasSuffix(s.Name, ":"+name) {
			return &s
		}
	}
	return nil
}

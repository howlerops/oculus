package agent

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbeck018/claude-go/pkg/config"
)

// AgentDefinition describes a custom or built-in agent
type AgentDefinition struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Model        string   `json:"model,omitempty"`
	Tools        []string `json:"tools,omitempty"`
	SystemPrompt string   `json:"system_prompt,omitempty"`
	IsBuiltIn    bool     `json:"is_built_in"`
	FilePath     string   `json:"file_path,omitempty"`
}

// BuiltInAgents returns the default agent definitions
func BuiltInAgents() []AgentDefinition {
	return []AgentDefinition{
		{Name: "general-purpose", Description: "General-purpose agent for research and multi-step tasks", IsBuiltIn: true},
		{Name: "Explore", Description: "Fast agent for exploring codebases", IsBuiltIn: true},
		{Name: "Plan", Description: "Software architect for designing implementation plans", IsBuiltIn: true},
		{Name: "code-reviewer", Description: "Code review specialist", IsBuiltIn: true},
		{Name: "debugger", Description: "Debugging specialist for errors and test failures", IsBuiltIn: true},
		{Name: "executor", Description: "Task executor for implementation work", IsBuiltIn: true},
		{Name: "writer", Description: "Technical documentation writer", IsBuiltIn: true},
		{Name: "architect", Description: "Strategic architecture advisor", IsBuiltIn: true},
		{Name: "verifier", Description: "Verification and evidence-based completion checks", IsBuiltIn: true},
	}
}

// LoadAgentDefinitions loads custom agents from disk + built-ins
func LoadAgentDefinitions() []AgentDefinition {
	agents := BuiltInAgents()

	// Search paths for custom agents
	searchDirs := []string{
		filepath.Join(".claude", "agents"),
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		searchDirs = append(searchDirs, filepath.Join(config.GetClaudeConfigDir(), "agents"))
	}

	for _, dir := range searchDirs {
		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if !strings.HasSuffix(name, ".md") && !strings.HasSuffix(name, ".json") {
				continue
			}

			path := filepath.Join(dir, name)
			agentName := strings.TrimSuffix(strings.TrimSuffix(name, ".md"), ".json")

			if strings.HasSuffix(name, ".json") {
				data, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				var def AgentDefinition
				if json.Unmarshal(data, &def) == nil {
					def.FilePath = path
					agents = append(agents, def)
				}
			} else {
				// Markdown agent - content becomes system prompt
				content, err := os.ReadFile(path)
				if err != nil {
					continue
				}
				agents = append(agents, AgentDefinition{
					Name:         agentName,
					Description:  "Custom agent: " + agentName,
					SystemPrompt: string(content),
					FilePath:     path,
				})
			}
		}
	}

	return agents
}

// FindAgent looks up an agent by name
func FindAgent(name string) *AgentDefinition {
	for _, a := range LoadAgentDefinitions() {
		if strings.EqualFold(a.Name, name) {
			return &a
		}
	}
	return nil
}

// IsBuiltInAgent checks if an agent name is a built-in
func IsBuiltInAgent(name string) bool {
	for _, a := range BuiltInAgents() {
		if strings.EqualFold(a.Name, name) {
			return true
		}
	}
	return false
}

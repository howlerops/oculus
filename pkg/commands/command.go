package commands

import (
	"context"
	"fmt"
	"sort"
	"strings"
)

// Command represents a slash command
type Command struct {
	Name        string
	Aliases     []string
	Description string
	IsHidden    bool
	// Run executes the command. Args is the text after the command name.
	// Returns the output string and whether to continue the conversation.
	Run func(ctx context.Context, args string) (output string, continueConversation bool, err error)
}

// Registry holds all registered commands
type Registry struct {
	commands map[string]*Command
}

// NewRegistry creates an empty command registry
func NewRegistry() *Registry {
	return &Registry{commands: make(map[string]*Command)}
}

// Register adds a command to the registry
func (r *Registry) Register(cmd *Command) {
	r.commands[cmd.Name] = cmd
	for _, alias := range cmd.Aliases {
		r.commands[alias] = cmd
	}
}

// Find looks up a command by name or prefix
func (r *Registry) Find(name string) *Command {
	name = strings.TrimPrefix(name, "/")
	if cmd, ok := r.commands[name]; ok {
		return cmd
	}
	// Prefix match
	var match *Command
	var matchCount int
	for k, cmd := range r.commands {
		if strings.HasPrefix(k, name) {
			match = cmd
			matchCount++
		}
	}
	if matchCount == 1 {
		return match
	}
	return nil
}

// List returns all non-hidden commands sorted by name
func (r *Registry) List() []*Command {
	seen := make(map[string]bool)
	var result []*Command
	for _, cmd := range r.commands {
		if !seen[cmd.Name] && !cmd.IsHidden {
			seen[cmd.Name] = true
			result = append(result, cmd)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Name < result[j].Name
	})
	return result
}

// IsSlashCommand checks if text starts with /
func IsSlashCommand(text string) bool {
	return strings.HasPrefix(strings.TrimSpace(text), "/")
}

// ParseCommand splits "/cmd args" into command name and args
func ParseCommand(text string) (name, args string) {
	text = strings.TrimSpace(text)
	text = strings.TrimPrefix(text, "/")
	parts := strings.SplitN(text, " ", 2)
	name = parts[0]
	if len(parts) > 1 {
		args = parts[1]
	}
	return
}

// FormatHelp returns a formatted help string for all commands
func (r *Registry) FormatHelp() string {
	cmds := r.List()
	var sb strings.Builder
	sb.WriteString("Available commands:\n\n")
	for _, cmd := range cmds {
		sb.WriteString(fmt.Sprintf("  /%s - %s\n", cmd.Name, cmd.Description))
	}
	return sb.String()
}

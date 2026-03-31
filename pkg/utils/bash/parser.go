package bash

import "strings"

// ParsedCommand represents a parsed shell command
type ParsedCommand struct {
	Command    string
	Args       []string
	Operator   string // "", "&&", "||", ";", "|"
	IsSubshell bool
}

// SplitCommands splits a command string into individual commands
func SplitCommands(input string) []ParsedCommand {
	var commands []ParsedCommand
	var current strings.Builder
	inSingleQuote := false
	inDoubleQuote := false
	inBacktick := false
	parenDepth := 0

	runes := []rune(input)
	for i := 0; i < len(runes); i++ {
		ch := runes[i]

		switch {
		case ch == '\'' && !inDoubleQuote && !inBacktick:
			inSingleQuote = !inSingleQuote
			current.WriteRune(ch)
		case ch == '"' && !inSingleQuote && !inBacktick:
			inDoubleQuote = !inDoubleQuote
			current.WriteRune(ch)
		case ch == '`' && !inSingleQuote:
			inBacktick = !inBacktick
			current.WriteRune(ch)
		case ch == '(' && !inSingleQuote && !inDoubleQuote:
			parenDepth++
			current.WriteRune(ch)
		case ch == ')' && !inSingleQuote && !inDoubleQuote:
			parenDepth--
			current.WriteRune(ch)
		case (ch == '&' || ch == ';') && !inSingleQuote && !inDoubleQuote && !inBacktick && parenDepth == 0:
			op := string(ch)
			if i+1 < len(runes) && runes[i+1] == ch && ch == '&' {
				op = "&&"
				i++
			}
			cmd := strings.TrimSpace(current.String())
			if cmd != "" {
				commands = append(commands, ParsedCommand{Command: cmd, Operator: op})
			}
			current.Reset()
		case ch == '|' && !inSingleQuote && !inDoubleQuote && !inBacktick && parenDepth == 0:
			op := "|"
			if i+1 < len(runes) && runes[i+1] == '|' {
				op = "||"
				i++
			}
			cmd := strings.TrimSpace(current.String())
			if cmd != "" {
				commands = append(commands, ParsedCommand{Command: cmd, Operator: op})
			}
			current.Reset()
		default:
			current.WriteRune(ch)
		}
	}

	cmd := strings.TrimSpace(current.String())
	if cmd != "" {
		commands = append(commands, ParsedCommand{Command: cmd})
	}

	// Extract base command names and args
	for i := range commands {
		parts := strings.Fields(commands[i].Command)
		if len(parts) > 0 {
			// Strip variable assignments
			base := parts[0]
			for strings.Contains(base, "=") && len(parts) > 1 {
				parts = parts[1:]
				base = parts[0]
			}
			commands[i].Args = parts[1:]
			if strings.HasPrefix(base, "(") || strings.HasPrefix(base, "$(") {
				commands[i].IsSubshell = true
			}
		}
	}

	return commands
}

// GetBaseCommand extracts the command name from a full command string
func GetBaseCommand(command string) string {
	command = strings.TrimSpace(command)
	// Strip env vars
	for strings.Contains(command, "=") {
		parts := strings.SplitN(command, " ", 2)
		if len(parts) < 2 || !strings.Contains(parts[0], "=") {
			break
		}
		command = strings.TrimSpace(parts[1])
	}
	parts := strings.Fields(command)
	if len(parts) == 0 {
		return ""
	}
	return parts[0]
}

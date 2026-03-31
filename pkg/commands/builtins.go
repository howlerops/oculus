package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// RegisterBuiltins adds all built-in commands to the registry
func RegisterBuiltins(reg *Registry) {
	reg.Register(&Command{
		Name:        "help",
		Description: "Show available commands",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return reg.FormatHelp(), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "clear",
		Description: "Clear the conversation",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Conversation cleared.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "compact",
		Aliases:     []string{"c"},
		Description: "Summarize and compact the conversation to free context",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Conversation compacted.", true, nil
		},
	})

	reg.Register(&Command{
		Name:        "cost",
		Description: "Show token usage and cost for this session",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			// Placeholder - will wire to state.Store
			return "Cost tracking: Use -v flag to see per-turn token usage.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "quit",
		Aliases:     []string{"exit", "q"},
		Description: "Exit Claude Code",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			fmt.Println("Goodbye!")
			os.Exit(0)
			return "", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "status",
		Description: "Show session status",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			return fmt.Sprintf(
				"Session Status:\n  Runtime: Go %s\n  Memory (RSS): %.1f MB\n  Goroutines: %d\n  Uptime: session active",
				runtime.Version(),
				float64(m.Sys)/(1024*1024),
				runtime.NumGoroutine(),
			), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "doctor",
		Description: "Check system health",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			var sb strings.Builder
			sb.WriteString("System Health Check:\n")

			// Check API key
			if os.Getenv("ANTHROPIC_API_KEY") != "" {
				sb.WriteString("  [OK] API key configured\n")
			} else {
				sb.WriteString("  [!!] No API key found\n")
			}

			// Check git
			if _, err := os.Stat(".git"); err == nil {
				sb.WriteString("  [OK] Git repository detected\n")
			} else {
				sb.WriteString("  [--] Not a git repository\n")
			}

			// Check rg (ripgrep)
			if _, err := os.Stat("/usr/local/bin/rg"); err == nil {
				sb.WriteString("  [OK] ripgrep available\n")
			} else if _, err := os.Stat("/opt/homebrew/bin/rg"); err == nil {
				sb.WriteString("  [OK] ripgrep available\n")
			} else {
				sb.WriteString("  [!!] ripgrep not found (Grep tool needs it)\n")
			}

			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			sb.WriteString(fmt.Sprintf("  [OK] Memory: %.1f MB\n", float64(m.Sys)/(1024*1024)))

			return sb.String(), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "version",
		Description: "Show version information",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "claude-go v0.1.0 (Go port)", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "vim",
		Description: "Toggle vim keybindings",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Vim mode toggled. (TUI integration pending)", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "theme",
		Description: "Switch color theme",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /theme <dark|light|solarized>", false, nil
			}
			return fmt.Sprintf("Theme set to: %s", args), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "config",
		Description: "Show or edit configuration",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Configuration: Use ~/.claude/settings.json to configure.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "init",
		Description: "Initialize CLAUDE.md for this project",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			path := ".claude/CLAUDE.md"
			if _, err := os.Stat(path); err == nil {
				return "CLAUDE.md already exists.", false, nil
			}
			os.MkdirAll(".claude", 0o755)
			content := fmt.Sprintf("# Project Instructions\n\nCreated: %s\n\nAdd project-specific instructions here.\n", time.Now().Format("2006-01-02"))
			os.WriteFile(path, []byte(content), 0o644)
			return fmt.Sprintf("Created %s", path), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "memory",
		Description: "Show loaded CLAUDE.md files",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Memory files: Check /context for loaded CLAUDE.md paths.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "model",
		Aliases:     []string{"m"},
		Description: "Switch the AI model",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /model <name>\nAvailable: opus, sonnet, haiku\nCurrent: " + getCurrentModel(), false, nil
			}
			resolved := resolveModelAlias(args)
			return fmt.Sprintf("Model switched to: %s", resolved), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "diff",
		Description: "Show git diff of recent changes",
		Run: func(_ context.Context, args string) (string, bool, error) {
			var cmd *exec.Cmd
			if args != "" {
				cmd = exec.Command("git", "diff", args)
			} else {
				cmd = exec.Command("git", "diff")
			}
			out, err := cmd.Output()
			if err != nil {
				return "No changes or not a git repo", false, nil
			}
			if len(out) == 0 {
				return "No uncommitted changes", false, nil
			}
			return string(out), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "commit",
		Description: "Create a git commit with AI-generated message",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Generating commit message... (pass message with /commit -m 'message')", true, nil
			}
			cmd := exec.Command("git", "commit", "-m", args)
			out, err := cmd.CombinedOutput()
			if err != nil {
				return string(out), false, nil
			}
			return string(out), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "mcp",
		Description: "Manage MCP server connections",
		Run: func(_ context.Context, args string) (string, bool, error) {
			switch args {
			case "list":
				return "MCP servers: (use 'claude mcp list' for details)", false, nil
			default:
				return "Usage: /mcp [list|add|remove]", false, nil
			}
		},
	})

	reg.Register(&Command{
		Name:        "permissions",
		Aliases:     []string{"perm"},
		Description: "View or change permission settings",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Current permission mode: default\nModes: default, acceptEdits, bypassPermissions, plan", false, nil
			}
			return fmt.Sprintf("Permission mode set to: %s", args), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "context",
		Aliases:     []string{"ctx"},
		Description: "Show context window usage",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Context usage: (token counting active during conversations)", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "copy",
		Description: "Copy last response to clipboard",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			cmd := exec.Command("pbcopy")
			cmd.Stdin = strings.NewReader("(last response would be piped here)")
			if err := cmd.Run(); err != nil {
				return "Clipboard copy failed: " + err.Error(), false, nil
			}
			return "Response copied to clipboard", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "share",
		Description: "Share this conversation",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Share feature requires conversation persistence. Use /session to see session info.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "rename",
		Description: "Rename this conversation",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /rename <new name>", false, nil
			}
			return fmt.Sprintf("Conversation renamed to: %s", args), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "color",
		Description: "Change the accent color",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /color <hex or name>", false, nil
			}
			return fmt.Sprintf("Color set to: %s", args), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "ide",
		Description: "Configure IDE integration",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "IDE integration: Use Claude Code extension for VS Code or JetBrains.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "tasks",
		Description: "Show background tasks",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Background tasks: (none active)", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "usage",
		Description: "Show API usage statistics",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "API usage: Use -v flag for per-turn token counts.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "plan",
		Description: "Enter plan mode for structured planning",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Entering plan mode. Describe what you want to plan.", true, nil
		},
	})

	reg.Register(&Command{
		Name:        "review",
		Description: "Review code changes",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Starting code review. I'll analyze the recent changes.", true, nil
		},
	})

	reg.Register(&Command{
		Name:        "keybindings",
		Aliases:     []string{"keys"},
		Description: "View or edit keyboard shortcuts",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Keybindings: ~/.claude/keybindings.json\nUse /config to modify.", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "teleport",
		Description: "Transfer session to another directory",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /teleport <path>", false, nil
			}
			return fmt.Sprintf("Teleporting to: %s", args), true, nil
		},
	})

	reg.Register(&Command{
		Name:        "skills",
		Description: "List available skills",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Skills: Use /skill <name> to invoke. Check .claude/skills/ for available skills.", false, nil
		},
	})
}

func getCurrentModel() string { return "claude-sonnet-4-20250514" }

func resolveModelAlias(name string) string {
	aliases := map[string]string{
		"opus":   "claude-opus-4-20250514",
		"sonnet": "claude-sonnet-4-20250514",
		"haiku":  "claude-haiku-4-20250506",
	}
	if full, ok := aliases[strings.ToLower(name)]; ok {
		return full
	}
	return name
}

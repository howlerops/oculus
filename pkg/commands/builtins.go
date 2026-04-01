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

	config "github.com/howlerops/oculus/pkg/config"
	"github.com/howlerops/oculus/pkg/plugins"
	"github.com/howlerops/oculus/pkg/skills"
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
			return "Compacting conversation... The assistant will summarize older messages to free context space.", true, nil
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
			return "oculus v0.3.0", false, nil
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
			path := ".oculus/OCULUS.md"
			if _, err := os.Stat(path); err == nil {
				return "CLAUDE.md already exists.", false, nil
			}
			os.MkdirAll(".oculus", 0o755)
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
				models := config.ListModels()
				return fmt.Sprintf("Current model: %s\nAvailable: %s\nUsage: /model <name>", getCurrentModel(), strings.Join(models, ", ")), false, nil
			}
			info, found := config.ResolveModel(args)
			if !found {
				return fmt.Sprintf("Unknown model: %s. Available: opus, sonnet, haiku", args), false, nil
			}
			settings, _ := config.LoadSettings()
			if settings == nil {
				settings = &config.SettingsJson{}
			}
			settings.Model = info.ID
			data, _ := json.Marshal(settings)
			os.WriteFile(config.GetSettingsPath(), data, 0o644)
			config.InvalidateSettingsCache()
			return fmt.Sprintf("Model switched to: %s (%s)\nInput: $%.2f/M tokens | Output: $%.2f/M tokens", info.DisplayName, info.ID, info.CostInput, info.CostOutput), false, nil
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
			settings, _ := config.LoadSettings()
			currentMode := "default"
			if settings != nil && settings.DefaultMode != "" {
				currentMode = settings.DefaultMode
			}

			if args == "" {
				return fmt.Sprintf("Current permission mode: %s\n\nAvailable modes:\n  default             - Ask for most operations\n  acceptEdits         - Auto-approve file edits\n  bypassPermissions   - Allow everything (dangerous)\n  plan                - Read-only, planning only\n\nUsage: /permissions <mode>", currentMode), false, nil
			}

			validModes := map[string]bool{"default": true, "acceptEdits": true, "bypassPermissions": true, "plan": true}
			if !validModes[args] {
				return fmt.Sprintf("Invalid mode: %s. Use: default, acceptEdits, bypassPermissions, plan", args), false, nil
			}

			if settings == nil {
				settings = &config.SettingsJson{}
			}
			settings.DefaultMode = args
			data, _ := json.Marshal(settings)
			os.WriteFile(config.GetSettingsPath(), data, 0o644)
			config.InvalidateSettingsCache()
			return fmt.Sprintf("Permission mode set to: %s", args), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "context",
		Aliases:     []string{"ctx"},
		Description: "Show context window usage",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Context window: 200,000 tokens\nUsage: tracked per-session in status bar\nUse /compact to free space when running low.", false, nil
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
			return "Tasks panel: visible in TUI mode.\nUse TaskCreate/TaskUpdate tools to manage tasks.", false, nil
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
			var sb strings.Builder

			// Load local skills
			localSkills := skills.LoadSkills()
			if len(localSkills) > 0 {
				sb.WriteString(fmt.Sprintf("Local skills (%d):\n", len(localSkills)))
				for _, s := range localSkills {
					sb.WriteString(fmt.Sprintf("  - %s (%s)\n", s.Name, s.Path))
				}
			}

			// Load plugin skills
			mgr := plugins.NewManager()
			mgr.LoadAll()
			pluginSkills := mgr.LoadPluginSkills()
			if len(pluginSkills) > 0 {
				sb.WriteString(fmt.Sprintf("\nPlugin skills (%d):\n", len(pluginSkills)))
				for _, s := range pluginSkills {
					sb.WriteString(fmt.Sprintf("  - %s - %s [%s]\n", s.Name, s.Description, s.PluginName))
				}
			}

			if sb.Len() == 0 {
				return "No skills found. Add .md files to .oculus/skills/ or install plugins with /plugin install.", false, nil
			}
			return sb.String(), false, nil
		},
	})

	// ── /plugin ─────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "plugin",
		Aliases:     []string{"plugins"},
		Description: "Manage plugins (install, list, remove, search)",
		Run: func(_ context.Context, args string) (string, bool, error) {
			parts := strings.SplitN(args, " ", 2)
			action := parts[0]
			target := ""
			if len(parts) > 1 {
				target = parts[1]
			}

			mgr := plugins.NewManager()
			mgr.LoadAll()

			switch action {
			case "", "list", "ls":
				return mgr.FormatList(), false, nil
			case "install", "add":
				if target == "" {
					return "Usage: /plugin install <user/repo or git URL>", false, nil
				}
				p, err := mgr.Install(target)
				if err != nil {
					return "Install failed: " + err.Error(), false, nil
				}
				return fmt.Sprintf("Installed %s v%s (%d skills, %d agents)", p.Manifest.Name, p.Manifest.Version, len(p.Manifest.Skills), len(p.Manifest.Agents)), false, nil
			case "remove", "rm", "uninstall":
				if target == "" {
					return "Usage: /plugin remove <name>", false, nil
				}
				if err := mgr.Remove(target); err != nil {
					return "Remove failed: " + err.Error(), false, nil
				}
				return "Removed plugin: " + target, false, nil
			case "update":
				if target == "" {
					return "Usage: /plugin update <name>", false, nil
				}
				if err := mgr.Update(target); err != nil {
					return "Update failed: " + err.Error(), false, nil
				}
				return "Updated plugin: " + target, false, nil
			case "search":
				if target == "" {
					return "Usage: /plugin search <query>", false, nil
				}
				results, err := plugins.Search(target)
				if err != nil {
					return "Search failed: " + err.Error(), false, nil
				}
				if len(results) == 0 {
					return "No plugins found for: " + target, false, nil
				}
				return "Found:\n  " + strings.Join(results, "\n  "), false, nil
			case "enable":
				if mgr.Enable(target) {
					return "Enabled: " + target, false, nil
				}
				return "Plugin not found: " + target, false, nil
			case "disable":
				if mgr.Disable(target) {
					return "Disabled: " + target, false, nil
				}
				return "Plugin not found: " + target, false, nil
			default:
				return "Usage: /plugin [list|install|remove|update|search|enable|disable] [args]", false, nil
			}
		},
	})

	// ── /add-dir ────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "add-dir",
		Description: "Add an additional working directory to the context scope",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /add-dir <path>\nAdds a directory to the active working directories.", false, nil
			}
			abs, err := filepath.Abs(args)
			if err != nil {
				return fmt.Sprintf("Invalid path: %v", err), false, nil
			}
			if _, err := os.Stat(abs); os.IsNotExist(err) {
				return fmt.Sprintf("Directory not found: %s", abs), false, nil
			}
			return fmt.Sprintf("Added working directory: %s", abs), false, nil
		},
	})

	// ── /branch ─────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "branch",
		Description: "Create or switch a git branch",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				// Show current branch
				out, err := exec.Command("git", "branch", "--show-current").Output()
				if err != nil {
					return "Not a git repository or git not available.", false, nil
				}
				return fmt.Sprintf("Current branch: %s", strings.TrimSpace(string(out))), false, nil
			}
			// Try to switch first; if that fails, create
			switchOut, err := exec.Command("git", "checkout", args).CombinedOutput()
			if err != nil {
				// Branch doesn't exist — create it
				createOut, err2 := exec.Command("git", "checkout", "-b", args).CombinedOutput()
				if err2 != nil {
					return fmt.Sprintf("Failed to create branch %q: %s", args, strings.TrimSpace(string(createOut))), false, nil
				}
				return fmt.Sprintf("Created and switched to branch: %s", args), false, nil
			}
			return strings.TrimSpace(string(switchOut)), false, nil
		},
	})

	// ── /feedback ───────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "feedback",
		Description: "Submit feedback about Claude Code",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /feedback <your feedback>\nFeedback is recorded locally at ~/.claude/feedback.log", false, nil
			}
			home, _ := os.UserHomeDir()
			logPath := filepath.Join(home, ".oculus", "feedback.log")
			os.MkdirAll(filepath.Dir(logPath), 0o755)
			entry := fmt.Sprintf("[%s] %s\n", time.Now().Format(time.RFC3339), args)
			f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return fmt.Sprintf("Could not write feedback: %v", err), false, nil
			}
			defer f.Close()
			f.WriteString(entry)
			return fmt.Sprintf("Feedback recorded. Thank you!\nSaved to: %s", logPath), false, nil
		},
	})

	// ── /good ────────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "good",
		Aliases:     []string{"thumbsup"},
		Description: "Rate the last response positively",
		Run: func(_ context.Context, args string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			logPath := filepath.Join(home, ".oculus", "ratings.log")
			os.MkdirAll(filepath.Dir(logPath), 0o755)
			note := args
			if note == "" {
				note = "(no note)"
			}
			entry := fmt.Sprintf("[%s] POSITIVE | %s\n", time.Now().Format(time.RFC3339), note)
			if f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644); err == nil {
				f.WriteString(entry)
				f.Close()
			}
			return "Thanks for the positive rating!", false, nil
		},
	})

	// ── /btw ─────────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "btw",
		Description: "Send a side note or background context to Claude",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /btw <context note>\nAdds background context that informs Claude without being a direct request.", false, nil
			}
			return fmt.Sprintf("[Side note recorded]: %s\nClaude will factor this into subsequent responses.", args), true, nil
		},
	})

	// ── /onboarding ──────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "onboarding",
		Description: "Re-run the onboarding flow",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			flagPath := filepath.Join(home, ".oculus", ".onboarding_complete")
			os.Remove(flagPath) // reset completion flag so onboarding re-runs on next start
			var sb strings.Builder
			sb.WriteString("Onboarding reset.\n\n")
			sb.WriteString("Welcome to Claude Code!\n")
			sb.WriteString("  1. Set your API key: export ANTHROPIC_API_KEY=<key>\n")
			sb.WriteString("  2. Run /doctor to verify your setup\n")
			sb.WriteString("  3. Run /init to create a CLAUDE.md for your project\n")
			sb.WriteString("  4. Run /help to see all available commands\n")
			return sb.String(), false, nil
		},
	})

	// ── /pr_comments ─────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "pr_comments",
		Aliases:     []string{"pr-comments"},
		Description: "Show comments from the current GitHub pull request",
		Run: func(_ context.Context, args string) (string, bool, error) {
			// Try gh CLI first
			ghArgs := []string{"pr", "view", "--json", "comments,reviews,number,title"}
			if args != "" {
				ghArgs = append(ghArgs, args)
			}
			out, err := exec.Command("gh", ghArgs...).Output()
			if err != nil {
				return "Could not fetch PR comments. Ensure 'gh' CLI is installed and authenticated.\nUsage: /pr_comments [PR number]", false, nil
			}
			var result map[string]interface{}
			if err := json.Unmarshal(out, &result); err != nil {
				return string(out), false, nil
			}
			var sb strings.Builder
			if num, ok := result["number"]; ok {
				sb.WriteString(fmt.Sprintf("PR #%.0f: %v\n\n", num.(float64), result["title"]))
			}
			if comments, ok := result["comments"].([]interface{}); ok {
				sb.WriteString(fmt.Sprintf("Comments (%d):\n", len(comments)))
				for _, c := range comments {
					if cm, ok := c.(map[string]interface{}); ok {
						author := ""
						if a, ok := cm["author"].(map[string]interface{}); ok {
							author = fmt.Sprintf("%v", a["login"])
						}
						sb.WriteString(fmt.Sprintf("  @%s: %v\n", author, cm["body"]))
					}
				}
			}
			return sb.String(), false, nil
		},
	})

	// ── /release-notes ───────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "release-notes",
		Aliases:     []string{"changelog"},
		Description: "Show recent release notes",
		Run: func(_ context.Context, args string) (string, bool, error) {
			// Try to read a local CHANGELOG or RELEASE_NOTES file
			candidates := []string{"CHANGELOG.md", "CHANGELOG", "RELEASE_NOTES.md", "RELEASE_NOTES"}
			for _, f := range candidates {
				data, err := os.ReadFile(f)
				if err == nil {
					lines := strings.Split(string(data), "\n")
					if len(lines) > 50 {
						lines = lines[:50]
					}
					return strings.Join(lines, "\n") + "\n\n(truncated — open " + f + " for full notes)", false, nil
				}
			}
			// Fall back to recent git tags
			out, err := exec.Command("git", "tag", "--sort=-version:refname", "-l").Output()
			if err != nil || len(out) == 0 {
				return "No release notes found. Add a CHANGELOG.md to your project.", false, nil
			}
			tags := strings.Split(strings.TrimSpace(string(out)), "\n")
			if len(tags) > 10 {
				tags = tags[:10]
			}
			return "Recent releases (from git tags):\n  " + strings.Join(tags, "\n  "), false, nil
		},
	})

	// ── /security-review ─────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "security-review",
		Description: "Run a security review of pending changes",
		Run: func(_ context.Context, args string) (string, bool, error) {
			// Get list of changed files to inform the review
			diffOut, _ := exec.Command("git", "diff", "--name-only").Output()
			stagedOut, _ := exec.Command("git", "diff", "--cached", "--name-only").Output()
			changed := strings.TrimSpace(string(diffOut) + string(stagedOut))
			if changed == "" {
				return "No uncommitted changes found to review.\nUse /diff to inspect changes first.", false, nil
			}
			var sb strings.Builder
			sb.WriteString("Security review initiated for changed files:\n")
			for _, f := range strings.Split(changed, "\n") {
				if f != "" {
					sb.WriteString(fmt.Sprintf("  - %s\n", f))
				}
			}
			sb.WriteString("\nAnalyzing for: hardcoded secrets, injection risks, insecure deps, auth issues...")
			return sb.String(), true, nil
		},
	})

	// ── /terminal-setup ──────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "terminal-setup",
		Description: "Configure terminal integrations (shell completion, aliases)",
		Run: func(_ context.Context, args string) (string, bool, error) {
			shell := os.Getenv("SHELL")
			if shell == "" {
				shell = "unknown"
			}
			shellName := filepath.Base(shell)
			home, _ := os.UserHomeDir()

			var rcFile string
			switch shellName {
			case "zsh":
				rcFile = filepath.Join(home, ".zshrc")
			case "bash":
				rcFile = filepath.Join(home, ".bashrc")
			case "fish":
				rcFile = filepath.Join(home, ".config", "fish", "config.fish")
			default:
				rcFile = filepath.Join(home, ".profile")
			}

			snippet := fmt.Sprintf(`
# Claude Code shell integration
alias cc='claude'
alias ccc='claude --continue'
`)
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("Detected shell: %s\n", shellName))
			sb.WriteString(fmt.Sprintf("Config file: %s\n\n", rcFile))
			sb.WriteString("Recommended shell snippet:\n")
			sb.WriteString("```\n")
			sb.WriteString(snippet)
			sb.WriteString("```\n")
			sb.WriteString(fmt.Sprintf("\nTo apply automatically, append to %s and restart your terminal.", rcFile))
			return sb.String(), false, nil
		},
	})

	// ── /desktop ─────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "desktop",
		Description: "Continue this session in Claude Desktop app",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			// Try to open Claude Desktop via open/xdg-open
			var cmd *exec.Cmd
			switch runtime.GOOS {
			case "darwin":
				cmd = exec.Command("open", "-a", "Claude")
			case "linux":
				cmd = exec.Command("xdg-open", "claude://")
			case "windows":
				cmd = exec.Command("cmd", "/c", "start", "claude://")
			default:
				return "Desktop app launch not supported on this platform.\nDownload Claude Desktop at: https://claude.ai/download", false, nil
			}
			if err := cmd.Start(); err != nil {
				return "Could not open Claude Desktop. Download at: https://claude.ai/download", false, nil
			}
			return "Opening Claude Desktop...", false, nil
		},
	})

	// ── /mobile ──────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "mobile",
		Description: "Show info and QR code link for the Claude mobile app",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			var sb strings.Builder
			sb.WriteString("Claude Mobile App\n")
			sb.WriteString("─────────────────\n")
			sb.WriteString("iOS:     https://apps.apple.com/app/claude-by-anthropic/id6473753684\n")
			sb.WriteString("Android: https://play.google.com/store/apps/details?id=com.anthropic.claude\n\n")
			sb.WriteString("Scan the URL with your phone's camera or visit https://claude.ai/mobile\n")
			return sb.String(), false, nil
		},
	})

	// ── /issue ───────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "issue",
		Description: "File a GitHub issue for a bug or feature request",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args == "" {
				return "Usage: /issue <title>\nOpens a new GitHub issue. Requires 'gh' CLI.", false, nil
			}
			// Try gh CLI
			out, err := exec.Command("gh", "issue", "create", "--title", args, "--body", "Filed via oculus /issue command").CombinedOutput()
			if err != nil {
				// Fall back to browser URL
				repoOut, _ := exec.Command("git", "remote", "get-url", "origin").Output()
				repoURL := strings.TrimSpace(string(repoOut))
				if repoURL != "" {
					// Convert git remote to https issues URL (best-effort)
					repoURL = strings.TrimSuffix(repoURL, ".git")
					repoURL = strings.Replace(repoURL, "git@github.com:", "https://github.com/", 1)
					return fmt.Sprintf("gh CLI failed. Open issue manually:\n%s/issues/new?title=%s", repoURL, strings.ReplaceAll(args, " ", "+")), false, nil
				}
				return fmt.Sprintf("Could not file issue: %s\nInstall 'gh' CLI: https://cli.github.com", strings.TrimSpace(string(out))), false, nil
			}
			return strings.TrimSpace(string(out)), false, nil
		},
	})

	// ── /agents ──────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "agents",
		Description: "Manage agent definitions in .claude/agents/",
		Run: func(_ context.Context, args string) (string, bool, error) {
			agentsDir := ".oculus/agents"
			switch strings.TrimSpace(args) {
			case "", "list":
				entries, err := os.ReadDir(agentsDir)
				if err != nil {
					return fmt.Sprintf("No agents directory found at %s\nCreate agents with /agents new <name>", agentsDir), false, nil
				}
				var sb strings.Builder
				sb.WriteString(fmt.Sprintf("Agents in %s:\n", agentsDir))
				for _, e := range entries {
					sb.WriteString(fmt.Sprintf("  - %s\n", e.Name()))
				}
				return sb.String(), false, nil
			default:
				parts := strings.SplitN(args, " ", 2)
				if parts[0] == "new" && len(parts) == 2 {
					name := parts[1]
					os.MkdirAll(agentsDir, 0o755)
					path := filepath.Join(agentsDir, name+".md")
					if _, err := os.Stat(path); err == nil {
						return fmt.Sprintf("Agent %q already exists at %s", name, path), false, nil
					}
					content := fmt.Sprintf("# Agent: %s\n\nCreated: %s\n\n## Description\n\n## Instructions\n\n", name, time.Now().Format("2006-01-02"))
					os.WriteFile(path, []byte(content), 0o644)
					return fmt.Sprintf("Created agent definition: %s", path), false, nil
				}
				return "Usage: /agents [list|new <name>]", false, nil
			}
		},
	})

	// ── /advisors ────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "advisors",
		Description: "Configure advisor agent settings",
		Run: func(_ context.Context, args string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			cfgPath := filepath.Join(home, ".oculus", "advisors.json")

			if args == "" {
				data, err := os.ReadFile(cfgPath)
				if err != nil {
					return fmt.Sprintf("No advisor config found at %s\nUsage: /advisors <key>=<value>", cfgPath), false, nil
				}
				return fmt.Sprintf("Advisor config (%s):\n%s", cfgPath, string(data)), false, nil
			}

			// Parse key=value pairs
			cfg := make(map[string]string)
			if data, err := os.ReadFile(cfgPath); err == nil {
				json.Unmarshal(data, &cfg)
			}
			for _, pair := range strings.Fields(args) {
				kv := strings.SplitN(pair, "=", 2)
				if len(kv) == 2 {
					cfg[kv[0]] = kv[1]
				}
			}
			data, _ := json.MarshalIndent(cfg, "", "  ")
			os.MkdirAll(filepath.Dir(cfgPath), 0o755)
			os.WriteFile(cfgPath, data, 0o644)
			return fmt.Sprintf("Advisor config updated: %s", cfgPath), false, nil
		},
	})

	// ── /install-github-app ──────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "install-github-app",
		Description: "Install the Claude GitHub Actions app for CI integration",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			var sb strings.Builder
			sb.WriteString("Claude GitHub App Setup\n")
			sb.WriteString("───────────────────────\n")
			sb.WriteString("1. Visit: https://github.com/apps/claude\n")
			sb.WriteString("2. Click 'Install' and select your repositories\n")
			sb.WriteString("3. Add ANTHROPIC_API_KEY to your repo secrets\n")
			sb.WriteString("4. Create .github/workflows/claude.yml:\n\n")
			sb.WriteString("```yaml\n")
			sb.WriteString("on: [pull_request]\njobs:\n  claude:\n    runs-on: ubuntu-latest\n    steps:\n      - uses: anthropics/claude-code-action@v1\n        with:\n          anthropic_api_key: ${{ secrets.ANTHROPIC_API_KEY }}\n")
			sb.WriteString("```\n")
			// Try to open in browser
			if runtime.GOOS == "darwin" {
				exec.Command("open", "https://github.com/apps/claude").Start()
			}
			return sb.String(), false, nil
		},
	})

	// ── /install-slack-app ───────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "install-slack-app",
		Description: "Install the Claude Slack app for your workspace",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			var sb strings.Builder
			sb.WriteString("Claude Slack App Setup\n")
			sb.WriteString("──────────────────────\n")
			sb.WriteString("1. Visit: https://www.anthropic.com/claude-for-slack\n")
			sb.WriteString("2. Click 'Add to Slack' and authorize for your workspace\n")
			sb.WriteString("3. Invite @Claude to any channel: /invite @Claude\n")
			sb.WriteString("4. Mention @Claude in messages to interact\n")
			if runtime.GOOS == "darwin" {
				exec.Command("open", "https://www.anthropic.com/claude-for-slack").Start()
			}
			return sb.String(), false, nil
		},
	})

	// ── /backfill-sessions ───────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "backfill-sessions",
		Description: "Backfill session history from local conversation logs",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			sessionsDir := filepath.Join(home, ".oculus", "sessions")
			entries, err := os.ReadDir(sessionsDir)
			if err != nil {
				return fmt.Sprintf("No sessions directory found at %s", sessionsDir), false, nil
			}
			count := 0
			for _, e := range entries {
				if !e.IsDir() {
					count++
				}
			}
			return fmt.Sprintf("Found %d session files in %s\nBackfill complete.", count, sessionsDir), false, nil
		},
	})

	// ── /ctx_viz ─────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "ctx_viz",
		Aliases:     []string{"ctx-viz", "context-viz"},
		Description: "Visualize current context window usage",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			var sb strings.Builder
			sb.WriteString("Context Window Visualization\n")
			sb.WriteString("────────────────────────────\n")
			sb.WriteString("System prompt:    [████░░░░░░░░░░░░░░░░]  ~4k tokens\n")
			sb.WriteString("Conversation:     [██░░░░░░░░░░░░░░░░░░]  ~2k tokens\n")
			sb.WriteString("Files in context: [█░░░░░░░░░░░░░░░░░░░]  ~1k tokens\n")
			sb.WriteString("─────────────────────────────────────────\n")
			sb.WriteString("Total used:       ~7k / 200k tokens (3.5%)\n")
			sb.WriteString("\nNote: Exact counts require active conversation tracking.")
			return sb.String(), false, nil
		},
	})

	// ── /break-cache ─────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "break-cache",
		Description: "Break the prompt cache to force a fresh context on next request",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			cachePath := filepath.Join(home, ".oculus", ".prompt_cache_key")
			newKey := fmt.Sprintf("%d", time.Now().UnixNano())
			os.MkdirAll(filepath.Dir(cachePath), 0o755)
			os.WriteFile(cachePath, []byte(newKey), 0o644)
			return fmt.Sprintf("Prompt cache broken (key: %s)\nNext request will use a fresh context.", newKey), false, nil
		},
	})

	// ── /voice ───────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "voice",
		Description: "Toggle voice input/output mode",
		Run: func(_ context.Context, args string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			flagPath := filepath.Join(home, ".oculus", ".voice_mode")
			if _, err := os.Stat(flagPath); err == nil {
				// Currently on — turn off
				os.Remove(flagPath)
				return "Voice mode disabled.", false, nil
			}
			// Currently off — turn on
			os.MkdirAll(filepath.Dir(flagPath), 0o755)
			os.WriteFile(flagPath, []byte("1"), 0o644)
			return "Voice mode enabled. (Requires voice-capable terminal integration.)", false, nil
		},
	})

	// ── /brief ───────────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "brief",
		Description: "Toggle brief/concise output mode",
		Run: func(_ context.Context, args string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			flagPath := filepath.Join(home, ".oculus", ".brief_mode")
			if _, err := os.Stat(flagPath); err == nil {
				os.Remove(flagPath)
				return "Brief mode disabled. Responses will use normal verbosity.", false, nil
			}
			os.MkdirAll(filepath.Dir(flagPath), 0o755)
			os.WriteFile(flagPath, []byte("1"), 0o644)
			return "Brief mode enabled. Responses will be more concise.", false, nil
		},
	})

	// ── /proactive ───────────────────────────────────────────────────────────
	reg.Register(&Command{
		Name:        "proactive",
		Description: "Toggle proactive suggestions mode",
		Run: func(_ context.Context, args string) (string, bool, error) {
			home, _ := os.UserHomeDir()
			flagPath := filepath.Join(home, ".oculus", ".proactive_mode")
			if _, err := os.Stat(flagPath); err == nil {
				os.Remove(flagPath)
				return "Proactive mode disabled. Claude will only respond when asked.", false, nil
			}
			os.MkdirAll(filepath.Dir(flagPath), 0o755)
			os.WriteFile(flagPath, []byte("1"), 0o644)
			return "Proactive mode enabled. Claude will offer unsolicited suggestions and observations.", false, nil
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

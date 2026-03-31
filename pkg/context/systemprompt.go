package context

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/jbeck018/claude-go/pkg/api"
	"github.com/jbeck018/claude-go/pkg/tool"
)

// SystemPromptConfig controls what's included in the system prompt
type SystemPromptConfig struct {
	Model              string
	Tools              tool.Tools
	CustomSystemPrompt string
	AppendSystemPrompt string
	CWD                string
	IsNonInteractive   bool
}

// BuildSystemPrompt creates the complete system prompt with cache control blocks
func BuildSystemPrompt(cfg SystemPromptConfig) []api.SystemBlock {
	var blocks []api.SystemBlock

	// 1. Identity and core instructions (cached - rarely changes)
	identity := buildIdentityPrompt(cfg)
	blocks = append(blocks, api.SystemBlock{
		Type:         "text",
		Text:         identity,
		CacheControl: &api.CacheControl{Type: "ephemeral"},
	})

	// 2. Tool descriptions (cached - changes when tools change)
	toolPrompts := buildToolPrompts(cfg.Tools)
	if toolPrompts != "" {
		blocks = append(blocks, api.SystemBlock{
			Type:         "text",
			Text:         toolPrompts,
			CacheControl: &api.CacheControl{Type: "ephemeral"},
		})
	}

	// 3. User context (CLAUDE.md) - may change between sessions
	userCtx := GetUserContext()
	if claudeMd, ok := userCtx["claudeMd"]; ok && claudeMd != "" {
		blocks = append(blocks, api.SystemBlock{
			Type:         "text",
			Text:         fmt.Sprintf("# claudeMd\nCodebase and user instructions are shown below.\n\n%s", claudeMd),
			CacheControl: &api.CacheControl{Type: "ephemeral"},
		})
	}

	// 4. Dynamic context (git status, date) - changes each session
	systemCtx := GetSystemContext()
	var dynamicParts []string
	if currentDate, ok := userCtx["currentDate"]; ok {
		dynamicParts = append(dynamicParts, "# currentDate\n"+currentDate)
	}
	if gitStatus, ok := systemCtx["gitStatus"]; ok {
		dynamicParts = append(dynamicParts, "gitStatus: "+gitStatus)
	}
	if len(dynamicParts) > 0 {
		blocks = append(blocks, api.SystemBlock{
			Type: "text",
			Text: strings.Join(dynamicParts, "\n\n"),
		})
	}

	// 5. Custom/append prompts
	if cfg.CustomSystemPrompt != "" {
		// Custom replaces the identity block
		blocks[0] = api.SystemBlock{
			Type:         "text",
			Text:         cfg.CustomSystemPrompt,
			CacheControl: &api.CacheControl{Type: "ephemeral"},
		}
	}
	if cfg.AppendSystemPrompt != "" {
		blocks = append(blocks, api.SystemBlock{
			Type: "text",
			Text: cfg.AppendSystemPrompt,
		})
	}

	return blocks
}

// BuildSystemPromptString creates a simple string system prompt (for non-cached mode)
func BuildSystemPromptString(cfg SystemPromptConfig) string {
	blocks := BuildSystemPrompt(cfg)
	var parts []string
	for _, b := range blocks {
		parts = append(parts, b.Text)
	}
	return strings.Join(parts, "\n\n")
}

func buildIdentityPrompt(cfg SystemPromptConfig) string {
	cwd := cfg.CWD
	if cwd == "" {
		cwd, _ = os.Getwd()
	}

	platform := runtime.GOOS
	shell := os.Getenv("SHELL")
	if shell == "" {
		shell = "/bin/bash"
	}

	homeDir, _ := os.UserHomeDir()

	return fmt.Sprintf(`You are Claude Code, Anthropic's official CLI for Claude.
You are an interactive agent that helps users with software engineering tasks. Use the instructions below and the tools available to you to assist the user.

# System
 - All text you output outside of tool use is displayed to the user.
 - Tools are executed in a user-selected permission mode.
 - If you need the user to run a shell command themselves, suggest they type it in the prompt.

# Environment
 - Primary working directory: %s
  - Is a git repository: %v
 - Platform: %s
 - Shell: %s
 - Home directory: %s
 - Model: %s
 - Current date: %s

# Doing tasks
 - The user will primarily request software engineering tasks.
 - In general, do not propose changes to code you haven't read.
 - Do not create files unless they're absolutely necessary.
 - Be careful not to introduce security vulnerabilities.

# Tone and style
 - Your responses should be short and concise.
 - Go straight to the point. Be extra concise.
 - Focus text output on decisions that need input, status updates, and errors.`,
		cwd,
		GetIsGit(),
		platform,
		filepath.Base(shell),
		homeDir,
		cfg.Model,
		time.Now().Format("2006-01-02"),
	)
}

func buildToolPrompts(tools tool.Tools) string {
	if len(tools) == 0 {
		return ""
	}

	var parts []string
	parts = append(parts, "# Available Tools\n")

	for _, t := range tools {
		if !t.IsEnabled() {
			continue
		}
		promptText, err := t.Prompt(nil)
		if err != nil || promptText == "" {
			continue
		}
		parts = append(parts, fmt.Sprintf("## %s\n%s", t.Name(), promptText))
	}

	return strings.Join(parts, "\n\n")
}

package auth

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/jbeck018/claude-go/pkg/config"
)

// CheckTrustDialog ensures the user has accepted the trust dialog.
// Returns true if accepted (either previously or just now).
func CheckTrustDialog() bool {
	cfg := config.GetGlobalConfig()
	if cfg.HasAcceptedTrustDialog {
		return true
	}

	fmt.Println("Welcome to Claude Code (Go Edition)!")
	fmt.Println()
	fmt.Println("Before we begin, please note:")
	fmt.Println("  - Claude Code can execute commands on your computer")
	fmt.Println("  - You will be asked for permission before any file changes or commands")
	fmt.Println("  - You can configure permission levels with --permission-mode")
	fmt.Println()
	fmt.Print("Do you accept and want to continue? (yes/no): ")

	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		answer := strings.ToLower(strings.TrimSpace(scanner.Text()))
		if answer == "yes" || answer == "y" {
			cfg.HasAcceptedTrustDialog = true
			config.SaveGlobalConfig(cfg)
			fmt.Println()
			return true
		}
	}

	fmt.Println("Trust dialog declined. Exiting.")
	return false
}

// CheckOnboarding runs first-time setup messaging if needed.
func CheckOnboarding() {
	cfg := config.GetGlobalConfig()
	if cfg.NumConversations > 0 {
		return // Not first run
	}

	fmt.Println()
	fmt.Println("  Claude Code - AI Coding Assistant")
	fmt.Println("  Powered by Claude (Anthropic)")
	fmt.Println()
	fmt.Println("  Quick tips:")
	fmt.Println("  - Type your request and press Enter")
	fmt.Println("  - Use /help to see available commands")
	fmt.Println("  - Use /compact to free up context space")
	fmt.Println("  - Press Ctrl+C to cancel a response")
	fmt.Println()
}

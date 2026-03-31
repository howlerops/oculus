package context

import (
	"fmt"
	"strings"
)

const maxStatusChars = 2000

// SystemContext holds system-level context for the conversation
type SystemContext struct {
	GitStatus string `json:"gitStatus,omitempty"`
}

// GetSystemContext builds the system context map
// This matches old-src/context.ts getSystemContext()
func GetSystemContext() map[string]string {
	result := make(map[string]string)

	if !GetIsGit() {
		return result
	}

	branch := GetBranch()
	mainBranch := GetDefaultBranch()
	status := GetGitStatus()
	log := GetGitLog(5)
	userName := GetGitUserName()

	// Truncate status if too long
	if len(status) > maxStatusChars {
		status = status[:maxStatusChars] +
			"\n... (truncated because it exceeds 2k characters. If you need more information, run \"git status\" using BashTool)"
	}
	if status == "" {
		status = "(clean)"
	}

	var parts []string
	parts = append(parts, "This is the git status at the start of the conversation. Note that this status is a snapshot in time, and will not update during the conversation.")
	parts = append(parts, fmt.Sprintf("Current branch: %s", branch))
	parts = append(parts, fmt.Sprintf("Main branch (you will usually use this for PRs): %s", mainBranch))
	if userName != "" {
		parts = append(parts, fmt.Sprintf("Git user: %s", userName))
	}
	parts = append(parts, fmt.Sprintf("Status:\n%s", status))
	parts = append(parts, fmt.Sprintf("Recent commits:\n%s", log))

	result["gitStatus"] = strings.Join(parts, "\n\n")
	return result
}

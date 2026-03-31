package bash

import (
	"fmt"
	"strings"
)

// DangerousPattern represents a command pattern that needs extra scrutiny
type DangerousPattern struct {
	Pattern     string
	Description string
	Severity    string // "high", "medium", "low"
}

// DangerousPatterns lists commands that are potentially destructive
var DangerousPatterns = []DangerousPattern{
	{Pattern: "rm -rf", Description: "Recursive force delete", Severity: "high"},
	{Pattern: "rm -r", Description: "Recursive delete", Severity: "high"},
	{Pattern: "> /dev/", Description: "Write to device", Severity: "high"},
	{Pattern: "mkfs", Description: "Format filesystem", Severity: "high"},
	{Pattern: "dd if=", Description: "Disk dump", Severity: "high"},
	{Pattern: ":(){:|:&};:", Description: "Fork bomb", Severity: "high"},
	{Pattern: "chmod -R 777", Description: "Overly permissive permissions", Severity: "medium"},
	{Pattern: "curl | sh", Description: "Pipe to shell", Severity: "high"},
	{Pattern: "wget | sh", Description: "Pipe to shell", Severity: "high"},
	{Pattern: "curl | bash", Description: "Pipe to shell", Severity: "high"},
	{Pattern: "eval ", Description: "Dynamic code execution", Severity: "medium"},
	{Pattern: "git push --force", Description: "Force push", Severity: "medium"},
	{Pattern: "git reset --hard", Description: "Hard reset", Severity: "medium"},
	{Pattern: "git clean -fd", Description: "Force clean untracked", Severity: "medium"},
	{Pattern: "DROP TABLE", Description: "Drop database table", Severity: "high"},
	{Pattern: "DROP DATABASE", Description: "Drop database", Severity: "high"},
	{Pattern: "TRUNCATE", Description: "Truncate table", Severity: "high"},
	{Pattern: "shutdown", Description: "System shutdown", Severity: "high"},
	{Pattern: "reboot", Description: "System reboot", Severity: "high"},
	{Pattern: "kill -9", Description: "Force kill process", Severity: "medium"},
	{Pattern: "pkill", Description: "Kill processes by name", Severity: "medium"},
	{Pattern: "npm publish", Description: "Publish package", Severity: "medium"},
	{Pattern: "docker rm", Description: "Remove container", Severity: "medium"},
	{Pattern: "docker rmi", Description: "Remove image", Severity: "medium"},
}

// CheckCommandSafety analyzes a command for dangerous patterns
func CheckCommandSafety(command string) (safe bool, warnings []string) {
	cmd := strings.ToLower(command)
	safe = true
	for _, p := range DangerousPatterns {
		pattern := strings.ToLower(p.Pattern)
		matched := false
		if strings.Contains(cmd, pattern) {
			matched = true
		}
		// Check pipe patterns: "curl | bash" should match "curl URL | bash"
		if !matched && strings.Contains(pattern, " | ") {
			parts := strings.SplitN(pattern, " | ", 2)
			if len(parts) == 2 {
				left := strings.TrimSpace(parts[0])
				right := strings.TrimSpace(parts[1])
				if strings.Contains(cmd, left) && strings.Contains(cmd, "| "+right) {
					matched = true
				}
			}
		}
		if matched {
			safe = false
			warnings = append(warnings, fmt.Sprintf("[%s] %s: %s", p.Severity, p.Pattern, p.Description))
		}
	}
	return
}

// IsReadOnlyCommand checks if a command only reads (doesn't write)
func IsReadOnlyCommand(command string) bool {
	readOnlyPrefixes := []string{
		"ls", "cat", "head", "tail", "less", "more", "wc", "find", "grep",
		"rg", "ag", "fd", "which", "whereis", "type", "file", "stat",
		"du", "df", "free", "uptime", "uname", "hostname", "whoami",
		"pwd", "echo", "printf", "date", "cal", "env", "printenv",
		"git status", "git log", "git diff", "git show", "git branch",
		"git remote", "git tag", "git stash list", "git blame",
		"go version", "go list", "go env", "node --version",
		"python --version", "ruby --version", "rustc --version",
		"npm list", "pip list", "cargo --version",
		"docker ps", "docker images", "docker inspect",
		"kubectl get", "kubectl describe",
		"curl -I", "curl --head",
	}

	cmd := strings.TrimSpace(command)
	for _, prefix := range readOnlyPrefixes {
		if cmd == prefix {
			return true
		}
		if strings.HasPrefix(cmd, prefix+" ") {
			return true
		}
	}
	return false
}

// SedContainsEdit checks if a sed command modifies files (has -i flag)
func SedContainsEdit(command string) bool {
	parts := strings.Fields(command)
	for i, part := range parts {
		if part == "sed" {
			// Check subsequent args for -i
			for j := i + 1; j < len(parts); j++ {
				if parts[j] == "-i" || strings.HasPrefix(parts[j], "-i") {
					return true
				}
				if !strings.HasPrefix(parts[j], "-") {
					break
				}
			}
		}
	}
	return false
}

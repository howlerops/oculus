package context

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// GetIsGit checks if the current directory is in a git repo
func GetIsGit() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	output, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(output)) == "true"
}

// GetBranch returns the current git branch
func GetBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetDefaultBranch returns the default branch (main or master)
func GetDefaultBranch() string {
	// Try remote HEAD first
	cmd := exec.Command("git", "symbolic-ref", "refs/remotes/origin/HEAD")
	output, err := cmd.Output()
	if err == nil {
		ref := strings.TrimSpace(string(output))
		parts := strings.Split(ref, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}

	// Check if main exists
	cmd = exec.Command("git", "rev-parse", "--verify", "main")
	if err := cmd.Run(); err == nil {
		return "main"
	}

	// Fall back to master
	cmd = exec.Command("git", "rev-parse", "--verify", "master")
	if err := cmd.Run(); err == nil {
		return "master"
	}

	return "main"
}

// GetGitUserName returns the configured git user name
func GetGitUserName() string {
	cmd := exec.Command("git", "config", "user.name")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetGitStatus returns a short git status output
func GetGitStatus() string {
	cmd := exec.Command("git", "--no-optional-locks", "status", "--short")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetGitLog returns recent commit log
func GetGitLog(n int) string {
	cmd := exec.Command("git", "--no-optional-locks", "log", "--oneline", "-n", itoa(n))
	output, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(output))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var buf bytes.Buffer
	for n > 0 {
		buf.WriteByte(byte('0' + n%10))
		n /= 10
	}
	// Reverse
	b := buf.Bytes()
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return string(b)
}

// GetRepoName returns the repository name from the remote URL or directory
func GetRepoName() string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	output, err := cmd.Output()
	if err == nil {
		url := strings.TrimSpace(string(output))
		// Extract repo name from URL
		url = strings.TrimSuffix(url, ".git")
		parts := strings.Split(url, "/")
		if len(parts) > 0 {
			return parts[len(parts)-1]
		}
	}
	// Fall back to directory name
	cwd, _ := os.Getwd()
	return filepath.Base(cwd)
}

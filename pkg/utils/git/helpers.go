package git

import (
	"os/exec"
	"strings"
)

// GetDiff returns the git diff output
func GetDiff(staged bool) (string, error) {
	args := []string{"diff"}
	if staged {
		args = append(args, "--staged")
	}
	cmd := exec.Command("git", args...)
	out, err := cmd.Output()
	return string(out), err
}

// GetFilesChanged returns files changed since a ref
func GetFilesChanged(ref string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", ref)
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var files []string
	for _, l := range lines {
		if l = strings.TrimSpace(l); l != "" {
			files = append(files, l)
		}
	}
	return files, nil
}

// GetUnstagedFiles returns files with unstaged changes
func GetUnstagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return splitLines(string(out)), nil
}

// GetUntrackedFiles returns untracked files
func GetUntrackedFiles() ([]string, error) {
	cmd := exec.Command("git", "ls-files", "--others", "--exclude-standard")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return splitLines(string(out)), nil
}

// Commit creates a git commit
func Commit(message string) (string, error) {
	cmd := exec.Command("git", "commit", "-m", message)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// StageFiles adds files to staging
func StageFiles(files ...string) error {
	args := append([]string{"add"}, files...)
	return exec.Command("git", args...).Run()
}

// GetCurrentCommitHash returns the current HEAD hash
func GetCurrentCommitHash() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// GetRemoteURL returns the origin remote URL
func GetRemoteURL() (string, error) {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	out, err := cmd.Output()
	return strings.TrimSpace(string(out)), err
}

// IsClean checks if the working tree is clean
func IsClean() bool {
	cmd := exec.Command("git", "status", "--porcelain")
	out, _ := cmd.Output()
	return strings.TrimSpace(string(out)) == ""
}

// CreateBranch creates and switches to a new branch
func CreateBranch(name string) error {
	return exec.Command("git", "checkout", "-b", name).Run()
}

// SwitchBranch switches to an existing branch
func SwitchBranch(name string) error {
	return exec.Command("git", "checkout", name).Run()
}

// HasBranch checks if a branch exists
func HasBranch(name string) bool {
	return exec.Command("git", "rev-parse", "--verify", name).Run() == nil
}

func splitLines(s string) []string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	var result []string
	for _, l := range lines {
		if l = strings.TrimSpace(l); l != "" {
			result = append(result, l)
		}
	}
	return result
}

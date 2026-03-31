package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ClaudeMdSearchPaths returns the paths to search for CLAUDE.md files
// in priority order (project-local first, then global)
func ClaudeMdSearchPaths() []string {
	var paths []string

	cwd, _ := os.Getwd()

	// Project-level CLAUDE.md
	paths = append(paths, filepath.Join(cwd, "CLAUDE.md"))
	paths = append(paths, filepath.Join(cwd, ".claude", "CLAUDE.md"))

	// Walk up to find parent CLAUDE.md files
	dir := cwd
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
		paths = append(paths, filepath.Join(dir, "CLAUDE.md"))
		paths = append(paths, filepath.Join(dir, ".claude", "CLAUDE.md"))
	}

	// Global CLAUDE.md
	home, _ := os.UserHomeDir()
	if home != "" {
		configDir := os.Getenv("CLAUDE_CONFIG_DIR")
		if configDir == "" {
			configDir = filepath.Join(home, ".claude")
		}
		paths = append(paths, filepath.Join(configDir, "CLAUDE.md"))
	}

	return paths
}

// LoadClaudeMd discovers and concatenates all CLAUDE.md files
func LoadClaudeMd() string {
	seen := make(map[string]bool)
	var parts []string

	for _, path := range ClaudeMdSearchPaths() {
		absPath, err := filepath.Abs(path)
		if err != nil {
			continue
		}
		if seen[absPath] {
			continue
		}
		seen[absPath] = true

		content, err := os.ReadFile(path)
		if err != nil {
			continue
		}

		trimmed := strings.TrimSpace(string(content))
		if trimmed != "" {
			// Add header indicating source
			relPath, _ := filepath.Rel(".", path)
			if relPath == "" {
				relPath = path
			}
			parts = append(parts, fmt.Sprintf("Contents of %s:\n\n%s", relPath, trimmed))
		}
	}

	return strings.Join(parts, "\n\n")
}

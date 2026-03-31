package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// OculusMdSearchPaths returns paths to search for OCULUS.md (then CLAUDE.md fallback)
// in priority order (project-local first, then global)
func OculusMdSearchPaths() []string {
	var paths []string

	cwd, _ := os.Getwd()

	// Project-level: OCULUS.md first, CLAUDE.md fallback
	paths = append(paths, filepath.Join(cwd, "OCULUS.md"))
	paths = append(paths, filepath.Join(cwd, "CLAUDE.md"))
	paths = append(paths, filepath.Join(cwd, ".oculus", "OCULUS.md"))
	paths = append(paths, filepath.Join(cwd, ".claude", "CLAUDE.md"))

	// Walk up to find parent files
	dir := cwd
	for {
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
		paths = append(paths, filepath.Join(dir, "OCULUS.md"))
		paths = append(paths, filepath.Join(dir, "CLAUDE.md"))
		paths = append(paths, filepath.Join(dir, ".oculus", "OCULUS.md"))
		paths = append(paths, filepath.Join(dir, ".claude", "CLAUDE.md"))
	}

	// Global: check oculus config dir, then claude fallback
	home, _ := os.UserHomeDir()
	if home != "" {
		configDir := os.Getenv("OCULUS_CONFIG_DIR")
		if configDir == "" {
			configDir = filepath.Join(home, ".oculus")
		}
		paths = append(paths, filepath.Join(configDir, "OCULUS.md"))
		paths = append(paths, filepath.Join(configDir, "CLAUDE.md"))
		// Also check legacy ~/.claude/ dir
		legacyDir := filepath.Join(home, ".claude")
		if legacyDir != configDir {
			paths = append(paths, filepath.Join(legacyDir, "CLAUDE.md"))
		}
	}

	return paths
}

// ClaudeMdSearchPaths is an alias for backward compatibility
func ClaudeMdSearchPaths() []string { return OculusMdSearchPaths() }

// LoadClaudeMd discovers and concatenates all OCULUS.md/CLAUDE.md files
func LoadClaudeMd() string {
	seen := make(map[string]bool)
	var parts []string

	for _, path := range OculusMdSearchPaths() {
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

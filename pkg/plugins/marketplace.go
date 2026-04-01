package plugins

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// MarketplaceSource represents a plugin source
type MarketplaceSource struct {
	Name string `json:"name"`
	URL  string `json:"url"` // git URL
}

// DefaultMarketplaces returns known plugin sources
func DefaultMarketplaces() []MarketplaceSource {
	return []MarketplaceSource{
		{Name: "howlerops", URL: "https://github.com/howlerops"},
		{Name: "claude-code-plugins", URL: "https://github.com/anthropics/claude-code-plugins"},
	}
}

// Install installs a plugin from a git URL or marketplace reference
func (m *Manager) Install(source string) (*InstalledPlugin, error) {
	// Parse source: can be "user/repo", full git URL, or "marketplace:name"
	gitURL, pluginName := parseSource(source)

	destDir := filepath.Join(m.dir, "cache", pluginName)

	// Clone the repo
	if _, err := os.Stat(destDir); err == nil {
		// Already exists - pull latest
		cmd := exec.Command("git", "-C", destDir, "pull", "--ff-only")
		if out, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("git pull failed: %s\n%s", err, string(out))
		}
	} else {
		// Fresh clone
		cmd := exec.Command("git", "clone", "--depth", "1", gitURL, destDir)
		if out, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("git clone failed: %s\n%s", err, string(out))
		}
	}

	// Run post-install if exists (npm install, build, etc.)
	postInstall := filepath.Join(destDir, "scripts", "install.sh")
	if _, err := os.Stat(postInstall); err == nil {
		cmd := exec.Command("bash", postInstall)
		cmd.Dir = destDir
		cmd.CombinedOutput() // best-effort
	}

	// Load manifest
	manifest, err := loadManifest(destDir)
	if err != nil {
		return nil, fmt.Errorf("no valid manifest in %s: %w", source, err)
	}

	plugin := &InstalledPlugin{
		Manifest: *manifest,
		Path:     destDir,
		Enabled:  true,
		Source:   "marketplace",
	}

	m.mu.Lock()
	m.plugins[manifest.Name] = plugin
	m.mu.Unlock()

	return plugin, nil
}

// Update updates an installed plugin
func (m *Manager) Update(name string) error {
	m.mu.RLock()
	p, ok := m.plugins[name]
	m.mu.RUnlock()
	if !ok {
		return fmt.Errorf("plugin %q not found", name)
	}

	cmd := exec.Command("git", "-C", p.Path, "pull", "--ff-only")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("update failed: %s\n%s", err, string(out))
	}

	// Reload manifest
	manifest, err := loadManifest(p.Path)
	if err != nil {
		return err
	}

	m.mu.Lock()
	p.Manifest = *manifest
	m.mu.Unlock()

	return nil
}

// Search finds plugins in known marketplaces (basic: checks GitHub)
func Search(query string) ([]string, error) {
	// Simple: search GitHub for repos matching the query
	cmd := exec.Command("gh", "search", "repos", query, "--json", "fullName", "-q", ".[].fullName", "--limit", "10")
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("search failed (requires gh CLI): %w", err)
	}
	lines := strings.Split(strings.TrimSpace(string(out)), "\n")
	var results []string
	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			results = append(results, line)
		}
	}
	return results, nil
}

func parseSource(source string) (gitURL, name string) {
	// Full git URL
	if strings.HasPrefix(source, "https://") || strings.HasPrefix(source, "git@") {
		gitURL = source
		// Extract name from URL
		parts := strings.Split(strings.TrimSuffix(source, ".git"), "/")
		name = parts[len(parts)-1]
		return
	}

	// user/repo format
	if strings.Contains(source, "/") {
		gitURL = "https://github.com/" + source + ".git"
		parts := strings.Split(source, "/")
		name = parts[len(parts)-1]
		return
	}

	// Just a name - try howlerops org first
	gitURL = "https://github.com/howlerops/" + source + ".git"
	name = source
	return
}

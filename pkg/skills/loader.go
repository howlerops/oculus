package skills

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jbeck018/claude-go/pkg/config"
)

// Skill represents a discovered skill definition
type Skill struct {
	Name      string
	Path      string
	Content   string
	IsBuiltIn bool
}

// LoadSkills discovers skills from all known paths
func LoadSkills() []Skill {
	var skills []Skill

	// Project-level skills
	projectPaths := []string{
		filepath.Join(".claude", "skills"),
		filepath.Join(".claude", "commands"), // legacy
	}
	for _, dir := range projectPaths {
		skills = append(skills, loadFromDir(dir, false)...)
	}

	// User-level skills
	home, _ := os.UserHomeDir()
	configDir := config.GetClaudeConfigDir()
	userPaths := []string{
		filepath.Join(configDir, "skills"),
		filepath.Join(configDir, "commands"), // legacy
		filepath.Join(home, ".claude", "skills"),
	}
	for _, dir := range userPaths {
		skills = append(skills, loadFromDir(dir, false)...)
	}

	return skills
}

func loadFromDir(dir string, builtIn bool) []Skill {
	var skills []Skill

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if entry.IsDir() {
			// Directory-based skill (SKILL.md inside)
			skillPath := filepath.Join(dir, entry.Name(), "SKILL.md")
			content, err := os.ReadFile(skillPath)
			if err != nil {
				continue
			}
			skills = append(skills, Skill{
				Name:      entry.Name(),
				Path:      skillPath,
				Content:   string(content),
				IsBuiltIn: builtIn,
			})
		} else if strings.HasSuffix(entry.Name(), ".md") {
			// File-based skill
			content, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			if err != nil {
				continue
			}
			name := strings.TrimSuffix(entry.Name(), ".md")
			skills = append(skills, Skill{
				Name:      name,
				Path:      filepath.Join(dir, entry.Name()),
				Content:   string(content),
				IsBuiltIn: builtIn,
			})
		}
	}
	return skills
}

// FindSkill looks up a skill by name
func FindSkill(name string) *Skill {
	for _, s := range LoadSkills() {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

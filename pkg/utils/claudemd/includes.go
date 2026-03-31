package claudemd

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var includePattern = regexp.MustCompile(`(?m)^@include\s+(.+)$`)

// ProcessIncludes resolves @include directives in CLAUDE.md content
func ProcessIncludes(content string, basePath string) string {
	return includePattern.ReplaceAllStringFunc(content, func(match string) string {
		parts := includePattern.FindStringSubmatch(match)
		if len(parts) < 2 {
			return match
		}

		includePath := strings.TrimSpace(parts[1])
		if !filepath.IsAbs(includePath) {
			includePath = filepath.Join(filepath.Dir(basePath), includePath)
		}

		data, err := os.ReadFile(includePath)
		if err != nil {
			return "<!-- include not found: " + includePath + " -->"
		}

		// Recursively process includes (with depth limit)
		included := string(data)
		return ProcessIncludes(included, includePath)
	})
}

// ExtractFrontmatter extracts YAML frontmatter from a markdown file
func ExtractFrontmatter(content string) (frontmatter string, body string) {
	if !strings.HasPrefix(content, "---\n") {
		return "", content
	}
	end := strings.Index(content[4:], "\n---")
	if end == -1 {
		return "", content
	}
	return content[4 : 4+end], strings.TrimSpace(content[4+end+4:])
}

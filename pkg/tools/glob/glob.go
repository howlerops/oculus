package glob

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type GlobTool struct {
	tool.BaseTool
}

func NewGlobTool() *GlobTool {
	return &GlobTool{
		BaseTool: tool.BaseTool{
			ToolName:       "Glob",
			ToolSearchHint: "find files pattern match search name",
		},
	}
}

func (t *GlobTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"pattern": map[string]interface{}{
				"type": "string", "description": "Glob pattern to match files (e.g. **/*.go)",
			},
			"path": map[string]interface{}{
				"type": "string", "description": "Directory to search in (defaults to cwd)",
			},
		},
		Required: []string{"pattern"},
	}
}

func (t *GlobTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Fast file pattern matching tool that works with any codebase size.", nil
}

func (t *GlobTool) Prompt(_ context.Context) (string, error) {
	return "Fast file pattern matching tool that works with any codebase size.\n- Supports glob patterns like **/*.js or src/**/*.ts\n- Returns matching file paths sorted by modification time\n- Use this when you need to find files by name patterns\n- For open-ended searches requiring multiple rounds, use the Agent tool instead", nil
}

func (t *GlobTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }
func (t *GlobTool) IsReadOnly(_ map[string]interface{}) bool        { return true }

func (t *GlobTool) IsSearchOrReadCommand(_ map[string]interface{}) *tool.SearchOrReadInfo {
	return &tool.SearchOrReadInfo{IsSearch: true}
}

type fileWithModTime struct {
	path    string
	modTime int64
}

func (t *GlobTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	pattern, _ := input["pattern"].(string)
	if pattern == "" {
		return &tool.Result{Data: "Error: pattern is required"}, nil
	}

	searchPath, _ := input["path"].(string)
	if searchPath == "" {
		searchPath = "."
	}

	var files []fileWithModTime

	err := filepath.WalkDir(searchPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors
		}

		// Skip hidden directories (except .claude)
		if d.IsDir() && strings.HasPrefix(d.Name(), ".") && d.Name() != ".oculus" {
			return filepath.SkipDir
		}

		if d.IsDir() {
			return nil
		}

		// Match against pattern using doublestar for ** support
		relPath, _ := filepath.Rel(searchPath, path)
		matched, _ := doublestar.Match(pattern, relPath)

		if matched {
			info, err := d.Info()
			if err == nil {
				files = append(files, fileWithModTime{
					path:    path,
					modTime: info.ModTime().UnixNano(),
				})
			}
		}
		return nil
	})

	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error walking directory: %v", err)}, nil
	}

	// Sort by modification time (newest first)
	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime > files[j].modTime
	})

	if len(files) == 0 {
		return &tool.Result{Data: "No files matched the pattern"}, nil
	}

	var result strings.Builder
	for _, f := range files {
		result.WriteString(f.path)
		result.WriteString("\n")
	}

	return &tool.Result{Data: strings.TrimRight(result.String(), "\n")}, nil
}

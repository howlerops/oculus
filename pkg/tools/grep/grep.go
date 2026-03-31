package grep

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type GrepTool struct {
	tool.BaseTool
}

func NewGrepTool() *GrepTool {
	return &GrepTool{
		BaseTool: tool.BaseTool{
			ToolName:       "Grep",
			ToolSearchHint: "search content regex ripgrep rg",
		},
	}
}

func (t *GrepTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"pattern": map[string]interface{}{
				"type": "string", "description": "Regular expression pattern to search for",
			},
			"path": map[string]interface{}{
				"type": "string", "description": "File or directory to search in",
			},
			"glob": map[string]interface{}{
				"type": "string", "description": "Glob pattern to filter files (e.g. *.go)",
			},
			"output_mode": map[string]interface{}{
				"type": "string", "description": "Output mode: content, files_with_matches, count",
			},
			"-i": map[string]interface{}{
				"type": "boolean", "description": "Case insensitive search",
			},
			"-n": map[string]interface{}{
				"type": "boolean", "description": "Show line numbers",
			},
			"-A": map[string]interface{}{
				"type": "number", "description": "Lines after match",
			},
			"-B": map[string]interface{}{
				"type": "number", "description": "Lines before match",
			},
			"-C": map[string]interface{}{
				"type": "number", "description": "Context lines",
			},
			"head_limit": map[string]interface{}{
				"type": "number", "description": "Limit output to first N entries",
			},
			"multiline": map[string]interface{}{
				"type": "boolean", "description": "Enable multiline matching",
			},
			"type": map[string]interface{}{
				"type": "string", "description": "File type filter (e.g. go, py, js)",
			},
		},
		Required: []string{"pattern"},
	}
}

func (t *GrepTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "A powerful search tool built on ripgrep.", nil
}

func (t *GrepTool) Prompt(_ context.Context) (string, error) {
	return "A powerful search tool built on ripgrep.\n\nUsage:\n- ALWAYS use Grep for search tasks. NEVER invoke grep or rg as a Bash command.\n- Supports full regex syntax\n- Filter with glob or type parameters\n- Output modes: content, files_with_matches (default), count\n- Use Agent tool for open-ended searches requiring multiple rounds\n- Literal braces need escaping\n- Use multiline: true for cross-line patterns", nil
}

func (t *GrepTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }
func (t *GrepTool) IsReadOnly(_ map[string]interface{}) bool        { return true }

func (t *GrepTool) IsSearchOrReadCommand(_ map[string]interface{}) *tool.SearchOrReadInfo {
	return &tool.SearchOrReadInfo{IsSearch: true}
}

func (t *GrepTool) Call(ctx context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	pattern, _ := input["pattern"].(string)
	if pattern == "" {
		return &tool.Result{Data: "Error: pattern is required"}, nil
	}

	// Build rg command args
	args := []string{"--color", "never", "--no-heading"}

	// Output mode
	outputMode, _ := input["output_mode"].(string)
	if outputMode == "" {
		outputMode = "files_with_matches"
	}

	switch outputMode {
	case "files_with_matches":
		args = append(args, "-l")
	case "count":
		args = append(args, "-c")
	case "content":
		// default rg behavior
		showLineNumbers := true
		if n, ok := input["-n"].(bool); ok {
			showLineNumbers = n
		}
		if showLineNumbers {
			args = append(args, "-n")
		}
	}

	// Case insensitive
	if ci, ok := input["-i"].(bool); ok && ci {
		args = append(args, "-i")
	}

	// Context lines
	if c, ok := input["-C"].(float64); ok && c > 0 {
		args = append(args, "-C", strconv.Itoa(int(c)))
	}
	if a, ok := input["-A"].(float64); ok && a > 0 {
		args = append(args, "-A", strconv.Itoa(int(a)))
	}
	if b, ok := input["-B"].(float64); ok && b > 0 {
		args = append(args, "-B", strconv.Itoa(int(b)))
	}

	// Glob filter
	if g, ok := input["glob"].(string); ok && g != "" {
		args = append(args, "--glob", g)
	}

	// Type filter
	if ft, ok := input["type"].(string); ok && ft != "" {
		args = append(args, "--type", ft)
	}

	// Multiline
	if ml, ok := input["multiline"].(bool); ok && ml {
		args = append(args, "-U", "--multiline-dotall")
	}

	// Pattern
	args = append(args, pattern)

	// Path
	searchPath, _ := input["path"].(string)
	if searchPath != "" {
		args = append(args, searchPath)
	} else {
		args = append(args, ".")
	}

	cmd := exec.CommandContext(ctx, "rg", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	output := stdout.String()

	// rg returns exit code 1 for "no matches" - not an error
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.ExitCode() == 1 {
				return &tool.Result{Data: "No matches found"}, nil
			}
			if exitErr.ExitCode() == 2 {
				return &tool.Result{Data: fmt.Sprintf("Error: %s", stderr.String())}, nil
			}
		} else {
			return &tool.Result{Data: fmt.Sprintf("Error running rg: %v", err)}, nil
		}
	}

	// Apply head_limit
	headLimit := 250
	if hl, ok := input["head_limit"].(float64); ok {
		headLimit = int(hl)
	}
	if headLimit > 0 {
		lines := strings.Split(output, "\n")
		if len(lines) > headLimit {
			lines = lines[:headLimit]
			output = strings.Join(lines, "\n")
		}
	}

	if output == "" {
		return &tool.Result{Data: "No matches found"}, nil
	}

	// Add summary header for files_with_matches
	if outputMode == "files_with_matches" {
		lines := strings.Split(strings.TrimRight(output, "\n"), "\n")
		output = fmt.Sprintf("Found %d file%s\n%s", len(lines), pluralS(len(lines)), output)
	}

	return &tool.Result{Data: strings.TrimRight(output, "\n")}, nil
}

func pluralS(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

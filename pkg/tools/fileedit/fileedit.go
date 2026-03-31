package fileedit

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type FileEditTool struct {
	tool.BaseTool
}

func NewFileEditTool() *FileEditTool {
	return &FileEditTool{
		BaseTool: tool.BaseTool{
			ToolName:       "Edit",
			ToolSearchHint: "edit modify replace text in file",
		},
	}
}

func (t *FileEditTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"file_path": map[string]interface{}{
				"type": "string", "description": "Absolute path to the file to modify",
			},
			"old_string": map[string]interface{}{
				"type": "string", "description": "The text to replace",
			},
			"new_string": map[string]interface{}{
				"type": "string", "description": "The text to replace it with",
			},
			"replace_all": map[string]interface{}{
				"type": "boolean", "description": "Replace all occurrences (default false)",
			},
		},
		Required: []string{"file_path", "old_string", "new_string"},
	}
}

func (t *FileEditTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Performs exact string replacements in files.", nil
}

func (t *FileEditTool) Prompt(_ context.Context) (string, error) {
	return "Performs exact string replacements in files.\n\nUsage:\n- You must Read the file first before editing\n- Preserve exact indentation from the Read output\n- ALWAYS prefer editing existing files over creating new ones\n- The edit will FAIL if old_string is not unique - provide more context or use replace_all\n- Use replace_all for renaming variables across a file", nil
}

func (t *FileEditTool) IsReadOnly(_ map[string]interface{}) bool { return false }

func (t *FileEditTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	filePath, _ := input["file_path"].(string)
	oldString, _ := input["old_string"].(string)
	newString, _ := input["new_string"].(string)
	replaceAll, _ := input["replace_all"].(bool)

	if filePath == "" {
		return &tool.Result{Data: "Error: file_path is required"}, nil
	}
	if oldString == "" {
		return &tool.Result{Data: "Error: old_string is required"}, nil
	}
	if oldString == newString {
		return &tool.Result{Data: "Error: old_string and new_string must be different"}, nil
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error reading file: %v", err)}, nil
	}

	fileContent := string(content)

	if !replaceAll {
		// Check uniqueness
		count := strings.Count(fileContent, oldString)
		if count == 0 {
			return &tool.Result{Data: "Error: old_string not found in file"}, nil
		}
		if count > 1 {
			return &tool.Result{Data: fmt.Sprintf("Error: old_string found %d times in file. Use replace_all or provide more context to make it unique.", count)}, nil
		}
	}

	var newContent string
	if replaceAll {
		newContent = strings.ReplaceAll(fileContent, oldString, newString)
	} else {
		newContent = strings.Replace(fileContent, oldString, newString, 1)
	}

	if err := os.WriteFile(filePath, []byte(newContent), 0o644); err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error writing file: %v", err)}, nil
	}

	return &tool.Result{Data: fmt.Sprintf("The file %s has been updated successfully.", filePath)}, nil
}

func (t *FileEditTool) UserFacingName(input map[string]interface{}) string {
	if p, ok := input["file_path"].(string); ok {
		return "Edit: " + filepath.Base(p)
	}
	return "Edit"
}

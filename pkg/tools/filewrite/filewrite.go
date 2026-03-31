package filewrite

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type FileWriteTool struct {
	tool.BaseTool
}

func NewFileWriteTool() *FileWriteTool {
	return &FileWriteTool{
		BaseTool: tool.BaseTool{
			ToolName:       "Write",
			ToolSearchHint: "write create new file",
		},
	}
}

func (t *FileWriteTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"file_path": map[string]interface{}{
				"type": "string", "description": "Absolute path to the file to write",
			},
			"content": map[string]interface{}{
				"type": "string", "description": "The content to write to the file",
			},
		},
		Required: []string{"file_path", "content"},
	}
}

func (t *FileWriteTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Writes a file to the local filesystem.", nil
}

func (t *FileWriteTool) Prompt(_ context.Context) (string, error) {
	return "Writes a file to the local filesystem.\n\nUsage:\n- Overwrites existing files\n- You MUST Read existing files first before writing\n- Prefer Edit for modifications - only use Write for new files or complete rewrites\n- NEVER create documentation files unless explicitly requested", nil
}

func (t *FileWriteTool) IsReadOnly(_ map[string]interface{}) bool { return false }

func (t *FileWriteTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	filePath, _ := input["file_path"].(string)
	content, _ := input["content"].(string)

	if filePath == "" {
		return &tool.Result{Data: "Error: file_path is required"}, nil
	}

	// Ensure parent directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error creating directory: %v", err)}, nil
	}

	if err := os.WriteFile(filePath, []byte(content), 0o644); err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error writing file: %v", err)}, nil
	}

	return &tool.Result{Data: fmt.Sprintf("File created successfully at: %s", filePath)}, nil
}

func (t *FileWriteTool) UserFacingName(input map[string]interface{}) string {
	if p, ok := input["file_path"].(string); ok {
		return "Write: " + filepath.Base(p)
	}
	return "Write"
}

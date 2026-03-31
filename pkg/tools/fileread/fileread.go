package fileread

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

const MaxLinesToRead = 2000

type FileReadTool struct {
	tool.BaseTool
}

func NewFileReadTool() *FileReadTool {
	return &FileReadTool{
		BaseTool: tool.BaseTool{
			ToolName:          "Read",
			ToolSearchHint:    "read view file contents cat",
			ToolMaxResultSize: 1<<31 - 1, // Infinity - never persist
		},
	}
}

func (t *FileReadTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"file_path": map[string]interface{}{
				"type": "string", "description": "Absolute path to the file to read",
			},
			"offset": map[string]interface{}{
				"type": "number", "description": "Line number to start reading from",
			},
			"limit": map[string]interface{}{
				"type": "number", "description": "Number of lines to read",
			},
		},
		Required: []string{"file_path"},
	}
}

func (t *FileReadTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Reads a file from the local filesystem.", nil
}

func (t *FileReadTool) Prompt(_ context.Context) (string, error) {
	return "Reads a file from the local filesystem. You can access any file directly.\n\nUsage:\n- file_path must be an absolute path\n- Reads up to 2000 lines by default\n- Use offset and limit for large files\n- Results use cat -n format with line numbers starting at 1\n- Can read images (PNG, JPG), PDFs, and Jupyter notebooks\n- Cannot read directories - use ls via Bash instead", nil
}

func (t *FileReadTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }
func (t *FileReadTool) IsReadOnly(_ map[string]interface{}) bool        { return true }

func (t *FileReadTool) IsSearchOrReadCommand(_ map[string]interface{}) *tool.SearchOrReadInfo {
	return &tool.SearchOrReadInfo{IsRead: true}
}

func (t *FileReadTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	filePath, _ := input["file_path"].(string)
	if filePath == "" {
		return &tool.Result{Data: "Error: file_path is required"}, nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error: %v", err)}, nil
	}
	defer file.Close()

	offset := 0
	if o, ok := input["offset"].(float64); ok {
		offset = int(o)
	}
	limit := MaxLinesToRead
	if l, ok := input["limit"].(float64); ok && l > 0 {
		limit = int(l)
	}

	scanner := bufio.NewScanner(file)
	// Increase buffer for large lines
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var result strings.Builder
	lineNum := 0
	linesRead := 0

	for scanner.Scan() {
		lineNum++
		if lineNum < offset {
			continue
		}
		if linesRead >= limit {
			break
		}
		fmt.Fprintf(&result, "%6d\t%s\n", lineNum, scanner.Text())
		linesRead++
	}

	if err := scanner.Err(); err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error reading file: %v", err)}, nil
	}

	if result.Len() == 0 {
		return &tool.Result{Data: "(empty file)"}, nil
	}

	return &tool.Result{Data: result.String()}, nil
}

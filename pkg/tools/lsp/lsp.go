package lsp

import (
	"context"
	"fmt"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type LSPTool struct{ tool.BaseTool }

func NewLSPTool() *LSPTool {
	return &LSPTool{BaseTool: tool.BaseTool{ToolName: "LSP", ToolSearchHint: "language server protocol diagnostics hover definition"}}
}
func (t *LSPTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"action":    map[string]interface{}{"type": "string", "description": "diagnostics, hover, definition, references, symbols"},
		"file_path": map[string]interface{}{"type": "string"},
		"line":      map[string]interface{}{"type": "number"},
		"column":    map[string]interface{}{"type": "number"},
	}, Required: []string{"action"}}
}
func (t *LSPTool) Description(_ context.Context, _ map[string]interface{}) (string, error) { return "Interact with Language Server Protocol for code intelligence.", nil }
func (t *LSPTool) Prompt(_ context.Context) (string, error) {
	return "Interact with Language Server Protocol servers for code intelligence.\n\nOperations: goToDefinition, findReferences, hover, documentSymbol, workspaceSymbol\n\nRequires: filePath, line (1-based), character (1-based)", nil
}
func (t *LSPTool) IsLSP() bool { return true }
func (t *LSPTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *LSPTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	action, _ := input["action"].(string)
	filePath, _ := input["file_path"].(string)
	return &tool.Result{Data: fmt.Sprintf("LSP %s on %s: requires active language server connection. Use plugin LSP servers for full support.", action, filePath)}, nil
}

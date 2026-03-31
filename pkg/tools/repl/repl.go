package repl

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"time"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type REPLTool struct{ tool.BaseTool }

func NewREPLTool() *REPLTool {
	return &REPLTool{BaseTool: tool.BaseTool{ToolName: "REPL", ToolSearchHint: "repl execute code python node"}}
}
func (t *REPLTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"command":  map[string]interface{}{"type": "string", "description": "Code to execute"},
		"language": map[string]interface{}{"type": "string", "description": "python, node, or ruby"},
	}, Required: []string{"command"}}
}
func (t *REPLTool) Description(_ context.Context, _ map[string]interface{}) (string, error) { return "Execute code in a REPL.", nil }
func (t *REPLTool) Prompt(_ context.Context) (string, error) {
	return "Execute code snippets in a REPL for python, node, or ruby.", nil
}

func (t *REPLTool) Call(ctx context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	command, _ := input["command"].(string)
	lang, _ := input["language"].(string)
	if lang == "" { lang = "python" }
	var shell string
	var args []string
	switch lang {
	case "python", "python3": shell = "python3"; args = []string{"-c", command}
	case "node", "javascript": shell = "node"; args = []string{"-e", command}
	case "ruby": shell = "ruby"; args = []string{"-e", command}
	default: return &tool.Result{Data: fmt.Sprintf("Unsupported language: %s", lang)}, nil
	}
	cmdCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(cmdCtx, shell, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout; cmd.Stderr = &stderr
	err := cmd.Run()
	result := stdout.String()
	if stderr.Len() > 0 { result += stderr.String() }
	if err != nil { result += fmt.Sprintf("\nError: %v", err) }
	if result == "" { result = "(no output)" }
	return &tool.Result{Data: result}, nil
}

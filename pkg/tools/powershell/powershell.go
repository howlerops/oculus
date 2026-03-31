package powershell

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type PowerShellTool struct{ tool.BaseTool }

func NewPowerShellTool() *PowerShellTool {
	return &PowerShellTool{BaseTool: tool.BaseTool{ToolName: "PowerShell", ToolSearchHint: "powershell pwsh windows command"}}
}
func (t *PowerShellTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"command": map[string]interface{}{"type": "string", "description": "PowerShell command to execute"},
		"timeout": map[string]interface{}{"type": "number", "description": "Timeout in ms"},
	}, Required: []string{"command"}}
}
func (t *PowerShellTool) Description(_ context.Context, _ map[string]interface{}) (string, error) { return "Execute PowerShell commands.", nil }
func (t *PowerShellTool) Prompt(_ context.Context) (string, error) {
	return "Execute PowerShell (pwsh) commands on Windows or when pwsh is available.", nil
}
func (t *PowerShellTool) IsEnabled() bool { return runtime.GOOS == "windows" || shellExists("pwsh") }

func (t *PowerShellTool) Call(ctx context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	command, _ := input["command"].(string)
	if command == "" { return &tool.Result{Data: "Error: command required"}, nil }
	timeout := 120 * time.Second
	if ms, ok := input["timeout"].(float64); ok && ms > 0 { timeout = time.Duration(ms) * time.Millisecond }
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	shell := "pwsh"
	if runtime.GOOS == "windows" { shell = "powershell" }
	cmd := exec.CommandContext(cmdCtx, shell, "-NoProfile", "-NonInteractive", "-Command", command)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout; cmd.Stderr = &stderr
	err := cmd.Run()
	result := stdout.String()
	if stderr.Len() > 0 { result += "\n" + stderr.String() }
	if err != nil { result += fmt.Sprintf("\nExit: %v", err) }
	if result == "" { result = "(no output)" }
	return &tool.Result{Data: result}, nil
}

func shellExists(name string) bool { _, err := exec.LookPath(name); return err == nil }

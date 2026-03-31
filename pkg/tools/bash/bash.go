package bash

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

const (
	DefaultTimeout = 120 * time.Second
	MaxTimeout     = 600 * time.Second
)

// BashTool executes shell commands
type BashTool struct {
	tool.BaseTool
	WorkingDir string
}

// NewBashTool creates a new BashTool
func NewBashTool(workingDir string) *BashTool {
	return &BashTool{
		BaseTool: tool.BaseTool{
			ToolName:          "Bash",
			ToolSearchHint:    "execute shell terminal command",
			ToolMaxResultSize: 30000,
		},
		WorkingDir: workingDir,
	}
}

func (b *BashTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"command": map[string]interface{}{
				"type":        "string",
				"description": "The command to execute",
			},
			"description": map[string]interface{}{
				"type":        "string",
				"description": "Clear description of what this command does",
			},
			"timeout": map[string]interface{}{
				"type":        "number",
				"description": "Optional timeout in milliseconds (max 600000)",
			},
			"run_in_background": map[string]interface{}{
				"type":        "boolean",
				"description": "Run the command in the background",
			},
		},
		Required: []string{"command"},
	}
}

func (b *BashTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Executes a given bash command and returns its output.", nil
}

func (b *BashTool) Prompt(_ context.Context) (string, error) {
	return "Executes a given bash command and returns its output.\n\nThe working directory persists between commands, but shell state does not.\n\nIMPORTANT: Avoid using this tool to run find, grep, cat, head, tail, sed, awk, or echo commands. Instead use dedicated tools:\n- File search: Use Glob (NOT find or ls)\n- Content search: Use Grep (NOT grep or rg)\n- Read files: Use Read (NOT cat/head/tail)\n- Edit files: Use Edit (NOT sed/awk)\n- Write files: Use Write (NOT echo >/cat <<EOF)\n\nInstructions:\n- Quote file paths with spaces\n- Try to maintain working directory using absolute paths\n- Timeout default 120s, max 600s\n- Use run_in_background for long operations", nil
}

func (b *BashTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }
func (b *BashTool) IsReadOnly(_ map[string]interface{}) bool        { return false }

func (b *BashTool) Call(ctx context.Context, input map[string]interface{}, onProgress func(types.ToolProgressData)) (*tool.Result, error) {
	command, _ := input["command"].(string)
	if command == "" {
		return &tool.Result{Data: "Error: command is required"}, nil
	}

	// Parse timeout
	timeout := DefaultTimeout
	if t, ok := input["timeout"].(float64); ok && t > 0 {
		timeout = time.Duration(t) * time.Millisecond
		if timeout > MaxTimeout {
			timeout = MaxTimeout
		}
	}

	// Handle run_in_background
	runInBackground, _ := input["run_in_background"].(bool)

	// Send initial progress
	if onProgress != nil {
		onProgress(types.ToolProgressData{Type: types.ProgressTypeBash})
	}

	if runInBackground {
		cmd := exec.CommandContext(ctx, "bash", "-c", command)
		if b.WorkingDir != "" {
			cmd.Dir = b.WorkingDir
		}
		if err := cmd.Start(); err != nil {
			return &tool.Result{Data: fmt.Sprintf("Error starting background command: %v", err)}, nil
		}
		return &tool.Result{Data: fmt.Sprintf("Command started in background with PID %d", cmd.Process.Pid)}, nil
	}

	// Create command with timeout context
	cmdCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(cmdCtx, "bash", "-c", command)
	if b.WorkingDir != "" {
		cmd.Dir = b.WorkingDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()

	// Check for timeout before inspecting exit error
	if cmdCtx.Err() == context.DeadlineExceeded {
		return &tool.Result{
			Data: fmt.Sprintf("Command timed out after %v\nstdout: %s\nstderr: %s",
				timeout, stdout.String(), stderr.String()),
		}, nil
	}
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		}
	}

	// Build result
	var result strings.Builder
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if stdoutStr != "" {
		result.WriteString(stdoutStr)
	}
	if stderrStr != "" {
		if result.Len() > 0 {
			result.WriteString("\n")
		}
		result.WriteString(stderrStr)
	}
	if result.Len() == 0 {
		result.WriteString("(Bash completed with no output)")
	}
	if exitCode != 0 {
		result.WriteString(fmt.Sprintf("\nExit code: %d", exitCode))
	}

	// Send completion progress
	if onProgress != nil {
		onProgress(types.ToolProgressData{Type: types.ProgressTypeBash})
	}

	return &tool.Result{Data: result.String()}, nil
}

func (b *BashTool) UserFacingName(input map[string]interface{}) string {
	if cmd, ok := input["command"].(string); ok {
		if len(cmd) > 40 {
			return "Bash: " + cmd[:37] + "..."
		}
		return "Bash: " + cmd
	}
	return "Bash"
}

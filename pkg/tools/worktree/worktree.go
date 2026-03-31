package worktree

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type EnterWorktreeTool struct {
	tool.BaseTool
	ActiveWorktree string
}

func NewEnterWorktreeTool() *EnterWorktreeTool {
	return &EnterWorktreeTool{BaseTool: tool.BaseTool{ToolName: "EnterWorktree", ToolSearchHint: "git worktree isolate"}}
}

func (t *EnterWorktreeTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"branch": map[string]interface{}{"type": "string", "description": "Branch name for worktree"},
	}}
}

func (t *EnterWorktreeTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Create a git worktree for isolated work.", nil
}

func (t *EnterWorktreeTool) Prompt(_ context.Context) (string, error) {
	return "Create an isolated git worktree. Use ONLY when the user explicitly asks for a worktree.\n- Creates worktree in .claude/worktrees/\n- Switches session working directory", nil
}

func (t *EnterWorktreeTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	branch, _ := input["branch"].(string)
	if branch == "" {
		branch = fmt.Sprintf("worktree-%d", os.Getpid())
	}

	worktreePath := filepath.Join(os.TempDir(), "oculus-worktree-"+branch)
	cmd := exec.Command("git", "worktree", "add", "-b", branch, worktreePath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Try without -b if branch exists
		cmd = exec.Command("git", "worktree", "add", worktreePath, branch)
		out, err = cmd.CombinedOutput()
		if err != nil {
			return &tool.Result{Data: fmt.Sprintf("Error creating worktree: %s\n%s", err, string(out))}, nil
		}
	}
	t.ActiveWorktree = worktreePath
	return &tool.Result{Data: fmt.Sprintf("Worktree created at %s (branch: %s)", worktreePath, branch)}, nil
}

type ExitWorktreeTool struct {
	tool.BaseTool
	WorktreeRef *string
}

func NewExitWorktreeTool(ref *string) *ExitWorktreeTool {
	return &ExitWorktreeTool{BaseTool: tool.BaseTool{ToolName: "ExitWorktree"}, WorktreeRef: ref}
}

func (t *ExitWorktreeTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{}}
}

func (t *ExitWorktreeTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Remove the git worktree and clean up.", nil
}

func (t *ExitWorktreeTool) Prompt(_ context.Context) (string, error) {
	return "Exit a worktree session and return to original directory.\n- action: keep or remove\n- Only operates on worktrees from this session", nil
}

func (t *ExitWorktreeTool) Call(_ context.Context, _ map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	if t.WorktreeRef == nil || *t.WorktreeRef == "" {
		return &tool.Result{Data: "No active worktree to exit."}, nil
	}
	path := *t.WorktreeRef
	cmd := exec.Command("git", "worktree", "remove", path, "--force")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Error removing worktree: %s\n%s", err, strings.TrimSpace(string(out)))}, nil
	}
	*t.WorktreeRef = ""
	return &tool.Result{Data: fmt.Sprintf("Worktree at %s removed.", path)}, nil
}

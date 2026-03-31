package plan

import (
	"context"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type EnterPlanModeTool struct {
	tool.BaseTool
	InPlanMode *bool
}

func NewEnterPlanModeTool(planMode *bool) *EnterPlanModeTool {
	return &EnterPlanModeTool{BaseTool: tool.BaseTool{ToolName: "EnterPlanMode"}, InPlanMode: planMode}
}

func (t *EnterPlanModeTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{}}
}

func (t *EnterPlanModeTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Enter plan mode for structured planning.", nil
}

func (t *EnterPlanModeTool) Prompt(_ context.Context) (string, error) {
	return "Use this tool when you're about to start a non-trivial implementation task. Getting user sign-off prevents wasted effort.\n\nUse when:\n- New feature implementation\n- Multiple valid approaches\n- Multi-file changes\n- Unclear requirements", nil
}

func (t *EnterPlanModeTool) Call(_ context.Context, _ map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	*t.InPlanMode = true
	return &tool.Result{Data: "Entered plan mode. Present your plan for approval."}, nil
}

type ExitPlanModeTool struct {
	tool.BaseTool
	InPlanMode *bool
}

func NewExitPlanModeTool(planMode *bool) *ExitPlanModeTool {
	return &ExitPlanModeTool{BaseTool: tool.BaseTool{ToolName: "ExitPlanMode"}, InPlanMode: planMode}
}

func (t *ExitPlanModeTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{}}
}

func (t *ExitPlanModeTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Exit plan mode after plan approval.", nil
}

func (t *ExitPlanModeTool) Prompt(_ context.Context) (string, error) {
	return "Use when you have finished writing your plan and are ready for user approval.\n- Write plan to plan file first\n- This tool reads the plan from the file\n- Do NOT use AskUserQuestion to ask 'Is this plan okay?' - that's what THIS tool does", nil
}

func (t *ExitPlanModeTool) Call(_ context.Context, _ map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	*t.InPlanMode = false
	return &tool.Result{Data: "Exited plan mode. Proceeding with implementation."}, nil
}

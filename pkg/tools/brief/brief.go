package brief

import (
	"context"
	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type BriefTool struct {
	tool.BaseTool
	IsBrief *bool
}

func NewBriefTool(isBrief *bool) *BriefTool {
	return &BriefTool{BaseTool: tool.BaseTool{ToolName: "Brief", ToolSearchHint: "brief concise output mode"}, IsBrief: isBrief}
}
func (t *BriefTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"enabled": map[string]interface{}{"type": "boolean", "description": "Enable or disable brief mode"},
	}}
}
func (t *BriefTool) Description(_ context.Context, _ map[string]interface{}) (string, error) { return "Toggle brief output mode.", nil }
func (t *BriefTool) Prompt(_ context.Context) (string, error) {
	return "Toggle brief output mode to reduce verbosity of responses.", nil
}
func (t *BriefTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	if enabled, ok := input["enabled"].(bool); ok {
		*t.IsBrief = enabled
	} else {
		*t.IsBrief = !*t.IsBrief
	}
	if *t.IsBrief { return &tool.Result{Data: "Brief mode enabled."}, nil }
	return &tool.Result{Data: "Brief mode disabled."}, nil
}

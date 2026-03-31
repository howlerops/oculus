package syntheticoutput

import (
	"context"
	"encoding/json"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type SyntheticOutputTool struct{ tool.BaseTool }

func NewSyntheticOutputTool() *SyntheticOutputTool {
	return &SyntheticOutputTool{BaseTool: tool.BaseTool{ToolName: "SyntheticOutput", ToolSearchHint: "structured output json schema"}}
}
func (t *SyntheticOutputTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"data": map[string]interface{}{"description": "Structured data to output"},
	}, Required: []string{"data"}}
}
func (t *SyntheticOutputTool) Description(_ context.Context, _ map[string]interface{}) (string, error) { return "Generate structured output.", nil }
func (t *SyntheticOutputTool) Prompt(_ context.Context) (string, error) {
	return "Emit structured JSON output for programmatic consumption by callers.", nil
}

func (t *SyntheticOutputTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	data := input["data"]
	out, _ := json.MarshalIndent(data, "", "  ")
	return &tool.Result{Data: string(out)}, nil
}

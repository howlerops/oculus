package sleep

import (
	"context"
	"fmt"
	"time"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type SleepTool struct {
	tool.BaseTool
}

func NewSleepTool() *SleepTool {
	return &SleepTool{
		BaseTool: tool.BaseTool{ToolName: "Sleep", ToolSearchHint: "wait pause delay timer"},
	}
}

func (t *SleepTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"duration_ms": map[string]interface{}{"type": "number", "description": "Duration to sleep in milliseconds"},
		},
		Required: []string{"duration_ms"},
	}
}

func (t *SleepTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Sleep for a specified duration. Can be interrupted.", nil
}

func (t *SleepTool) Prompt(_ context.Context) (string, error) {
	return "Wait for a specified duration. User can interrupt at any time.\n\nPrefer this over Bash(sleep ...) - it doesn't hold a shell process.", nil
}

func (t *SleepTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *SleepTool) GetInterruptBehavior() tool.InterruptBehavior { return tool.InterruptCancel }

func (t *SleepTool) Call(ctx context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	durationMs, _ := input["duration_ms"].(float64)
	if durationMs <= 0 {
		return &tool.Result{Data: "Error: duration_ms must be positive"}, nil
	}
	if durationMs > 3600000 { // 1 hour max
		durationMs = 3600000
	}

	duration := time.Duration(durationMs) * time.Millisecond
	select {
	case <-time.After(duration):
		return &tool.Result{Data: fmt.Sprintf("Slept for %v", duration)}, nil
	case <-ctx.Done():
		return &tool.Result{Data: "Sleep interrupted"}, nil
	}
}

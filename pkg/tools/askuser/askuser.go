package askuser

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type AskUserQuestionTool struct {
	tool.BaseTool
}

func NewAskUserQuestionTool() *AskUserQuestionTool {
	return &AskUserQuestionTool{
		BaseTool: tool.BaseTool{
			ToolName:       "AskUserQuestion",
			ToolSearchHint: "ask user question prompt interactive",
		},
	}
}

func (t *AskUserQuestionTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"questions": map[string]interface{}{
				"type":        "array",
				"description": "Questions to ask (1-4)",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"question":    map[string]interface{}{"type": "string"},
						"header":      map[string]interface{}{"type": "string"},
						"options":     map[string]interface{}{"type": "array"},
						"multiSelect": map[string]interface{}{"type": "boolean"},
					},
				},
			},
			"answers": map[string]interface{}{
				"type":        "object",
				"description": "Pre-filled answers",
			},
		},
		Required: []string{"questions"},
	}
}

func (t *AskUserQuestionTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Ask the user questions during execution.", nil
}

func (t *AskUserQuestionTool) Prompt(_ context.Context) (string, error) {
	return "Ask the user questions during execution.\n\nUse to:\n1. Gather user preferences or requirements\n2. Clarify ambiguous instructions\n3. Get decisions on implementation choices\n4. Offer choices about direction\n\nNotes:\n- Users can always select 'Other' for custom input\n- Use multiSelect: true for multiple selections\n- Put recommended option first with '(Recommended)' label", nil
}

func (t *AskUserQuestionTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	questions, _ := input["questions"]
	answers, hasAnswers := input["answers"].(map[string]interface{})

	if !hasAnswers || len(answers) == 0 {
		// In non-interactive mode, format questions for display and prompt via stdout
		questionsJSON, _ := json.MarshalIndent(questions, "", "  ")
		return &tool.Result{
			Data: fmt.Sprintf("Questions presented to user:\n%s\n\nWaiting for answers...", string(questionsJSON)),
		}, nil
	}

	// Format answers
	var parts []string
	for q, a := range answers {
		parts = append(parts, fmt.Sprintf("%q=%q", q, a))
	}

	return &tool.Result{
		Data: fmt.Sprintf("User has answered your questions: %s", strings.Join(parts, ", ")),
	}, nil
}

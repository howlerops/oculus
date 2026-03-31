package todowrite

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type TodoItem struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Status  string `json:"status"` // "pending", "in_progress", "completed"
}

type TodoWriteTool struct {
	tool.BaseTool
	// Todos stored in memory (would be in AppState in full impl)
	Todos []TodoItem
}

func NewTodoWriteTool() *TodoWriteTool {
	return &TodoWriteTool{
		BaseTool: tool.BaseTool{ToolName: "TodoWrite", ToolSearchHint: "todo task list manage write"},
	}
}

func (t *TodoWriteTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"todos": map[string]interface{}{
				"type": "array", "description": "Complete replacement todo list",
				"items": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"id":      map[string]interface{}{"type": "string"},
						"content": map[string]interface{}{"type": "string"},
						"status":  map[string]interface{}{"type": "string"},
					},
				},
			},
		},
		Required: []string{"todos"},
	}
}

func (t *TodoWriteTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Write and manage todo items for tracking progress.", nil
}

func (t *TodoWriteTool) Prompt(_ context.Context) (string, error) {
	return "Create and manage a structured task list for your coding session.\n\nUse when:\n- Complex multi-step tasks (3+ steps)\n- User provides multiple tasks\n- After receiving new instructions\n- Mark in_progress BEFORE starting work\n- Mark completed after finishing", nil
}

func (t *TodoWriteTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	todosRaw, _ := input["todos"]
	data, _ := json.Marshal(todosRaw)
	var newTodos []TodoItem
	json.Unmarshal(data, &newTodos)

	oldTodos := t.Todos
	// If all completed, clear
	allComplete := len(newTodos) > 0
	for _, td := range newTodos {
		if td.Status != "completed" {
			allComplete = false
			break
		}
	}
	if allComplete {
		t.Todos = nil
	} else {
		t.Todos = newTodos
	}

	completed := 0
	for _, td := range newTodos {
		if td.Status == "completed" {
			completed++
		}
	}

	return &tool.Result{
		Data: fmt.Sprintf("Updated todos: %d items (%d completed, was %d items). %s",
			len(newTodos), completed, len(oldTodos),
			func() string {
				if allComplete {
					return "All complete - list cleared."
				}
				return ""
			}()),
	}, nil
}

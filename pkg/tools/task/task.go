package task

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	apptask "github.com/howlerops/oculus/pkg/task"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

// TaskCreateTool creates background tasks
type TaskCreateTool struct {
	tool.BaseTool
	Tasks map[string]*apptask.TaskState
}

func NewTaskCreateTool() *TaskCreateTool {
	return &TaskCreateTool{
		BaseTool: tool.BaseTool{ToolName: "TaskCreate", ToolSearchHint: "create background task spawn"},
		Tasks:    make(map[string]*apptask.TaskState),
	}
}

func (t *TaskCreateTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"subject":     map[string]interface{}{"type": "string", "description": "Brief task title"},
			"description": map[string]interface{}{"type": "string", "description": "What to do"},
		},
		Required: []string{"subject", "description"},
	}
}

func (t *TaskCreateTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Create a new task for tracking work.", nil
}

func (t *TaskCreateTool) Prompt(_ context.Context) (string, error) {
	return "Create structured tasks to track progress.\n- Use for complex multi-step work\n- All tasks created as pending", nil
}

func (t *TaskCreateTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	subject, _ := input["subject"].(string)
	description, _ := input["description"].(string)
	id := apptask.GenerateTaskId(apptask.TaskTypeLocalBash)
	state := apptask.NewTaskState(id, apptask.TaskTypeLocalBash, description, "")
	t.Tasks[id] = &state
	return &tool.Result{Data: fmt.Sprintf("Task created: %s (%s)", id, subject)}, nil
}

// TaskGetTool retrieves task state
type TaskGetTool struct {
	tool.BaseTool
	Tasks map[string]*apptask.TaskState
}

func NewTaskGetTool(tasks map[string]*apptask.TaskState) *TaskGetTool {
	return &TaskGetTool{
		BaseTool: tool.BaseTool{ToolName: "TaskGet"},
		Tasks:    tasks,
	}
}

func (t *TaskGetTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"task_id": map[string]interface{}{"type": "string"},
	}, Required: []string{"task_id"}}
}

func (t *TaskGetTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Get task details by ID.", nil
}

func (t *TaskGetTool) Prompt(_ context.Context) (string, error) {
	return "Retrieve a task by ID. Returns full details including dependencies.", nil
}

func (t *TaskGetTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *TaskGetTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	id, _ := input["task_id"].(string)
	task, ok := t.Tasks[id]
	if !ok {
		return &tool.Result{Data: fmt.Sprintf("Task %s not found", id)}, nil
	}
	return &tool.Result{Data: fmt.Sprintf("Task %s: status=%s desc=%q started=%s", task.ID, task.Status, task.Description, task.StartTime.Format(time.RFC3339))}, nil
}

// TaskUpdateTool updates task status
type TaskUpdateTool struct {
	tool.BaseTool
	Tasks map[string]*apptask.TaskState
}

func NewTaskUpdateTool(tasks map[string]*apptask.TaskState) *TaskUpdateTool {
	return &TaskUpdateTool{BaseTool: tool.BaseTool{ToolName: "TaskUpdate"}, Tasks: tasks}
}

func (t *TaskUpdateTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"task_id": map[string]interface{}{"type": "string"},
		"status":  map[string]interface{}{"type": "string", "description": "pending/in_progress/completed"},
	}, Required: []string{"task_id", "status"}}
}

func (t *TaskUpdateTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Update a task's status.", nil
}

func (t *TaskUpdateTool) Prompt(_ context.Context) (string, error) {
	return "Update a task status. Mark in_progress before starting, completed when done. ONLY mark completed when FULLY accomplished.", nil
}

func (t *TaskUpdateTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	id, _ := input["task_id"].(string)
	status, _ := input["status"].(string)
	task, ok := t.Tasks[id]
	if !ok {
		return &tool.Result{Data: fmt.Sprintf("Task %s not found", id)}, nil
	}
	task.Status = apptask.TaskStatus(status)
	if apptask.IsTerminalTaskStatus(task.Status) {
		now := time.Now()
		task.EndTime = &now
	}
	return &tool.Result{Data: fmt.Sprintf("Task %s updated to %s", id, status)}, nil
}

// TaskListTool lists all tasks
type TaskListTool struct {
	tool.BaseTool
	Tasks map[string]*apptask.TaskState
}

func NewTaskListTool(tasks map[string]*apptask.TaskState) *TaskListTool {
	return &TaskListTool{BaseTool: tool.BaseTool{ToolName: "TaskList"}, Tasks: tasks}
}

func (t *TaskListTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{}}
}

func (t *TaskListTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "List all tasks.", nil
}

func (t *TaskListTool) Prompt(_ context.Context) (string, error) {
	return "List all tasks. Check for available work and overall progress.", nil
}

func (t *TaskListTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *TaskListTool) Call(_ context.Context, _ map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	if len(t.Tasks) == 0 {
		return &tool.Result{Data: "No tasks."}, nil
	}
	var lines []string
	for _, task := range t.Tasks {
		lines = append(lines, fmt.Sprintf("  %s [%s] %s", task.ID, task.Status, task.Description))
	}
	return &tool.Result{Data: fmt.Sprintf("Tasks (%d):\n%s", len(t.Tasks), strings.Join(lines, "\n"))}, nil
}

// TaskStopTool stops a running task
type TaskStopTool struct {
	tool.BaseTool
	Tasks map[string]*apptask.TaskState
}

func NewTaskStopTool(tasks map[string]*apptask.TaskState) *TaskStopTool {
	return &TaskStopTool{BaseTool: tool.BaseTool{ToolName: "TaskStop"}, Tasks: tasks}
}

func (t *TaskStopTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"task_id": map[string]interface{}{"type": "string"},
	}, Required: []string{"task_id"}}
}

func (t *TaskStopTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Stop a running task.", nil
}

func (t *TaskStopTool) Prompt(_ context.Context) (string, error) {
	return "Stop a running background task by ID.", nil
}

func (t *TaskStopTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	id, _ := input["task_id"].(string)
	task, ok := t.Tasks[id]
	if !ok {
		return &tool.Result{Data: fmt.Sprintf("Task %s not found", id)}, nil
	}
	task.Status = apptask.TaskStatusKilled
	now := time.Now()
	task.EndTime = &now
	return &tool.Result{Data: fmt.Sprintf("Task %s stopped", id)}, nil
}

// TaskOutputTool reads task output
type TaskOutputTool struct {
	tool.BaseTool
	Tasks map[string]*apptask.TaskState
}

func NewTaskOutputTool(tasks map[string]*apptask.TaskState) *TaskOutputTool {
	return &TaskOutputTool{BaseTool: tool.BaseTool{ToolName: "TaskOutput", ToolSearchHint: "read background task output"}, Tasks: tasks}
}

func (t *TaskOutputTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"task_id": map[string]interface{}{"type": "string"},
	}, Required: []string{"task_id"}}
}

func (t *TaskOutputTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Read output from a background task.", nil
}

func (t *TaskOutputTool) Prompt(_ context.Context) (string, error) {
	return "Read output from a background task.", nil
}

func (t *TaskOutputTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *TaskOutputTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	id, _ := input["task_id"].(string)
	task, ok := t.Tasks[id]
	if !ok {
		return &tool.Result{Data: fmt.Sprintf("Task %s not found", id)}, nil
	}
	data, err := os.ReadFile(task.OutputFile)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Task %s: no output file (status: %s)", id, task.Status)}, nil
	}
	return &tool.Result{Data: string(data)}, nil
}

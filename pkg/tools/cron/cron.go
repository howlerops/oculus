package cron

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

// cronJob holds an in-memory cron job record.
type cronJob struct {
	ID       string
	Schedule string
	Command  string
	Created  time.Time
}

var (
	cronMu   sync.Mutex
	cronJobs = make(map[string]*cronJob)
)

// ScheduleCronTool is the base cron scheduling tool (legacy name, kept for compatibility).
type ScheduleCronTool struct {
	tool.BaseTool
}

func NewScheduleCronTool() *ScheduleCronTool {
	return &ScheduleCronTool{BaseTool: tool.BaseTool{ToolName: "ScheduleCron", ToolSearchHint: "schedule cron recurring task"}}
}

func (t *ScheduleCronTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"schedule": map[string]interface{}{"type": "string", "description": "Cron expression (e.g. '0 * * * *')"},
			"command":  map[string]interface{}{"type": "string", "description": "Shell command to run on schedule"},
		},
		Required: []string{"schedule", "command"},
	}
}

func (t *ScheduleCronTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Schedule a recurring command using a cron expression.", nil
}

func (t *ScheduleCronTool) Prompt(_ context.Context) (string, error) {
	return "Schedule a recurring shell command using a cron expression.", nil
}

func (t *ScheduleCronTool) Call(ctx context.Context, input map[string]interface{}, progress func(types.ToolProgressData)) (*tool.Result, error) {
	ct := &CronCreateTool{BaseTool: tool.BaseTool{ToolName: "ScheduleCron"}}
	return ct.Call(ctx, input, progress)
}

// CronCreateTool creates a new cron job entry.
type CronCreateTool struct {
	tool.BaseTool
}

func NewCronCreateTool() *CronCreateTool {
	return &CronCreateTool{
		BaseTool: tool.BaseTool{
			ToolName:       "CronCreate",
			ToolSearchHint: "create schedule recurring cron job",
		},
	}
}

func (t *CronCreateTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"schedule": map[string]interface{}{
				"type":        "string",
				"description": "Cron expression (e.g. '0 9 * * 1' for every Monday at 9am)",
			},
			"command": map[string]interface{}{
				"type":        "string",
				"description": "Shell command to run",
			},
		},
		Required: []string{"schedule", "command"},
	}
}

func (t *CronCreateTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Create a new cron job that runs a command on a schedule.", nil
}

func (t *CronCreateTool) Prompt(_ context.Context) (string, error) {
	return "Create a new cron job that runs a shell command on a recurring schedule.", nil
}

func (t *CronCreateTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	schedule, _ := input["schedule"].(string)
	command, _ := input["command"].(string)
	if schedule == "" {
		return &tool.Result{Data: "Error: schedule is required"}, nil
	}
	if command == "" {
		return &tool.Result{Data: "Error: command is required"}, nil
	}

	id := fmt.Sprintf("cron_%d", time.Now().UnixNano())
	cronMu.Lock()
	cronJobs[id] = &cronJob{
		ID:       id,
		Schedule: schedule,
		Command:  command,
		Created:  time.Now(),
	}
	cronMu.Unlock()

	return &tool.Result{
		Data: fmt.Sprintf("Cron job created: id=%s schedule=%q command=%q", id, schedule, command),
	}, nil
}

// CronDeleteTool removes an existing cron job by ID.
type CronDeleteTool struct {
	tool.BaseTool
}

func NewCronDeleteTool() *CronDeleteTool {
	return &CronDeleteTool{
		BaseTool: tool.BaseTool{
			ToolName:       "CronDelete",
			ToolSearchHint: "delete remove cron job",
		},
	}
}

func (t *CronDeleteTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"id": map[string]interface{}{
				"type":        "string",
				"description": "ID of the cron job to delete",
			},
		},
		Required: []string{"id"},
	}
}

func (t *CronDeleteTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Delete a cron job by its ID.", nil
}

func (t *CronDeleteTool) Prompt(_ context.Context) (string, error) {
	return "Delete a registered cron job by its ID.", nil
}

func (t *CronDeleteTool) IsDestructive(_ map[string]interface{}) bool { return true }

func (t *CronDeleteTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	id, _ := input["id"].(string)
	if id == "" {
		return &tool.Result{Data: "Error: id is required"}, nil
	}

	cronMu.Lock()
	_, ok := cronJobs[id]
	if ok {
		delete(cronJobs, id)
	}
	cronMu.Unlock()

	if !ok {
		return &tool.Result{Data: fmt.Sprintf("Error: cron job %q not found", id)}, nil
	}
	return &tool.Result{Data: fmt.Sprintf("Cron job %q deleted", id)}, nil
}

// CronListTool lists all current cron jobs.
type CronListTool struct {
	tool.BaseTool
}

func NewCronListTool() *CronListTool {
	return &CronListTool{
		BaseTool: tool.BaseTool{
			ToolName:       "CronList",
			ToolSearchHint: "list cron jobs scheduled tasks",
		},
	}
}

func (t *CronListTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{}}
}

func (t *CronListTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "List all currently registered cron jobs.", nil
}

func (t *CronListTool) Prompt(_ context.Context) (string, error) {
	return "List all currently registered cron jobs and their schedules.", nil
}

func (t *CronListTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *CronListTool) Call(_ context.Context, _ map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	cronMu.Lock()
	jobs := make([]*cronJob, 0, len(cronJobs))
	for _, j := range cronJobs {
		jobs = append(jobs, j)
	}
	cronMu.Unlock()

	if len(jobs) == 0 {
		return &tool.Result{Data: "No cron jobs registered."}, nil
	}

	var lines []string
	lines = append(lines, fmt.Sprintf("Cron jobs (%d):", len(jobs)))
	for _, j := range jobs {
		lines = append(lines, fmt.Sprintf("  %s  schedule=%q  command=%q  created=%s",
			j.ID, j.Schedule, j.Command, j.Created.Format(time.RFC3339)))
	}
	return &tool.Result{Data: strings.Join(lines, "\n")}, nil
}

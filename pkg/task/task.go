package task

import (
	"crypto/rand"
	"path/filepath"
	"time"
)

// TaskType identifies the kind of background task
type TaskType string

const (
	TaskTypeLocalBash         TaskType = "local_bash"
	TaskTypeLocalAgent        TaskType = "local_agent"
	TaskTypeRemoteAgent       TaskType = "remote_agent"
	TaskTypeInProcessTeammate TaskType = "in_process_teammate"
	TaskTypeLocalWorkflow     TaskType = "local_workflow"
	TaskTypeMonitorMCP        TaskType = "monitor_mcp"
	TaskTypeDream             TaskType = "dream"
)

// TaskStatus is the lifecycle state of a task
type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusKilled    TaskStatus = "killed"
)

// IsTerminalTaskStatus returns true when a task won't transition further
func IsTerminalTaskStatus(status TaskStatus) bool {
	return status == TaskStatusCompleted || status == TaskStatusFailed || status == TaskStatusKilled
}

// TaskState holds the full state of a task
type TaskState struct {
	ID            string     `json:"id"`
	Type          TaskType   `json:"type"`
	Status        TaskStatus `json:"status"`
	Description   string     `json:"description"`
	ToolUseID     string     `json:"toolUseId,omitempty"`
	StartTime     time.Time  `json:"startTime"`
	EndTime       *time.Time `json:"endTime,omitempty"`
	TotalPausedMs int64      `json:"totalPausedMs,omitempty"`
	OutputFile    string     `json:"outputFile"`
	OutputOffset  int64      `json:"outputOffset"`
	Notified      bool       `json:"notified"`
}

// TaskHandle is a reference to a running task
type TaskHandle struct {
	TaskID  string
	Cleanup func()
}

var taskIDPrefixes = map[TaskType]string{
	TaskTypeLocalBash:         "b",
	TaskTypeLocalAgent:        "a",
	TaskTypeRemoteAgent:       "r",
	TaskTypeInProcessTeammate: "t",
	TaskTypeLocalWorkflow:     "w",
	TaskTypeMonitorMCP:        "m",
	TaskTypeDream:             "d",
}

const taskIDAlphabet = "0123456789abcdefghijklmnopqrstuvwxyz"

// GenerateTaskId creates a new task ID with type-based prefix
func GenerateTaskId(taskType TaskType) string {
	prefix, ok := taskIDPrefixes[taskType]
	if !ok {
		prefix = "x"
	}
	bytes := make([]byte, 8)
	_, _ = rand.Read(bytes)
	id := prefix
	for i := 0; i < 8; i++ {
		id += string(taskIDAlphabet[bytes[i]%byte(len(taskIDAlphabet))])
	}
	return id
}

// GetTaskOutputPath returns the file path for task output
func GetTaskOutputPath(taskID string) string {
	return filepath.Join(".oculus", "tasks", taskID+".output")
}

// NewTaskState creates a new task in pending status
func NewTaskState(id string, taskType TaskType, description string, toolUseID string) TaskState {
	return TaskState{
		ID:          id,
		Type:        taskType,
		Status:      TaskStatusPending,
		Description: description,
		ToolUseID:   toolUseID,
		StartTime:   time.Now(),
		OutputFile:  GetTaskOutputPath(id),
	}
}

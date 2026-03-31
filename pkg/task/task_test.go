package task

import (
	"strings"
	"testing"
)

func TestGenerateTaskId(t *testing.T) {
	tests := []struct {
		taskType       TaskType
		expectedPrefix string
	}{
		{TaskTypeLocalBash, "b"},
		{TaskTypeLocalAgent, "a"},
		{TaskTypeRemoteAgent, "r"},
		{TaskTypeInProcessTeammate, "t"},
		{TaskTypeLocalWorkflow, "w"},
		{TaskTypeMonitorMCP, "m"},
		{TaskTypeDream, "d"},
		{TaskType("unknown"), "x"},
	}

	for _, tt := range tests {
		t.Run(string(tt.taskType), func(t *testing.T) {
			id := GenerateTaskId(tt.taskType)
			if !strings.HasPrefix(id, tt.expectedPrefix) {
				t.Errorf("GenerateTaskId(%s) = %s, want prefix %s", tt.taskType, id, tt.expectedPrefix)
			}
			if len(id) != 9 { // 1 prefix + 8 random
				t.Errorf("GenerateTaskId(%s) length = %d, want 9", tt.taskType, len(id))
			}
		})
	}

	// Test uniqueness
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := GenerateTaskId(TaskTypeLocalBash)
		if ids[id] {
			t.Errorf("GenerateTaskId produced duplicate: %s", id)
		}
		ids[id] = true
	}
}

func TestIsTerminalTaskStatus(t *testing.T) {
	tests := []struct {
		status   TaskStatus
		expected bool
	}{
		{TaskStatusPending, false},
		{TaskStatusRunning, false},
		{TaskStatusCompleted, true},
		{TaskStatusFailed, true},
		{TaskStatusKilled, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := IsTerminalTaskStatus(tt.status); got != tt.expected {
				t.Errorf("IsTerminalTaskStatus(%s) = %v, want %v", tt.status, got, tt.expected)
			}
		})
	}
}

func TestNewTaskState(t *testing.T) {
	state := NewTaskState("b12345678", TaskTypeLocalBash, "test task", "tool-123")

	if state.ID != "b12345678" {
		t.Errorf("ID = %s, want b12345678", state.ID)
	}
	if state.Type != TaskTypeLocalBash {
		t.Errorf("Type = %s, want local_bash", state.Type)
	}
	if state.Status != TaskStatusPending {
		t.Errorf("Status = %s, want pending", state.Status)
	}
	if state.Description != "test task" {
		t.Errorf("Description = %s, want 'test task'", state.Description)
	}
	if state.ToolUseID != "tool-123" {
		t.Errorf("ToolUseID = %s, want tool-123", state.ToolUseID)
	}
	if state.StartTime.IsZero() {
		t.Error("StartTime should not be zero")
	}
	if state.EndTime != nil {
		t.Error("EndTime should be nil for new task")
	}
}

func TestGetTaskOutputPath(t *testing.T) {
	path := GetTaskOutputPath("b12345678")
	if !strings.Contains(path, "b12345678.output") {
		t.Errorf("GetTaskOutputPath = %s, want to contain b12345678.output", path)
	}
}

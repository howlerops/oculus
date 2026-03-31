package state

import (
	"github.com/howlerops/oculus/pkg/task"
	"github.com/howlerops/oculus/pkg/types"
)

// GetRunningTasks returns all tasks that are currently running
func GetRunningTasks(state AppState) []task.TaskState {
	var running []task.TaskState
	for _, t := range state.Tasks {
		if t.Status == task.TaskStatusRunning {
			running = append(running, t)
		}
	}
	return running
}

// GetPendingTasks returns all tasks that are pending
func GetPendingTasks(state AppState) []task.TaskState {
	var pending []task.TaskState
	for _, t := range state.Tasks {
		if t.Status == task.TaskStatusPending {
			pending = append(pending, t)
		}
	}
	return pending
}

// GetActiveTasks returns all non-terminal tasks
func GetActiveTasks(state AppState) []task.TaskState {
	var active []task.TaskState
	for _, t := range state.Tasks {
		if !task.IsTerminalTaskStatus(t.Status) {
			active = append(active, t)
		}
	}
	return active
}

// HasInProgressTools returns true if any tool is currently executing
func HasInProgressTools(state AppState) bool {
	return len(state.InProgressToolUseIDs) > 0
}

// GetMessageCount returns the number of messages in the conversation
func GetMessageCount(state AppState) int {
	return len(state.Messages)
}

// GetLastMessage returns the last message, or nil if empty
func GetLastMessage(state AppState) *types.Message {
	if len(state.Messages) == 0 {
		return nil
	}
	msg := state.Messages[len(state.Messages)-1]
	return &msg
}

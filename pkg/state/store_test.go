package state

import (
	"sync"
	"testing"

	"github.com/howlerops/oculus/pkg/task"
	"github.com/howlerops/oculus/pkg/types"
)

func TestNewStore(t *testing.T) {
	initial := NewAppState("claude-sonnet-4-20250514")
	store := NewStore(initial)

	state := store.Get()
	if state.MainLoopModel != "claude-sonnet-4-20250514" {
		t.Errorf("MainLoopModel = %s, want claude-sonnet-4-20250514", state.MainLoopModel)
	}
}

func TestStoreUpdate(t *testing.T) {
	store := NewStore(NewAppState("test-model"))

	store.Update(func(prev AppState) AppState {
		prev.Verbose = true
		prev.TotalInputTokens = 100
		return prev
	})

	state := store.Get()
	if !state.Verbose {
		t.Error("Verbose should be true after update")
	}
	if state.TotalInputTokens != 100 {
		t.Errorf("TotalInputTokens = %d, want 100", state.TotalInputTokens)
	}
}

func TestStoreConcurrentAccess(t *testing.T) {
	store := NewStore(NewAppState("test-model"))

	var wg sync.WaitGroup

	// Concurrent writers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			store.Update(func(prev AppState) AppState {
				prev.TotalInputTokens += 1
				return prev
			})
		}(i)
	}

	// Concurrent readers
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_ = store.Get()
		}()
	}

	wg.Wait()

	state := store.Get()
	if state.TotalInputTokens != 100 {
		t.Errorf("TotalInputTokens = %d, want 100 after 100 concurrent increments", state.TotalInputTokens)
	}
}

func TestStoreSubscribe(t *testing.T) {
	store := NewStore(NewAppState("test-model"))

	called := false
	unsub := store.Subscribe(func(state AppState) {
		called = true
	})
	defer unsub()

	store.Update(func(prev AppState) AppState {
		prev.Verbose = true
		return prev
	})

	if !called {
		t.Error("Subscriber should have been called after update")
	}
}

func TestGetRunningTasks(t *testing.T) {
	state := NewAppState("test")
	state.Tasks["task1"] = task.TaskState{ID: "task1", Status: task.TaskStatusRunning}
	state.Tasks["task2"] = task.TaskState{ID: "task2", Status: task.TaskStatusCompleted}
	state.Tasks["task3"] = task.TaskState{ID: "task3", Status: task.TaskStatusRunning}

	running := GetRunningTasks(state)
	if len(running) != 2 {
		t.Errorf("GetRunningTasks returned %d tasks, want 2", len(running))
	}
}

func TestGetLastMessage(t *testing.T) {
	state := NewAppState("test")

	// Empty
	if msg := GetLastMessage(state); msg != nil {
		t.Error("GetLastMessage should return nil for empty messages")
	}

	// With messages
	state.Messages = append(state.Messages, types.NewUserMessage("hello"))
	state.Messages = append(state.Messages, types.NewUserMessage("world"))

	msg := GetLastMessage(state)
	if msg == nil {
		t.Fatal("GetLastMessage should not be nil")
	}
	if msg.User.Content[0].Text != "world" {
		t.Errorf("GetLastMessage text = %s, want 'world'", msg.User.Content[0].Text)
	}
}

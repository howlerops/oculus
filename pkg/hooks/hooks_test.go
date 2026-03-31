package hooks

import (
	"context"
	"testing"

	"github.com/howlerops/oculus/pkg/types"
)

func TestNewRegistry(t *testing.T) {
	reg := NewRegistry()
	hooks := reg.GetHooks(types.HookEventPreToolUse)
	if len(hooks) != 0 {
		t.Errorf("expected 0 hooks, got %d", len(hooks))
	}
}

func TestRegister(t *testing.T) {
	reg := NewRegistry()
	reg.Register(types.HookEventPreToolUse, HookConfig{
		Command: "echo test",
	})
	hooks := reg.GetHooks(types.HookEventPreToolUse)
	if len(hooks) != 1 {
		t.Errorf("expected 1 hook, got %d", len(hooks))
	}
}

func TestExecuteNoHooks(t *testing.T) {
	reg := NewRegistry()
	result, err := reg.Execute(context.Background(), types.HookEventPreToolUse, HookInput{})
	if err != nil {
		t.Fatal(err)
	}
	if result.PreventContinuation {
		t.Error("expected no prevention with no hooks")
	}
}

func TestMatchesHookPattern(t *testing.T) {
	tests := []struct {
		matcher  string
		toolName string
		expected bool
	}{
		{"Bash", "Bash", true},
		{"Bash", "Read", false},
		{"*", "Bash", true},
		{"Bash(git *)", "Bash", true},
	}
	for _, tt := range tests {
		got := matchesHookPattern(tt.matcher, tt.toolName, nil)
		if got != tt.expected {
			t.Errorf("matchesHookPattern(%q, %q) = %v, want %v", tt.matcher, tt.toolName, got, tt.expected)
		}
	}
}

package bash

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/howlerops/oculus/pkg/types"
)

func newTool() *BashTool {
	return NewBashTool("")
}

func TestSimpleEcho(t *testing.T) {
	bt := newTool()
	result, err := bt.Call(context.Background(), map[string]interface{}{
		"command": "echo hello",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Data.(string), "hello") {
		t.Errorf("expected 'hello' in output, got: %s", result.Data)
	}
}

func TestLsCommand(t *testing.T) {
	bt := newTool()
	result, err := bt.Call(context.Background(), map[string]interface{}{
		"command": "ls /tmp",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Data.(string) == "" {
		t.Error("expected non-empty output from ls")
	}
}

func TestExitCodeCapture(t *testing.T) {
	bt := newTool()
	result, err := bt.Call(context.Background(), map[string]interface{}{
		"command": "exit 42",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Non-zero exit code should appear in the result
	if !strings.Contains(result.Data.(string), "42") {
		t.Errorf("expected exit code 42 in output, got: %s", result.Data)
	}
}

func TestStderrCapture(t *testing.T) {
	bt := newTool()
	result, err := bt.Call(context.Background(), map[string]interface{}{
		"command": "echo error-output >&2",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Data.(string), "error-output") {
		t.Errorf("expected stderr in output, got: %s", result.Data)
	}
}

func TestTimeout(t *testing.T) {
	bt := newTool()
	start := time.Now()
	result, err := bt.Call(context.Background(), map[string]interface{}{
		"command": "sleep 10",
		"timeout": float64(200), // 200ms
	}, nil)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed > 2*time.Second {
		t.Errorf("timeout did not fire in time, elapsed: %v", elapsed)
	}
	if !strings.Contains(result.Data.(string), "timed out") {
		t.Errorf("expected timeout message, got: %s", result.Data)
	}
}

func TestEmptyCommand(t *testing.T) {
	bt := newTool()
	result, err := bt.Call(context.Background(), map[string]interface{}{
		"command": "",
	}, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(result.Data.(string), "Error") {
		t.Errorf("expected error message for empty command, got: %s", result.Data)
	}
}

func TestProgressCallback(t *testing.T) {
	bt := newTool()
	var progressEvents []types.ToolProgressData

	_, err := bt.Call(context.Background(), map[string]interface{}{
		"command": "echo progress-test",
	}, func(p types.ToolProgressData) {
		progressEvents = append(progressEvents, p)
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Expect at least 2 events: initial + completion
	if len(progressEvents) < 2 {
		t.Errorf("expected at least 2 progress events, got %d", len(progressEvents))
	}
	// All events should have bash progress type
	for _, ev := range progressEvents {
		if ev.Type != types.ProgressTypeBash {
			t.Errorf("expected progress type %s, got %s", types.ProgressTypeBash, ev.Type)
		}
	}
}

func TestToolMetadata(t *testing.T) {
	bt := newTool()

	if bt.Name() != "Bash" {
		t.Errorf("expected name 'Bash', got %s", bt.Name())
	}
	if !bt.IsConcurrencySafe(nil) {
		t.Error("expected IsConcurrencySafe to be true")
	}
	if bt.IsReadOnly(nil) {
		t.Error("expected IsReadOnly to be false")
	}
	if bt.MaxResultSizeChars() != 30000 {
		t.Errorf("expected MaxResultSizeChars 30000, got %d", bt.MaxResultSizeChars())
	}
}

func TestUserFacingName(t *testing.T) {
	bt := newTool()

	name := bt.UserFacingName(map[string]interface{}{"command": "echo hi"})
	if name != "Bash: echo hi" {
		t.Errorf("unexpected name: %s", name)
	}

	longCmd := strings.Repeat("x", 50)
	name = bt.UserFacingName(map[string]interface{}{"command": longCmd})
	if !strings.HasSuffix(name, "...") {
		t.Errorf("expected truncated name to end with '...', got: %s", name)
	}
}

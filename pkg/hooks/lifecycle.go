package hooks

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/howlerops/oculus/pkg/types"
)

// SessionStartHook holds data passed to session start hooks
type SessionStartHook struct {
	SessionID string
	CWD       string
	Model     string
	StartTime time.Time
}

// RunSessionStartHooks executes all session start hooks
func (r *Registry) RunSessionStartHooks(ctx context.Context, session SessionStartHook) error {
	input := HookInput{
		SessionID: session.SessionID,
		Event:     types.HookEventSessionStart,
	}
	_, err := r.Execute(ctx, types.HookEventSessionStart, input)
	return err
}

// StopHookResult is the outcome of running stop hooks
type StopHookResult struct {
	ShouldStop bool
	Reason     string
	Summary    string
}

// RunStopHooks executes post-sampling stop hooks
func (r *Registry) RunStopHooks(ctx context.Context, sessionID string, lastAssistantText string) (*StopHookResult, error) {
	input := HookInput{
		SessionID: sessionID,
		Event:     types.HookEventStop,
	}
	result, err := r.Execute(ctx, types.HookEventStop, input)
	if err != nil {
		return nil, err
	}

	return &StopHookResult{
		ShouldStop: !result.PreventContinuation,
		Reason:     result.StopReason,
	}, nil
}

// FileChangedWatcher watches a set of paths for modification and triggers a callback
type FileChangedWatcher struct {
	paths    []string
	interval time.Duration
	modTimes map[string]time.Time
	onChange func(path string)
	done     chan struct{}
}

// NewFileChangedWatcher creates a watcher for the given paths, polling at interval
func NewFileChangedWatcher(paths []string, interval time.Duration, onChange func(string)) *FileChangedWatcher {
	return &FileChangedWatcher{
		paths:    paths,
		interval: interval,
		modTimes: make(map[string]time.Time),
		onChange: onChange,
		done:     make(chan struct{}),
	}
}

// Start begins polling in a background goroutine
func (w *FileChangedWatcher) Start() {
	// Seed initial mod times so we only fire on changes after Start
	for _, path := range w.paths {
		info, err := os.Stat(path)
		if err == nil {
			w.modTimes[path] = info.ModTime()
		}
	}

	go func() {
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				for _, path := range w.paths {
					info, err := os.Stat(path)
					if err != nil {
						continue
					}
					if prev, ok := w.modTimes[path]; ok {
						if info.ModTime().After(prev) {
							w.modTimes[path] = info.ModTime()
							if w.onChange != nil {
								w.onChange(path)
							}
						}
					} else {
						w.modTimes[path] = info.ModTime()
					}
				}
			case <-w.done:
				return
			}
		}
	}()
}

// Stop halts the polling goroutine
func (w *FileChangedWatcher) Stop() { close(w.done) }

// InputProcessor processes user input before it is sent to the API
type InputProcessor struct{}

// ProcessTextPrompt handles slash commands, ! shell shortcuts, and plain text.
// Returns (processedText, isCommand, commandName).
func (p *InputProcessor) ProcessTextPrompt(text string) (processedText string, isCommand bool, commandName string) {
	text = strings.TrimSpace(text)

	// Slash command: /name [rest]
	if strings.HasPrefix(text, "/") {
		parts := strings.SplitN(text, " ", 2)
		commandName = strings.TrimPrefix(parts[0], "/")
		isCommand = true
		if len(parts) > 1 {
			processedText = parts[1]
		}
		return
	}

	// Shell shorthand: !<cmd>
	if strings.HasPrefix(text, "!") {
		cmd := strings.TrimSpace(strings.TrimPrefix(text, "!"))
		if cmd != "" {
			out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
			if err != nil {
				processedText = fmt.Sprintf("Shell command output:\n```\n%s\nError: %v\n```", string(out), err)
			} else {
				processedText = fmt.Sprintf("Shell command output:\n```\n%s```", string(out))
			}
			return processedText, false, ""
		}
	}

	return text, false, ""
}

package hooks

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"golang.org/x/term"
)

// TerminalSize represents terminal dimensions
type TerminalSize struct {
	Columns int
	Rows    int
}

// GetTerminalSize returns current terminal dimensions
func GetTerminalSize() TerminalSize {
	cols, rows, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return TerminalSize{Columns: 80, Rows: 24}
	}
	return TerminalSize{Columns: cols, Rows: rows}
}

// OnTerminalResize calls the callback when terminal size changes
func OnTerminalResize(callback func(TerminalSize)) func() {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-ch:
				callback(GetTerminalSize())
			case <-done:
				return
			}
		}
	}()

	return func() {
		signal.Stop(ch)
		close(done)
	}
}

// ElapsedTimer tracks elapsed time for display
type ElapsedTimer struct {
	start time.Time
	mu    sync.Mutex
}

// NewElapsedTimer creates a new ElapsedTimer starting now
func NewElapsedTimer() *ElapsedTimer {
	return &ElapsedTimer{start: time.Now()}
}

// Reset resets the timer to now
func (t *ElapsedTimer) Reset() {
	t.mu.Lock()
	t.start = time.Now()
	t.mu.Unlock()
}

// Elapsed returns the duration since the timer was started or last reset
func (t *ElapsedTimer) Elapsed() time.Duration {
	t.mu.Lock()
	defer t.mu.Unlock()
	return time.Since(t.start)
}

// BlinkState manages a boolean that toggles on an interval for cursor animation
type BlinkState struct {
	visible  bool
	interval time.Duration
	mu       sync.Mutex
	done     chan struct{}
}

// NewBlinkState creates a BlinkState that toggles at the given interval
func NewBlinkState(interval time.Duration) *BlinkState {
	bs := &BlinkState{
		visible:  true,
		interval: interval,
		done:     make(chan struct{}),
	}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				bs.mu.Lock()
				bs.visible = !bs.visible
				bs.mu.Unlock()
			case <-bs.done:
				return
			}
		}
	}()
	return bs
}

// IsVisible reports whether the blink state is currently visible
func (bs *BlinkState) IsVisible() bool {
	bs.mu.Lock()
	defer bs.mu.Unlock()
	return bs.visible
}

// Stop stops the background ticker goroutine
func (bs *BlinkState) Stop() { close(bs.done) }

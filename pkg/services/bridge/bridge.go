package bridge

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// SessionState tracks the remote bridge session
type SessionState struct {
	SessionID   string     `json:"session_id"`
	Status      string     `json:"status"` // "disconnected", "connecting", "connected"
	ConnectedAt *time.Time `json:"connected_at,omitempty"`
	LastPing    *time.Time `json:"last_ping,omitempty"`
}

// Message is a bridge protocol message
type Message struct {
	Type    string          `json:"type"`
	ID      string          `json:"id,omitempty"`
	Payload json.RawMessage `json:"payload,omitempty"`
}

// Bridge manages a remote session connection
type Bridge struct {
	mu     sync.Mutex
	state  SessionState
	inbox  chan Message
	outbox chan Message
	done   chan struct{}
}

// NewBridge creates a new bridge instance
func NewBridge(sessionID string) *Bridge {
	return &Bridge{
		state: SessionState{
			SessionID: sessionID,
			Status:    "disconnected",
		},
		inbox:  make(chan Message, 100),
		outbox: make(chan Message, 100),
		done:   make(chan struct{}),
	}
}

// GetState returns the current bridge state
func (b *Bridge) GetState() SessionState {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state
}

// SetStatus updates the bridge connection status
func (b *Bridge) SetStatus(status string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.state.Status = status
	if status == "connected" {
		now := time.Now()
		b.state.ConnectedAt = &now
	}
}

// Send queues a message for sending
func (b *Bridge) Send(msg Message) error {
	select {
	case b.outbox <- msg:
		return nil
	default:
		return fmt.Errorf("outbox full")
	}
}

// Receive returns the next incoming message (non-blocking)
func (b *Bridge) Receive() (Message, bool) {
	select {
	case msg := <-b.inbox:
		return msg, true
	default:
		return Message{}, false
	}
}

// Close shuts down the bridge
func (b *Bridge) Close() {
	b.mu.Lock()
	b.state.Status = "disconnected"
	b.mu.Unlock()
	close(b.done)
}

// IsConnected returns whether the bridge is active
func (b *Bridge) IsConnected() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.state.Status == "connected"
}

package bridge

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sync"
	"time"
)

// SessionManager manages bridge sessions
type SessionManager struct {
	mu       sync.Mutex
	sessions map[string]*BridgeSession
}

// BridgeSession holds all state for a single bridge session
type BridgeSession struct {
	ID        string
	Bridge    *Bridge
	Transport Transport
	Config    BridgeConfig
	CreatedAt time.Time
	LastPing  time.Time
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return &SessionManager{sessions: make(map[string]*BridgeSession)}
}

// CreateSession initialises a new bridge session with the given config
func (m *SessionManager) CreateSession(ctx context.Context, cfg BridgeConfig) (*BridgeSession, error) {
	if cfg.SessionID == "" {
		b := make([]byte, 16)
		if _, err := rand.Read(b); err != nil {
			return nil, fmt.Errorf("generate session id: %w", err)
		}
		cfg.SessionID = hex.EncodeToString(b)
	}

	// Choose transport
	var transport Transport
	switch cfg.Transport {
	case "polling":
		transport = NewPollingTransport()
	default:
		transport = NewWebSocketTransport()
	}

	if err := transport.Connect(ctx, cfg); err != nil {
		return nil, fmt.Errorf("connect transport: %w", err)
	}

	session := &BridgeSession{
		ID:        cfg.SessionID,
		Bridge:    NewBridge(cfg.SessionID),
		Transport: transport,
		Config:    cfg,
		CreatedAt: time.Now(),
	}

	m.mu.Lock()
	m.sessions[cfg.SessionID] = session
	m.mu.Unlock()

	// Start ping loop
	go m.pingLoop(session)

	return session, nil
}

// GetSession retrieves a session by ID, returning nil if not found
func (m *SessionManager) GetSession(id string) *BridgeSession {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.sessions[id]
}

// CloseSession tears down a session and removes it from the manager
func (m *SessionManager) CloseSession(id string) {
	m.mu.Lock()
	session, ok := m.sessions[id]
	if ok {
		delete(m.sessions, id)
	}
	m.mu.Unlock()

	if session != nil {
		session.Transport.Close() //nolint:errcheck
		session.Bridge.Close()
	}
}

// CloseAll tears down every active session
func (m *SessionManager) CloseAll() {
	m.mu.Lock()
	ids := make([]string, 0, len(m.sessions))
	for id := range m.sessions {
		ids = append(ids, id)
	}
	m.mu.Unlock()

	for _, id := range ids {
		m.CloseSession(id)
	}
}

func (m *SessionManager) pingLoop(session *BridgeSession) {
	interval := time.Duration(session.Config.PollInterval) * time.Millisecond
	if interval <= 0 {
		interval = time.Duration(DefaultPollInterval) * time.Millisecond
	}
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			session.Transport.Send(ProtocolMessage{ //nolint:errcheck
				Type:      MsgTypePing,
				SessionID: session.ID,
			})
			session.LastPing = time.Now()
		}
	}
}

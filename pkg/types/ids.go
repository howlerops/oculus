package types

import (
	"crypto/rand"
	"fmt"
	"regexp"
)

// SessionId uniquely identifies a Claude Code session
type SessionId string

// AgentId uniquely identifies a subagent within a session
type AgentId string

// AsSessionId casts a raw string to SessionId
func AsSessionId(id string) SessionId {
	return SessionId(id)
}

// AsAgentId casts a raw string to AgentId
func AsAgentId(id string) AgentId {
	return AgentId(id)
}

var agentIDPattern = regexp.MustCompile(`^a(?:.+-)?[0-9a-f]{16}$`)

// ToAgentId validates and brands a string as AgentId
// Returns empty string if invalid
func ToAgentId(s string) (AgentId, bool) {
	if agentIDPattern.MatchString(s) {
		return AgentId(s), true
	}
	return "", false
}

// NewAgentId creates a new random agent ID with optional label
func NewAgentId(label string) AgentId {
	b := make([]byte, 8)
	rand.Read(b)
	hex := fmt.Sprintf("%x", b)
	if label != "" {
		return AgentId("a" + label + "-" + hex)
	}
	return AgentId("a" + hex)
}

package state

import (
	"sync"
)

// Store provides thread-safe access to AppState
type Store struct {
	mu          sync.RWMutex
	state       AppState
	subscribers []func(AppState)
}

// NewStore creates a new state store with the given initial state
func NewStore(initial AppState) *Store {
	return &Store{
		state: initial,
	}
}

// Get returns a snapshot of the current state
func (s *Store) Get() AppState {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.state
}

// Set replaces the entire state
func (s *Store) Set(state AppState) {
	s.mu.Lock()
	s.state = state
	subs := s.subscribers
	s.mu.Unlock()

	for _, sub := range subs {
		if sub != nil {
			sub(state)
		}
	}
}

// Update applies a mutation function to the state
func (s *Store) Update(fn func(prev AppState) AppState) {
	s.mu.Lock()
	s.state = fn(s.state)
	newState := s.state
	subs := s.subscribers
	s.mu.Unlock()

	for _, sub := range subs {
		if sub != nil {
			sub(newState)
		}
	}
}

// Subscribe registers a callback for state changes.
// Returns an unsubscribe function.
func (s *Store) Subscribe(fn func(AppState)) func() {
	s.mu.Lock()
	s.subscribers = append(s.subscribers, fn)
	idx := len(s.subscribers) - 1
	s.mu.Unlock()

	return func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		// Set to nil instead of removing to preserve indices
		if idx < len(s.subscribers) {
			s.subscribers[idx] = nil
		}
	}
}

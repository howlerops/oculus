package analytics

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Event struct {
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	SessionID  string                 `json:"session_id,omitempty"`
}

type Sink interface {
	Send(event Event) error
	Flush() error
}

// FileSink writes events to a JSONL file
type FileSink struct {
	mu   sync.Mutex
	path string
	file *os.File
}

func NewFileSink(dir string) (*FileSink, error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}
	path := filepath.Join(dir, "events.jsonl")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, err
	}
	return &FileSink{path: path, file: f}, nil
}

func (s *FileSink) Send(event Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	data, err := json.Marshal(event)
	if err != nil {
		return err
	}
	_, err = s.file.Write(append(data, '\n'))
	return err
}

func (s *FileSink) Flush() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.file.Sync()
}

func (s *FileSink) Close() error {
	s.Flush()
	return s.file.Close()
}

// NullSink discards all events (for when analytics is disabled)
type NullSink struct{}

func (s *NullSink) Send(_ Event) error { return nil }
func (s *NullSink) Flush() error       { return nil }

// Logger is the main analytics interface
type Logger struct {
	sink      Sink
	sessionID string
	disabled  bool
}

func NewLogger(sink Sink, sessionID string) *Logger {
	return &Logger{sink: sink, sessionID: sessionID}
}

func NewDisabledLogger() *Logger {
	return &Logger{sink: &NullSink{}, disabled: true}
}

func (l *Logger) LogEvent(name string, properties map[string]interface{}) {
	if l.disabled {
		return
	}
	l.sink.Send(Event{
		Name:       name,
		Properties: properties,
		Timestamp:  time.Now(),
		SessionID:  l.sessionID,
	})
}

func (l *Logger) Flush() {
	if l.sink != nil {
		l.sink.Flush()
	}
}

// IsAnalyticsDisabled checks if the user opted out
func IsAnalyticsDisabled() bool {
	return os.Getenv("CLAUDE_CODE_DISABLE_ANALYTICS") == "1"
}

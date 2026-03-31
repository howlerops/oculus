package filestate

import (
	"os"
	"sync"
	"time"
)

type FileState struct {
	Path    string
	ModTime time.Time
	Size    int64
	Exists  bool
}

type Cache struct {
	mu      sync.RWMutex
	entries map[string]*FileState
	maxSize int
}

func NewCache(maxSize int) *Cache {
	return &Cache{entries: make(map[string]*FileState), maxSize: maxSize}
}

func (c *Cache) Get(path string) (*FileState, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	s, ok := c.entries[path]
	return s, ok
}

func (c *Cache) Refresh(path string) *FileState {
	info, err := os.Stat(path)
	state := &FileState{Path: path}
	if err == nil {
		state.Exists = true
		state.ModTime = info.ModTime()
		state.Size = info.Size()
	}
	c.mu.Lock()
	if len(c.entries) >= c.maxSize {
		for k := range c.entries {
			delete(c.entries, k)
			break
		}
	}
	c.entries[path] = state
	c.mu.Unlock()
	return state
}

func (c *Cache) Clear() {
	c.mu.Lock()
	c.entries = make(map[string]*FileState)
	c.mu.Unlock()
}

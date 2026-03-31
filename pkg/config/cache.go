package config

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// SettingsCache caches loaded settings with change detection
type SettingsCache struct {
	mu            sync.RWMutex
	settings      *SettingsJson
	lastHash      string
	lastCheck     time.Time
	checkInterval time.Duration
}

var globalSettingsCache = &SettingsCache{
	checkInterval: 5 * time.Second,
}

// GetCachedSettings returns settings, reloading if stale
func GetCachedSettings() *SettingsJson {
	globalSettingsCache.mu.RLock()
	if globalSettingsCache.settings != nil && time.Since(globalSettingsCache.lastCheck) < globalSettingsCache.checkInterval {
		s := globalSettingsCache.settings
		globalSettingsCache.mu.RUnlock()
		return s
	}
	globalSettingsCache.mu.RUnlock()

	// Reload
	globalSettingsCache.mu.Lock()
	defer globalSettingsCache.mu.Unlock()

	// Double-check after acquiring write lock
	if globalSettingsCache.settings != nil && time.Since(globalSettingsCache.lastCheck) < globalSettingsCache.checkInterval {
		return globalSettingsCache.settings
	}

	settings, _ := LoadSettings()
	if settings == nil {
		settings = &SettingsJson{}
	}
	globalSettingsCache.settings = settings
	globalSettingsCache.lastCheck = time.Now()

	// Check hash for change detection
	data, _ := json.Marshal(settings)
	hash := fmt.Sprintf("%x", sha256.Sum256(data))
	if hash != globalSettingsCache.lastHash {
		globalSettingsCache.lastHash = hash
		// Settings changed - could notify watchers here
	}

	return settings
}

// InvalidateSettingsCache forces reload on next access
func InvalidateSettingsCache() {
	globalSettingsCache.mu.Lock()
	globalSettingsCache.settings = nil
	globalSettingsCache.lastCheck = time.Time{}
	globalSettingsCache.mu.Unlock()
}

// WatchSettings watches settings.json for changes and calls onChange when they occur.
// It returns a stop function that cancels the watcher.
func WatchSettings(onChange func(*SettingsJson)) func() {
	done := make(chan struct{})
	go func() {
		var lastHash string
		ticker := time.NewTicker(2 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				data, err := os.ReadFile(GetSettingsPath())
				if err != nil {
					continue
				}
				hash := fmt.Sprintf("%x", sha256.Sum256(data))
				if hash != lastHash && lastHash != "" {
					lastHash = hash
					var s SettingsJson
					if json.Unmarshal(data, &s) == nil {
						InvalidateSettingsCache()
						onChange(&s)
					}
				} else {
					lastHash = hash
				}
			case <-done:
				return
			}
		}
	}()
	return func() { close(done) }
}

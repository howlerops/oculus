package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/howlerops/oculus/pkg/config"
)

// Entry represents a single conversation history entry
type Entry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Model     string    `json:"model"`
	StartedAt time.Time `json:"startedAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Turns     int       `json:"turns"`
	CWD       string    `json:"cwd"`
}

// GetHistoryDir returns the history directory path
func GetHistoryDir() string {
	return filepath.Join(config.GetOculusDir(), "history")
}

// AddEntry adds a conversation to history
func AddEntry(entry Entry) error {
	dir := GetHistoryDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}

	path := filepath.Join(dir, entry.ID+".json")
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0o644)
}

// ListEntries returns recent history entries
func ListEntries(limit int) ([]Entry, error) {
	dir := GetHistoryDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var result []Entry
	for i := len(entries) - 1; i >= 0 && len(result) < limit; i-- {
		e := entries[i]
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var entry Entry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}
		result = append(result, entry)
	}
	return result, nil
}

// FormatEntry returns a display string for an entry
func FormatEntry(e Entry) string {
	return fmt.Sprintf("[%s] %s (%d turns, %s)", e.ID[:8], e.Title, e.Turns, e.UpdatedAt.Format("2006-01-02 15:04"))
}

package memory

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/jbeck018/claude-go/pkg/config"
	"github.com/jbeck018/claude-go/pkg/types"
)

type MemoryEntry struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	Source    string    `json:"source"` // "user", "auto", "tool"
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags,omitempty"`
}

func GetMemoryDir() string {
	return filepath.Join(config.GetClaudeConfigDir(), "memory")
}

func SaveMemory(entry MemoryEntry) error {
	dir := GetMemoryDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, entry.ID+".json"), data, 0o644)
}

func LoadMemories() ([]MemoryEntry, error) {
	dir := GetMemoryDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil
	}
	var result []MemoryEntry
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var mem MemoryEntry
		if json.Unmarshal(data, &mem) == nil {
			result = append(result, mem)
		}
	}
	return result, nil
}

// ExtractMemoriesFromConversation scans for <remember> tags in assistant messages.
func ExtractMemoriesFromConversation(messages []types.Message) []MemoryEntry {
	var memories []MemoryEntry
	for _, msg := range messages {
		if msg.Kind != "assistant" || msg.Assistant == nil {
			continue
		}
		for _, block := range msg.Assistant.Content {
			if block.Type != types.ContentBlockText {
				continue
			}
			text := block.Text
			for {
				start := strings.Index(text, "<remember>")
				if start == -1 {
					break
				}
				end := strings.Index(text[start:], "</remember>")
				if end == -1 {
					break
				}
				content := text[start+10 : start+end]
				memories = append(memories, MemoryEntry{
					ID:        generateMemoryID(),
					Content:   strings.TrimSpace(content),
					Source:    "auto",
					CreatedAt: time.Now(),
				})
				text = text[start+end+11:]
			}
		}
	}
	return memories
}

func generateMemoryID() string {
	return time.Now().Format("20060102-150405")
}

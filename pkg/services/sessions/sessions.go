package sessions

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/jbeck018/claude-go/pkg/config"
	"github.com/jbeck018/claude-go/pkg/types"
)

type SessionMetadata struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Model     string    `json:"model"`
	CWD       string    `json:"cwd"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Turns     int       `json:"turns"`
}

type SavedSession struct {
	Metadata SessionMetadata `json:"metadata"`
	Messages []types.Message `json:"messages"`
}

func GetSessionsDir() string {
	return filepath.Join(config.GetClaudeConfigDir(), "conversations")
}

func GenerateSessionID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func Save(session SavedSession) error {
	dir := GetSessionsDir()
	os.MkdirAll(dir, 0o755)
	path := filepath.Join(dir, session.Metadata.ID+".json")
	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func Load(id string) (*SavedSession, error) {
	path := filepath.Join(GetSessionsDir(), id+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var session SavedSession
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, err
	}
	return &session, nil
}

func ListRecent(limit int) ([]SessionMetadata, error) {
	dir := GetSessionsDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var sessions []SessionMetadata
	for _, e := range entries {
		if filepath.Ext(e.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		var saved SavedSession
		if err := json.Unmarshal(data, &saved); err != nil {
			continue
		}
		sessions = append(sessions, saved.Metadata)
	}

	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].UpdatedAt.After(sessions[j].UpdatedAt)
	})

	if limit > 0 && len(sessions) > limit {
		sessions = sessions[:limit]
	}
	return sessions, nil
}

func Delete(id string) error {
	return os.Remove(filepath.Join(GetSessionsDir(), id+".json"))
}

func FormatSessionList(sessions []SessionMetadata) string {
	if len(sessions) == 0 {
		return "No previous sessions found."
	}
	result := fmt.Sprintf("Recent sessions (%d):\n", len(sessions))
	for _, s := range sessions {
		age := time.Since(s.UpdatedAt).Round(time.Minute)
		result += fmt.Sprintf("  [%s] %s (%d turns, %s ago, %s)\n", s.ID[:8], s.Title, s.Turns, age, s.CWD)
	}
	return result
}

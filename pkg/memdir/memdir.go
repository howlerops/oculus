package memdir

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/jbeck018/claude-go/pkg/config"
)

const MemoryDirName = "memory"

// GetMemoryDir returns the auto-memory directory path
func GetMemoryDir() string {
	return filepath.Join(config.GetClaudeConfigDir(), "projects", getCurrentProjectHash(), MemoryDirName)
}

// IsAutoMemoryEnabled checks if auto-memory feature is on
func IsAutoMemoryEnabled() bool {
	return os.Getenv("CLAUDE_CODE_DISABLE_AUTO_MEMORY") != "1"
}

// MemoryFile represents a single memory file
type MemoryFile struct {
	Name    string
	Path    string
	Content string
}

// ListMemoryFiles returns all memory files
func ListMemoryFiles() ([]MemoryFile, error) {
	dir := GetMemoryDir()
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	var files []MemoryFile
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		content, err := os.ReadFile(filepath.Join(dir, e.Name()))
		if err != nil {
			continue
		}
		files = append(files, MemoryFile{
			Name:    strings.TrimSuffix(e.Name(), ".md"),
			Path:    filepath.Join(dir, e.Name()),
			Content: string(content),
		})
	}
	return files, nil
}

// SaveMemoryFile writes a memory file
func SaveMemoryFile(name, content string) error {
	dir := GetMemoryDir()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, name+".md"), []byte(content), 0o644)
}

// DeleteMemoryFile removes a memory file
func DeleteMemoryFile(name string) error {
	return os.Remove(filepath.Join(GetMemoryDir(), name+".md"))
}

func getCurrentProjectHash() string {
	cwd, _ := os.Getwd()
	// Simple hash: use last 2 path components
	parts := strings.Split(filepath.Clean(cwd), string(filepath.Separator))
	if len(parts) >= 2 {
		return parts[len(parts)-2] + "-" + parts[len(parts)-1]
	}
	return filepath.Base(cwd)
}

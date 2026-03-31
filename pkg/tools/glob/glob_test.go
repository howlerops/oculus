package glob

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGlobMatchesFiles(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte("hello"), 0644)
	os.MkdirAll(filepath.Join(dir, "sub"), 0755)
	os.WriteFile(filepath.Join(dir, "sub", "nested.go"), []byte("package sub"), 0644)

	tool := NewGlobTool()
	result, err := tool.Call(context.Background(), map[string]interface{}{
		"pattern": "*.go",
		"path":    dir,
	}, nil)

	if err != nil {
		t.Fatal(err)
	}

	data := result.Data.(string)
	if !strings.Contains(data, "test.go") {
		t.Errorf("expected test.go in results, got: %s", data)
	}
}

func TestGlobNoMatches(t *testing.T) {
	dir := t.TempDir()

	tool := NewGlobTool()
	result, _ := tool.Call(context.Background(), map[string]interface{}{
		"pattern": "*.xyz",
		"path":    dir,
	}, nil)

	data := result.Data.(string)
	if !strings.Contains(data, "No files matched") {
		t.Errorf("expected no match message, got: %s", data)
	}
}

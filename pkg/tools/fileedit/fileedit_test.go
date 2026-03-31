package fileedit

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

func TestFileEditReplace(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello world\nfoo bar\n"), 0644)

	tool := NewFileEditTool()
	result, err := tool.Call(context.Background(), map[string]interface{}{
		"file_path":  path,
		"old_string": "hello",
		"new_string": "goodbye",
	}, nil)

	if err != nil {
		t.Fatal(err)
	}
	if result == nil {
		t.Fatal("result is nil")
	}

	content, _ := os.ReadFile(path)
	if string(content) != "goodbye world\nfoo bar\n" {
		t.Errorf("got %q, want 'goodbye world\\nfoo bar\\n'", string(content))
	}
}

func TestFileEditNonUnique(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("aaa\naaa\n"), 0644)

	tool := NewFileEditTool()
	result, _ := tool.Call(context.Background(), map[string]interface{}{
		"file_path":  path,
		"old_string": "aaa",
		"new_string": "bbb",
	}, nil)

	data := result.Data.(string)
	if data == "" || !contains(data, "found 2 times") {
		t.Errorf("expected non-unique error, got: %s", data)
	}
}

func TestFileEditReplaceAll(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("aaa\naaa\n"), 0644)

	tool := NewFileEditTool()
	tool.Call(context.Background(), map[string]interface{}{
		"file_path":   path,
		"old_string":  "aaa",
		"new_string":  "bbb",
		"replace_all": true,
	}, nil)

	content, _ := os.ReadFile(path)
	if string(content) != "bbb\nbbb\n" {
		t.Errorf("got %q, want 'bbb\\nbbb\\n'", string(content))
	}
}

func TestFileEditNotFound(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.txt")
	os.WriteFile(path, []byte("hello"), 0644)

	tool := NewFileEditTool()
	result, _ := tool.Call(context.Background(), map[string]interface{}{
		"file_path":  path,
		"old_string": "xyz",
		"new_string": "abc",
	}, nil)

	data := result.Data.(string)
	if !contains(data, "not found") {
		t.Errorf("expected not found error, got: %s", data)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsHelper(s, sub))
}

func containsHelper(s, sub string) bool {
	for i := 0; i <= len(s)-len(sub); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}

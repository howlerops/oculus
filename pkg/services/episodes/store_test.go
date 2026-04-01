package episodes

import (
	"os"
	"testing"
)

func TestCreateAndCompactEpisode(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	ep := store.CreateEpisode()
	if ep.Status != StatusActive {
		t.Errorf("expected active, got %s", ep.Status)
	}

	store.AddMessage(ep.ID, "user", "hello world", "", 10)
	store.AddMessage(ep.ID, "assistant", "hi there", "Bash", 15)

	if store.GetTotalTokens() != 25 {
		t.Errorf("expected 25 tokens, got %d", store.GetTotalTokens())
	}

	store.CompactEpisode(ep.ID, "Greeting exchange", []string{"hello", "greeting"})

	compacted := store.episodes[ep.ID]
	if compacted.Status != StatusCompacted {
		t.Errorf("expected compacted, got %s", compacted.Status)
	}
	if compacted.Summary != "Greeting exchange" {
		t.Errorf("wrong summary")
	}
	if len(compacted.Messages) != 0 {
		t.Errorf("messages should be cleared after compaction")
	}

	// Verify persisted to disk
	files, _ := os.ReadDir(dir)
	if len(files) != 1 {
		t.Errorf("expected 1 file on disk, got %d", len(files))
	}
}

func TestSearchEpisodes(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)

	ep1 := store.CreateEpisode()
	store.AddMessage(ep1.ID, "user", "fix the auth bug", "", 10)
	store.CompactEpisode(ep1.ID, "Fixed authentication bug in login flow", []string{"auth", "bug", "login"})

	ep2 := store.CreateEpisode()
	store.AddMessage(ep2.ID, "user", "add database migration", "", 10)
	store.CompactEpisode(ep2.ID, "Created database migration for users table", []string{"database", "migration", "users"})

	results := store.SearchEpisodes("auth login", 5)
	if len(results) != 1 {
		t.Errorf("expected 1 result, got %d", len(results))
	}
	if results[0].ID != ep1.ID {
		t.Errorf("expected ep1, got %s", results[0].ID)
	}
}

func TestContextLoop(t *testing.T) {
	dir := t.TempDir()
	store := NewStore(dir)
	config := LCMConfig{Enabled: true, SoftThreshold: 100, HardThreshold: 200, MaxEpisodeTokens: 50}
	loop := NewContextLoop(config, store)

	if loop.CheckThreshold(50) != "ok" {
		t.Error("expected ok at 50 tokens")
	}
	if loop.CheckThreshold(150) != "soft" {
		t.Error("expected soft at 150 tokens")
	}
	if loop.CheckThreshold(250) != "hard" {
		t.Error("expected hard at 250 tokens")
	}
}

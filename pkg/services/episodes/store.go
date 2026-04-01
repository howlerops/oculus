package episodes

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

// EpisodeStatus describes lifecycle state
type EpisodeStatus string

const (
	StatusActive    EpisodeStatus = "active"
	StatusCompacted EpisodeStatus = "compacted"
	StatusArchived  EpisodeStatus = "archived"
)

// Episode represents a bounded unit of conversation
type Episode struct {
	ID         string            `json:"id"`
	Status     EpisodeStatus     `json:"status"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	TokenCount int               `json:"token_count"`
	Summary    string            `json:"summary,omitempty"`
	Keywords   []string          `json:"keywords,omitempty"`
	Messages   []EpisodeMessage  `json:"messages,omitempty"` // only for active episodes
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// EpisodeMessage is a simplified message within an episode
type EpisodeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	ToolUse string `json:"tool_use,omitempty"`
}

// Store manages episodes on disk
type Store struct {
	mu        sync.RWMutex
	episodes  map[string]*Episode
	dir       string
	maxActive int
}

// NewStore creates an episode store at the given directory
func NewStore(dir string) *Store {
	os.MkdirAll(dir, 0o755)
	s := &Store{
		episodes:  make(map[string]*Episode),
		dir:       dir,
		maxActive: 10, // keep last 10 active episodes
	}
	s.loadFromDisk()
	return s
}

// CreateEpisode starts a new episode
func (s *Store) CreateEpisode() *Episode {
	s.mu.Lock()
	defer s.mu.Unlock()

	ep := &Episode{
		ID:        fmt.Sprintf("ep_%d", time.Now().UnixNano()),
		Status:    StatusActive,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	s.episodes[ep.ID] = ep
	return ep
}

// AddMessage adds a message to the active episode
func (s *Store) AddMessage(episodeID, role, content, toolUse string, tokens int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ep, ok := s.episodes[episodeID]
	if !ok || ep.Status != StatusActive {
		return
	}

	ep.Messages = append(ep.Messages, EpisodeMessage{
		Role: role, Content: content, ToolUse: toolUse,
	})
	ep.TokenCount += tokens
	ep.UpdatedAt = time.Now()
}

// CompactEpisode transitions an episode from active to compacted with a summary
func (s *Store) CompactEpisode(episodeID, summary string, keywords []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ep, ok := s.episodes[episodeID]
	if !ok {
		return
	}

	ep.Status = StatusCompacted
	ep.Summary = summary
	ep.Keywords = keywords
	ep.Messages = nil // free message memory
	ep.UpdatedAt = time.Now()

	s.saveToDisk(ep)
}

// ArchiveEpisode transitions to archived (deeply compressed)
func (s *Store) ArchiveEpisode(episodeID string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	ep, ok := s.episodes[episodeID]
	if !ok {
		return
	}

	ep.Status = StatusArchived
	ep.UpdatedAt = time.Now()
	s.saveToDisk(ep)
}

// GetActiveEpisodes returns all active episodes ordered by creation time
func (s *Store) GetActiveEpisodes() []*Episode {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var active []*Episode
	for _, ep := range s.episodes {
		if ep.Status == StatusActive {
			active = append(active, ep)
		}
	}
	sort.Slice(active, func(i, j int) bool {
		return active[i].CreatedAt.Before(active[j].CreatedAt)
	})
	return active
}

// SearchEpisodes finds episodes matching keywords using TF-IDF scoring
func (s *Store) SearchEpisodes(query string, maxResults int) []*Episode {
	s.mu.RLock()
	defer s.mu.RUnlock()

	queryTerms := tokenize(strings.ToLower(query))
	type scored struct {
		ep    *Episode
		score float64
	}

	var results []scored
	for _, ep := range s.episodes {
		if ep.Status == StatusActive {
			continue
		} // only search compacted/archived

		score := 0.0
		epText := strings.ToLower(ep.Summary + " " + strings.Join(ep.Keywords, " "))
		epTerms := tokenize(epText)
		termSet := make(map[string]bool)
		for _, t := range epTerms {
			termSet[t] = true
		}

		for _, qt := range queryTerms {
			if termSet[qt] {
				score += 1.0
			}
		}

		if score > 0 {
			results = append(results, scored{ep: ep, score: score})
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i].score > results[j].score })

	var out []*Episode
	for i, r := range results {
		if i >= maxResults {
			break
		}
		out = append(out, r.ep)
	}
	return out
}

// GetTotalTokens returns total tokens across all active episodes
func (s *Store) GetTotalTokens() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	total := 0
	for _, ep := range s.episodes {
		if ep.Status == StatusActive {
			total += ep.TokenCount
		}
	}
	return total
}

// Prune removes oldest archived episodes beyond a limit
func (s *Store) Prune(maxArchived int) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	var archived []*Episode
	for _, ep := range s.episodes {
		if ep.Status == StatusArchived {
			archived = append(archived, ep)
		}
	}

	sort.Slice(archived, func(i, j int) bool {
		return archived[i].CreatedAt.Before(archived[j].CreatedAt)
	})

	pruned := 0
	for i := 0; i < len(archived)-maxArchived; i++ {
		delete(s.episodes, archived[i].ID)
		os.Remove(filepath.Join(s.dir, archived[i].ID+".json"))
		pruned++
	}
	return pruned
}

func (s *Store) saveToDisk(ep *Episode) {
	data, _ := json.MarshalIndent(ep, "", "  ")
	os.WriteFile(filepath.Join(s.dir, ep.ID+".json"), data, 0o644)
}

func (s *Store) loadFromDisk() {
	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return
	}
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		data, err := os.ReadFile(filepath.Join(s.dir, e.Name()))
		if err != nil {
			continue
		}
		var ep Episode
		if json.Unmarshal(data, &ep) == nil {
			s.episodes[ep.ID] = &ep
		}
	}
}

func tokenize(text string) []string {
	var words []string
	word := ""
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			word += string(r)
		} else {
			if len(word) > 2 {
				words = append(words, word)
			}
			word = ""
		}
	}
	if len(word) > 2 {
		words = append(words, word)
	}
	return words
}

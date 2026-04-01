package episodes

import (
	"fmt"
	"strings"
	"time"
)

// LCMConfig configures the context management engine
type LCMConfig struct {
	Enabled          bool `json:"enabled"`
	SoftThreshold    int  `json:"soft_threshold"`    // τ_soft: async compaction begins
	HardThreshold    int  `json:"hard_threshold"`    // τ_hard: blocking compaction
	MaxEpisodeTokens int  `json:"max_episode_tokens"`
}

// DefaultLCMConfig returns sensible defaults for 200k context
func DefaultLCMConfig() LCMConfig {
	return LCMConfig{
		Enabled:          true,
		SoftThreshold:    140000, // 70% of 200k
		HardThreshold:    180000, // 90% of 200k
		MaxEpisodeTokens: 4000,
	}
}

// ContextState tracks the current context window state
type ContextState struct {
	TotalTokens    int
	ActiveEpisodes int
	Level          string // "ok", "soft", "hard"
	LastCompacted  time.Time
}

// ContextLoop manages the dual-threshold context control
type ContextLoop struct {
	config LCMConfig
	store  *Store
	state  ContextState
}

// NewContextLoop creates an LCM context loop
func NewContextLoop(config LCMConfig, store *Store) *ContextLoop {
	return &ContextLoop{
		config: config,
		store:  store,
	}
}

// CheckThreshold evaluates context usage and returns the action needed
func (cl *ContextLoop) CheckThreshold(totalTokens int) string {
	cl.state.TotalTokens = totalTokens
	cl.state.ActiveEpisodes = len(cl.store.GetActiveEpisodes())

	if totalTokens >= cl.config.HardThreshold {
		cl.state.Level = "hard"
		return "hard" // blocking compaction needed NOW
	}
	if totalTokens >= cl.config.SoftThreshold {
		cl.state.Level = "soft"
		return "soft" // async compaction between turns
	}
	cl.state.Level = "ok"
	return "ok"
}

// CompactOldest compacts the oldest active episode
func (cl *ContextLoop) CompactOldest() (*Episode, error) {
	active := cl.store.GetActiveEpisodes()
	if len(active) == 0 {
		return nil, fmt.Errorf("no active episodes to compact")
	}

	oldest := active[0]

	// Generate summary from messages
	summary := generateEpisodeSummary(oldest)
	keywords := extractKeywords(oldest)

	cl.store.CompactEpisode(oldest.ID, summary, keywords)
	cl.state.LastCompacted = time.Now()

	return oldest, nil
}

// CompactUntilBelow keeps compacting until tokens drop below threshold
func (cl *ContextLoop) CompactUntilBelow(targetTokens int) int {
	compacted := 0
	for cl.store.GetTotalTokens() > targetTokens {
		_, err := cl.CompactOldest()
		if err != nil {
			break
		}
		compacted++
	}
	return compacted
}

// RetrieveRelevant finds episodes relevant to a query
func (cl *ContextLoop) RetrieveRelevant(query string, maxResults int) []*Episode {
	return cl.store.SearchEpisodes(query, maxResults)
}

// GetState returns the current context state
func (cl *ContextLoop) GetState() ContextState {
	return cl.state
}

func generateEpisodeSummary(ep *Episode) string {
	var parts []string
	for _, msg := range ep.Messages {
		switch msg.Role {
		case "user":
			text := msg.Content
			if len(text) > 80 {
				text = text[:77] + "..."
			}
			parts = append(parts, "User: "+text)
		case "assistant":
			if msg.ToolUse != "" {
				parts = append(parts, "Tool: "+msg.ToolUse)
			} else {
				text := msg.Content
				if len(text) > 80 {
					text = text[:77] + "..."
				}
				parts = append(parts, "Assistant: "+text)
			}
		}
	}
	return strings.Join(parts, "\n")
}

func extractKeywords(ep *Episode) []string {
	// Simple keyword extraction: most frequent non-trivial words
	freq := make(map[string]int)
	for _, msg := range ep.Messages {
		words := tokenize(strings.ToLower(msg.Content + " " + msg.ToolUse))
		for _, w := range words {
			freq[w]++
		}
	}

	type kv struct {
		word  string
		count int
	}
	var sorted []kv
	for w, c := range freq {
		sorted = append(sorted, kv{w, c})
	}
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[j].count > sorted[i].count {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}

	var keywords []string
	for i, kv := range sorted {
		if i >= 10 {
			break
		}
		keywords = append(keywords, kv.word)
	}
	return keywords
}

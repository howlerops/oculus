package compact

import (
	"fmt"

	"github.com/howlerops/oculus/pkg/types"
)

// LCMConfig mirrors the RLM module's LCM dual-threshold configuration.
// Thresholds are expressed in tokens.
type LCMConfig struct {
	Enabled          bool
	SoftThreshold    int // async compaction between turns
	HardThreshold    int // blocking compaction before next request
	MaxEpisodeTokens int
}

// DefaultLCMConfig returns thresholds calibrated for a 200k context window.
func DefaultLCMConfig() LCMConfig {
	return LCMConfig{
		Enabled:          true,
		SoftThreshold:    140000, // 70% of 200k
		HardThreshold:    180000, // 90% of 200k
		MaxEpisodeTokens: 4000,
	}
}

// RLMCompactor uses TF-IDF scoring (inspired by the Recursive LLM module)
// for intelligent context compaction.
type RLMCompactor struct {
	config LCMConfig
}

// NewRLMCompactor creates an RLM-backed compactor with default thresholds.
func NewRLMCompactor() (*RLMCompactor, error) {
	return &RLMCompactor{
		config: DefaultLCMConfig(),
	}, nil
}

// NewRLMCompactorWithConfig creates an RLM-backed compactor with the given config.
func NewRLMCompactorWithConfig(cfg LCMConfig) *RLMCompactor {
	return &RLMCompactor{config: cfg}
}

// CompactionLevel describes the urgency returned by ShouldCompact.
type CompactionLevel string

const (
	CompactionNone CompactionLevel = "none"
	CompactionSoft CompactionLevel = "soft" // async compaction between turns
	CompactionHard CompactionLevel = "hard" // blocking compaction needed
)

// ShouldCompact checks if compaction is needed based on dual thresholds.
func (c *RLMCompactor) ShouldCompact(totalTokens int) CompactionLevel {
	if totalTokens >= c.config.HardThreshold {
		return CompactionHard
	}
	if totalTokens >= c.config.SoftThreshold {
		return CompactionSoft
	}
	return CompactionNone
}

// CompactWithRLM performs intelligent compaction using TF-IDF episode summarization.
func (c *RLMCompactor) CompactWithRLM(messages []types.Message, opts CompactOptions) []types.Message {
	if len(messages) <= opts.PreserveRecent {
		return messages
	}

	// Split into episodes based on turn boundaries
	oldMessages := messages[:len(messages)-opts.PreserveRecent]
	recentMessages := messages[len(messages)-opts.PreserveRecent:]

	// Use TF-IDF to extract key sentences from old messages
	summary := extractKeyContent(oldMessages)

	// Create compacted boundary
	boundaryMsg := types.NewSystemMessage(
		types.SystemMsgCompactBoundary,
		fmt.Sprintf("--- Context compacted (RLM) ---\n\nEpisode summary (%d messages):\n%s", len(oldMessages), summary),
	)

	result := []types.Message{boundaryMsg}
	result = append(result, recentMessages...)
	return result
}

// extractKeyContent uses TF-IDF scoring to find the most important content.
func extractKeyContent(messages []types.Message) string {
	var texts []string
	for _, msg := range messages {
		switch msg.Kind {
		case "user":
			if msg.User != nil {
				for _, b := range msg.User.Content {
					if b.Type == types.ContentBlockText && b.Text != "" {
						texts = append(texts, b.Text)
					}
				}
			}
		case "assistant":
			if msg.Assistant != nil {
				for _, b := range msg.Assistant.Content {
					if b.Type == types.ContentBlockText && b.Text != "" {
						texts = append(texts, b.Text)
					}
					if b.Type == types.ContentBlockToolUse {
						texts = append(texts, fmt.Sprintf("Used tool: %s", b.Name))
					}
				}
			}
		}
	}

	if len(texts) == 0 {
		return "(no content to summarize)"
	}

	scored := tfidfScore(texts)

	maxSentences := 20
	if len(scored) < maxSentences {
		maxSentences = len(scored)
	}

	var result string
	for i := 0; i < maxSentences; i++ {
		result += "- " + scored[i].text + "\n"
	}
	return result
}

type scoredText struct {
	text  string
	score float64
}

// tfidfScore scores texts by term frequency-inverse document frequency.
func tfidfScore(texts []string) []scoredText {
	// Build document frequency map
	termFreq := make(map[string]int)
	for _, t := range texts {
		words := tokenize(t)
		seen := make(map[string]bool)
		for _, w := range words {
			if !seen[w] {
				termFreq[w]++
				seen[w] = true
			}
		}
	}

	var scored []scoredText
	for _, t := range texts {
		if len(t) < 10 {
			continue // skip trivial
		}
		words := tokenize(t)
		score := 0.0
		for _, w := range words {
			if df, ok := termFreq[w]; ok && df > 0 {
				// IDF = log(N/df), approximated as 1/df for simplicity
				idf := 1.0 / float64(df)
				score += idf
			}
		}
		// Normalize by length
		if len(words) > 0 {
			score /= float64(len(words))
		}
		// Boost tool use mentions
		if containsToolMention(t) {
			score *= 1.5
		}
		scored = append(scored, scoredText{text: truncateText(t, 150), score: score})
	}

	// Sort by score descending (selection sort for small N)
	for i := 0; i < len(scored); i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[i].score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	return scored
}

// tokenize splits text into lowercase words, discarding tokens shorter than 3 chars.
func tokenize(text string) []string {
	var words []string
	word := ""
	for _, r := range text {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			if r >= 'A' && r <= 'Z' {
				r = r + 32 // lowercase
			}
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

// containsToolMention returns true if text references a known tool name.
func containsToolMention(text string) bool {
	tools := []string{"Bash", "Read", "Edit", "Write", "Glob", "Grep", "Agent"}
	for _, t := range tools {
		for i := 0; i <= len(text)-len(t); i++ {
			if text[i:i+len(t)] == t {
				return true
			}
		}
	}
	return false
}

// truncateText shortens s to max runes, appending "..." if truncated.
func truncateText(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

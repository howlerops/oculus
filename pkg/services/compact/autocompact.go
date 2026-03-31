package compact

import (
	"github.com/jbeck018/claude-go/pkg/types"
)

// AutoCompactConfig controls automatic compaction behaviour
type AutoCompactConfig struct {
	Enabled           bool
	WarningThreshold  float64 // context % at which to show a warning (default 70)
	CompactThreshold  float64 // context % at which to auto-compact (default 80)
	CriticalThreshold float64 // context % considered urgent (default 90)
	MaxContextTokens  int     // maximum context window size
}

// DefaultAutoCompactConfig returns sensible defaults
func DefaultAutoCompactConfig() AutoCompactConfig {
	return AutoCompactConfig{
		Enabled:           true,
		WarningThreshold:  70.0,
		CompactThreshold:  80.0,
		CriticalThreshold: 90.0,
		MaxContextTokens:  200000,
	}
}

// TokenWarningLevel describes severity of context usage
type TokenWarningLevel string

const (
	TokenWarningLevelOK       TokenWarningLevel = "ok"
	TokenWarningLevelWarning  TokenWarningLevel = "warning"
	TokenWarningLevelCompact  TokenWarningLevel = "compact"
	TokenWarningLevelCritical TokenWarningLevel = "critical"
)

// TokenWarningState captures the current context window usage assessment
type TokenWarningState struct {
	TotalTokens  int
	MaxTokens    int
	UsagePercent float64
	Level        TokenWarningLevel
	ShouldCompact bool
	ShouldWarn   bool
}

// CalculateTokenWarningState determines the compaction urgency for the given token count
func CalculateTokenWarningState(totalTokens int, config AutoCompactConfig) TokenWarningState {
	maxTokens := config.MaxContextTokens
	if maxTokens == 0 {
		maxTokens = 200000
	}

	pct := float64(totalTokens) / float64(maxTokens) * 100

	state := TokenWarningState{
		TotalTokens:  totalTokens,
		MaxTokens:    maxTokens,
		UsagePercent: pct,
		Level:        TokenWarningLevelOK,
	}

	switch {
	case pct >= config.CriticalThreshold:
		state.Level = TokenWarningLevelCritical
		state.ShouldCompact = true
		state.ShouldWarn = true
	case pct >= config.CompactThreshold:
		state.Level = TokenWarningLevelCompact
		state.ShouldCompact = config.Enabled
		state.ShouldWarn = true
	case pct >= config.WarningThreshold:
		state.Level = TokenWarningLevelWarning
		state.ShouldWarn = true
	}

	return state
}

// ShouldAutoCompact returns true when the accumulated message tokens exceed
// the configured compact threshold
func ShouldAutoCompact(messages []types.Message, config AutoCompactConfig) bool {
	total := 0
	for _, msg := range messages {
		total += estimateMessageTokens(msg)
	}
	return CalculateTokenWarningState(total, config).ShouldCompact
}

// estimateMessageTokens approximates token count for a message (chars/4 heuristic)
func estimateMessageTokens(msg types.Message) int {
	chars := 0
	switch msg.Kind {
	case "user":
		if msg.User != nil {
			for _, b := range msg.User.Content {
				chars += len(b.Text) + len(b.Thinking)
			}
		}
	case "assistant":
		if msg.Assistant != nil {
			for _, b := range msg.Assistant.Content {
				chars += len(b.Text) + len(b.Thinking)
			}
		}
	case "system":
		if msg.System != nil {
			chars += len(msg.System.Text)
		}
	case "attachment":
		if msg.Attachment != nil {
			chars += len(msg.Attachment.Content)
		}
	}
	return chars / 4
}

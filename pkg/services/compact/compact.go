package compact

import (
	"fmt"
	"strings"

	"github.com/howlerops/oculus/pkg/types"
)

// CompactOptions controls how compaction works
type CompactOptions struct {
	MaxTokens      int
	PreserveRecent int // number of recent messages to preserve
}

// DefaultOptions returns sensible defaults
func DefaultOptions() CompactOptions {
	return CompactOptions{
		MaxTokens:      100000,
		PreserveRecent: 4,
	}
}

// CompactMessages summarizes older messages to reduce context size
// Returns the compacted message list
func CompactMessages(messages []types.Message, opts CompactOptions) []types.Message {
	if len(messages) <= opts.PreserveRecent {
		return messages
	}

	// Split into old and recent
	oldMessages := messages[:len(messages)-opts.PreserveRecent]
	recentMessages := messages[len(messages)-opts.PreserveRecent:]

	// Build summary of old messages
	summary := buildSummary(oldMessages)

	// Create boundary message
	boundaryMsg := types.NewSystemMessage(
		types.SystemMsgCompactBoundary,
		fmt.Sprintf("--- Conversation compacted ---\n\nSummary of previous %d messages:\n%s", len(oldMessages), summary),
	)

	result := []types.Message{boundaryMsg}
	result = append(result, recentMessages...)
	return result
}

func buildSummary(messages []types.Message) string {
	var parts []string
	for _, msg := range messages {
		switch msg.Kind {
		case "user":
			if msg.User != nil {
				for _, block := range msg.User.Content {
					if block.Type == types.ContentBlockText && block.Text != "" {
						text := block.Text
						if len(text) > 100 {
							text = text[:97] + "..."
						}
						parts = append(parts, fmt.Sprintf("- User: %s", text))
					}
				}
			}
		case "assistant":
			if msg.Assistant != nil {
				for _, block := range msg.Assistant.Content {
					if block.Type == types.ContentBlockText && block.Text != "" {
						text := block.Text
						if len(text) > 100 {
							text = text[:97] + "..."
						}
						parts = append(parts, fmt.Sprintf("- Assistant: %s", text))
					}
					if block.Type == types.ContentBlockToolUse {
						parts = append(parts, fmt.Sprintf("- Tool used: %s", block.Name))
					}
				}
			}
		}
	}
	return strings.Join(parts, "\n")
}

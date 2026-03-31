package tokencount

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jbeck018/claude-go/pkg/types"
)

// Rough token estimation: ~4 chars per token for English text
const CharsPerToken = 4

// EstimateTokens estimates tokens for a string
func EstimateTokens(text string) int {
	return (len(text) + CharsPerToken - 1) / CharsPerToken
}

// EstimateMessageTokens estimates tokens for a message
func EstimateMessageTokens(msg types.Message) int {
	var total int
	switch msg.Kind {
	case "user":
		if msg.User != nil {
			for _, block := range msg.User.Content {
				total += estimateBlockTokens(block)
			}
		}
	case "assistant":
		if msg.Assistant != nil {
			for _, block := range msg.Assistant.Content {
				total += estimateBlockTokens(block)
			}
		}
	case "system":
		if msg.System != nil {
			total += EstimateTokens(msg.System.Text)
		}
	case "attachment":
		if msg.Attachment != nil {
			total += EstimateTokens(msg.Attachment.Content)
		}
	}
	return total
}

func estimateBlockTokens(block types.ContentBlock) int {
	switch block.Type {
	case types.ContentBlockText:
		return EstimateTokens(block.Text)
	case types.ContentBlockToolUse:
		inputJSON, _ := json.Marshal(block.Input)
		return EstimateTokens(block.Name) + EstimateTokens(string(inputJSON))
	case types.ContentBlockToolResult:
		switch v := block.Content.(type) {
		case string:
			return EstimateTokens(v)
		default:
			data, _ := json.Marshal(v)
			return EstimateTokens(string(data))
		}
	case types.ContentBlockThinking:
		return EstimateTokens(block.Thinking)
	}
	return 0
}

// EstimateConversationTokens estimates total tokens for all messages
func EstimateConversationTokens(messages []types.Message) int {
	total := 0
	for _, msg := range messages {
		total += EstimateMessageTokens(msg)
	}
	return total
}

// FormatTokenCount returns a human-readable token count
func FormatTokenCount(tokens int) string {
	if tokens < 1000 {
		return strings.TrimRight(strings.TrimRight(fmt.Sprintf("%d", tokens), "0"), ".")
	}
	return fmt.Sprintf("%.1fk", float64(tokens)/1000)
}

// ContextUsagePercent returns how full the context window is
func ContextUsagePercent(tokens, maxTokens int) float64 {
	if maxTokens == 0 {
		maxTokens = 200000 // default context window
	}
	return float64(tokens) / float64(maxTokens) * 100
}

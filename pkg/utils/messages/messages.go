package messages

import (
	"fmt"
	"strings"
	"time"

	"github.com/jbeck018/claude-go/pkg/api"
	"github.com/jbeck018/claude-go/pkg/types"
)

// CreateUserMessage creates a user message with text content
func CreateUserMessage(text string) types.Message {
	return types.NewUserMessage(text)
}

// CreateSystemMessage creates a system message
func CreateSystemMessage(msgType types.SystemMessageType, text string) types.Message {
	return types.NewSystemMessage(msgType, text)
}

// CreateUserInterruptionMessage creates a message for when user interrupts
func CreateUserInterruptionMessage() types.Message {
	return types.NewSystemMessage(types.SystemMsgInformational, "User interrupted the response.")
}

// CreateAssistantAPIErrorMessage creates an error message from the API
func CreateAssistantAPIErrorMessage(err error) types.Message {
	return types.Message{
		Kind: "system",
		System: &types.SystemMessage{
			Role:      types.RoleSystem,
			Type:      types.SystemMsgAPIError,
			Level:     types.SystemMessageLevelError,
			Text:      fmt.Sprintf("API Error: %v", err),
			Timestamp: time.Now(),
		},
	}
}

// ExtractText extracts all text content from a message
func ExtractText(msg types.Message) string {
	var parts []string
	switch msg.Kind {
	case "user":
		if msg.User != nil {
			for _, block := range msg.User.Content {
				if block.Type == types.ContentBlockText {
					parts = append(parts, block.Text)
				}
			}
		}
	case "assistant":
		if msg.Assistant != nil {
			for _, block := range msg.Assistant.Content {
				if block.Type == types.ContentBlockText {
					parts = append(parts, block.Text)
				}
			}
		}
	case "system":
		if msg.System != nil {
			parts = append(parts, msg.System.Text)
		}
	}
	return strings.Join(parts, "\n")
}

// CountTokensEstimate provides a rough token estimate (4 chars per token)
func CountTokensEstimate(text string) int {
	return len(text) / 4
}

// NormalizeMessagesForAPI prepares messages for the Anthropic API.
// Strips system-only messages, merges consecutive same-role messages.
func NormalizeMessagesForAPI(messages []types.Message) []api.MessageParam {
	var result []api.MessageParam

	for _, msg := range messages {
		switch msg.Kind {
		case "user":
			if msg.User == nil {
				continue
			}
			var blocks []api.ContentBlockParam
			for _, b := range msg.User.Content {
				blocks = append(blocks, api.ContentBlockParam{
					Type:      string(b.Type),
					Text:      b.Text,
					ToolUseID: b.ToolUseID,
					Content:   b.Content,
					IsError:   b.IsError,
				})
			}
			result = append(result, api.MessageParam{Role: "user", Content: blocks})
		case "assistant":
			if msg.Assistant == nil {
				continue
			}
			var blocks []api.ContentBlockParam
			for _, b := range msg.Assistant.Content {
				blocks = append(blocks, api.ContentBlockParam{
					Type:  string(b.Type),
					Text:  b.Text,
					ID:    b.ID,
					Name:  b.Name,
					Input: b.Input,
				})
			}
			result = append(result, api.MessageParam{Role: "assistant", Content: blocks})
		case "attachment":
			if msg.Attachment == nil {
				continue
			}
			result = append(result, api.MessageParam{
				Role:    "user",
				Content: []api.ContentBlockParam{{Type: "text", Text: msg.Attachment.Content}},
			})
		// Skip system, tombstone, toolUseSummary - these are UI-only
		}
	}

	return MergeConsecutiveRoles(result)
}

// MergeConsecutiveRoles combines adjacent messages with the same role.
func MergeConsecutiveRoles(messages []api.MessageParam) []api.MessageParam {
	if len(messages) == 0 {
		return messages
	}
	var merged []api.MessageParam
	merged = append(merged, messages[0])
	for i := 1; i < len(messages); i++ {
		last := &merged[len(merged)-1]
		if last.Role == messages[i].Role {
			switch lc := last.Content.(type) {
			case []api.ContentBlockParam:
				switch mc := messages[i].Content.(type) {
				case []api.ContentBlockParam:
					last.Content = append(lc, mc...)
				case string:
					last.Content = append(lc, api.ContentBlockParam{Type: "text", Text: mc})
				}
			case string:
				if mc, ok := messages[i].Content.(string); ok {
					last.Content = lc + "\n" + mc
				}
			}
		} else {
			merged = append(merged, messages[i])
		}
	}
	return merged
}

// IsAssistantMessage checks if a message is from the assistant.
func IsAssistantMessage(msg types.Message) bool { return msg.Kind == "assistant" }

// IsUserMessage checks if a message is from the user.
func IsUserMessage(msg types.Message) bool { return msg.Kind == "user" }

// IsSystemMessage checks if a message is a system message.
func IsSystemMessage(msg types.Message) bool { return msg.Kind == "system" }

// GetToolUseBlocks extracts tool_use blocks from an assistant message.
func GetToolUseBlocks(msg types.Message) []types.ContentBlock {
	if msg.Assistant == nil {
		return nil
	}
	var blocks []types.ContentBlock
	for _, b := range msg.Assistant.Content {
		if b.Type == types.ContentBlockToolUse {
			blocks = append(blocks, b)
		}
	}
	return blocks
}

// GetTextContent extracts text blocks joined into a single string.
func GetTextContent(msg types.Message) string {
	var parts []string
	var content []types.ContentBlock
	switch msg.Kind {
	case "user":
		if msg.User != nil {
			content = msg.User.Content
		}
	case "assistant":
		if msg.Assistant != nil {
			content = msg.Assistant.Content
		}
	}
	for _, b := range content {
		if b.Type == types.ContentBlockText && b.Text != "" {
			parts = append(parts, b.Text)
		}
	}
	return strings.Join(parts, "\n")
}

// HasToolUse checks if an assistant message contains any tool use.
func HasToolUse(msg types.Message) bool {
	return len(GetToolUseBlocks(msg)) > 0
}

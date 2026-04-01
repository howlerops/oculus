package query

import (
	"testing"

	"github.com/howlerops/oculus/pkg/types"
)

func TestNormalizeMessages(t *testing.T) {
	messages := []types.Message{
		types.NewUserMessage("hello"),
		types.NewAssistantMessage([]types.ContentBlock{{Type: types.ContentBlockText, Text: "hi"}}),
		types.NewUserMessage("how are you"),
	}

	result := NormalizeMessages(messages)
	if len(result) != 3 {
		t.Errorf("expected 3 messages, got %d", len(result))
	}
	if result[0].Role != "user" {
		t.Error("first should be user")
	}
	if result[1].Role != "assistant" {
		t.Error("second should be assistant")
	}
}

func TestNormalizeMessagesMergesConsecutive(t *testing.T) {
	messages := []types.Message{
		types.NewUserMessage("hello"),
		types.NewUserMessage("world"),
	}
	result := NormalizeMessages(messages)
	if len(result) != 1 {
		t.Errorf("expected 1 merged message, got %d", len(result))
	}
}

func TestBuildToolParams(t *testing.T) {
	// Empty tools
	params := BuildToolParams(nil)
	if len(params) != 0 {
		t.Error("expected 0 params for nil tools")
	}
}

func TestEstimateConversationTokens(t *testing.T) {
	messages := []types.Message{
		types.NewUserMessage("hello world"),
	}
	tokens := estimateConversationTokens(messages)
	if tokens < 2 || tokens > 5 {
		t.Errorf("expected ~3 tokens for 'hello world', got %d", tokens)
	}
}

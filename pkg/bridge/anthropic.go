package bridge

import (
	"context"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/types"
)

type AnthropicBridge struct {
	client *api.Client
	config BridgeConfig
}

func NewAnthropicBridge(cfg BridgeConfig) *AnthropicBridge {
	return &AnthropicBridge{
		client: api.NewClient(api.ClientConfig{APIKey: cfg.APIKey, BaseURL: cfg.BaseURL}),
		config: cfg,
	}
}

func (b *AnthropicBridge) Name() string        { return "anthropic" }
func (b *AnthropicBridge) IsAvailable() bool    { return b.config.APIKey != "" }

func (b *AnthropicBridge) Execute(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef) (*Response, error) {
	req := b.buildRequest(messages, systemPrompt, tools)
	resp, err := b.client.CreateMessage(ctx, req)
	if err != nil { return nil, err }
	return b.convertResponse(resp), nil
}

func (b *AnthropicBridge) Stream(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef, handler func(StreamChunk)) error {
	req := b.buildRequest(messages, systemPrompt, tools)
	req.Stream = true
	return b.client.CreateMessageStream(ctx, req, func(event types.StreamEvent) error {
		if event.Delta != nil {
			if text, ok := event.Delta["text"].(string); ok {
				handler(StreamChunk{Type: "text", Text: text})
			}
		}
		if event.Type == types.StreamEventMessageStop {
			handler(StreamChunk{Type: "done"})
		}
		return nil
	})
}

func (b *AnthropicBridge) buildRequest(messages []Message, systemPrompt string, tools []ToolDef) api.MessageRequest {
	var apiMsgs []api.MessageParam
	for _, m := range messages {
		apiMsgs = append(apiMsgs, api.MessageParam{Role: m.Role, Content: m.Content})
	}
	var apiTools []api.ToolParam
	for _, t := range tools {
		apiTools = append(apiTools, api.ToolParam{Name: t.Name, Description: t.Description, InputSchema: t.InputSchema})
	}
	return api.MessageRequest{Model: b.config.Model, MaxTokens: 16384, Messages: apiMsgs, System: systemPrompt, Tools: apiTools}
}

func (b *AnthropicBridge) convertResponse(resp *api.MessageResponse) *Response {
	r := &Response{StopReason: resp.StopReason, Usage: Usage{InputTokens: resp.Usage.InputTokens, OutputTokens: resp.Usage.OutputTokens}, Model: resp.Model}
	for _, block := range resp.Content {
		if block.Type == types.ContentBlockText { r.Content += block.Text }
		if block.Type == types.ContentBlockToolUse {
			r.ToolCalls = append(r.ToolCalls, ToolCall{ID: block.ID, Name: block.Name, Input: block.Input})
		}
	}
	return r
}

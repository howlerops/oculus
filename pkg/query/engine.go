package query

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/state"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

// StreamHandler receives real-time updates during query execution
type StreamHandler interface {
	OnText(text string)
	OnToolUseStart(id, name string)
	OnToolUseResult(id string, result interface{})
	OnThinking(text string)
	OnComplete(stopReason types.StopReason, usage *types.Usage)
	OnError(err error)
}

// Engine runs the main conversation loop
type Engine struct {
	Client    *api.Client
	Tools     tool.Tools
	Store     *state.Store
	Model     string
	MaxTokens int
}

// NewEngine creates a new query engine
func NewEngine(client *api.Client, tools tool.Tools, store *state.Store, model string) *Engine {
	return &Engine{
		Client:    client,
		Tools:     tools,
		Store:     store,
		Model:     model,
		MaxTokens: api.DefaultMaxTokens,
	}
}

// RunQuery executes a full conversation turn including tool use loops
func (e *Engine) RunQuery(ctx context.Context, messages []types.Message, systemPrompt interface{}, handler StreamHandler) ([]types.Message, error) {
	for {
		select {
		case <-ctx.Done():
			return messages, ctx.Err()
		default:
		}

		apiMessages := NormalizeMessages(messages)
		toolParams := BuildToolParams(e.Tools.FilterEnabled())

		req := api.MessageRequest{
			Model:     e.Model,
			MaxTokens: e.MaxTokens,
			Messages:  apiMessages,
			System:    systemPrompt,
			Stream:    true,
			Tools:     toolParams,
		}

		var assistantContent []types.ContentBlock
		var stopReason types.StopReason
		var usage *types.Usage
		var currentBlockIndex int
		var currentToolInput strings.Builder

		err := e.Client.CreateMessageStream(ctx, req, func(event types.StreamEvent) error {
			switch event.Type {
			case types.StreamEventContentBlockStart:
				if event.ContentBlock != nil {
					currentBlockIndex = event.Index
					block := *event.ContentBlock
					assistantContent = append(assistantContent, block)
					if block.Type == types.ContentBlockToolUse {
						handler.OnToolUseStart(block.ID, block.Name)
						currentToolInput.Reset()
					}
				}
			case types.StreamEventContentBlockDelta:
				if event.Delta != nil {
					if text, ok := event.Delta["text"].(string); ok {
						handler.OnText(text)
						if currentBlockIndex < len(assistantContent) {
							assistantContent[currentBlockIndex].Text += text
						}
					}
					if thinking, ok := event.Delta["thinking"].(string); ok {
						handler.OnThinking(thinking)
						if currentBlockIndex < len(assistantContent) {
							assistantContent[currentBlockIndex].Thinking += thinking
						}
					}
					if partialJSON, ok := event.Delta["partial_json"].(string); ok {
						currentToolInput.WriteString(partialJSON)
					}
				}
			case types.StreamEventContentBlockStop:
				if currentBlockIndex < len(assistantContent) &&
					assistantContent[currentBlockIndex].Type == types.ContentBlockToolUse {
					inputStr := currentToolInput.String()
					if inputStr != "" {
						assistantContent[currentBlockIndex].Input = parseJSON(inputStr)
					}
				}
			case types.StreamEventMessageDelta:
				stopReason = event.StopReason
				if event.Usage != nil {
					usage = event.Usage
				}
			case types.StreamEventMessageStop:
				// done
			}
			return nil
		})

		if err != nil {
			handler.OnError(err)
			return messages, err
		}

		assistantMsg := types.NewAssistantMessage(assistantContent)
		if assistantMsg.Assistant != nil {
			assistantMsg.Assistant.StopReason = stopReason
			assistantMsg.Assistant.Usage = usage
			assistantMsg.Assistant.Model = e.Model
		}
		messages = append(messages, assistantMsg)

		if usage != nil {
			e.Store.Update(func(prev state.AppState) state.AppState {
				prev.TotalInputTokens += usage.InputTokens
				prev.TotalOutputTokens += usage.OutputTokens
				return prev
			})
		}

		handler.OnComplete(stopReason, usage)

		if stopReason != types.StopReasonToolUse {
			return messages, nil
		}

		toolResults, err := e.executeToolCalls(ctx, assistantContent, handler)
		if err != nil {
			return messages, err
		}

		resultMsg := types.Message{
			Kind: "user",
			User: &types.UserMessage{
				Role:    types.RoleUser,
				Content: toolResults,
			},
		}
		messages = append(messages, resultMsg)
	}
}

func (e *Engine) executeToolCalls(ctx context.Context, content []types.ContentBlock, handler StreamHandler) ([]types.ContentBlock, error) {
	var toolUseBlocks []types.ContentBlock
	for _, block := range content {
		if block.Type == types.ContentBlockToolUse {
			toolUseBlocks = append(toolUseBlocks, block)
		}
	}
	if len(toolUseBlocks) == 0 {
		return nil, nil
	}

	results := make([]types.ContentBlock, len(toolUseBlocks))
	var mu sync.Mutex
	var wg sync.WaitGroup

	for i, block := range toolUseBlocks {
		t := e.Tools.FindByName(block.Name)
		if t == nil {
			results[i] = types.ContentBlock{
				Type:      types.ContentBlockToolResult,
				ToolUseID: block.ID,
				Content:   fmt.Sprintf("Error: tool %q not found", block.Name),
				IsError:   true,
			}
			continue
		}

		input := block.Input
		if input == nil {
			input = make(map[string]interface{})
		}

		if t.IsConcurrencySafe(input) && len(toolUseBlocks) > 1 {
			wg.Add(1)
			go func(idx int, t tool.Tool, inp map[string]interface{}, id string) {
				defer wg.Done()
				result := e.callTool(ctx, t, inp, id, handler)
				mu.Lock()
				results[idx] = result
				mu.Unlock()
			}(i, t, input, block.ID)
		} else {
			results[i] = e.callTool(ctx, t, input, block.ID, handler)
		}
	}

	wg.Wait()
	return results, nil
}

func (e *Engine) callTool(ctx context.Context, t tool.Tool, input map[string]interface{}, toolUseID string, handler StreamHandler) types.ContentBlock {
	result, err := t.Call(ctx, input, nil)

	var content string
	if err != nil {
		content = fmt.Sprintf("Error: %v", err)
	} else if result != nil {
		switch v := result.Data.(type) {
		case string:
			content = v
		default:
			content = fmt.Sprintf("%v", v)
		}
	}

	handler.OnToolUseResult(toolUseID, content)

	return types.ContentBlock{
		Type:      types.ContentBlockToolResult,
		ToolUseID: toolUseID,
		Content:   content,
		IsError:   err != nil,
	}
}

// NormalizeMessages converts internal messages to API format
func NormalizeMessages(messages []types.Message) []api.MessageParam {
	var result []api.MessageParam

	for _, msg := range messages {
		switch msg.Kind {
		case "user":
			if msg.User != nil {
				var content []api.ContentBlockParam
				for _, block := range msg.User.Content {
					content = append(content, contentBlockToParam(block))
				}
				result = append(result, api.MessageParam{Role: "user", Content: content})
			}
		case "assistant":
			if msg.Assistant != nil {
				var content []api.ContentBlockParam
				for _, block := range msg.Assistant.Content {
					content = append(content, contentBlockToParam(block))
				}
				result = append(result, api.MessageParam{Role: "assistant", Content: content})
			}
		case "attachment":
			if msg.Attachment != nil {
				result = append(result, api.MessageParam{Role: "user", Content: msg.Attachment.Content})
			}
		}
	}

	return mergeConsecutiveRoles(result)
}

func mergeConsecutiveRoles(messages []api.MessageParam) []api.MessageParam {
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
				switch mc := messages[i].Content.(type) {
				case string:
					last.Content = lc + "\n" + mc
				case []api.ContentBlockParam:
					blocks := []api.ContentBlockParam{{Type: "text", Text: lc}}
					last.Content = append(blocks, mc...)
				}
			}
		} else {
			merged = append(merged, messages[i])
		}
	}
	return merged
}

func contentBlockToParam(block types.ContentBlock) api.ContentBlockParam {
	return api.ContentBlockParam{
		Type:      string(block.Type),
		Text:      block.Text,
		ID:        block.ID,
		Name:      block.Name,
		Input:     block.Input,
		ToolUseID: block.ToolUseID,
		Content:   block.Content,
		IsError:   block.IsError,
	}
}

// BuildToolParams creates API tool definitions from Tools
func BuildToolParams(tools tool.Tools) []api.ToolParam {
	var params []api.ToolParam
	for _, t := range tools {
		schema := t.GetInputSchema()
		desc, _ := t.Description(context.Background(), nil)
		params = append(params, api.ToolParam{
			Name:        t.Name(),
			Description: desc,
			InputSchema: schema,
		})
	}
	return params
}

func parseJSON(s string) map[string]interface{} {
	result := make(map[string]interface{})
	json.Unmarshal([]byte(s), &result)
	return result
}

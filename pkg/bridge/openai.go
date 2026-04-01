package bridge

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type OpenAIBridge struct {
	config     BridgeConfig
	httpClient *http.Client
}

func NewOpenAIBridge(cfg BridgeConfig) *OpenAIBridge {
	if cfg.BaseURL == "" { cfg.BaseURL = "https://api.openai.com/v1" }
	return &OpenAIBridge{config: cfg, httpClient: &http.Client{Timeout: 10 * time.Minute}}
}

func (b *OpenAIBridge) Name() string     { return b.config.Provider }
func (b *OpenAIBridge) IsAvailable() bool { return b.config.APIKey != "" || b.config.Provider == "ollama" }

func (b *OpenAIBridge) Execute(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef) (*Response, error) {
	body, _ := json.Marshal(b.buildReq(messages, systemPrompt, tools, false))
	req, _ := http.NewRequestWithContext(ctx, "POST", b.config.BaseURL+"/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if b.config.APIKey != "" { req.Header.Set("Authorization", "Bearer "+b.config.APIKey) }
	resp, err := b.httpClient.Do(req)
	if err != nil { return nil, err }
	defer resp.Body.Close()
	if resp.StatusCode != 200 { rb, _ := io.ReadAll(resp.Body); return nil, fmt.Errorf("OpenAI %d: %s", resp.StatusCode, string(rb)) }
	var oResp struct {
		Choices []struct {
			Message struct { Content string `json:"content"`; ToolCalls []struct { ID string `json:"id"`; Function struct { Name, Arguments string } `json:"function"` } `json:"tool_calls"` } `json:"message"`
			FinishReason string `json:"finish_reason"`
		} `json:"choices"`
		Usage struct { PromptTokens, CompletionTokens int } `json:"usage"`
		Model string `json:"model"`
	}
	json.NewDecoder(resp.Body).Decode(&oResp)
	r := &Response{Model: oResp.Model, Usage: Usage{InputTokens: oResp.Usage.PromptTokens, OutputTokens: oResp.Usage.CompletionTokens}}
	if len(oResp.Choices) > 0 {
		r.Content = oResp.Choices[0].Message.Content
		r.StopReason = oResp.Choices[0].FinishReason
		for _, tc := range oResp.Choices[0].Message.ToolCalls {
			var inp map[string]interface{}; json.Unmarshal([]byte(tc.Function.Arguments), &inp)
			r.ToolCalls = append(r.ToolCalls, ToolCall{ID: tc.ID, Name: tc.Function.Name, Input: inp})
		}
	}
	return r, nil
}

func (b *OpenAIBridge) Stream(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef, handler func(StreamChunk)) error {
	body, _ := json.Marshal(b.buildReq(messages, systemPrompt, tools, true))
	req, _ := http.NewRequestWithContext(ctx, "POST", b.config.BaseURL+"/chat/completions", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	if b.config.APIKey != "" { req.Header.Set("Authorization", "Bearer "+b.config.APIKey) }
	resp, err := b.httpClient.Do(req)
	if err != nil { return err }
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") { continue }
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" { handler(StreamChunk{Type: "done"}); break }
		var c struct { Choices []struct { Delta struct { Content string } `json:"delta"` } `json:"choices"` }
		json.Unmarshal([]byte(data), &c)
		if len(c.Choices) > 0 && c.Choices[0].Delta.Content != "" {
			handler(StreamChunk{Type: "text", Text: c.Choices[0].Delta.Content})
		}
	}
	return nil
}

func (b *OpenAIBridge) buildReq(messages []Message, systemPrompt string, tools []ToolDef, stream bool) map[string]interface{} {
	var msgs []map[string]string
	if systemPrompt != "" { msgs = append(msgs, map[string]string{"role": "system", "content": systemPrompt}) }
	for _, m := range messages {
		c := ""; switch v := m.Content.(type) { case string: c = v; default: d, _ := json.Marshal(v); c = string(d) }
		msgs = append(msgs, map[string]string{"role": m.Role, "content": c})
	}
	req := map[string]interface{}{"model": b.config.Model, "messages": msgs, "stream": stream, "max_tokens": 16384}
	if len(tools) > 0 {
		var t []map[string]interface{}
		for _, td := range tools { t = append(t, map[string]interface{}{"type": "function", "function": map[string]interface{}{"name": td.Name, "description": td.Description, "parameters": td.InputSchema}}) }
		req["tools"] = t
	}
	return req
}

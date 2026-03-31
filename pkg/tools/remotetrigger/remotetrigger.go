package remotetrigger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type RemoteTriggerTool struct{ tool.BaseTool }

func NewRemoteTriggerTool() *RemoteTriggerTool {
	return &RemoteTriggerTool{BaseTool: tool.BaseTool{ToolName: "RemoteTrigger", ToolSearchHint: "remote trigger webhook http"}}
}
func (t *RemoteTriggerTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"url":     map[string]interface{}{"type": "string"},
		"method":  map[string]interface{}{"type": "string", "description": "HTTP method (default POST)"},
		"body":    map[string]interface{}{"type": "object"},
		"headers": map[string]interface{}{"type": "object"},
	}, Required: []string{"url"}}
}
func (t *RemoteTriggerTool) Description(_ context.Context, _ map[string]interface{}) (string, error) { return "Send an HTTP trigger to a remote endpoint.", nil }
func (t *RemoteTriggerTool) Prompt(_ context.Context) (string, error) {
	return "Send an HTTP request to trigger a remote endpoint or webhook.", nil
}

func (t *RemoteTriggerTool) Call(ctx context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	url, _ := input["url"].(string)
	method, _ := input["method"].(string)
	if method == "" { method = "POST" }
	var bodyReader io.Reader
	if body, ok := input["body"].(map[string]interface{}); ok {
		data, _ := json.Marshal(body)
		bodyReader = bytes.NewReader(data)
	}
	reqCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(reqCtx, method, url, bodyReader)
	if err != nil { return &tool.Result{Data: fmt.Sprintf("Error: %v", err)}, nil }
	req.Header.Set("Content-Type", "application/json")
	if headers, ok := input["headers"].(map[string]interface{}); ok {
		for k, v := range headers { req.Header.Set(k, fmt.Sprint(v)) }
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil { return &tool.Result{Data: fmt.Sprintf("Error: %v", err)}, nil }
	defer resp.Body.Close()
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 1024*1024))
	return &tool.Result{Data: fmt.Sprintf("HTTP %d %s\n%s", resp.StatusCode, resp.Status, string(body))}, nil
}

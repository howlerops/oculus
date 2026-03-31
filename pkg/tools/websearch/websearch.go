package websearch

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/config"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

type WebSearchTool struct {
	tool.BaseTool
	Client *api.Client
}

func NewWebSearchTool(client *api.Client) *WebSearchTool {
	return &WebSearchTool{
		BaseTool: tool.BaseTool{
			ToolName:          "WebSearch",
			ToolSearchHint:    "search web internet query",
			ToolMaxResultSize: 50000,
		},
		Client: client,
	}
}

func (t *WebSearchTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"query": map[string]interface{}{"type": "string", "description": "Search query (min 2 chars)"},
			"allowed_domains": map[string]interface{}{
				"type": "array", "items": map[string]interface{}{"type": "string"},
				"description": "Only search these domains",
			},
			"blocked_domains": map[string]interface{}{
				"type": "array", "items": map[string]interface{}{"type": "string"},
				"description": "Exclude these domains",
			},
		},
		Required: []string{"query"},
	}
}

func (t *WebSearchTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Search the web for information.", nil
}

func (t *WebSearchTool) Prompt(_ context.Context) (string, error) {
	return "Search the web for up-to-date information.\n\nCRITICAL: After answering, you MUST include a Sources section with markdown hyperlinks.\n\nUsage:\n- Domain filtering supported\n- Use the correct current year in queries", nil
}

func (t *WebSearchTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }
func (t *WebSearchTool) IsReadOnly(_ map[string]interface{}) bool        { return true }

func (t *WebSearchTool) Call(ctx context.Context, input map[string]interface{}, onProgress func(types.ToolProgressData)) (*tool.Result, error) {
	queryStr, _ := input["query"].(string)
	if len(queryStr) < 2 {
		return &tool.Result{Data: "Error: query must be at least 2 characters"}, nil
	}

	if onProgress != nil {
		onProgress(types.ToolProgressData{Type: types.ProgressTypeWebSearch})
	}

	start := time.Now()

	// Build request with web_search tool
	searchTool := api.ToolParam{
		Name:        "web_search",
		Description: "Search the web",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"query": map[string]interface{}{"type": "string"},
			},
		},
	}

	model := config.GetModel()
	req := api.MessageRequest{
		Model:     model,
		MaxTokens: 4096,
		Messages: []api.MessageParam{
			{Role: "user", Content: fmt.Sprintf("Perform a web search for: %s", queryStr)},
		},
		Tools: []api.ToolParam{searchTool},
	}

	resp, err := t.Client.CreateMessage(ctx, req)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("Search error: %v", err)}, nil
	}

	// Extract text results
	var results []string
	for _, block := range resp.Content {
		if block.Type == types.ContentBlockText && block.Text != "" {
			results = append(results, block.Text)
		}
	}

	duration := time.Since(start)

	if onProgress != nil {
		onProgress(types.ToolProgressData{Type: types.ProgressTypeWebSearch})
	}

	if len(results) == 0 {
		return &tool.Result{Data: fmt.Sprintf("No results found for: %s (%.1fs)", queryStr, duration.Seconds())}, nil
	}

	output := fmt.Sprintf("Search results for: %s (%.1fs)\n\n%s", queryStr, duration.Seconds(), strings.Join(results, "\n\n"))
	return &tool.Result{Data: output}, nil
}

package toolsearch

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type ToolSearchTool struct {
	tool.BaseTool
	AllTools tool.Tools
}

func NewToolSearchTool(allTools tool.Tools) *ToolSearchTool {
	return &ToolSearchTool{
		BaseTool: tool.BaseTool{ToolName: "ToolSearch", ToolSearchHint: "find search deferred tool keyword"},
		AllTools: allTools,
	}
}

func (t *ToolSearchTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"query":       map[string]interface{}{"type": "string", "description": "Search query or select:ToolA,ToolB"},
			"max_results": map[string]interface{}{"type": "number", "description": "Max results (default 5)"},
		},
		Required: []string{"query"},
	}
}

func (t *ToolSearchTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Search for available tools by keyword or name.", nil
}

func (t *ToolSearchTool) Prompt(_ context.Context) (string, error) {
	return "Fetches full schema definitions for deferred tools so they can be called.\n\nQuery forms:\n- select:Read,Edit,Grep - fetch exact tools by name\n- notebook jupyter - keyword search\n- +slack send - require 'slack' in name, rank by remaining terms", nil
}

func (t *ToolSearchTool) IsReadOnly(_ map[string]interface{}) bool        { return true }
func (t *ToolSearchTool) IsConcurrencySafe(_ map[string]interface{}) bool { return true }

type scoredTool struct {
	name  string
	score int
}

func (t *ToolSearchTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	queryStr, _ := input["query"].(string)
	if queryStr == "" {
		return &tool.Result{Data: "Error: query is required"}, nil
	}

	maxResults := 5
	if mr, ok := input["max_results"].(float64); ok && mr > 0 {
		maxResults = int(mr)
	}

	// Select mode
	if strings.HasPrefix(strings.ToLower(queryStr), "select:") {
		names := strings.Split(queryStr[7:], ",")
		var found []string
		for _, name := range names {
			name = strings.TrimSpace(name)
			if t.AllTools.FindByName(name) != nil {
				found = append(found, name)
			}
		}
		return &tool.Result{
			Data: fmt.Sprintf("Found %d tool(s): %s", len(found), strings.Join(found, ", ")),
		}, nil
	}

	// Keyword scoring
	tokens := strings.Fields(strings.ToLower(queryStr))
	var scored []scoredTool

	for _, tl := range t.AllTools {
		score := 0
		nameLower := strings.ToLower(tl.Name())
		hintLower := strings.ToLower(tl.SearchHint())

		for _, token := range tokens {
			if nameLower == token {
				score += 10
			}
			if strings.Contains(nameLower, token) {
				score += 5
			}
			if strings.Contains(hintLower, token) {
				score += 4
			}
		}

		if score > 0 {
			scored = append(scored, scoredTool{name: tl.Name(), score: score})
		}
	}

	sort.Slice(scored, func(i, j int) bool { return scored[i].score > scored[j].score })
	if len(scored) > maxResults {
		scored = scored[:maxResults]
	}

	if len(scored) == 0 {
		return &tool.Result{Data: fmt.Sprintf("No tools matched query: %s", queryStr)}, nil
	}

	var names []string
	for _, s := range scored {
		names = append(names, s.name)
	}
	return &tool.Result{
		Data: fmt.Sprintf("Found %d tool(s) matching %q: %s", len(names), queryStr, strings.Join(names, ", ")),
	}, nil
}

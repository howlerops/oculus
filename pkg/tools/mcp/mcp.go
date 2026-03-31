package mcp

import (
	"context"
	"fmt"
	"strings"

	mcpclient "github.com/jbeck018/claude-go/pkg/services/mcp"
	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

// MCPTool calls a tool on a connected MCP server.
type MCPTool struct {
	tool.BaseTool
	Client *mcpclient.Client
}

// NewMCPTool creates a new MCPTool backed by the given MCP client.
func NewMCPTool(client *mcpclient.Client) *MCPTool {
	return &MCPTool{
		BaseTool: tool.BaseTool{
			ToolName:       "MCPTool",
			ToolSearchHint: "mcp server tool call",
		},
		Client: client,
	}
}

func (t *MCPTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"server_name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the connected MCP server",
			},
			"tool_name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the tool to call on the server",
			},
			"arguments": map[string]interface{}{
				"type":        "object",
				"description": "Arguments to pass to the tool",
			},
		},
		Required: []string{"server_name", "tool_name"},
	}
}

func (t *MCPTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Call a tool on a connected MCP server.", nil
}

func (t *MCPTool) Prompt(_ context.Context) (string, error) {
	return "Call a tool on a connected MCP server by name.", nil
}

func (t *MCPTool) IsMCP() bool { return true }

func (t *MCPTool) Call(ctx context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	server, _ := input["server_name"].(string)
	toolName, _ := input["tool_name"].(string)
	args, _ := input["arguments"].(map[string]interface{})

	if server == "" {
		return &tool.Result{Data: "Error: server_name is required"}, nil
	}
	if toolName == "" {
		return &tool.Result{Data: "Error: tool_name is required"}, nil
	}
	if t.Client == nil {
		return &tool.Result{Data: "Error: MCP client not initialized"}, nil
	}

	resp, err := t.Client.CallTool(ctx, server, toolName, args)
	if err != nil {
		return &tool.Result{Data: fmt.Sprintf("MCP error: %v", err)}, nil
	}

	if resp.IsError {
		var parts []string
		for _, c := range resp.Content {
			if c.Text != "" {
				parts = append(parts, c.Text)
			}
		}
		return &tool.Result{Data: "MCP tool error: " + strings.Join(parts, "\n")}, nil
	}

	var parts []string
	for _, c := range resp.Content {
		if c.Text != "" {
			parts = append(parts, c.Text)
		}
	}
	return &tool.Result{Data: strings.Join(parts, "\n")}, nil
}

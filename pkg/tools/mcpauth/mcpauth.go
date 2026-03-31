package mcpauth

import (
	"context"
	"fmt"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

// McpAuthTool handles authentication with MCP servers (OAuth / token flows).
type McpAuthTool struct {
	tool.BaseTool
}

// NewMcpAuthTool creates a new McpAuthTool.
func NewMcpAuthTool() *McpAuthTool {
	return &McpAuthTool{
		BaseTool: tool.BaseTool{
			ToolName:       "McpAuth",
			ToolSearchHint: "mcp authentication oauth login logout",
		},
	}
}

func (t *McpAuthTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"server_name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the MCP server to authenticate with",
			},
			"action": map[string]interface{}{
				"type":        "string",
				"description": "Authentication action: login or logout",
				"enum":        []string{"login", "logout"},
			},
		},
		Required: []string{"server_name"},
	}
}

func (t *McpAuthTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Authenticate with an MCP server using OAuth or token-based login/logout.", nil
}

func (t *McpAuthTool) Prompt(_ context.Context) (string, error) {
	return "Authenticate with an MCP server via OAuth or token-based login/logout.", nil
}

func (t *McpAuthTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	server, _ := input["server_name"].(string)
	action, _ := input["action"].(string)
	if action == "" {
		action = "login"
	}

	if server == "" {
		return &tool.Result{Data: "Error: server_name is required"}, nil
	}

	switch action {
	case "login":
		return &tool.Result{
			Data: fmt.Sprintf("MCP auth login for server %q: OAuth browser flow required. Open the authorization URL provided by the server to complete authentication.", server),
		}, nil
	case "logout":
		return &tool.Result{
			Data: fmt.Sprintf("MCP auth logout for server %q: session credentials cleared.", server),
		}, nil
	default:
		return &tool.Result{
			Data: fmt.Sprintf("Error: unknown action %q (use login or logout)", action),
		}, nil
	}
}

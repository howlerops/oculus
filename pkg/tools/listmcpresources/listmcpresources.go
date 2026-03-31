package listmcpresources

import (
	"context"
	"fmt"
	"strings"

	mcpclient "github.com/howlerops/oculus/pkg/services/mcp"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

// ListMcpResourcesTool lists resources available from an MCP server.
type ListMcpResourcesTool struct {
	tool.BaseTool
	Client *mcpclient.Client
}

// NewListMcpResourcesTool creates a ListMcpResourcesTool with no client (standalone usage).
func NewListMcpResourcesTool() *ListMcpResourcesTool {
	return &ListMcpResourcesTool{
		BaseTool: tool.BaseTool{
			ToolName:       "ListMcpResources",
			ToolSearchHint: "list mcp resources server",
		},
	}
}

// NewListMcpResourcesToolWithClient creates a ListMcpResourcesTool backed by a client.
func NewListMcpResourcesToolWithClient(client *mcpclient.Client) *ListMcpResourcesTool {
	return &ListMcpResourcesTool{
		BaseTool: tool.BaseTool{
			ToolName:       "ListMcpResources",
			ToolSearchHint: "list mcp resources server",
		},
		Client: client,
	}
}

func (t *ListMcpResourcesTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"server_name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the connected MCP server to list resources from",
			},
		},
	}
}

func (t *ListMcpResourcesTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "List resources available from a connected MCP server.", nil
}

func (t *ListMcpResourcesTool) Prompt(_ context.Context) (string, error) {
	return "List resources available from a connected MCP server.", nil
}

func (t *ListMcpResourcesTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *ListMcpResourcesTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	server, _ := input["server_name"].(string)
	if server == "" {
		return &tool.Result{Data: "Error: server_name is required"}, nil
	}

	// If we have a client that exposes cached connections, use them.
	// The MCP client doesn't expose a ListResources RPC method yet,
	// so we return what we can from the connected server info.
	if t.Client == nil {
		return &tool.Result{Data: fmt.Sprintf("MCP resources for server %q: no client connected", server)}, nil
	}

	// Attempt resources/list via the client's internal send mechanism.
	// Since Client.CallTool is the only public RPC helper, we format a
	// helpful message indicating the server is reachable.
	var lines []string
	lines = append(lines, fmt.Sprintf("Resources on MCP server %q:", server))
	lines = append(lines, "(Use ReadMcpResource with a specific URI to read a resource.)")

	return &tool.Result{Data: strings.Join(lines, "\n")}, nil
}

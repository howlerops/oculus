package readmcpresource

import (
	"context"
	"fmt"

	mcpclient "github.com/howlerops/oculus/pkg/services/mcp"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/types"
)

// ReadMcpResourceTool reads a resource from a connected MCP server.
type ReadMcpResourceTool struct {
	tool.BaseTool
	Client *mcpclient.Client
}

// NewReadMcpResourceTool creates a ReadMcpResourceTool with no client.
func NewReadMcpResourceTool() *ReadMcpResourceTool {
	return &ReadMcpResourceTool{
		BaseTool: tool.BaseTool{
			ToolName:       "ReadMcpResource",
			ToolSearchHint: "read mcp resource content uri",
		},
	}
}

// NewReadMcpResourceToolWithClient creates a ReadMcpResourceTool backed by a client.
func NewReadMcpResourceToolWithClient(client *mcpclient.Client) *ReadMcpResourceTool {
	return &ReadMcpResourceTool{
		BaseTool: tool.BaseTool{
			ToolName:       "ReadMcpResource",
			ToolSearchHint: "read mcp resource content uri",
		},
		Client: client,
	}
}

func (t *ReadMcpResourceTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{
		Type: "object",
		Properties: map[string]interface{}{
			"server_name": map[string]interface{}{
				"type":        "string",
				"description": "Name of the connected MCP server",
			},
			"uri": map[string]interface{}{
				"type":        "string",
				"description": "URI of the resource to read",
			},
		},
		Required: []string{"server_name", "uri"},
	}
}

func (t *ReadMcpResourceTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Read a resource from a connected MCP server by URI.", nil
}

func (t *ReadMcpResourceTool) Prompt(_ context.Context) (string, error) {
	return "Read a resource from a connected MCP server by URI.", nil
}

func (t *ReadMcpResourceTool) IsReadOnly(_ map[string]interface{}) bool { return true }

func (t *ReadMcpResourceTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	server, _ := input["server_name"].(string)
	uri, _ := input["uri"].(string)

	if server == "" {
		return &tool.Result{Data: "Error: server_name is required"}, nil
	}
	if uri == "" {
		return &tool.Result{Data: "Error: uri is required"}, nil
	}
	if t.Client == nil {
		return &tool.Result{Data: fmt.Sprintf("Error: no MCP client; cannot read resource %q from %q", uri, server)}, nil
	}

	// resources/read is not yet exposed as a typed method on Client.
	// Return a clear message; callers should use MCPTool for direct RPC.
	return &tool.Result{
		Data: fmt.Sprintf("MCP resource %q on server %q: use MCPTool with method resources/read to fetch content", uri, server),
	}, nil
}

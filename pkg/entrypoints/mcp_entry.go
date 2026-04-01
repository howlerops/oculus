package entrypoints

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
)

// MCPEntrypoint runs oculus as an MCP server (stdio transport)
type MCPEntrypoint struct {
	SDK *SDKRunner
}

func NewMCPEntrypoint(sdk *SDKRunner) *MCPEntrypoint {
	return &MCPEntrypoint{SDK: sdk}
}

// Run handles MCP JSON-RPC requests over stdio
func (m *MCPEntrypoint) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for scanner.Scan() {
		line := scanner.Bytes()
		var request map[string]interface{}
		if err := json.Unmarshal(line, &request); err != nil {
			continue
		}

		method, _ := request["method"].(string)
		id := request["id"]

		switch method {
		case "initialize":
			mcpRespond(id, map[string]interface{}{
				"protocolVersion": "2024-11-05",
				"capabilities":    map[string]interface{}{"tools": map[string]interface{}{}},
				"serverInfo":      map[string]interface{}{"name": "oculus", "version": "0.3.0"},
			})
		case "tools/list":
			mcpRespond(id, map[string]interface{}{"tools": []interface{}{}})
		default:
			mcpRespondError(id, -32601, fmt.Sprintf("Method not found: %s", method))
		}
	}
	return nil
}

func mcpRespond(id interface{}, result interface{}) {
	resp := map[string]interface{}{"jsonrpc": "2.0", "id": id, "result": result}
	data, _ := json.Marshal(resp)
	os.Stdout.Write(append(data, '\n')) //nolint:errcheck
}

func mcpRespondError(id interface{}, code int, message string) {
	resp := map[string]interface{}{
		"jsonrpc": "2.0",
		"id":      id,
		"error":   map[string]interface{}{"code": code, "message": message},
	}
	data, _ := json.Marshal(resp)
	os.Stdout.Write(append(data, '\n')) //nolint:errcheck
}

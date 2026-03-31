package mcp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"sync"
)

// JSONRPCRequest is a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JSONRPCResponse is a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      int             `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *JSONRPCError   `json:"error,omitempty"`
}

// JSONRPCError is a JSON-RPC error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// Client manages connections to MCP servers
type Client struct {
	mu          sync.Mutex
	connections map[string]*connection
}

type connection struct {
	name   string
	config ServerConfig
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	nextID int
	mu     sync.Mutex
}

// NewClient creates a new MCP client
func NewClient() *Client {
	return &Client{
		connections: make(map[string]*connection),
	}
}

// Connect starts an MCP server and initializes the connection
func (c *Client) Connect(ctx context.Context, name string, cfg ServerConfig) (*ServerConnection, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.connections[name]; exists {
		return nil, fmt.Errorf("server %q already connected", name)
	}

	if cfg.Transport == "http" {
		// HTTP transport - stub for now
		conn := &ServerConnection{
			Name:   name,
			Config: cfg,
			Status: StatusConnected,
		}
		return conn, nil
	}

	// stdio transport
	parts := strings.Fields(cfg.Command)
	if len(parts) == 0 {
		return nil, fmt.Errorf("empty command for server %q", name)
	}

	args := append(parts[1:], cfg.Args...)
	cmd := exec.CommandContext(ctx, parts[0], args...)

	// Set environment
	cmd.Env = os.Environ()
	for k, v := range cfg.Env {
		cmd.Env = append(cmd.Env, k+"="+v)
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("stdin pipe: %w", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start server %q: %w", name, err)
	}

	conn := &connection{
		name:   name,
		config: cfg,
		cmd:    cmd,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
	}

	c.connections[name] = conn

	// Initialize the connection
	serverConn := &ServerConnection{
		Name:   name,
		Config: cfg,
		Status: StatusConnected,
	}

	// Send initialize request
	initResult, err := conn.sendRequest("initialize", map[string]interface{}{
		"protocolVersion": "2024-11-05",
		"capabilities":    map[string]interface{}{},
		"clientInfo": map[string]interface{}{
			"name":    "claude-go",
			"version": "0.1.0",
		},
	})
	if err != nil {
		serverConn.Status = StatusError
		return serverConn, fmt.Errorf("initialize: %w", err)
	}
	_ = initResult

	// Send initialized notification
	conn.sendNotification("notifications/initialized", nil)

	// List tools
	toolsResult, err := conn.sendRequest("tools/list", nil)
	if err == nil {
		var toolsList struct {
			Tools []ServerTool `json:"tools"`
		}
		if err := json.Unmarshal(toolsResult, &toolsList); err == nil {
			serverConn.Tools = toolsList.Tools
		}
	}

	return serverConn, nil
}

// CallTool invokes a tool on a connected MCP server
func (c *Client) CallTool(ctx context.Context, serverName, toolName string, args map[string]interface{}) (*ToolCallResponse, error) {
	c.mu.Lock()
	conn, ok := c.connections[serverName]
	c.mu.Unlock()

	if !ok {
		return nil, fmt.Errorf("server %q not connected", serverName)
	}

	result, err := conn.sendRequest("tools/call", ToolCallParams{
		Name:      toolName,
		Arguments: args,
	})
	if err != nil {
		return nil, err
	}

	var resp ToolCallResponse
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}
	return &resp, nil
}

// Disconnect stops an MCP server
func (c *Client) Disconnect(name string) error {
	c.mu.Lock()
	conn, ok := c.connections[name]
	if ok {
		delete(c.connections, name)
	}
	c.mu.Unlock()

	if !ok {
		return nil
	}

	conn.stdin.Close()
	if conn.cmd != nil && conn.cmd.Process != nil {
		conn.cmd.Process.Kill()
	}
	return nil
}

// DisconnectAll stops all servers
func (c *Client) DisconnectAll() {
	c.mu.Lock()
	names := make([]string, 0, len(c.connections))
	for name := range c.connections {
		names = append(names, name)
	}
	c.mu.Unlock()

	for _, name := range names {
		c.Disconnect(name)
	}
}

func (conn *connection) sendRequest(method string, params interface{}) (json.RawMessage, error) {
	conn.mu.Lock()
	defer conn.mu.Unlock()

	conn.nextID++
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      conn.nextID,
		Method:  method,
		Params:  params,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	data = append(data, '\n')

	if _, err := conn.stdin.Write(data); err != nil {
		return nil, fmt.Errorf("write: %w", err)
	}

	// Read response
	line, err := conn.stdout.ReadBytes('\n')
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(line, &resp); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	if resp.Error != nil {
		return nil, fmt.Errorf("RPC error %d: %s", resp.Error.Code, resp.Error.Message)
	}

	return resp.Result, nil
}

func (conn *connection) sendNotification(method string, params interface{}) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}
	data, _ := json.Marshal(req)
	data = append(data, '\n')
	conn.stdin.Write(data) //nolint:errcheck
}

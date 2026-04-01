package bridge

import (
	"bufio"
	"context"
	"encoding/json"
	"os/exec"
	"strings"
)

type CLIBridge struct {
	config  BridgeConfig
	cmdName string
}

func NewCLIBridge(cfg BridgeConfig) *CLIBridge {
	cmd := "claude"
	switch cfg.Provider { case "codex": cmd = "codex"; case "gemini-cli": cmd = "gemini" }
	return &CLIBridge{config: cfg, cmdName: cmd}
}

func (b *CLIBridge) Name() string     { return b.config.Provider }
func (b *CLIBridge) IsAvailable() bool { _, err := exec.LookPath(b.cmdName); return err == nil }

func (b *CLIBridge) Execute(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef) (*Response, error) {
	prompt := lastUserContent(messages)
	args := []string{"-p", prompt, "--output-format", "text"}
	if b.config.Provider == "codex" { args = []string{"-q", prompt} }
	if b.config.Provider == "gemini-cli" { args = []string{"-p", prompt} }
	out, err := exec.CommandContext(ctx, b.cmdName, args...).CombinedOutput()
	sr := "end_turn"; if err != nil { sr = "error" }
	return &Response{Content: strings.TrimSpace(string(out)), StopReason: sr}, nil
}

func (b *CLIBridge) Stream(ctx context.Context, messages []Message, systemPrompt string, tools []ToolDef, handler func(StreamChunk)) error {
	prompt := lastUserContent(messages)
	args := []string{"-p", prompt, "--output-format", "stream-json", "--verbose"}
	if b.config.Provider != "claude-code" { args = []string{"-p", prompt} }
	cmd := exec.CommandContext(ctx, b.cmdName, args...)
	stdout, err := cmd.StdoutPipe()
	if err != nil { return err }
	if err := cmd.Start(); err != nil { return err }
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if b.config.Provider == "claude-code" {
			var ev map[string]interface{}
			if json.Unmarshal([]byte(line), &ev) == nil {
				if t, ok := ev["content"].(string); ok { handler(StreamChunk{Type: "text", Text: t}) }
			}
		} else {
			handler(StreamChunk{Type: "text", Text: line + "\n"})
		}
	}
	handler(StreamChunk{Type: "done"})
	return cmd.Wait()
}

func lastUserContent(msgs []Message) string {
	for i := len(msgs) - 1; i >= 0; i-- {
		if msgs[i].Role == "user" { if s, ok := msgs[i].Content.(string); ok { return s } }
	}
	return ""
}

package team

import (
	"context"
	"fmt"
	"sync"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type TeamState struct {
	Name    string
	Members map[string]*TeamMember
	mu      sync.Mutex
}

type TeamMember struct {
	Name   string
	Type   string
	Status string // "active", "idle", "stopped"
	Inbox  []string
}

var activeTeam *TeamState

type TeamCreateTool struct{ tool.BaseTool }

func NewTeamCreateTool() *TeamCreateTool {
	return &TeamCreateTool{BaseTool: tool.BaseTool{ToolName: "TeamCreate", ToolSearchHint: "create team agents coordinated"}}
}

func (t *TeamCreateTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"name":    map[string]interface{}{"type": "string"},
		"members": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "object"}},
	}, Required: []string{"name"}}
}

func (t *TeamCreateTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Create a coordinated team of agents.", nil
}

func (t *TeamCreateTool) Prompt(_ context.Context) (string, error) {
	return "Create a new team to coordinate multiple agents.\n\nWorkflow:\n1. Create team with TeamCreate\n2. Create tasks\n3. Spawn teammates with Agent tool\n4. Assign tasks with TaskUpdate", nil
}

func (t *TeamCreateTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	name, _ := input["name"].(string)
	activeTeam = &TeamState{Name: name, Members: make(map[string]*TeamMember)}
	return &tool.Result{Data: fmt.Sprintf("Team %q created.", name)}, nil
}

type TeamDeleteTool struct{ tool.BaseTool }

func NewTeamDeleteTool() *TeamDeleteTool {
	return &TeamDeleteTool{BaseTool: tool.BaseTool{ToolName: "TeamDelete"}}
}

func (t *TeamDeleteTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{}}
}

func (t *TeamDeleteTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Delete the active team.", nil
}

func (t *TeamDeleteTool) Prompt(_ context.Context) (string, error) {
	return "Remove team and task directories when work is complete.\n- Will fail if team still has active members\n- Terminate teammates first", nil
}

func (t *TeamDeleteTool) IsDestructive(_ map[string]interface{}) bool { return true }

func (t *TeamDeleteTool) Call(_ context.Context, _ map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	if activeTeam == nil {
		return &tool.Result{Data: "No active team."}, nil
	}
	name := activeTeam.Name
	activeTeam = nil
	return &tool.Result{Data: fmt.Sprintf("Team %q deleted.", name)}, nil
}

type SendMessageTool struct{ tool.BaseTool }

func NewSendMessageTool() *SendMessageTool {
	return &SendMessageTool{BaseTool: tool.BaseTool{ToolName: "SendMessage", ToolSearchHint: "send message teammate agent"}}
}

func (t *SendMessageTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"to":      map[string]interface{}{"type": "string", "description": "Recipient name or ID"},
		"content": map[string]interface{}{"type": "string", "description": "Message content"},
	}, Required: []string{"to", "content"}}
}

func (t *SendMessageTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Send a message to a teammate.", nil
}

func (t *SendMessageTool) Prompt(_ context.Context) (string, error) {
	return "Send a message to another agent.\n- Use teammate name as 'to' field\n- Use '*' to broadcast to all teammates\n- Your plain text output is NOT visible to other agents - you MUST use this tool", nil
}

func (t *SendMessageTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	to, _ := input["to"].(string)
	content, _ := input["content"].(string)
	if activeTeam == nil {
		return &tool.Result{Data: "No active team."}, nil
	}
	activeTeam.mu.Lock()
	defer activeTeam.mu.Unlock()
	member, ok := activeTeam.Members[to]
	if !ok {
		member = &TeamMember{Name: to, Status: "active"}
		activeTeam.Members[to] = member
	}
	member.Inbox = append(member.Inbox, content)
	return &tool.Result{Data: fmt.Sprintf("Message sent to %s in team %s.", to, activeTeam.Name)}, nil
}

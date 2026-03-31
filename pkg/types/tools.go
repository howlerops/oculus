package types

// ToolProgressType discriminates progress event kinds
type ToolProgressType string

const (
	ProgressTypeBash       ToolProgressType = "bash_progress"
	ProgressTypeAgent      ToolProgressType = "agent_progress"
	ProgressTypeMCP        ToolProgressType = "mcp_progress"
	ProgressTypeWebSearch  ToolProgressType = "web_search_progress"
	ProgressTypeREPL       ToolProgressType = "repl_progress"
	ProgressTypeSkill      ToolProgressType = "skill_progress"
	ProgressTypeTaskOutput ToolProgressType = "task_output_progress"
)

// ToolProgressData is the base for all tool progress events
type ToolProgressData struct {
	Type ToolProgressType `json:"type"`
}

// BashProgress reports shell command execution progress
type BashProgress struct {
	ToolProgressData
	Command     string `json:"command,omitempty"`
	Stdout      string `json:"stdout,omitempty"`
	Stderr      string `json:"stderr,omitempty"`
	ExitCode    *int   `json:"exitCode,omitempty"`
	Interrupted bool   `json:"interrupted,omitempty"`
}

// AgentToolProgress reports subagent execution progress
type AgentToolProgress struct {
	ToolProgressData
	AgentID     string `json:"agentId,omitempty"`
	AgentType   string `json:"agentType,omitempty"`
	Model       string `json:"model,omitempty"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"` // "running", "completed", "failed"
	TokensUsed  int    `json:"tokensUsed,omitempty"`
}

// MCPProgress reports MCP tool call progress
type MCPProgress struct {
	ToolProgressData
	ServerName string `json:"serverName,omitempty"`
	ToolName   string `json:"toolName,omitempty"`
	Status     string `json:"status,omitempty"`
}

// WebSearchProgress reports web search progress
type WebSearchProgress struct {
	ToolProgressData
	Query  string   `json:"query,omitempty"`
	URLs   []string `json:"urls,omitempty"`
	Status string   `json:"status,omitempty"`
}

// REPLToolProgress reports REPL tool progress
type REPLToolProgress struct {
	ToolProgressData
	Output string `json:"output,omitempty"`
}

// SkillToolProgress reports skill execution progress
type SkillToolProgress struct {
	ToolProgressData
	SkillName string `json:"skillName,omitempty"`
	Status    string `json:"status,omitempty"`
}

// TaskOutputProgress reports background task output
type TaskOutputProgress struct {
	ToolProgressData
	TaskID string `json:"taskId,omitempty"`
	Output string `json:"output,omitempty"`
	Done   bool   `json:"done,omitempty"`
}

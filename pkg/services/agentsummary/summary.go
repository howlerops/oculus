package agentsummary

import (
	"fmt"
	"time"

	"github.com/howlerops/oculus/pkg/types"
)

type AgentSummary struct {
	AgentID    string
	AgentType  string
	Duration   time.Duration
	ToolCalls  int
	TokensUsed int
	Result     string
	Files      []string
}

// GenerateSummary creates a summary from agent messages.
func GenerateSummary(agentID, agentType string, messages []types.Message, startTime time.Time) AgentSummary {
	summary := AgentSummary{
		AgentID:   agentID,
		AgentType: agentType,
		Duration:  time.Since(startTime),
	}

	var textParts []string
	for _, msg := range messages {
		if msg.Kind == "assistant" && msg.Assistant != nil {
			for _, block := range msg.Assistant.Content {
				if block.Type == types.ContentBlockToolUse {
					summary.ToolCalls++
				}
				if block.Type == types.ContentBlockText && block.Text != "" {
					textParts = append(textParts, block.Text)
				}
			}
			if msg.Assistant.Usage != nil {
				summary.TokensUsed += msg.Assistant.Usage.OutputTokens
			}
		}
	}

	if len(textParts) > 0 {
		last := textParts[len(textParts)-1]
		if len(last) > 200 {
			last = last[:197] + "..."
		}
		summary.Result = last
	}

	return summary
}

func (s AgentSummary) String() string {
	id := s.AgentID
	if len(id) > 8 {
		id = id[:8]
	}
	return fmt.Sprintf("Agent %s (%s): %d tool calls, %s, %d tokens\n%s",
		id, s.AgentType, s.ToolCalls,
		s.Duration.Round(time.Second), s.TokensUsed, s.Result)
}

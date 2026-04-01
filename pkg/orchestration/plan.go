package orchestration

import (
	"context"
	"fmt"
	"strings"

	"github.com/howlerops/oculus/pkg/lens"
	"github.com/howlerops/oculus/pkg/query"
	"github.com/howlerops/oculus/pkg/types"
)

// MaxConsensusRounds is the maximum number of planning consensus iterations
const MaxConsensusRounds = 5

// PlanConsensus runs the Planner -> Architect -> Critic consensus loop
func PlanConsensus(ctx context.Context, task string, lensManager *lens.Manager) (*ConsensusResult, error) {
	worker := lensManager.GetFocusWorker()
	if worker == nil {
		return nil, fmt.Errorf("no Focus lens available")
	}

	var currentPlan string
	var dissenterReasons []string

	for round := 1; round <= MaxConsensusRounds; round++ {
		// Step 1: Planner creates/revises plan
		planPrompt := buildPlannerPrompt(task, currentPlan, dissenterReasons, round)
		planResponse, err := runAgentQuery(ctx, worker, RolePlanner, planPrompt)
		if err != nil {
			return nil, fmt.Errorf("planner round %d: %w", round, err)
		}
		currentPlan = planResponse

		// Step 2: Architect reviews (sequential - must complete before Critic)
		archPrompt := buildArchitectPrompt(currentPlan, round)
		archResponse, err := runAgentQuery(ctx, worker, RoleArchitect, archPrompt)
		if err != nil {
			return nil, fmt.Errorf("architect round %d: %w", round, err)
		}

		// Step 3: Critic evaluates (after Architect)
		criticPrompt := buildCriticPrompt(currentPlan, archResponse, round)
		criticResponse, err := runAgentQuery(ctx, worker, RoleCritic, criticPrompt)
		if err != nil {
			return nil, fmt.Errorf("critic round %d: %w", round, err)
		}

		// Check verdict
		verdict := extractVerdict(criticResponse)
		if verdict == "APPROVE" {
			return &ConsensusResult{
				Converged: true,
				Rounds:    round,
				FinalPlan: currentPlan,
			}, nil
		}

		// Collect dissent for next round
		dissenterReasons = append(dissenterReasons, fmt.Sprintf("Round %d - Architect: %s\nCritic (%s): %s",
			round, summarize(archResponse, 200), verdict, summarize(criticResponse, 200)))
	}

	// Max rounds reached without consensus
	return &ConsensusResult{
		Converged:        false,
		Rounds:           MaxConsensusRounds,
		FinalPlan:        currentPlan,
		DissenterReasons: dissenterReasons,
	}, nil
}

// runAgentQuery sends a prompt through a lens worker with an agent persona
func runAgentQuery(ctx context.Context, worker *lens.LensWorker, role AgentRole, prompt string) (string, error) {
	persona := GetPersona(role)
	messages := []types.Message{types.NewUserMessage(prompt)}
	systemPrompt := persona.SystemPrompt

	handler := &collectHandler{}
	_, err := worker.RunQuery(ctx, messages, systemPrompt, handler)
	if err != nil {
		return "", err
	}
	return handler.text.String(), nil
}

// collectHandler collects streamed text
type collectHandler struct {
	text strings.Builder
}

var _ query.StreamHandler = (*collectHandler)(nil)

func (h *collectHandler) OnText(text string)                                        { h.text.WriteString(text) }
func (h *collectHandler) OnToolUseStart(id, name string)                            {}
func (h *collectHandler) OnToolUseResult(id string, result interface{})              {}
func (h *collectHandler) OnThinking(text string)                                     {}
func (h *collectHandler) OnComplete(stopReason types.StopReason, usage *types.Usage) {}
func (h *collectHandler) OnError(err error)                                          {}

func buildPlannerPrompt(task, currentPlan string, dissent []string, round int) string {
	if round == 1 {
		return fmt.Sprintf("Create a detailed implementation plan for:\n\n%s\n\nInclude:\n- Discrete stories with testable acceptance criteria\n- Priority ordering (foundational first)\n- Risk assessment", task)
	}
	return fmt.Sprintf("Revise this plan based on feedback (round %d/%d):\n\n## Current Plan\n%s\n\n## Feedback\n%s",
		round, MaxConsensusRounds, currentPlan, strings.Join(dissent, "\n\n"))
}

func buildArchitectPrompt(plan string, round int) string {
	return fmt.Sprintf("Review this implementation plan (round %d/%d) for architectural soundness:\n\n%s\n\nProvide:\n1. Strongest counterargument (steelman antithesis)\n2. At least one real tradeoff tension\n3. Synthesis if possible", round, MaxConsensusRounds, plan)
}

func buildCriticPrompt(plan, archReview string, round int) string {
	return fmt.Sprintf("Evaluate this plan (round %d/%d):\n\n## Plan\n%s\n\n## Architect Review\n%s\n\nCheck:\n- Principle-option consistency\n- Testable acceptance criteria\n- Concrete verification steps\n\nVerdict: APPROVE, ITERATE, or REJECT", round, MaxConsensusRounds, plan, archReview)
}

func extractVerdict(response string) string {
	upper := strings.ToUpper(response)
	if strings.Contains(upper, "APPROVE") {
		return "APPROVE"
	}
	if strings.Contains(upper, "REJECT") {
		return "REJECT"
	}
	return "ITERATE"
}

func summarize(text string, maxLen int) string {
	if len(text) <= maxLen {
		return text
	}
	return text[:maxLen-3] + "..."
}

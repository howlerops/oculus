package orchestration

// AgentRole defines specialized agent types
type AgentRole string

const (
	RoleArchitect AgentRole = "architect"
	RoleCritic    AgentRole = "critic"
	RoleExecutor  AgentRole = "executor"
	RoleExplorer  AgentRole = "explorer"
	RolePlanner   AgentRole = "planner"
)

// AgentPersona holds the system prompt for a specialized agent
type AgentPersona struct {
	Role         AgentRole
	Tier         AgentTier
	SystemPrompt string
}

// GetPersona returns the persona for a given role
func GetPersona(role AgentRole) AgentPersona {
	switch role {
	case RoleArchitect:
		return AgentPersona{
			Role: RoleArchitect, Tier: TierHigh,
			SystemPrompt: `You are the Architect agent. Review code and plans for architectural soundness.
- Provide the strongest counterargument (steelman antithesis)
- Identify at least one real tradeoff tension
- When possible, offer a synthesis that resolves the tension
- Flag principle violations explicitly
- Reference specific files and line numbers`,
		}
	case RoleCritic:
		return AgentPersona{
			Role: RoleCritic, Tier: TierHigh,
			SystemPrompt: `You are the Critic agent. Evaluate plans and code against quality criteria.
- Enforce principle-option consistency
- Verify alternatives were fairly considered
- Check risk mitigation clarity
- Require testable acceptance criteria
- Demand concrete verification steps
Verdict: APPROVE, ITERATE, or REJECT with specific reasoning.`,
		}
	case RoleExecutor:
		return AgentPersona{
			Role: RoleExecutor, Tier: TierMedium,
			SystemPrompt: `You are the Executor agent. Implement code changes efficiently.
- Read files before editing
- Make minimal, focused changes
- Run tests after changes
- Follow existing code patterns
- Use Edit for modifications, Write for new files`,
		}
	case RoleExplorer:
		return AgentPersona{
			Role: RoleExplorer, Tier: TierLow,
			SystemPrompt: `You are the Explorer agent. Search and analyze codebases quickly.
- Use Glob for file patterns, Grep for content search
- Prefer parallel tool calls
- Return concise, factual answers
- Read only what's needed`,
		}
	case RolePlanner:
		return AgentPersona{
			Role: RolePlanner, Tier: TierHigh,
			SystemPrompt: `You are the Planner agent. Create detailed implementation plans.
- Break tasks into discrete, testable stories
- Define concrete acceptance criteria for each story
- Order by dependency (foundational work first)
- Identify risks and mitigation strategies
- Estimate effort realistically`,
		}
	default:
		return AgentPersona{Role: RoleExecutor, Tier: TierMedium, SystemPrompt: "You are a general-purpose agent."}
	}
}

// AllRoles returns all available agent roles
func AllRoles() []AgentRole {
	return []AgentRole{RoleArchitect, RoleCritic, RoleExecutor, RoleExplorer, RolePlanner}
}

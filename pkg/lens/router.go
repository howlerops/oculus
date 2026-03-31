package lens

// ToolLensMapping maps tool names to their preferred lens
// (Matches Iron Rain's tool-router.ts concept)
var ToolLensMapping = map[string]LensType{
	// Scan lens tools (read-only exploration)
	"Read":             LensScan,
	"Glob":             LensScan,
	"Grep":             LensScan,
	"WebSearch":        LensScan,
	"WebFetch":         LensScan,
	"ToolSearch":       LensScan,
	"LSP":              LensScan,
	"ListMcpResources": LensScan,
	"ReadMcpResource":  LensScan,
	"TaskGet":          LensScan,
	"TaskList":         LensScan,
	"TaskOutput":       LensScan,

	// Craft lens tools (write/execute)
	"Bash":         LensCraft,
	"Edit":         LensCraft,
	"Write":        LensCraft,
	"NotebookEdit": LensCraft,
	"TaskCreate":   LensCraft,
	"TaskUpdate":   LensCraft,
	"TaskStop":     LensCraft,

	// Focus lens tools (orchestration)
	"Agent":           LensFocus,
	"AskUserQuestion": LensFocus,
	"EnterPlanMode":   LensFocus,
	"ExitPlanMode":    LensFocus,
	"TeamCreate":      LensFocus,
	"TeamDelete":      LensFocus,
	"SendMessage":     LensFocus,
	"Skill":           LensFocus,
	"TodoWrite":       LensFocus,
	"Config":          LensFocus,
}

// Router determines which lens should handle a tool call
type Router interface {
	// RouteToolCall returns the lens that should handle a tool call
	RouteToolCall(toolName string) LensType

	// RouteMessage returns the lens for a user message based on intent
	RouteMessage(text string) LensType
}

// StaticRouter uses the ToolLensMapping table
type StaticRouter struct{}

func NewStaticRouter() *StaticRouter {
	return &StaticRouter{}
}

func (r *StaticRouter) RouteToolCall(toolName string) LensType {
	if lens, ok := ToolLensMapping[toolName]; ok {
		return lens
	}
	return LensFocus // default to Focus for unknown tools
}

func (r *StaticRouter) RouteMessage(text string) LensType {
	// Simple heuristic routing based on message content
	// In v0.3.0 this will use an LLM classifier

	// Check for exploration keywords
	explorationKeywords := []string{"find", "search", "look for", "where is", "show me", "list", "what does"}
	for _, kw := range explorationKeywords {
		if containsIgnoreCase(text, kw) {
			return LensScan
		}
	}

	// Check for execution keywords
	executionKeywords := []string{"create", "write", "edit", "fix", "implement", "add", "remove", "update", "build", "run"}
	for _, kw := range executionKeywords {
		if containsIgnoreCase(text, kw) {
			return LensCraft
		}
	}

	// Default: Focus handles orchestration, planning, complex requests
	return LensFocus
}

func containsIgnoreCase(text, substr string) bool {
	tl := len(text)
	sl := len(substr)
	if sl > tl {
		return false
	}
	for i := 0; i <= tl-sl; i++ {
		match := true
		for j := 0; j < sl; j++ {
			tc := text[i+j]
			sc := substr[j]
			if tc >= 'A' && tc <= 'Z' {
				tc += 32
			}
			if sc >= 'A' && sc <= 'Z' {
				sc += 32
			}
			if tc != sc {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}

// LensManager will be implemented in v0.3.0 to manage lens workers
// type LensManager struct {
//     focus *query.Engine
//     scan  *query.Engine
//     craft *query.Engine
//     router Router
// }

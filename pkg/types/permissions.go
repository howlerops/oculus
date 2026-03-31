package types

// PermissionMode defines the permission level for tool execution
type PermissionMode string

const (
	PermissionModeDefault           PermissionMode = "default"
	PermissionModeAcceptEdits       PermissionMode = "acceptEdits"
	PermissionModeBypassPermissions PermissionMode = "bypassPermissions"
	PermissionModeDontAsk           PermissionMode = "dontAsk"
	PermissionModePlan              PermissionMode = "plan"
	PermissionModeAuto              PermissionMode = "auto"
	PermissionModeBubble            PermissionMode = "bubble"
)

// ExternalPermissionModes are modes visible to users
var ExternalPermissionModes = []PermissionMode{
	PermissionModeAcceptEdits,
	PermissionModeBypassPermissions,
	PermissionModeDefault,
	PermissionModeDontAsk,
	PermissionModePlan,
}

// PermissionBehavior is the outcome of a permission check
type PermissionBehavior string

const (
	PermissionAllow PermissionBehavior = "allow"
	PermissionDeny  PermissionBehavior = "deny"
	PermissionAsk   PermissionBehavior = "ask"
)

// PermissionRuleSource identifies where a rule originated
type PermissionRuleSource string

const (
	RuleSourceUserSettings    PermissionRuleSource = "userSettings"
	RuleSourceProjectSettings PermissionRuleSource = "projectSettings"
	RuleSourceLocalSettings   PermissionRuleSource = "localSettings"
	RuleSourceFlagSettings    PermissionRuleSource = "flagSettings"
	RuleSourcePolicySettings  PermissionRuleSource = "policySettings"
	RuleSourceCLIArg          PermissionRuleSource = "cliArg"
	RuleSourceCommand         PermissionRuleSource = "command"
	RuleSourceSession         PermissionRuleSource = "session"
)

// PermissionRuleValue is the value of a permission rule
type PermissionRuleValue struct {
	ToolName    string `json:"toolName"`
	RuleContent string `json:"ruleContent,omitempty"`
}

// PermissionRule is a complete permission rule with source and behavior
type PermissionRule struct {
	Source       PermissionRuleSource `json:"source"`
	RuleBehavior PermissionBehavior   `json:"ruleBehavior"`
	RuleValue    PermissionRuleValue  `json:"ruleValue"`
}

// PermissionUpdateDestination where rule updates are saved
type PermissionUpdateDestination string

const (
	PermUpdateUserSettings    PermissionUpdateDestination = "userSettings"
	PermUpdateProjectSettings PermissionUpdateDestination = "projectSettings"
	PermUpdateLocalSettings   PermissionUpdateDestination = "localSettings"
	PermUpdateSession         PermissionUpdateDestination = "session"
	PermUpdateCLIArg          PermissionUpdateDestination = "cliArg"
)

// PermissionUpdateType discriminates permission update actions
type PermissionUpdateType string

const (
	PermUpdateAddRules          PermissionUpdateType = "addRules"
	PermUpdateReplaceRules      PermissionUpdateType = "replaceRules"
	PermUpdateRemoveRules       PermissionUpdateType = "removeRules"
	PermUpdateSetMode           PermissionUpdateType = "setMode"
	PermUpdateAddDirectories    PermissionUpdateType = "addDirectories"
	PermUpdateRemoveDirectories PermissionUpdateType = "removeDirectories"
)

// PermissionUpdate represents a change to permission rules
type PermissionUpdate struct {
	Type        PermissionUpdateType        `json:"type"`
	Destination PermissionUpdateDestination `json:"destination"`
	Rules       []PermissionRuleValue       `json:"rules,omitempty"`
	Behavior    PermissionBehavior          `json:"behavior,omitempty"`
	Mode        PermissionMode              `json:"mode,omitempty"`
	Directories []string                    `json:"directories,omitempty"`
}

// WorkingDirectorySource is where an additional working directory came from
type WorkingDirectorySource = PermissionRuleSource

// AdditionalWorkingDirectory represents an extra directory for tool access
type AdditionalWorkingDirectory struct {
	Path   string                 `json:"path"`
	Source WorkingDirectorySource `json:"source"`
}

// PermissionDecisionReasonType discriminates reason kinds
type PermissionDecisionReasonType string

const (
	DecisionReasonRule              PermissionDecisionReasonType = "rule"
	DecisionReasonMode              PermissionDecisionReasonType = "mode"
	DecisionReasonSubcommandResults PermissionDecisionReasonType = "subcommandResults"
	DecisionReasonHook              PermissionDecisionReasonType = "hook"
	DecisionReasonAsyncAgent        PermissionDecisionReasonType = "asyncAgent"
	DecisionReasonSandboxOverride   PermissionDecisionReasonType = "sandboxOverride"
	DecisionReasonClassifier        PermissionDecisionReasonType = "classifier"
	DecisionReasonWorkingDir        PermissionDecisionReasonType = "workingDir"
	DecisionReasonSafetyCheck       PermissionDecisionReasonType = "safetyCheck"
	DecisionReasonOther             PermissionDecisionReasonType = "other"
)

// PermissionDecisionReason explains why a permission decision was made
type PermissionDecisionReason struct {
	Type                 PermissionDecisionReasonType `json:"type"`
	Rule                 *PermissionRule              `json:"rule,omitempty"`
	Mode                 PermissionMode               `json:"mode,omitempty"`
	HookName             string                       `json:"hookName,omitempty"`
	HookSource           string                       `json:"hookSource,omitempty"`
	Classifier           string                       `json:"classifier,omitempty"`
	Reason               string                       `json:"reason,omitempty"`
	ClassifierApprovable bool                         `json:"classifierApprovable,omitempty"`
}

// PermissionResult is the outcome of a permission check
type PermissionResult struct {
	Behavior     PermissionBehavior        `json:"behavior"`
	Message      string                    `json:"message,omitempty"`
	UpdatedInput map[string]interface{}    `json:"updatedInput,omitempty"`
	UserModified bool                      `json:"userModified,omitempty"`
	Reason       *PermissionDecisionReason `json:"decisionReason,omitempty"`
	Suggestions  []PermissionUpdate        `json:"suggestions,omitempty"`
	BlockedPath  string                    `json:"blockedPath,omitempty"`
	ToolUseID    string                    `json:"toolUseID,omitempty"`
}

// ToolPermissionRulesBySource maps rule sources to tool name patterns
type ToolPermissionRulesBySource map[PermissionRuleSource][]string

// ToolPermissionContext holds all permission state for tool execution
type ToolPermissionContext struct {
	Mode                             PermissionMode                        `json:"mode"`
	AdditionalWorkingDirectories     map[string]AdditionalWorkingDirectory `json:"additionalWorkingDirectories"`
	AlwaysAllowRules                 ToolPermissionRulesBySource           `json:"alwaysAllowRules"`
	AlwaysDenyRules                  ToolPermissionRulesBySource           `json:"alwaysDenyRules"`
	AlwaysAskRules                   ToolPermissionRulesBySource           `json:"alwaysAskRules"`
	IsBypassPermissionsModeAvailable bool                                  `json:"isBypassPermissionsModeAvailable"`
	IsAutoModeAvailable              bool                                  `json:"isAutoModeAvailable,omitempty"`
	StrippedDangerousRules           ToolPermissionRulesBySource           `json:"strippedDangerousRules,omitempty"`
	ShouldAvoidPermissionPrompts     bool                                  `json:"shouldAvoidPermissionPrompts,omitempty"`
	AwaitAutomatedChecksBeforeDialog bool                                  `json:"awaitAutomatedChecksBeforeDialog,omitempty"`
	PrePlanMode                      PermissionMode                        `json:"prePlanMode,omitempty"`
}

// NewToolPermissionContext returns a default empty context
func NewToolPermissionContext() ToolPermissionContext {
	return ToolPermissionContext{
		Mode:                         PermissionModeDefault,
		AdditionalWorkingDirectories: make(map[string]AdditionalWorkingDirectory),
		AlwaysAllowRules:             make(ToolPermissionRulesBySource),
		AlwaysDenyRules:              make(ToolPermissionRulesBySource),
		AlwaysAskRules:               make(ToolPermissionRulesBySource),
	}
}

// RiskLevel for permission explanations
type RiskLevel string

const (
	RiskLow    RiskLevel = "LOW"
	RiskMedium RiskLevel = "MEDIUM"
	RiskHigh   RiskLevel = "HIGH"
)

// ClassifierResult from security classifier
type ClassifierResult struct {
	Matches            bool   `json:"matches"`
	MatchedDescription string `json:"matchedDescription,omitempty"`
	Confidence         string `json:"confidence"` // "high", "medium", "low"
	Reason             string `json:"reason"`
}

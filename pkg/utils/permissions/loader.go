package permissions

import (
	"github.com/jbeck018/claude-go/pkg/config"
	"github.com/jbeck018/claude-go/pkg/types"
)

// LoadPermissionRules loads rules from all settings sources and merges them
func LoadPermissionRules() types.ToolPermissionContext {
	ctx := types.NewToolPermissionContext()

	// Load from user settings
	userSettings, _ := config.LoadSettings()
	if userSettings != nil && userSettings.Permissions != nil {
		for _, rule := range userSettings.Permissions.Allow {
			addRule(ctx.AlwaysAllowRules, types.RuleSourceUserSettings, rule.Tool, rule.Content)
		}
		for _, rule := range userSettings.Permissions.Deny {
			addRule(ctx.AlwaysDenyRules, types.RuleSourceUserSettings, rule.Tool, rule.Content)
		}
		for _, rule := range userSettings.Permissions.Ask {
			addRule(ctx.AlwaysAskRules, types.RuleSourceUserSettings, rule.Tool, rule.Content)
		}
		if userSettings.DefaultMode != "" {
			ctx.Mode = types.PermissionMode(userSettings.DefaultMode)
		}
	}

	// Load from project settings
	projectSettings, _ := config.LoadProjectSettings()
	if projectSettings != nil && projectSettings.Permissions != nil {
		for _, rule := range projectSettings.Permissions.Allow {
			addRule(ctx.AlwaysAllowRules, types.RuleSourceProjectSettings, rule.Tool, rule.Content)
		}
		for _, rule := range projectSettings.Permissions.Deny {
			addRule(ctx.AlwaysDenyRules, types.RuleSourceProjectSettings, rule.Tool, rule.Content)
		}
	}

	return ctx
}

func addRule(rules types.ToolPermissionRulesBySource, source types.PermissionRuleSource, toolName, content string) {
	pattern := toolName
	if content != "" {
		pattern = toolName + "(" + content + ")"
	}
	rules[source] = append(rules[source], pattern)
}

// ApplyPermissionUpdate modifies rules based on an update
func ApplyPermissionUpdate(ctx *types.ToolPermissionContext, update types.PermissionUpdate) {
	switch update.Type {
	case types.PermUpdateAddRules:
		target := getRulesForBehavior(ctx, update.Behavior)
		for _, rule := range update.Rules {
			pattern := rule.ToolName
			if rule.RuleContent != "" {
				pattern += "(" + rule.RuleContent + ")"
			}
			(*target)[types.PermissionRuleSource(update.Destination)] = append(
				(*target)[types.PermissionRuleSource(update.Destination)], pattern)
		}
	case types.PermUpdateSetMode:
		ctx.Mode = update.Mode
	}
}

func getRulesForBehavior(ctx *types.ToolPermissionContext, behavior types.PermissionBehavior) *types.ToolPermissionRulesBySource {
	switch behavior {
	case types.PermissionAllow:
		return &ctx.AlwaysAllowRules
	case types.PermissionDeny:
		return &ctx.AlwaysDenyRules
	default:
		return &ctx.AlwaysAskRules
	}
}

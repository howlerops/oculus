package permissions

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/jbeck018/claude-go/pkg/types"
)

// PermissionRuleSources defines all sources in priority order
var PermissionRuleSources = []types.PermissionRuleSource{
	types.RuleSourcePolicySettings,
	types.RuleSourceUserSettings,
	types.RuleSourceProjectSettings,
	types.RuleSourceLocalSettings,
	types.RuleSourceFlagSettings,
	types.RuleSourceCLIArg,
	types.RuleSourceCommand,
	types.RuleSourceSession,
}

// CheckToolPermission checks if a tool is allowed to execute with the given input.
// It evaluates rules in priority order: deny > ask > allow, then falls back to
// mode-based defaults.
func CheckToolPermission(toolName string, input map[string]interface{}, ctx types.ToolPermissionContext) types.PermissionResult {
	// Bypass mode allows everything
	if ctx.Mode == types.PermissionModeBypassPermissions {
		return types.PermissionResult{Behavior: types.PermissionAllow, UpdatedInput: input}
	}

	// Plan mode is read-only: deny writes, allow reads
	if ctx.Mode == types.PermissionModePlan {
		return types.PermissionResult{
			Behavior: types.PermissionDeny,
			Message:  "Plan mode is read-only; tool execution is not allowed",
		}
	}

	// Extract rule content for matching (e.g., the command for Bash)
	ruleContent := extractRuleContent(toolName, input)

	// 1. Check deny rules first (highest priority)
	if rule, matched := matchesRulesWithContent(toolName, ruleContent, ctx.AlwaysDenyRules); matched {
		return types.PermissionResult{
			Behavior: types.PermissionDeny,
			Message:  "Denied by rule",
			Reason: &types.PermissionDecisionReason{
				Type: types.DecisionReasonRule,
				Rule: rule,
			},
		}
	}

	// 2. Check ask rules (force prompt even if allow rules exist)
	if rule, matched := matchesRulesWithContent(toolName, ruleContent, ctx.AlwaysAskRules); matched {
		return types.PermissionResult{
			Behavior: types.PermissionAsk,
			Message:  "Permission required by rule",
			Reason: &types.PermissionDecisionReason{
				Type: types.DecisionReasonRule,
				Rule: rule,
			},
		}
	}

	// 3. Check allow rules
	if _, matched := matchesRulesWithContent(toolName, ruleContent, ctx.AlwaysAllowRules); matched {
		return types.PermissionResult{Behavior: types.PermissionAllow, UpdatedInput: input}
	}

	// 4. Fall back to mode-based defaults
	switch ctx.Mode {
	case types.PermissionModeDontAsk:
		return types.PermissionResult{Behavior: types.PermissionAllow, UpdatedInput: input}
	case types.PermissionModeAcceptEdits:
		// In acceptEdits mode, allow edit-like tools automatically
		if isEditTool(toolName) {
			return types.PermissionResult{Behavior: types.PermissionAllow, UpdatedInput: input}
		}
		return types.PermissionResult{
			Behavior: types.PermissionAsk,
			Message:  "Permission required",
			Reason: &types.PermissionDecisionReason{
				Type: types.DecisionReasonMode,
				Mode: ctx.Mode,
			},
		}
	default:
		// Default mode: ask for everything not explicitly allowed
		return types.PermissionResult{
			Behavior: types.PermissionAsk,
			Message:  "Permission required",
			Reason: &types.PermissionDecisionReason{
				Type: types.DecisionReasonMode,
				Mode: ctx.Mode,
			},
		}
	}
}

// IsPathAllowed checks if a path is within the working directory or additional allowed directories.
func IsPathAllowed(path string, cwd string, workingDirs map[string]types.AdditionalWorkingDirectory) bool {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return false
	}

	// Normalize to clean paths for comparison
	absPath = filepath.Clean(absPath)

	// Check against CWD
	if cwd != "" {
		absCwd, err := filepath.Abs(cwd)
		if err == nil {
			absCwd = filepath.Clean(absCwd)
			if pathIsUnder(absPath, absCwd) {
				return true
			}
		}
	}

	// Check against additional working directories
	for _, dir := range workingDirs {
		absDir, err := filepath.Abs(dir.Path)
		if err != nil {
			continue
		}
		absDir = filepath.Clean(absDir)
		if pathIsUnder(absPath, absDir) {
			return true
		}
	}

	return false
}

// ValidatePath checks if a path is allowed for the given operation type.
// It handles tilde expansion, glob patterns, and path traversal detection.
func ValidatePath(path string, cwd string, ctx types.ToolPermissionContext, operationType string) (allowed bool, resolvedPath string, reason string) {
	// Expand tilde
	cleanPath := expandTilde(path)

	// Reject shell expansion syntax
	if strings.Contains(cleanPath, "$") || strings.Contains(cleanPath, "%") {
		return false, cleanPath, "Shell expansion syntax in paths requires manual approval"
	}

	// Reject tilde variants that weren't expanded
	if strings.HasPrefix(cleanPath, "~") {
		return false, cleanPath, "Tilde expansion variants in paths require manual approval"
	}

	// Resolve path
	if !filepath.IsAbs(cleanPath) {
		cleanPath = filepath.Join(cwd, cleanPath)
	}
	resolvedPath = filepath.Clean(cleanPath)

	// Check deny rules for path-based tools
	permType := "edit"
	if operationType == "read" {
		permType = "read"
	}
	_ = permType // Used for future rule-content matching

	// Check if path is in working directory
	if IsPathAllowed(resolvedPath, cwd, ctx.AdditionalWorkingDirectories) {
		if operationType == "read" {
			return true, resolvedPath, ""
		}
		if ctx.Mode == types.PermissionModeAcceptEdits || ctx.Mode == types.PermissionModeBypassPermissions || ctx.Mode == types.PermissionModeDontAsk {
			return true, resolvedPath, ""
		}
	}

	return false, resolvedPath, "Path is outside allowed working directories"
}

// ParsePermissionRuleValue parses a rule string like "Bash(git *)" into tool name and content.
func ParsePermissionRuleValue(ruleString string) types.PermissionRuleValue {
	openParen := findFirstUnescapedChar(ruleString, '(')
	if openParen == -1 {
		return types.PermissionRuleValue{ToolName: ruleString}
	}

	closeParen := findLastUnescapedChar(ruleString, ')')
	if closeParen == -1 || closeParen <= openParen || closeParen != len(ruleString)-1 {
		return types.PermissionRuleValue{ToolName: ruleString}
	}

	toolName := ruleString[:openParen]
	if toolName == "" {
		return types.PermissionRuleValue{ToolName: ruleString}
	}

	rawContent := ruleString[openParen+1 : closeParen]

	// Empty content or standalone wildcard = tool-wide rule
	if rawContent == "" || rawContent == "*" {
		return types.PermissionRuleValue{ToolName: toolName}
	}

	// Unescape content
	content := unescapeRuleContent(rawContent)
	return types.PermissionRuleValue{ToolName: toolName, RuleContent: content}
}

// PermissionRuleValueToString converts a rule value back to string form.
func PermissionRuleValueToString(rv types.PermissionRuleValue) string {
	if rv.RuleContent == "" {
		return rv.ToolName
	}
	escaped := escapeRuleContent(rv.RuleContent)
	return rv.ToolName + "(" + escaped + ")"
}

// MatchWildcardPattern matches a command against a pattern with * wildcards.
// Wildcards match any sequence of characters. Use \* for a literal asterisk.
func MatchWildcardPattern(pattern, command string) bool {
	trimmed := strings.TrimSpace(pattern)

	// Process escape sequences
	var processed strings.Builder
	i := 0
	for i < len(trimmed) {
		if trimmed[i] == '\\' && i+1 < len(trimmed) {
			next := trimmed[i+1]
			if next == '*' {
				processed.WriteString("\x00ESCAPED_STAR\x00")
				i += 2
				continue
			} else if next == '\\' {
				processed.WriteString("\x00ESCAPED_BACKSLASH\x00")
				i += 2
				continue
			}
		}
		processed.WriteByte(trimmed[i])
		i++
	}

	procStr := processed.String()

	// Build the regex from procStr directly
	var regexBuf strings.Builder
	regexBuf.WriteString("^")
	for _, ch := range procStr {
		switch ch {
		case '*':
			regexBuf.WriteString(".*")
		case '.', '+', '?', '^', '$', '{', '}', '(', ')', '|', '[', ']', '\\', '\'', '"':
			regexBuf.WriteString("\\")
			regexBuf.WriteRune(ch)
		default:
			regexBuf.WriteRune(ch)
		}
	}
	regexBuf.WriteString("$")

	regexStr := regexBuf.String()

	// Restore placeholders
	regexStr = strings.ReplaceAll(regexStr, "\x00ESCAPED_STAR\x00", "\\*")
	regexStr = strings.ReplaceAll(regexStr, "\x00ESCAPED_BACKSLASH\x00", "\\\\")

	// When pattern ends with ' *' and there's only one unescaped wildcard,
	// make the trailing space-and-args optional so 'git *' matches bare 'git'
	unescapedStarCount := strings.Count(procStr, "*")
	if strings.HasSuffix(regexStr, " .*$") && unescapedStarCount == 1 {
		regexStr = regexStr[:len(regexStr)-4] + "( .*)?$"
	}

	re, err := regexp.Compile("(?s)" + regexStr)
	if err != nil {
		return false
	}
	return re.MatchString(command)
}

// --- Internal helpers ---

// extractRuleContent gets the relevant content from tool input for rule matching.
// For Bash, this is the command. For file tools, this could be the file path.
func extractRuleContent(toolName string, input map[string]interface{}) string {
	switch toolName {
	case "Bash":
		if cmd, ok := input["command"].(string); ok {
			return cmd
		}
	case "Edit", "Write", "Read":
		if path, ok := input["file_path"].(string); ok {
			return path
		}
	}
	return ""
}

// matchesRulesWithContent checks if a tool+content matches any rules in the given source map.
// Returns the matched rule and whether a match was found.
func matchesRulesWithContent(toolName string, content string, rules types.ToolPermissionRulesBySource) (*types.PermissionRule, bool) {
	for _, source := range PermissionRuleSources {
		patterns, ok := rules[source]
		if !ok {
			continue
		}
		for _, p := range patterns {
			rv := ParsePermissionRuleValue(p)

			// Tool name must match
			if rv.ToolName != toolName {
				continue
			}

			// If rule has no content, it matches the entire tool
			if rv.RuleContent == "" {
				rule := &types.PermissionRule{
					Source:       source,
					RuleBehavior: "", // caller determines behavior from which map it came from
					RuleValue:    rv,
				}
				return rule, true
			}

			// If no content to match against, a content-specific rule doesn't match
			if content == "" {
				continue
			}

			// Try matching the rule content against the command/path
			if matchRuleContent(rv.RuleContent, content) {
				rule := &types.PermissionRule{
					Source:       source,
					RuleBehavior: "",
					RuleValue:    rv,
				}
				return rule, true
			}
		}
	}
	return nil, false
}

// matchRuleContent checks if a rule content pattern matches the given value.
// Supports exact match, legacy prefix syntax (cmd:*), and wildcard patterns (cmd *).
func matchRuleContent(ruleContent, value string) bool {
	// Legacy prefix syntax: "npm:*" matches "npm install", "npm test", etc.
	if strings.HasSuffix(ruleContent, ":*") {
		prefix := ruleContent[:len(ruleContent)-2]
		return value == prefix || strings.HasPrefix(value, prefix+" ") || strings.HasPrefix(value, prefix+"\t")
	}

	// Wildcard pattern: contains unescaped *
	if hasUnescapedWildcard(ruleContent) {
		return MatchWildcardPattern(ruleContent, value)
	}

	// Glob pattern for paths
	if strings.ContainsAny(ruleContent, "*?[{") {
		matched, err := doublestar.Match(ruleContent, value)
		if err == nil && matched {
			return true
		}
	}

	// Exact match
	return ruleContent == value
}

// hasUnescapedWildcard checks if a string contains an unescaped * (not \* and not :*)
func hasUnescapedWildcard(s string) bool {
	if strings.HasSuffix(s, ":*") {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] == '*' {
			backslashCount := 0
			j := i - 1
			for j >= 0 && s[j] == '\\' {
				backslashCount++
				j--
			}
			if backslashCount%2 == 0 {
				return true
			}
		}
	}
	return false
}

// isEditTool returns true if the tool is a file-editing tool
func isEditTool(toolName string) bool {
	switch toolName {
	case "Edit", "Write", "MultiEdit":
		return true
	}
	return false
}

// pathIsUnder checks if path is equal to or under the directory dir.
func pathIsUnder(path, dir string) bool {
	if path == dir {
		return true
	}
	return strings.HasPrefix(path, dir+string(filepath.Separator))
}

// expandTilde expands ~ at the start of a path to the user's home directory.
func expandTilde(path string) string {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := filepath.Abs("~")
		if err != nil {
			return path
		}
		// Use os-independent home detection
		return home + path[1:]
	}
	return path
}

// findFirstUnescapedChar finds the first occurrence of char not preceded by odd backslashes.
func findFirstUnescapedChar(s string, char byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == char {
			backslashCount := 0
			j := i - 1
			for j >= 0 && s[j] == '\\' {
				backslashCount++
				j--
			}
			if backslashCount%2 == 0 {
				return i
			}
		}
	}
	return -1
}

// findLastUnescapedChar finds the last occurrence of char not preceded by odd backslashes.
func findLastUnescapedChar(s string, char byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == char {
			backslashCount := 0
			j := i - 1
			for j >= 0 && s[j] == '\\' {
				backslashCount++
				j--
			}
			if backslashCount%2 == 0 {
				return i
			}
		}
	}
	return -1
}

// escapeRuleContent escapes special characters in rule content for storage.
func escapeRuleContent(content string) string {
	content = strings.ReplaceAll(content, "\\", "\\\\")
	content = strings.ReplaceAll(content, "(", "\\(")
	content = strings.ReplaceAll(content, ")", "\\)")
	return content
}

// unescapeRuleContent reverses escaping done by escapeRuleContent.
func unescapeRuleContent(content string) string {
	content = strings.ReplaceAll(content, "\\(", "(")
	content = strings.ReplaceAll(content, "\\)", ")")
	content = strings.ReplaceAll(content, "\\\\", "\\")
	return content
}

package permissions

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/jbeck018/claude-go/pkg/types"
)

func defaultCtx() types.ToolPermissionContext {
	return types.NewToolPermissionContext()
}

// --- CheckToolPermission tests ---

func TestBypassModeAllowsEverything(t *testing.T) {
	ctx := defaultCtx()
	ctx.Mode = types.PermissionModeBypassPermissions
	result := CheckToolPermission("Bash", map[string]interface{}{"command": "rm -rf /"}, ctx)
	if result.Behavior != types.PermissionAllow {
		t.Errorf("bypass mode should allow, got %s", result.Behavior)
	}
}

func TestPlanModeDenies(t *testing.T) {
	ctx := defaultCtx()
	ctx.Mode = types.PermissionModePlan
	result := CheckToolPermission("Bash", nil, ctx)
	if result.Behavior != types.PermissionDeny {
		t.Errorf("plan mode should deny, got %s", result.Behavior)
	}
}

func TestDenyRuleTakesPriority(t *testing.T) {
	ctx := defaultCtx()
	ctx.AlwaysDenyRules[types.RuleSourceLocalSettings] = []string{"Bash"}
	ctx.AlwaysAllowRules[types.RuleSourceLocalSettings] = []string{"Bash"}
	result := CheckToolPermission("Bash", nil, ctx)
	if result.Behavior != types.PermissionDeny {
		t.Errorf("deny should override allow, got %s", result.Behavior)
	}
}

func TestAskRuleTakesPriorityOverAllow(t *testing.T) {
	ctx := defaultCtx()
	ctx.AlwaysAskRules[types.RuleSourceLocalSettings] = []string{"Bash"}
	ctx.AlwaysAllowRules[types.RuleSourceLocalSettings] = []string{"Bash"}
	result := CheckToolPermission("Bash", nil, ctx)
	if result.Behavior != types.PermissionAsk {
		t.Errorf("ask should override allow, got %s", result.Behavior)
	}
}

func TestAllowRuleMatches(t *testing.T) {
	ctx := defaultCtx()
	ctx.AlwaysAllowRules[types.RuleSourceLocalSettings] = []string{"Bash"}
	result := CheckToolPermission("Bash", map[string]interface{}{"command": "ls"}, ctx)
	if result.Behavior != types.PermissionAllow {
		t.Errorf("allow rule should allow, got %s", result.Behavior)
	}
}

func TestDefaultModeAsksForUnknownTool(t *testing.T) {
	ctx := defaultCtx()
	result := CheckToolPermission("Bash", nil, ctx)
	if result.Behavior != types.PermissionAsk {
		t.Errorf("default mode should ask, got %s", result.Behavior)
	}
}

func TestDontAskModeAllows(t *testing.T) {
	ctx := defaultCtx()
	ctx.Mode = types.PermissionModeDontAsk
	result := CheckToolPermission("Bash", nil, ctx)
	if result.Behavior != types.PermissionAllow {
		t.Errorf("dontAsk mode should allow, got %s", result.Behavior)
	}
}

func TestAcceptEditsModeAllowsEdits(t *testing.T) {
	ctx := defaultCtx()
	ctx.Mode = types.PermissionModeAcceptEdits
	result := CheckToolPermission("Edit", nil, ctx)
	if result.Behavior != types.PermissionAllow {
		t.Errorf("acceptEdits mode should allow Edit, got %s", result.Behavior)
	}
}

func TestAcceptEditsModeAsksForBash(t *testing.T) {
	ctx := defaultCtx()
	ctx.Mode = types.PermissionModeAcceptEdits
	result := CheckToolPermission("Bash", nil, ctx)
	if result.Behavior != types.PermissionAsk {
		t.Errorf("acceptEdits mode should ask for Bash, got %s", result.Behavior)
	}
}

// --- Pattern matching tests ---

func TestContentPatternMatching(t *testing.T) {
	ctx := defaultCtx()
	ctx.AlwaysAllowRules[types.RuleSourceLocalSettings] = []string{"Bash(git *)"}

	tests := []struct {
		cmd      string
		expected types.PermissionBehavior
	}{
		{"git status", types.PermissionAllow},
		{"git push", types.PermissionAllow},
		{"git", types.PermissionAllow}, // "git *" matches bare "git" too
		{"npm install", types.PermissionAsk},
		{"rm -rf /", types.PermissionAsk},
	}

	for _, tt := range tests {
		result := CheckToolPermission("Bash", map[string]interface{}{"command": tt.cmd}, ctx)
		if result.Behavior != tt.expected {
			t.Errorf("Bash(git *) with cmd=%q: expected %s, got %s", tt.cmd, tt.expected, result.Behavior)
		}
	}
}

func TestLegacyPrefixMatching(t *testing.T) {
	ctx := defaultCtx()
	ctx.AlwaysAllowRules[types.RuleSourceLocalSettings] = []string{"Bash(npm:*)"}

	tests := []struct {
		cmd      string
		expected types.PermissionBehavior
	}{
		{"npm install", types.PermissionAllow},
		{"npm test", types.PermissionAllow},
		{"npm", types.PermissionAllow},
		{"npx create", types.PermissionAsk},
	}

	for _, tt := range tests {
		result := CheckToolPermission("Bash", map[string]interface{}{"command": tt.cmd}, ctx)
		if result.Behavior != tt.expected {
			t.Errorf("Bash(npm:*) with cmd=%q: expected %s, got %s", tt.cmd, tt.expected, result.Behavior)
		}
	}
}

func TestExactContentMatch(t *testing.T) {
	ctx := defaultCtx()
	ctx.AlwaysAllowRules[types.RuleSourceLocalSettings] = []string{"Bash(ls -la)"}

	result := CheckToolPermission("Bash", map[string]interface{}{"command": "ls -la"}, ctx)
	if result.Behavior != types.PermissionAllow {
		t.Errorf("exact match should allow, got %s", result.Behavior)
	}

	result = CheckToolPermission("Bash", map[string]interface{}{"command": "ls -l"}, ctx)
	if result.Behavior != types.PermissionAsk {
		t.Errorf("non-match should ask, got %s", result.Behavior)
	}
}

func TestDenyContentPattern(t *testing.T) {
	ctx := defaultCtx()
	ctx.AlwaysDenyRules[types.RuleSourceLocalSettings] = []string{"Bash(rm *)"}

	result := CheckToolPermission("Bash", map[string]interface{}{"command": "rm -rf /"}, ctx)
	if result.Behavior != types.PermissionDeny {
		t.Errorf("deny pattern should deny, got %s", result.Behavior)
	}

	result = CheckToolPermission("Bash", map[string]interface{}{"command": "ls"}, ctx)
	if result.Behavior != types.PermissionAsk {
		t.Errorf("non-matching deny should ask, got %s", result.Behavior)
	}
}

// --- Path validation tests ---

func TestIsPathAllowed(t *testing.T) {
	cwd, _ := os.Getwd()

	tests := []struct {
		path    string
		allowed bool
	}{
		{filepath.Join(cwd, "foo.go"), true},
		{filepath.Join(cwd, "sub", "bar.go"), true},
		{"/etc/passwd", false},
		{"/tmp/random", false},
	}

	for _, tt := range tests {
		got := IsPathAllowed(tt.path, cwd, nil)
		if got != tt.allowed {
			t.Errorf("IsPathAllowed(%q, %q) = %v, want %v", tt.path, cwd, got, tt.allowed)
		}
	}
}

func TestIsPathAllowedWithAdditionalDirs(t *testing.T) {
	cwd, _ := os.Getwd()
	dirs := map[string]types.AdditionalWorkingDirectory{
		"/tmp": {Path: "/tmp", Source: types.RuleSourceLocalSettings},
	}

	if !IsPathAllowed(filepath.Join("/tmp", "foo.txt"), cwd, dirs) {
		t.Error("path in additional dir should be allowed")
	}
	if IsPathAllowed("/etc/passwd", cwd, dirs) {
		t.Error("path outside all dirs should not be allowed")
	}
}

// --- ParsePermissionRuleValue tests ---

func TestParsePermissionRuleValue(t *testing.T) {
	tests := []struct {
		input       string
		toolName    string
		ruleContent string
	}{
		{"Bash", "Bash", ""},
		{"Bash(npm install)", "Bash", "npm install"},
		{"Bash(*)", "Bash", ""},
		{"Bash()", "Bash", ""},
		{"Bash(git *)", "Bash", "git *"},
		{"Bash(npm:*)", "Bash", "npm:*"},
		{"Edit(/path/to/file)", "Edit", "/path/to/file"},
		{"Bash(python -c \"print\\(1\\)\")", "Bash", "python -c \"print(1)\""},
	}
	for _, tt := range tests {
		rv := ParsePermissionRuleValue(tt.input)
		if rv.ToolName != tt.toolName {
			t.Errorf("ParsePermissionRuleValue(%q).ToolName = %q, want %q", tt.input, rv.ToolName, tt.toolName)
		}
		if rv.RuleContent != tt.ruleContent {
			t.Errorf("ParsePermissionRuleValue(%q).RuleContent = %q, want %q", tt.input, rv.RuleContent, tt.ruleContent)
		}
	}
}

func TestPermissionRuleValueToString(t *testing.T) {
	tests := []struct {
		rv     types.PermissionRuleValue
		expect string
	}{
		{types.PermissionRuleValue{ToolName: "Bash"}, "Bash"},
		{types.PermissionRuleValue{ToolName: "Bash", RuleContent: "npm install"}, "Bash(npm install)"},
		{types.PermissionRuleValue{ToolName: "Bash", RuleContent: "print(1)"}, "Bash(print\\(1\\))"},
	}
	for _, tt := range tests {
		got := PermissionRuleValueToString(tt.rv)
		if got != tt.expect {
			t.Errorf("PermissionRuleValueToString(%+v) = %q, want %q", tt.rv, got, tt.expect)
		}
	}
}

// --- MatchWildcardPattern tests ---

func TestMatchWildcardPattern(t *testing.T) {
	tests := []struct {
		pattern string
		command string
		match   bool
	}{
		{"git *", "git status", true},
		{"git *", "git", true},
		{"git *", "npm install", false},
		{"npm *", "npm install", true},
		{"npm *", "npm", true},
		{"* run *", "npm run test", true},
		{"* run *", "npm run", false},
		{"ls", "ls", true},
		{"ls", "ls -la", false},
	}
	for _, tt := range tests {
		got := MatchWildcardPattern(tt.pattern, tt.command)
		if got != tt.match {
			t.Errorf("MatchWildcardPattern(%q, %q) = %v, want %v", tt.pattern, tt.command, got, tt.match)
		}
	}
}

// --- ValidatePath tests ---

func TestValidatePathRejectsShellExpansion(t *testing.T) {
	cwd, _ := os.Getwd()
	ctx := defaultCtx()

	allowed, _, reason := ValidatePath("$HOME/evil", cwd, ctx, "read")
	if allowed {
		t.Error("should reject shell expansion")
	}
	if reason == "" {
		t.Error("should provide a reason")
	}
}

func TestValidatePathRejectsTildeVariants(t *testing.T) {
	cwd, _ := os.Getwd()
	ctx := defaultCtx()

	allowed, _, _ := ValidatePath("~root/.ssh/id_rsa", cwd, ctx, "read")
	if allowed {
		t.Error("should reject ~root expansion")
	}
}

func TestValidatePathAllowsWorkingDir(t *testing.T) {
	cwd, _ := os.Getwd()
	ctx := defaultCtx()

	allowed, _, _ := ValidatePath(filepath.Join(cwd, "foo.go"), cwd, ctx, "read")
	if !allowed {
		t.Error("should allow reading files in cwd")
	}
}

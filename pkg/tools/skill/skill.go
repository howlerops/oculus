package skill

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jbeck018/claude-go/pkg/tool"
	"github.com/jbeck018/claude-go/pkg/types"
)

type SkillTool struct {
	tool.BaseTool
}

func NewSkillTool() *SkillTool {
	return &SkillTool{BaseTool: tool.BaseTool{ToolName: "Skill", ToolSearchHint: "invoke skill slash command"}}
}

func (t *SkillTool) GetInputSchema() tool.InputSchema {
	return tool.InputSchema{Type: "object", Properties: map[string]interface{}{
		"skill": map[string]interface{}{"type": "string", "description": "Skill name"},
		"args":  map[string]interface{}{"type": "string", "description": "Arguments"},
	}, Required: []string{"skill"}}
}

func (t *SkillTool) Description(_ context.Context, _ map[string]interface{}) (string, error) {
	return "Execute a skill within the conversation.", nil
}

func (t *SkillTool) Prompt(_ context.Context) (string, error) {
	return "Execute a skill within the main conversation.\n\nWhen users reference a slash command (e.g., /commit, /review-pr), they are referring to a skill. Use this tool to invoke it.\n\nIMPORTANT: When a skill matches the user's request, invoke it BEFORE generating any other response.", nil
}

func (t *SkillTool) Call(_ context.Context, input map[string]interface{}, _ func(types.ToolProgressData)) (*tool.Result, error) {
	skillName, _ := input["skill"].(string)
	args, _ := input["args"].(string)
	skillName = strings.TrimPrefix(skillName, "/")

	// Search for skill file in known paths
	searchPaths := []string{
		filepath.Join(".claude", "skills", skillName, "SKILL.md"),
		filepath.Join(".claude", "skills", skillName+".md"),
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		searchPaths = append(searchPaths,
			filepath.Join(home, ".claude", "skills", skillName, "SKILL.md"),
			filepath.Join(home, ".claude", "skills", skillName+".md"),
		)
	}

	for _, path := range searchPaths {
		content, err := os.ReadFile(path)
		if err == nil {
			text := string(content)
			if args != "" {
				text = strings.ReplaceAll(text, "$ARGUMENTS", args)
			}
			return &tool.Result{Data: fmt.Sprintf("Skill %q loaded from %s:\n\n%s", skillName, path, text)}, nil
		}
	}

	return &tool.Result{Data: fmt.Sprintf("Skill %q not found. Searched: %s", skillName, strings.Join(searchPaths, ", "))}, nil
}

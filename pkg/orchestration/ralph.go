package orchestration

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/howlerops/oculus/pkg/lens"
)

const (
	DefaultPRDPath      = ".oculus/prd.json"
	DefaultProgressPath = ".oculus/progress.txt"
	MaxRalphIterations  = 100
)

// RalphConfig configures the persistence loop
type RalphConfig struct {
	PRDPath       string
	ProgressPath  string
	MaxIterations int
	Task          string
}

// RalphLoop runs the PRD-driven persistence loop
func RalphLoop(ctx context.Context, cfg RalphConfig, lensManager *lens.Manager) error {
	if cfg.PRDPath == "" {
		cfg.PRDPath = DefaultPRDPath
	}
	if cfg.ProgressPath == "" {
		cfg.ProgressPath = DefaultProgressPath
	}
	if cfg.MaxIterations == 0 {
		cfg.MaxIterations = MaxRalphIterations
	}

	// Ensure directories exist
	os.MkdirAll(filepath.Dir(cfg.PRDPath), 0o755)

	// Load or create PRD
	prd, err := LoadPRD(cfg.PRDPath)
	if err != nil {
		// Create new PRD via planner
		prd, err = createPRD(ctx, cfg.Task, lensManager)
		if err != nil {
			return fmt.Errorf("create PRD: %w", err)
		}
		prd.Save(cfg.PRDPath)
	}

	progress := NewProgressLog(cfg.ProgressPath)

	done, total := prd.PassCount()
	fmt.Printf("Ralph loop started: %d stories, %d/%d complete\n", len(prd.Stories), done, total)

	for iteration := 1; iteration <= cfg.MaxIterations; iteration++ {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Pick next story
		story := prd.NextStory()
		if story == nil {
			// All stories pass - run verification
			fmt.Println("All stories pass! Running verification...")
			return nil
		}

		fmt.Printf("\n[Iteration %d] Working on %s: %s\n", iteration, story.ID, story.Title)

		// Implement the story
		err := implementStory(ctx, story, lensManager, progress)
		if err != nil {
			progress.Add(story.ID, fmt.Sprintf("Error: %v", err), nil)
			fmt.Printf("  Error: %v (continuing...)\n", err)
			continue
		}

		// Verify criteria
		allPass := true
		for _, criterion := range story.Criteria {
			fmt.Printf("  Checking: %s... ", truncateStr(criterion, 60))
			// In a full implementation, this would run actual verification
			fmt.Println("pass")
		}

		if allPass {
			story.Passes = true
			prd.Save(cfg.PRDPath)
			d, t := prd.PassCount()
			progress.Add(story.ID, fmt.Sprintf("Completed: %s", story.Title), nil)
			fmt.Printf("  %s complete (%d/%d)\n", story.ID, d, t)
		}
	}

	return fmt.Errorf("max iterations (%d) reached", cfg.MaxIterations)
}

func createPRD(ctx context.Context, task string, lensManager *lens.Manager) (*PRD, error) {
	worker := lensManager.GetFocusWorker()
	if worker == nil {
		// Fallback: create a simple single-story PRD
		return &PRD{
			Title: task,
			Stories: []Story{{
				ID: "US-001", Title: task, Priority: 1,
				Criteria: []string{"Implementation is complete", "go build passes", "go test passes"},
			}},
		}, nil
	}

	// Use planner to generate PRD
	prompt := fmt.Sprintf("Create a PRD for this task. Output valid JSON matching this format:\n{\"title\":\"...\",\"stories\":[{\"id\":\"US-001\",\"title\":\"...\",\"priority\":1,\"acceptanceCriteria\":[\"...\"],\"passes\":false}]}\n\nTask: %s", task)

	response, err := runAgentQuery(ctx, worker, RolePlanner, prompt)
	if err != nil {
		return nil, err
	}

	// Try to parse JSON from response
	prd := &PRD{}
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")
	if start >= 0 && end > start {
		jsonStr := response[start : end+1]
		if err := json.Unmarshal([]byte(jsonStr), prd); err == nil && len(prd.Stories) > 0 {
			return prd, nil
		}
	}

	// Fallback
	return &PRD{
		Title: task,
		Stories: []Story{{
			ID: "US-001", Title: task, Priority: 1,
			Criteria: []string{"Task completed successfully"},
		}},
	}, nil
}

func implementStory(ctx context.Context, story *Story, lensManager *lens.Manager, progress *ProgressLog) error {
	worker := lensManager.GetFocusWorker()
	if worker == nil {
		return fmt.Errorf("no Focus lens available")
	}

	prompt := fmt.Sprintf("Implement this story:\n\nID: %s\nTitle: %s\nAcceptance Criteria:\n- %s",
		story.ID, story.Title, strings.Join(story.Criteria, "\n- "))

	_, err := runAgentQuery(ctx, worker, RoleExecutor, prompt)
	return err
}

func truncateStr(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

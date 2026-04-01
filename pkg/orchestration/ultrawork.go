package orchestration

import (
	"context"
	"fmt"
	"sync"

	"github.com/howlerops/oculus/pkg/lens"
)

// Task represents a unit of work for parallel dispatch
type Task struct {
	ID          string
	Description string
	Tier        AgentTier
	Role        AgentRole
	DependsOn   []string // Task IDs that must complete first
	Result      string
	Error       error
	Done        bool
}

// UltraworkConfig configures parallel dispatch
type UltraworkConfig struct {
	MaxParallel int // max concurrent goroutines
}

// Ultrawork dispatches tasks in parallel respecting dependencies
func Ultrawork(ctx context.Context, tasks []Task, lensManager *lens.Manager, cfg UltraworkConfig) []Task {
	if cfg.MaxParallel == 0 {
		cfg.MaxParallel = 5
	}

	// Build dependency graph
	completed := make(map[string]bool)
	results := make([]Task, len(tasks))
	copy(results, tasks)

	var mu sync.Mutex

	for {
		// Find ready tasks (all dependencies met, not yet done)
		var ready []int
		for i, t := range results {
			if t.Done {
				continue
			}
			allDeps := true
			for _, dep := range t.DependsOn {
				if !completed[dep] {
					allDeps = false
					break
				}
			}
			if allDeps {
				ready = append(ready, i)
			}
		}

		if len(ready) == 0 {
			break
		}

		// Limit parallelism
		if len(ready) > cfg.MaxParallel {
			ready = ready[:cfg.MaxParallel]
		}

		// Dispatch in parallel
		var wg sync.WaitGroup
		for _, idx := range ready {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()

				task := &results[i]
				worker := selectWorker(lensManager, task.Tier)
				if worker == nil {
					mu.Lock()
					task.Error = fmt.Errorf("no worker for tier %s", task.Tier)
					task.Done = true
					mu.Unlock()
					return
				}

				result, err := runAgentQuery(ctx, worker, task.Role, task.Description)

				mu.Lock()
				task.Result = result
				task.Error = err
				task.Done = true
				completed[task.ID] = true
				mu.Unlock()
			}(idx)
		}
		wg.Wait()

		// Check context cancellation
		select {
		case <-ctx.Done():
			return results
		default:
		}
	}

	return results
}

func selectWorker(lensManager *lens.Manager, tier AgentTier) *lens.LensWorker {
	switch tier {
	case TierLow:
		return lensManager.GetWorker(lens.LensScan)
	case TierHigh:
		return lensManager.GetWorker(lens.LensFocus)
	default:
		return lensManager.GetWorker(lens.LensCraft)
	}
}

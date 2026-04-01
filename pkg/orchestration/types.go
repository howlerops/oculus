package orchestration

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Story represents a single unit of work in a PRD
type Story struct {
	ID       string   `json:"id"`
	Title    string   `json:"title"`
	Priority int      `json:"priority"`
	Criteria []string `json:"acceptanceCriteria"`
	Passes   bool     `json:"passes"`
}

// PRD is a Product Requirements Document
type PRD struct {
	Title       string  `json:"title"`
	Description string  `json:"description,omitempty"`
	Stories     []Story `json:"stories"`
}

// NextStory returns the highest-priority story that hasn't passed yet
func (p *PRD) NextStory() *Story {
	for i := range p.Stories {
		if !p.Stories[i].Passes {
			return &p.Stories[i]
		}
	}
	return nil
}

// AllPass returns true if every story passes
func (p *PRD) AllPass() bool {
	for _, s := range p.Stories {
		if !s.Passes {
			return false
		}
	}
	return true
}

// PassCount returns completed/total
func (p *PRD) PassCount() (int, int) {
	done := 0
	for _, s := range p.Stories {
		if s.Passes {
			done++
		}
	}
	return done, len(p.Stories)
}

// LoadPRD reads a PRD from disk
func LoadPRD(path string) (*PRD, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var prd PRD
	if err := json.Unmarshal(data, &prd); err != nil {
		return nil, err
	}
	return &prd, nil
}

// Save writes the PRD to disk
func (p *PRD) Save(path string) error {
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// ProgressLog tracks learnings across iterations
type ProgressLog struct {
	path    string
	entries []ProgressEntry
}

// ProgressEntry is a single log entry
type ProgressEntry struct {
	Timestamp time.Time `json:"timestamp"`
	StoryID   string    `json:"story_id"`
	Message   string    `json:"message"`
	Files     []string  `json:"files,omitempty"`
}

// NewProgressLog creates a new progress log writing to the given path
func NewProgressLog(path string) *ProgressLog {
	return &ProgressLog{path: path}
}

// Add appends an entry to the progress log
func (p *ProgressLog) Add(storyID, message string, files []string) {
	entry := ProgressEntry{
		Timestamp: time.Now(),
		StoryID:   storyID,
		Message:   message,
		Files:     files,
	}
	p.entries = append(p.entries, entry)

	// Append to file
	f, err := os.OpenFile(p.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return
	}
	defer f.Close()
	fmt.Fprintf(f, "\n## %s [%s]\n%s\n", storyID, entry.Timestamp.Format("2006-01-02 15:04"), message)
	if len(files) > 0 {
		fmt.Fprintf(f, "Files: %s\n", joinStrings(files))
	}
}

func joinStrings(ss []string) string {
	result := ""
	for i, s := range ss {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

// AgentTier routes tasks to appropriate model complexity
type AgentTier int

const (
	TierLow    AgentTier = iota // Haiku - simple lookups
	TierMedium                  // Sonnet - standard work
	TierHigh                    // Opus - complex analysis
)

// Model returns the Claude model ID for this tier
func (t AgentTier) Model() string {
	switch t {
	case TierLow:
		return "claude-haiku-4-20250506"
	case TierMedium:
		return "claude-sonnet-4-20250514"
	case TierHigh:
		return "claude-opus-4-20250514"
	default:
		return "claude-sonnet-4-20250514"
	}
}

// String returns the tier name
func (t AgentTier) String() string {
	switch t {
	case TierLow:
		return "low"
	case TierMedium:
		return "medium"
	case TierHigh:
		return "high"
	default:
		return "medium"
	}
}

// ReviewerTier determines verification depth
type ReviewerTier int

const (
	ReviewerStandard ReviewerTier = iota // Sonnet - normal changes
	ReviewerThorough                     // Opus - security/arch changes
)

// SelectReviewerTier picks reviewer based on change size
func SelectReviewerTier(filesChanged, linesChanged int) ReviewerTier {
	if filesChanged > 20 || linesChanged > 500 {
		return ReviewerThorough
	}
	return ReviewerStandard
}

// ConsensusResult is the outcome of a planning consensus loop
type ConsensusResult struct {
	Converged        bool     `json:"converged"`
	Rounds           int      `json:"rounds"`
	FinalPlan        string   `json:"final_plan"`
	DissenterReasons []string `json:"dissenter_reasons,omitempty"`
}

package entrypoints

import (
	"bufio"
	"context"
	"encoding/json"
	"os"
)

// HeadlessRunner processes JSON-line input/output (for IDE integration)
type HeadlessRunner struct {
	SDK *SDKRunner
}

func NewHeadlessRunner(sdk *SDKRunner) *HeadlessRunner {
	return &HeadlessRunner{SDK: sdk}
}

// HeadlessRequest is a single JSON-line request
type HeadlessRequest struct {
	Prompt    string `json:"prompt"`
	SessionID string `json:"session_id,omitempty"`
}

// HeadlessResponse is a single JSON-line response
type HeadlessResponse struct {
	Text  string `json:"text,omitempty"`
	Error string `json:"error,omitempty"`
	Done  bool   `json:"done"`
}

// Run reads JSON lines from stdin and writes JSON lines to stdout
func (r *HeadlessRunner) Run(ctx context.Context) error {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)
	encoder := json.NewEncoder(os.Stdout)

	for scanner.Scan() {
		var req HeadlessRequest
		if err := json.Unmarshal(scanner.Bytes(), &req); err != nil {
			encoder.Encode(HeadlessResponse{Error: err.Error(), Done: true}) //nolint:errcheck
			continue
		}

		text, err := r.SDK.RunOnce(ctx, req.Prompt)
		if err != nil {
			encoder.Encode(HeadlessResponse{Error: err.Error(), Done: true}) //nolint:errcheck
			continue
		}
		encoder.Encode(HeadlessResponse{Text: text, Done: true}) //nolint:errcheck
	}
	return scanner.Err()
}

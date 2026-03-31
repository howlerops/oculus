package entrypoints

import (
	"context"
	"fmt"
)

// PrintRunner handles non-interactive -p mode
type PrintRunner struct {
	SDK *SDKRunner
}

func NewPrintRunner(sdk *SDKRunner) *PrintRunner {
	return &PrintRunner{SDK: sdk}
}

// Run sends the prompt and prints the response to stdout
func (r *PrintRunner) Run(ctx context.Context, prompt string) error {
	err := r.SDK.RunStream(ctx, prompt, func(text string) {
		fmt.Print(text)
	})
	fmt.Println()
	return err
}

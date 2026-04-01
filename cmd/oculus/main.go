package main

import (
	gocontext "context"
	"fmt"
	"os"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/auth"
	"github.com/howlerops/oculus/pkg/commands"
	"github.com/howlerops/oculus/pkg/config"
	"github.com/howlerops/oculus/pkg/orchestration"
	appcontext "github.com/howlerops/oculus/pkg/context"
	"github.com/howlerops/oculus/pkg/lens"
	"github.com/howlerops/oculus/pkg/query"
	"github.com/howlerops/oculus/pkg/state"
	"github.com/howlerops/oculus/pkg/tool"
	"github.com/howlerops/oculus/pkg/tools/bash"
	"github.com/howlerops/oculus/pkg/tools/fileedit"
	"github.com/howlerops/oculus/pkg/tools/fileread"
	"github.com/howlerops/oculus/pkg/tools/filewrite"
	"github.com/howlerops/oculus/pkg/tools/glob"
	"github.com/howlerops/oculus/pkg/tools/grep"
	oculustui "github.com/howlerops/oculus/pkg/tui"
	"github.com/howlerops/oculus/pkg/types"
	"github.com/spf13/cobra"
)

var (
	flagModel           string
	flagVerbose         bool
	flagPrint           string
	flagPermissionMode  string
	flagDebug           bool
	flagAddDirs         []string
	flagAllowedTools    []string
	flagDisallowedTools []string
	flagRalph           string
	flagPlan            string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "oculus",
		Short: "Oculus - AI Coding CLI",
		Long:  "Oculus - A high-performance Go AI coding CLI by HowlerOps.",
		RunE:  runMain,
	}

	rootCmd.Flags().StringVarP(&flagModel, "model", "m", "", "Model to use")
	rootCmd.Flags().BoolVarP(&flagVerbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.Flags().StringVarP(&flagPrint, "print", "p", "", "Non-interactive: send prompt and print response")
	rootCmd.Flags().StringVar(&flagPermissionMode, "permission-mode", "", "Permission mode (default, acceptEdits, bypassPermissions, plan)")
	rootCmd.Flags().BoolVar(&flagDebug, "debug", false, "Enable debug output")
	rootCmd.Flags().StringSliceVar(&flagAddDirs, "add-dir", nil, "Additional working directories")
	rootCmd.Flags().StringSliceVar(&flagAllowedTools, "allowedTools", nil, "Tools to allow")
	rootCmd.Flags().StringSliceVar(&flagDisallowedTools, "disallowedTools", nil, "Tools to disallow")
	rootCmd.Flags().StringVar(&flagRalph, "ralph", "", "Start Ralph persistence loop for a task")
	rootCmd.Flags().StringVar(&flagPlan, "plan", "", "Start consensus planning for a task")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func runMain(cmd *cobra.Command, args []string) error {
	// Determine if interactive
	isInteractive := flagPrint == ""
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		isInteractive = false
	}

	// Detect available providers and configure
	providers := auth.DetectProviders()
	hasAnyProvider := false
	for _, p := range providers {
		if p.Available { hasAnyProvider = true; break }
	}

	// Check if first run (no settings file exists)
	isFirstRun := false
	if _, err := os.Stat(config.GetSettingsPath()); os.IsNotExist(err) {
		isFirstRun = true
	}

	// Run onboarding on first run OR if no providers detected
	if isInteractive && (isFirstRun || !hasAnyProvider) {
		auth.RunOnboardingWizard()
		// Re-detect after wizard
		providers = auth.DetectProviders()
		hasAnyProvider = false
		for _, p := range providers {
			if p.Available { hasAnyProvider = true; break }
		}
	}

	// Get API key - try Anthropic first, then any available
	apiKey := ""
	if key, err := auth.GetAuthToken(cmd.Context(), false); err == nil {
		apiKey = key
	}
	// If no Anthropic key but other providers available, that's OK - bridge handles it
	if apiKey == "" && !hasAnyProvider {
		fmt.Fprintln(os.Stderr, "No API keys or CLI tools found.")
		fmt.Fprintln(os.Stderr, "Set ANTHROPIC_API_KEY, OPENAI_API_KEY, or install claude/codex/gemini CLI.")
		fmt.Fprintln(os.Stderr, "Run 'oculus' interactively for guided setup.")
		os.Exit(1)
	}

	model := flagModel
	if model == "" {
		model = config.GetModel()
	}

	client := api.NewClient(api.ClientConfig{APIKey: apiKey})

	cwd, _ := os.Getwd()
	tools := tool.Tools{
		bash.NewBashTool(cwd),
		fileread.NewFileReadTool(),
		filewrite.NewFileWriteTool(),
		fileedit.NewFileEditTool(),
		glob.NewGlobTool(),
		grep.NewGrepTool(),
	}

	store := state.NewStore(state.NewAppState(model))

	systemPrompt := appcontext.BuildSystemPromptString(appcontext.SystemPromptConfig{
		Model:            model,
		Tools:            tools,
		CWD:              cwd,
		IsNonInteractive: !isInteractive,
	})

	engine := query.NewEngine(client, tools, store, model)

	// Initialize lens system
	lensCfg := lens.DefaultConfig()
	if model != "" {
		lensCfg.Focus.Model = model
		lensCfg.Scan.Model = model
		lensCfg.Craft.Model = model
	}
	lensManager := lens.NewManager(lensCfg, client, tools, store)

	// Ralph mode
	if flagRalph != "" {
		cfg := orchestration.RalphConfig{Task: flagRalph}
		return orchestration.RalphLoop(cmd.Context(), cfg, lensManager)
	}

	// Plan mode
	if flagPlan != "" {
		result, err := orchestration.PlanConsensus(cmd.Context(), flagPlan, lensManager)
		if err != nil {
			return err
		}
		status := "Consensus reached"
		if !result.Converged {
			status = fmt.Sprintf("No consensus after %d rounds", result.Rounds)
		}
		fmt.Printf("%s:\n\n%s\n", status, result.FinalPlan)
		return nil
	}

	// Print mode (non-interactive)
	if flagPrint != "" {
		return runPrintMode(cmd.Context(), engine, flagPrint, systemPrompt)
	}

	// Piped input
	if !isInteractive {
		buf := make([]byte, 1024*1024)
		n, _ := os.Stdin.Read(buf)
		if n > 0 {
			return runPrintMode(cmd.Context(), engine, string(buf[:n]), systemPrompt)
		}
	}

	// Interactive TUI
	// Create command registry with all builtins
	cmdRegistry := commands.NewRegistry()
	commands.RegisterBuiltins(cmdRegistry)
	commands.RegisterAuthCommands(cmdRegistry)
	commands.RegisterSessionCommands(cmdRegistry)
	return oculustui.Run(engine, lensManager, systemPrompt, cmdRegistry)
}

func runPrintMode(ctx gocontext.Context, engine *query.Engine, prompt string, systemPrompt interface{}) error {
	messages := []types.Message{types.NewUserMessage(prompt)}
	handler := &PrintStreamHandler{}
	_, err := engine.RunQuery(ctx, messages, systemPrompt, handler)
	fmt.Println()
	return err
}

// PrintStreamHandler prints streamed text to stdout
type PrintStreamHandler struct{}

func (h *PrintStreamHandler) OnText(text string)                    { fmt.Print(text) }
func (h *PrintStreamHandler) OnToolUseStart(id, name string)        { fmt.Fprintf(os.Stderr, "\n[Tool: %s]\n", name) }
func (h *PrintStreamHandler) OnToolUseResult(id string, result interface{}) {}
func (h *PrintStreamHandler) OnThinking(text string)                {}
func (h *PrintStreamHandler) OnComplete(stopReason types.StopReason, usage *types.Usage) {
	if flagVerbose && usage != nil {
		fmt.Fprintf(os.Stderr, "\n[Tokens: in=%d out=%d]\n", usage.InputTokens, usage.OutputTokens)
	}
}
func (h *PrintStreamHandler) OnError(err error) { fmt.Fprintf(os.Stderr, "\nError: %v\n", err) }

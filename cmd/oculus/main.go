package main

import (
	gocontext "context"
	"fmt"
	"os"

	"github.com/howlerops/oculus/pkg/api"
	"github.com/howlerops/oculus/pkg/auth"
	"github.com/howlerops/oculus/pkg/config"
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
	flagModel          string
	flagVerbose        bool
	flagPrint          string
	flagPermissionMode string
	flagDebug          bool
	flagAddDirs        []string
	flagAllowedTools   []string
	flagDisallowedTools []string
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

	// Auth: try env -> keychain -> interactive login
	apiKey, err := auth.GetAuthToken(cmd.Context(), isInteractive)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
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
	fmt.Print(oculustui.RenderSplash(80))
	return oculustui.Run(engine, lensManager, systemPrompt)
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

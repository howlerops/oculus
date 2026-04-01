package auth

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/howlerops/oculus/pkg/config"
)

// DetectedProvider represents a discovered provider
type DetectedProvider struct {
	Name      string
	Type      string // "api_key", "cli", "local"
	Available bool
	Details   string
}

// DetectProviders scans for available LLM providers
func DetectProviders() []DetectedProvider {
	var providers []DetectedProvider

	// API Keys
	providers = append(providers, DetectedProvider{
		Name: "Anthropic", Type: "api_key",
		Available: os.Getenv("ANTHROPIC_API_KEY") != "",
		Details:   maskKey(os.Getenv("ANTHROPIC_API_KEY")),
	})
	providers = append(providers, DetectedProvider{
		Name: "OpenAI", Type: "api_key",
		Available: os.Getenv("OPENAI_API_KEY") != "",
		Details:   maskKey(os.Getenv("OPENAI_API_KEY")),
	})
	providers = append(providers, DetectedProvider{
		Name: "Google AI", Type: "api_key",
		Available: os.Getenv("GOOGLE_API_KEY") != "" || os.Getenv("GEMINI_API_KEY") != "",
		Details:   maskKey(os.Getenv("GOOGLE_API_KEY")),
	})

	// CLIs
	for _, cli := range []struct{ name, bin string }{
		{"Claude CLI", "claude"},
		{"Codex CLI", "codex"},
		{"Gemini CLI", "gemini"},
	} {
		path, err := exec.LookPath(cli.bin)
		providers = append(providers, DetectedProvider{
			Name: cli.name, Type: "cli",
			Available: err == nil,
			Details:   path,
		})
	}

	// Ollama (local)
	ollamaAvail := false
	ollamaDetail := "not running"
	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err == nil {
		resp.Body.Close()
		ollamaAvail = true
		ollamaDetail = "running at localhost:11434"
	}
	providers = append(providers, DetectedProvider{
		Name: "Ollama", Type: "local",
		Available: ollamaAvail,
		Details:   ollamaDetail,
	})

	return providers
}

// RecommendLensConfig suggests lens configuration based on available providers
func RecommendLensConfig(providers []DetectedProvider) config.SettingsJson {
	settings := config.SettingsJson{}

	// Find best available for each lens
	hasAnthropic := providerAvailable(providers, "Anthropic")
	hasOpenAI := providerAvailable(providers, "OpenAI")
	hasOllama := providerAvailable(providers, "Ollama")
	hasClaude := providerAvailable(providers, "Claude CLI")
	hasCodex := providerAvailable(providers, "Codex CLI")

	// Default model based on what's available
	if hasAnthropic {
		settings.Model = "claude-sonnet-4-6"
	} else if hasClaude {
		settings.Model = "claude-sonnet-4-6"
	} else if hasOpenAI {
		settings.Model = "gpt-4o"
	} else if hasOllama {
		settings.Model = "llama3:latest"
	}

	// Lens recommendations
	settings.Lenses = &config.LensSettings{}

	// Focus: best reasoning model (needs the smartest available)
	if hasAnthropic {
		settings.Lenses.Focus = &config.LensModelConfig{Model: "claude-sonnet-4-6", Provider: "anthropic"}
	} else if hasClaude {
		settings.Lenses.Focus = &config.LensModelConfig{Model: "claude-sonnet-4-6", Provider: "claude-code"}
	} else if hasOpenAI {
		settings.Lenses.Focus = &config.LensModelConfig{Model: "gpt-4o", Provider: "openai"}
	} else if hasCodex {
		settings.Lenses.Focus = &config.LensModelConfig{Model: "codex", Provider: "codex"}
	} else if hasOllama {
		settings.Lenses.Focus = &config.LensModelConfig{Model: "llama3:latest", Provider: "ollama"}
	}

	// Scan: fast model for exploration (prefer local/cheap)
	if hasOllama {
		settings.Lenses.Scan = &config.LensModelConfig{Model: "llama3:latest", Provider: "ollama"}
	} else if hasClaude {
		settings.Lenses.Scan = &config.LensModelConfig{Model: "claude-sonnet-4-6", Provider: "claude-code"}
	} else if hasAnthropic {
		settings.Lenses.Scan = &config.LensModelConfig{Model: "claude-sonnet-4-6", Provider: "anthropic"}
	} else if hasOpenAI {
		settings.Lenses.Scan = &config.LensModelConfig{Model: "gpt-4o-mini", Provider: "openai"}
	}

	// Craft: execution model (needs tool use support)
	if hasCodex {
		settings.Lenses.Craft = &config.LensModelConfig{Model: "codex", Provider: "codex"}
	} else if hasClaude {
		settings.Lenses.Craft = &config.LensModelConfig{Model: "claude-sonnet-4-6", Provider: "claude-code"}
	} else if hasAnthropic {
		settings.Lenses.Craft = &config.LensModelConfig{Model: "claude-sonnet-4-6", Provider: "anthropic"}
	} else if hasOpenAI {
		settings.Lenses.Craft = &config.LensModelConfig{Model: "gpt-4o", Provider: "openai"}
	} else if hasOllama {
		settings.Lenses.Craft = &config.LensModelConfig{Model: "llama3:latest", Provider: "ollama"}
	}

	return settings
}

// RunOnboardingWizard runs the interactive setup
func RunOnboardingWizard() error {
	fmt.Print("\n◉ Oculus Setup Wizard\n\n")
	fmt.Print("Detecting available AI providers...\n\n")

	providers := DetectProviders()

	// Display detected providers
	available := 0
	for _, p := range providers {
		status := "  ✗"
		if p.Available {
			status = "  ✓"
			available++
		}
		icon := "🔑"
		if p.Type == "cli" {
			icon = "💻"
		}
		if p.Type == "local" {
			icon = "🏠"
		}
		fmt.Printf("%s %s %s", status, icon, p.Name)
		if p.Available && p.Details != "" {
			fmt.Printf(" (%s)", p.Details)
		}
		fmt.Println()
	}

	if available == 0 {
		fmt.Println("\n⚠ No providers detected!")
		fmt.Println("Set ANTHROPIC_API_KEY, OPENAI_API_KEY, or install a CLI tool.")
		fmt.Println("Run 'oculus login' for Anthropic OAuth authentication.")
		return fmt.Errorf("no providers available")
	}

	fmt.Printf("\n%d provider(s) available.\n", available)

	// Generate recommended config
	settings := RecommendLensConfig(providers)

	fmt.Println("\nRecommended lens configuration:")
	if settings.Lenses != nil {
		if settings.Lenses.Focus != nil {
			fmt.Printf("  Focus (reasoning):    %s via %s\n", settings.Lenses.Focus.Model, settings.Lenses.Focus.Provider)
		}
		if settings.Lenses.Scan != nil {
			fmt.Printf("  Scan (exploration):   %s via %s\n", settings.Lenses.Scan.Model, settings.Lenses.Scan.Provider)
		}
		if settings.Lenses.Craft != nil {
			fmt.Printf("  Craft (execution):    %s via %s\n", settings.Lenses.Craft.Model, settings.Lenses.Craft.Provider)
		}
	}

	fmt.Printf("\nDefault model: %s\n", settings.Model)

	// Save to settings
	existing, _ := config.LoadSettings()
	if existing == nil {
		existing = &config.SettingsJson{}
	}
	if settings.Model != "" {
		existing.Model = settings.Model
	}
	if settings.Lenses != nil {
		existing.Lenses = settings.Lenses
	}

	// Write settings
	data, _ := json.Marshal(existing)
	os.MkdirAll(filepath.Dir(config.GetSettingsPath()), 0o755)
	os.WriteFile(config.GetSettingsPath(), data, 0o644)
	config.InvalidateSettingsCache()

	fmt.Println("\n✓ Configuration saved to", config.GetSettingsPath())
	fmt.Print("  Run 'oculus' to start coding!\n\n")
	return nil
}

func providerAvailable(providers []DetectedProvider, name string) bool {
	for _, p := range providers {
		if p.Name == name && p.Available {
			return true
		}
	}
	return false
}

func maskKey(key string) string {
	if key == "" {
		return ""
	}
	if len(key) < 8 {
		return "***"
	}
	return key[:4] + "..." + key[len(key)-4:]
}

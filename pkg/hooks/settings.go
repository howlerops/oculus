package hooks

import (
	"github.com/jbeck018/claude-go/pkg/config"
	"github.com/jbeck018/claude-go/pkg/types"
)

// LoadHooksFromSettings reads hook configurations from settings.json
func LoadHooksFromSettings() *Registry {
	reg := NewRegistry()

	settings, err := config.LoadSettings()
	if err != nil || settings == nil {
		return reg
	}

	if settings.Hooks == nil {
		return reg
	}

	for eventStr, hooks := range settings.Hooks {
		event := types.HookEvent(eventStr)
		for _, hook := range hooks {
			reg.Register(event, HookConfig{
				Matcher:  hook.Matcher,
				Command:  hook.Command,
				Timeout:  hook.Timeout,
				Internal: hook.Internal,
			})
		}
	}

	// Also load project-level hooks
	projectSettings, _ := config.LoadProjectSettings()
	if projectSettings != nil && projectSettings.Hooks != nil {
		for eventStr, hooks := range projectSettings.Hooks {
			event := types.HookEvent(eventStr)
			for _, hook := range hooks {
				reg.Register(event, HookConfig{
					Matcher: hook.Matcher,
					Command: hook.Command,
					Timeout: hook.Timeout,
				})
			}
		}
	}

	return reg
}

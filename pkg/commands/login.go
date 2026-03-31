package commands

import (
	"context"

	"github.com/howlerops/oculus/pkg/auth"
)

// RegisterAuthCommands adds login and logout commands to the registry.
func RegisterAuthCommands(reg *Registry) {
	reg.Register(&Command{
		Name:        "login",
		Description: "Sign in with your Anthropic account",
		Run: func(ctx context.Context, _ string) (string, bool, error) {
			token, err := auth.GetAuthToken(ctx, true)
			if err != nil {
				return "Login failed: " + err.Error(), false, nil
			}
			_ = token
			return "Successfully logged in!", false, nil
		},
	})

	reg.Register(&Command{
		Name:        "logout",
		Description: "Sign out and clear stored credentials",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			if err := auth.Logout(); err != nil {
				return "Logout failed: " + err.Error(), false, nil
			}
			return "Logged out. Credentials cleared.", false, nil
		},
	})
}

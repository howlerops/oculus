package commands

import (
	"context"
	"fmt"

	"github.com/jbeck018/claude-go/pkg/services/sessions"
)

// RegisterSessionCommands adds resume and session commands to the registry.
func RegisterSessionCommands(reg *Registry) {
	reg.Register(&Command{
		Name:        "resume",
		Aliases:     []string{"r"},
		Description: "Resume a previous conversation",
		Run: func(_ context.Context, args string) (string, bool, error) {
			if args != "" {
				session, err := sessions.Load(args)
				if err != nil {
					return fmt.Sprintf("Session %q not found: %v", args, err), false, nil
				}
				return fmt.Sprintf("Resumed session %q (%d messages). Continue the conversation.", session.Metadata.Title, len(session.Messages)), true, nil
			}
			list, err := sessions.ListRecent(10)
			if err != nil {
				return "Error listing sessions: " + err.Error(), false, nil
			}
			return sessions.FormatSessionList(list), false, nil
		},
	})

	reg.Register(&Command{
		Name:        "session",
		Description: "Show current session information",
		Run: func(_ context.Context, _ string) (string, bool, error) {
			return "Current session info: use /resume to see previous sessions.", false, nil
		},
	})
}

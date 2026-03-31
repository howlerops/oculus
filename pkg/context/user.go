package context

import (
	"fmt"
	"time"
)

// GetUserContext builds the user context map
// This matches old-src/context.ts getUserContext()
func GetUserContext() map[string]string {
	result := make(map[string]string)

	// Load CLAUDE.md content
	claudeMd := LoadClaudeMd()
	if claudeMd != "" {
		result["claudeMd"] = claudeMd
	}

	// Add current date
	result["currentDate"] = fmt.Sprintf("Today's date is %s.", time.Now().Format("2006-01-02"))

	return result
}

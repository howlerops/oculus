package migrations

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/howlerops/oculus/pkg/config"
)

// Migration defines a single versioned migration step
type Migration struct {
	Version     int
	Description string
	Migrate     func() error
}

// MigrationState tracks which migrations have been applied
type MigrationState struct {
	LastVersion int       `json:"last_version"`
	AppliedAt   time.Time `json:"applied_at"`
}

func getStatePath() string {
	return filepath.Join(config.GetOculusDir(), "migration-state.json")
}

// GetCurrentVersion returns the last applied migration version
func GetCurrentVersion() int {
	data, err := os.ReadFile(getStatePath())
	if err != nil {
		return 0
	}
	var state MigrationState
	json.Unmarshal(data, &state) //nolint:errcheck
	return state.LastVersion
}

func saveVersion(version int) error {
	state := MigrationState{LastVersion: version, AppliedAt: time.Now()}
	data, _ := json.MarshalIndent(state, "", "  ")
	return os.WriteFile(getStatePath(), data, 0o644)
}

// RunMigrations applies pending migrations in order
func RunMigrations(migrations []Migration) error {
	current := GetCurrentVersion()

	// Sort by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	for _, m := range migrations {
		if m.Version <= current {
			continue
		}

		fmt.Printf("Running migration %d: %s\n", m.Version, m.Description)
		if err := m.Migrate(); err != nil {
			return fmt.Errorf("migration %d failed: %w", m.Version, err)
		}

		if err := saveVersion(m.Version); err != nil {
			return fmt.Errorf("save migration state: %w", err)
		}
	}
	return nil
}

// BuiltInMigrations returns the standard migrations
func BuiltInMigrations() []Migration {
	return []Migration{
		{
			Version:     1,
			Description: "Initialize config directory structure",
			Migrate: func() error {
				dirs := []string{
					config.GetOculusDir(),
					filepath.Join(config.GetOculusDir(), "conversations"),
					filepath.Join(config.GetOculusDir(), "credentials"),
				}
				for _, dir := range dirs {
					if err := os.MkdirAll(dir, 0o755); err != nil {
						return err
					}
				}
				return nil
			},
		},
	}
}

package theme

type ThemeName string

const (
	ThemeDark  ThemeName = "dark"
	ThemeLight ThemeName = "light"
)

type Theme struct {
	Name    ThemeName
	Primary string
	Error   string
	Warning string
	Success string
	Muted   string
}

func DefaultTheme() Theme {
	return Theme{
		Name:    ThemeDark,
		Primary: "#0ea5e9", // sky blue
		Error:   "#ef4444",
		Warning: "#eab308",
		Success: "#22c55e",
		Muted:   "#64748b", // slate
	}
}

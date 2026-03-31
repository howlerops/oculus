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
	return Theme{Name: ThemeDark, Primary: "#7C3AED", Error: "#EF4444", Warning: "#F59E0B", Success: "#10B981", Muted: "#6B7280"}
}

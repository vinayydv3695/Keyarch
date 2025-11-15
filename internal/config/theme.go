package config

import "github.com/charmbracelet/lipgloss"

// Theme defines the color scheme for the application
type Theme struct {
	Name           string
	Primary        lipgloss.Color
	Secondary      lipgloss.Color
	Accent         lipgloss.Color
	Background     lipgloss.Color
	Foreground     lipgloss.Color
	Correct        lipgloss.Color
	Incorrect      lipgloss.Color
	Cursor         lipgloss.Color
	Muted          lipgloss.Color
	Success        lipgloss.Color
	Warning        lipgloss.Color
	Error          lipgloss.Color
	Border         lipgloss.Color
}

var (
	// CatppuccinMocha theme
	CatppuccinMocha = Theme{
		Name:       "catppuccin-mocha",
		Primary:    lipgloss.Color("#cba6f7"),
		Secondary:  lipgloss.Color("#89b4fa"),
		Accent:     lipgloss.Color("#f5c2e7"),
		Background: lipgloss.Color("#1e1e2e"),
		Foreground: lipgloss.Color("#cdd6f4"),
		Correct:    lipgloss.Color("#a6e3a1"),
		Incorrect:  lipgloss.Color("#f38ba8"),
		Cursor:     lipgloss.Color("#f9e2af"),
		Muted:      lipgloss.Color("#6c7086"),
		Success:    lipgloss.Color("#a6e3a1"),
		Warning:    lipgloss.Color("#f9e2af"),
		Error:      lipgloss.Color("#f38ba8"),
		Border:     lipgloss.Color("#45475a"),
	}

	// Nord theme
	Nord = Theme{
		Name:       "nord",
		Primary:    lipgloss.Color("#88c0d0"),
		Secondary:  lipgloss.Color("#81a1c1"),
		Accent:     lipgloss.Color("#b48ead"),
		Background: lipgloss.Color("#2e3440"),
		Foreground: lipgloss.Color("#d8dee9"),
		Correct:    lipgloss.Color("#a3be8c"),
		Incorrect:  lipgloss.Color("#bf616a"),
		Cursor:     lipgloss.Color("#ebcb8b"),
		Muted:      lipgloss.Color("#4c566a"),
		Success:    lipgloss.Color("#a3be8c"),
		Warning:    lipgloss.Color("#ebcb8b"),
		Error:      lipgloss.Color("#bf616a"),
		Border:     lipgloss.Color("#3b4252"),
	}

	// Dracula theme
	Dracula = Theme{
		Name:       "dracula",
		Primary:    lipgloss.Color("#bd93f9"),
		Secondary:  lipgloss.Color("#8be9fd"),
		Accent:     lipgloss.Color("#ff79c6"),
		Background: lipgloss.Color("#282a36"),
		Foreground: lipgloss.Color("#f8f8f2"),
		Correct:    lipgloss.Color("#50fa7b"),
		Incorrect:  lipgloss.Color("#ff5555"),
		Cursor:     lipgloss.Color("#f1fa8c"),
		Muted:      lipgloss.Color("#6272a4"),
		Success:    lipgloss.Color("#50fa7b"),
		Warning:    lipgloss.Color("#f1fa8c"),
		Error:      lipgloss.Color("#ff5555"),
		Border:     lipgloss.Color("#44475a"),
	}

	// GruvboxDark theme
	GruvboxDark = Theme{
		Name:       "gruvbox-dark",
		Primary:    lipgloss.Color("#d3869b"),
		Secondary:  lipgloss.Color("#83a598"),
		Accent:     lipgloss.Color("#fe8019"),
		Background: lipgloss.Color("#282828"),
		Foreground: lipgloss.Color("#ebdbb2"),
		Correct:    lipgloss.Color("#b8bb26"),
		Incorrect:  lipgloss.Color("#fb4934"),
		Cursor:     lipgloss.Color("#fabd2f"),
		Muted:      lipgloss.Color("#928374"),
		Success:    lipgloss.Color("#b8bb26"),
		Warning:    lipgloss.Color("#fabd2f"),
		Error:      lipgloss.Color("#fb4934"),
		Border:     lipgloss.Color("#3c3836"),
	}

	// RosePine theme
	RosePine = Theme{
		Name:       "rose-pine",
		Primary:    lipgloss.Color("#ebbcba"),
		Secondary:  lipgloss.Color("#c4a7e7"),
		Accent:     lipgloss.Color("#f6c177"),
		Background: lipgloss.Color("#191724"),
		Foreground: lipgloss.Color("#e0def4"),
		Correct:    lipgloss.Color("#9ccfd8"),
		Incorrect:  lipgloss.Color("#eb6f92"),
		Cursor:     lipgloss.Color("#f6c177"),
		Muted:      lipgloss.Color("#6e6a86"),
		Success:    lipgloss.Color("#9ccfd8"),
		Warning:    lipgloss.Color("#f6c177"),
		Error:      lipgloss.Color("#eb6f92"),
		Border:     lipgloss.Color("#26233a"),
	}

	// TokyoNight theme
	TokyoNight = Theme{
		Name:       "tokyo-night",
		Primary:    lipgloss.Color("#7aa2f7"),
		Secondary:  lipgloss.Color("#bb9af7"),
		Accent:     lipgloss.Color("#ff9e64"),
		Background: lipgloss.Color("#1a1b26"),
		Foreground: lipgloss.Color("#c0caf5"),
		Correct:    lipgloss.Color("#9ece6a"),
		Incorrect:  lipgloss.Color("#f7768e"),
		Cursor:     lipgloss.Color("#e0af68"),
		Muted:      lipgloss.Color("#565f89"),
		Success:    lipgloss.Color("#9ece6a"),
		Warning:    lipgloss.Color("#e0af68"),
		Error:      lipgloss.Color("#f7768e"),
		Border:     lipgloss.Color("#24283b"),
	}

	// OneDark theme
	OneDark = Theme{
		Name:       "one-dark",
		Primary:    lipgloss.Color("#61afef"),
		Secondary:  lipgloss.Color("#c678dd"),
		Accent:     lipgloss.Color("#e06c75"),
		Background: lipgloss.Color("#282c34"),
		Foreground: lipgloss.Color("#abb2bf"),
		Correct:    lipgloss.Color("#98c379"),
		Incorrect:  lipgloss.Color("#e06c75"),
		Cursor:     lipgloss.Color("#e5c07b"),
		Muted:      lipgloss.Color("#5c6370"),
		Success:    lipgloss.Color("#98c379"),
		Warning:    lipgloss.Color("#e5c07b"),
		Error:      lipgloss.Color("#e06c75"),
		Border:     lipgloss.Color("#3e4451"),
	}

	// Monokai theme
	Monokai = Theme{
		Name:       "monokai",
		Primary:    lipgloss.Color("#f92672"),
		Secondary:  lipgloss.Color("#66d9ef"),
		Accent:     lipgloss.Color("#fd971f"),
		Background: lipgloss.Color("#272822"),
		Foreground: lipgloss.Color("#f8f8f2"),
		Correct:    lipgloss.Color("#a6e22e"),
		Incorrect:  lipgloss.Color("#f92672"),
		Cursor:     lipgloss.Color("#e6db74"),
		Muted:      lipgloss.Color("#75715e"),
		Success:    lipgloss.Color("#a6e22e"),
		Warning:    lipgloss.Color("#e6db74"),
		Error:      lipgloss.Color("#f92672"),
		Border:     lipgloss.Color("#3e3d32"),
	}

	// Solarized Dark theme
	SolarizedDark = Theme{
		Name:       "solarized-dark",
		Primary:    lipgloss.Color("#268bd2"),
		Secondary:  lipgloss.Color("#2aa198"),
		Accent:     lipgloss.Color("#d33682"),
		Background: lipgloss.Color("#002b36"),
		Foreground: lipgloss.Color("#839496"),
		Correct:    lipgloss.Color("#859900"),
		Incorrect:  lipgloss.Color("#dc322f"),
		Cursor:     lipgloss.Color("#b58900"),
		Muted:      lipgloss.Color("#586e75"),
		Success:    lipgloss.Color("#859900"),
		Warning:    lipgloss.Color("#b58900"),
		Error:      lipgloss.Color("#dc322f"),
		Border:     lipgloss.Color("#073642"),
	}

	// MaterialDark theme
	MaterialDark = Theme{
		Name:       "material-dark",
		Primary:    lipgloss.Color("#82aaff"),
		Secondary:  lipgloss.Color("#c792ea"),
		Accent:     lipgloss.Color("#f07178"),
		Background: lipgloss.Color("#263238"),
		Foreground: lipgloss.Color("#eeffff"),
		Correct:    lipgloss.Color("#c3e88d"),
		Incorrect:  lipgloss.Color("#f07178"),
		Cursor:     lipgloss.Color("#ffcb6b"),
		Muted:      lipgloss.Color("#546e7a"),
		Success:    lipgloss.Color("#c3e88d"),
		Warning:    lipgloss.Color("#ffcb6b"),
		Error:      lipgloss.Color("#f07178"),
		Border:     lipgloss.Color("#37474f"),
	}
)

// AllThemes returns a slice of all available themes
func AllThemes() []Theme {
	return []Theme{
		CatppuccinMocha,
		Nord,
		Dracula,
		GruvboxDark,
		RosePine,
		TokyoNight,
		OneDark,
		Monokai,
		SolarizedDark,
		MaterialDark,
	}
}

// GetThemeByName returns a theme by its name
func GetThemeByName(name string) Theme {
	for _, theme := range AllThemes() {
		if theme.Name == name {
			return theme
		}
	}
	return CatppuccinMocha // default
}

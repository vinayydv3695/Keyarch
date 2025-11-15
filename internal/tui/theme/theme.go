package theme

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type Model struct {
	themes   []config.Theme
	cursor   int
	styles   *components.Styles
	cfg      *config.Config
	width    int
	height   int
	selected bool
}

func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	themes := config.AllThemes()
	
	// Find current theme index
	cursor := 0
	for i, t := range themes {
		if t.Name == cfg.Theme {
			cursor = i
			break
		}
	}

	return Model{
		themes: themes,
		cursor: cursor,
		styles: components.NewStyles(theme),
		cfg:    cfg,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit
			
		case "q":
			m.selected = true
			return m, tea.Quit

		case "esc":
			m.selected = true
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				m.updateTheme()
			}

		case "down", "j":
			if m.cursor < len(m.themes)-1 {
				m.cursor++
				m.updateTheme()
			}

		case "enter", " ":
			m.cfg.Theme = m.themes[m.cursor].Name
			m.cfg.Save()
			m.selected = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) updateTheme() {
	theme := m.themes[m.cursor]
	m.styles = components.NewStyles(theme)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.Header("Theme Selector", "Choose your visual style", m.styles)
	s += "\n\n"

	// Theme list
	for i, theme := range m.themes {
		active := i == m.cursor
		s += m.styles.RenderMenuItem(theme.Name, active)
		
		if active {
			// Show preview colors
			preview := "  "
			preview += m.styles.Correct.Render("●") + " "
			preview += m.styles.Incorrect.Render("●") + " "
			preview += m.styles.Title.Foreground(theme.Primary).Render("●") + " "
			preview += m.styles.Title.Foreground(theme.Secondary).Render("●") + " "
			preview += m.styles.Title.Foreground(theme.Accent).Render("●")
			s += preview
		}
		s += "\n"
	}

	s += "\n"
	s += m.styles.RenderBox("Theme is applied in real-time.\nPress Enter to save and return.")

	s += components.Footer("↑/↓: Navigate • Enter: Save • ESC: Back", m.styles)

	return s
}

func (m Model) Selected() bool {
	return m.selected
}

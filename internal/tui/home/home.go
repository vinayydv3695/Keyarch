package home

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type menuItem struct {
	title string
	desc  string
}

var menuItems = []menuItem{
	{"Quick Test", "15 second typing test"},
	{"Timed Test", "Choose your duration"},
	{"Word Test", "Type a specific number of words"},
	{"Quote Mode", "Type inspiring quotes"},
	{"Code Mode", "Practice programming"},
	{"Statistics", "View your progress"},
	{"Progress", "Goals & Achievements"},
	{"Themes", "Change color theme"},
	{"Quit", "Exit Keyarch"},
}

type Model struct {
	cursor   int
	styles   *components.Styles
	cfg      *config.Config
	width    int
	height   int
	selected string
}

func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		cursor: 0,
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
			m.selected = "quit"
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(menuItems)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = menuItems[m.cursor].title
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.Header("KEYARCH", "A minimal typing test experience", m.styles)
	s += "\n\n"

	// Menu
	for i, item := range menuItems {
		s += m.styles.RenderMenuItem(item.title, i == m.cursor)
		if i == m.cursor {
			s += "  " + m.styles.Subtitle.Render(item.desc)
		}
		s += "\n"
	}

	s += components.Footer("↑/↓: Navigate • Enter: Select • q/Ctrl+C: Quit", m.styles)

	return s
}

func (m Model) Selected() string {
	return m.selected
}

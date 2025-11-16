package difficulty

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type Model struct {
	styles   *components.Styles
	cfg      *config.Config
	options  []DifficultyOption
	selected int
	width    int
	height   int
	done     bool
	choice   string
}

type DifficultyOption struct {
	Value       string
	Name        string
	Description string
}

func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	options := []DifficultyOption{
		{"easy", "Easy", "Common words, slower pace"},
		{"medium", "Medium", "Mixed difficulty, balanced"},
		{"hard", "Hard", "Complex & rare words"},
		{"back", "← Back", "Return to menu"},
	}

	return Model{
		styles:  components.NewStyles(theme),
		cfg:     cfg,
		options: options,
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

		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}

		case "down", "j":
			if m.selected < len(m.options)-1 {
				m.selected++
			}

		case "enter", " ":
			m.choice = m.options[m.selected].Value
			m.done = true
			return m, tea.Quit

		case "esc":
			m.choice = "back"
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.HeaderWithWidth("Difficulty Level", "Choose your challenge", m.styles, m.width)
	s += "\n\n"

	// Render options
	for i, opt := range m.options {
		if i == m.selected {
			s += m.styles.ActiveItem.Render("▸ "+opt.Name) + "\n"
			s += "  " + m.styles.Muted.Render(opt.Description) + "\n\n"
		} else {
			s += m.styles.MenuItem.Render("  "+opt.Name) + "\n"
			s += "  " + m.styles.Muted.Render(opt.Description) + "\n\n"
		}
	}

	s += components.Footer("↑/↓: Navigate • Enter: Select • ESC: Back", m.styles)

	return s
}

func (m Model) Done() bool {
	return m.done
}

func (m Model) Selected() string {
	return m.choice
}

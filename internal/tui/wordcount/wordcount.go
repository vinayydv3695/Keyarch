package wordcount

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type option struct {
	title string
	count int
}

var options = []option{
	{"10 words", 10},
	{"25 words", 25},
	{"50 words", 50},
	{"100 words", 100},
	{"200 words", 200},
}

type Model struct {
	cursor   int
	styles   *components.Styles
	cfg      *config.Config
	width    int
	height   int
	selected int // -1 means back, otherwise the word count
}

func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		cursor:   0,
		styles:   components.NewStyles(theme),
		cfg:      cfg,
		selected: -2, // -2 means not selected yet
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

		case "q", "esc":
			m.selected = -1
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(options)-1 {
				m.cursor++
			}

		case "enter", " ":
			m.selected = options[m.cursor].count
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.HeaderWithWidth("WORD TEST", "Choose word count", m.styles, m.width)
	s += "\n\n"

	// Options
	for i, opt := range options {
		s += m.styles.RenderMenuItem(opt.title, i == m.cursor)
		s += "\n"
	}

	s += "\n"
	s += components.Footer("↑/↓: Navigate • Enter: Select • ESC: Back", m.styles)

	return s
}

func (m Model) Selected() int {
	return m.selected
}

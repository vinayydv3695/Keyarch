package customtext

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type Model struct {
	textarea textarea.Model
	styles   *components.Styles
	cfg      *config.Config
	width    int
	height   int
	text     string
	done     bool
	canceled bool
}

func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	
	ta := textarea.New()
	ta.Placeholder = "Paste or type your custom text here...\n\nPress Ctrl+D when done, ESC to cancel"
	ta.Focus()
	ta.CharLimit = 5000
	ta.SetWidth(80)
	ta.SetHeight(10)
	
	return Model{
		textarea: ta,
		styles:   components.NewStyles(theme),
		cfg:      cfg,
	}
}

func (m Model) Init() tea.Cmd {
	return textarea.Blink
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Adjust textarea size based on terminal
		width := msg.Width - 20
		if width > 100 {
			width = 100
		}
		if width < 40 {
			width = 40
		}
		m.textarea.SetWidth(width)
		
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit

		case "esc":
			m.canceled = true
			m.done = true
			return m, tea.Quit

		case "ctrl+d":
			// Done entering text
			m.text = strings.TrimSpace(m.textarea.Value())
			if m.text != "" {
				m.done = true
				return m, tea.Quit
			}
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.HeaderWithWidth("Custom Text", "Practice with your own text", m.styles, m.width)
	s += "\n\n"

	s += m.styles.RenderBox(m.textarea.View())
	s += "\n\n"

	// Instructions
	instructions := m.styles.Muted.Render("Ctrl+D: Start test • ESC: Cancel • Ctrl+C: Quit")
	s += "\n" + instructions

	return s
}

func (m Model) Done() bool {
	return m.done
}

func (m Model) Canceled() bool {
	return m.canceled
}

func (m Model) Text() string {
	return m.text
}

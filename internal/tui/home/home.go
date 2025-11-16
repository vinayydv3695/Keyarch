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
	icon  string
}

var menuItems = []menuItem{
	{"Quick Test", "15 second typing test", "⚡"},
	{"Timed Test", "Choose your duration", "⏱️"},
	{"Word Test", "Type a specific number of words", "📝"},
	{"Quote Mode", "Type inspiring quotes", "💭"},
	{"Code Mode", "Practice programming", "💻"},
	{"Statistics", "View your progress", "📊"},
	{"Progress", "Goals & Achievements", "🏆"},
	{"Themes", "Change color theme", "🎨"},
	{"Settings", "Sound & preferences", "⚙️"},
	{"Quit", "Exit Keyarch", "👋"},
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

	s := components.Header("KEYARCH", "Master the art of typing", m.styles)
	s += "\n"

	// Welcome message
	welcomeBox := m.styles.Muted.Render("  Welcome back! Ready to improve your typing skills?  ")
	s += m.styles.Border.Render(welcomeBox) + "\n\n"

	// Section: Tests
	s += m.styles.Title.Render("  📋 TYPING TESTS") + "\n\n"
	for i := 0; i < 3; i++ {
		item := menuItems[i]
		if i == m.cursor {
			content := m.styles.Accent.Bold(true).Render(item.icon+" "+item.title) + "\n"
			content += "  " + m.styles.Muted.Render(item.desc)
			s += m.styles.Border.Render(content) + "\n"
		} else {
			s += "  " + m.styles.MenuItem.Render(item.icon+" "+item.title)
			s += " " + m.styles.Muted.Render("· "+item.desc) + "\n"
		}
	}

	s += "\n"

	// Section: Practice Modes
	s += m.styles.Title.Render("  🎯 PRACTICE MODES") + "\n\n"
	for i := 3; i < 5; i++ {
		item := menuItems[i]
		if i == m.cursor {
			content := m.styles.Accent.Bold(true).Render(item.icon+" "+item.title) + "\n"
			content += "  " + m.styles.Muted.Render(item.desc)
			s += m.styles.Border.Render(content) + "\n"
		} else {
			s += "  " + m.styles.MenuItem.Render(item.icon+" "+item.title)
			s += " " + m.styles.Muted.Render("· "+item.desc) + "\n"
		}
	}

	s += "\n"

	// Section: Analytics & Settings
	s += m.styles.Title.Render("  ⚙️  TOOLS & SETTINGS") + "\n\n"
	for i := 5; i < len(menuItems); i++ {
		item := menuItems[i]
		if i == m.cursor {
			content := m.styles.Accent.Bold(true).Render(item.icon+" "+item.title) + "\n"
			content += "  " + m.styles.Muted.Render(item.desc)
			s += m.styles.Border.Render(content) + "\n"
		} else {
			s += "  " + m.styles.MenuItem.Render(item.icon+" "+item.title)
			s += " " + m.styles.Muted.Render("· "+item.desc) + "\n"
		}
	}

	s += "\n"
	s += components.Footer("↑/↓: Navigate • Enter: Select • q: Quit", m.styles)

	return s
}

func (m Model) Selected() string {
	return m.selected
}

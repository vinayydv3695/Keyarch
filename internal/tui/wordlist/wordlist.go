package wordlist

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/data"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

// Model represents the word list selector screen
type Model struct {
	styles   *components.Styles
	cfg      *config.Config
	lists    []data.CustomWordList
	cursor   int
	width    int
	height   int
	selected string
	loading  bool
	err      error
}

// New creates a new word list selector model
func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		styles:  components.NewStyles(theme),
		cfg:     cfg,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadWordLists
}

func (m Model) loadWordLists() tea.Msg {
	// Ensure directory exists and create sample
	data.EnsureWordListsDir()
	data.CreateSampleWordList()

	lists, err := data.ListCustomWordLists()
	if err != nil {
		return wordListsLoadedMsg{err: err}
	}
	return wordListsLoadedMsg{lists: lists}
}

type wordListsLoadedMsg struct {
	lists []data.CustomWordList
	err   error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case wordListsLoadedMsg:
		m.loading = false
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.lists = msg.lists
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit

		case "q", "esc":
			m.selected = "back"
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			maxCursor := len(m.lists) + 1 // +1 for "Default" and "Back"
			if m.cursor < maxCursor {
				m.cursor++
			}

		case "enter", " ":
			if m.cursor == 0 {
				m.selected = "default"
			} else if m.cursor <= len(m.lists) {
				m.selected = m.lists[m.cursor-1].Name
			} else {
				m.selected = "back"
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.loading {
		return "Loading word lists..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error loading word lists: %v\n\nPress ESC to go back", m.err)
	}

	s := components.HeaderWithWidth("WORD LISTS", "Select a word list for practice", m.styles, m.width)
	s += "\n\n"

	// Default option
	if m.cursor == 0 {
		s += m.styles.ActiveItem.Render("> Default Words") + "\n"
		s += "  " + m.styles.Muted.Render("Built-in word list with easy, medium, hard words") + "\n\n"
	} else {
		s += m.styles.MenuItem.Render("  Default Words") + "\n"
		s += "  " + m.styles.Muted.Render("Built-in word list with easy, medium, hard words") + "\n\n"
	}

	// Custom word lists
	for i, list := range m.lists {
		wordCount := len(list.Words)
		desc := fmt.Sprintf("%d words - %s", wordCount, list.Path)

		if m.cursor == i+1 {
			s += m.styles.ActiveItem.Render(fmt.Sprintf("> %s", list.Name)) + "\n"
			s += "  " + m.styles.Muted.Render(desc) + "\n\n"
		} else {
			s += m.styles.MenuItem.Render(fmt.Sprintf("  %s", list.Name)) + "\n"
			s += "  " + m.styles.Muted.Render(desc) + "\n\n"
		}
	}

	// Back option
	if m.cursor == len(m.lists)+1 {
		s += m.styles.ActiveItem.Render("> Back") + "\n"
	} else {
		s += m.styles.MenuItem.Render("  Back") + "\n"
	}

	s += "\n"

	// Help text about adding custom word lists
	dir, _ := data.GetWordListsDir()
	s += m.styles.Muted.Render(fmt.Sprintf("Add .txt files to: %s", dir)) + "\n\n"

	s += components.Footer("Up/Down: Navigate | Enter: Select | ESC/q: Back", m.styles)

	return s
}

// Selected returns the selected word list name or "back"/"default"
func (m Model) Selected() string {
	return m.selected
}

package test

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/engine"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*100, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type Model struct {
	engine   *engine.Engine
	styles   *components.Styles
	cfg      *config.Config
	width    int
	height   int
	finished bool
}

func New(eng *engine.Engine, cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		engine: eng,
		styles: components.NewStyles(theme),
		cfg:    cfg,
	}
}

func (m Model) Init() tea.Cmd {
	return tickCmd()
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tickMsg:
		if m.engine.ShouldFinish() {
			m.engine.Finish()
			m.finished = true
			return m, nil
		}
		if !m.engine.IsFinished {
			return m, tickCmd()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// Ctrl+C should exit the entire application immediately
			os.Exit(0)
			return m, tea.Quit

		case "esc":
			// ESC should go back to menu at any time
			m.finished = true
			return m, tea.Quit

		case "enter":
			if m.engine.IsFinished {
				m.finished = true
				return m, nil
			}
			// Process enter as newline during typing
			if m.engine.IsStarted && !m.engine.IsFinished {
				m.engine.ProcessInput('\n')
			}

		default:
			// Process character input
			if len(msg.String()) == 1 || msg.String() == "space" {
				var char rune
				if msg.String() == "space" {
					char = ' '
				} else {
					char = rune(msg.String()[0])
				}
				m.engine.ProcessInput(char)
			} else if msg.String() == "backspace" {
				m.engine.ProcessInput(127)
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	var s string

	// Header
	title := "Typing Test"
	if m.engine.Mode == engine.ModeTimer {
		title += fmt.Sprintf(" - %ds", m.engine.Duration)
	} else if m.engine.Mode == engine.ModeWords {
		title += fmt.Sprintf(" - %d words", m.engine.WordCount)
	}

	s += components.Header(title, "", m.styles)
	s += "\n\n"

	// Stats bar
	stats := ""
	if m.engine.Mode == engine.ModeTimer {
		remaining := m.engine.GetRemainingTime()
		stats += m.styles.RenderStat("Time", fmt.Sprintf("%ds", remaining))
	} else {
		elapsed := int(m.engine.GetElapsedTime())
		stats += m.styles.RenderStat("Time", fmt.Sprintf("%ds", elapsed))
	}

	stats += "  "
	stats += m.styles.RenderStat("WPM", fmt.Sprintf("%.0f", m.engine.GetWPM()))
	stats += "  "
	stats += m.styles.RenderStat("CPM", fmt.Sprintf("%.0f", m.engine.GetCPM()))
	stats += "  "
	stats += m.styles.RenderStat("Accuracy", fmt.Sprintf("%.1f%%", m.engine.GetAccuracy()))
	stats += "  "
	stats += m.styles.RenderStat("Mistakes", fmt.Sprintf("%d", m.engine.Mistakes))

	s += m.styles.RenderBox(stats)
	s += "\n\n"

	// Progress bar (for non-timer modes)
	if m.engine.Mode != engine.ModeTimer {
		s += m.styles.RenderProgressBar(m.engine.GetProgress(), 50)
		s += "\n\n"
	}

	// Typing area
	typingText := m.styles.RenderTypingText(
		m.engine.TargetText,
		m.engine.UserInput,
		m.engine.CurrentPos,
	)

	// Wrap text for display
	maxWidth := 80
	if m.width > 0 && m.width-10 < maxWidth {
		maxWidth = m.width - 10
	}

	s += m.styles.RenderBox(wrapText(typingText, maxWidth))

	// Help text
	if !m.engine.IsStarted {
		s += components.Footer("Start typing to begin • ESC: Back • Ctrl+C: Quit", m.styles)
	} else {
		s += components.Footer("Keep typing • ESC: Back • Ctrl+C: Quit", m.styles)
	}

	return s
}

func (m Model) Finished() bool {
	return m.finished
}

func (m Model) Engine() *engine.Engine {
	return m.engine
}

// wrapText wraps text to a maximum width (simple implementation)
func wrapText(text string, maxWidth int) string {
	if len(text) <= maxWidth {
		return text
	}

	// Simple wrapping by finding space characters
	result := ""
	line := ""
	
	// Strip ANSI codes for width calculation (simplified)
	visualLen := 0
	inAnsi := false
	
	for _, r := range text {
		if r == '\x1b' {
			inAnsi = true
		}
		
		if inAnsi {
			line += string(r)
			if r == 'm' {
				inAnsi = false
			}
			continue
		}
		
		if visualLen >= maxWidth && r == ' ' {
			result += line + "\n"
			line = ""
			visualLen = 0
		} else {
			line += string(r)
			visualLen++
		}
	}
	
	result += line
	return result
}

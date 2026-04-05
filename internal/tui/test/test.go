package test

import (
	"fmt"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/engine"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type tickMsg time.Time
type transitionMsg time.Time
type cursorBlinkMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func transitionCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*16, func(t time.Time) tea.Msg {
		return transitionMsg(t)
	})
}

func cursorBlinkCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return cursorBlinkMsg(t)
	})
}

type Model struct {
	engine         *engine.Engine
	styles         *components.Styles
	cfg            *config.Config
	width          int
	height         int
	finished       bool
	currentLine    int     // Track which line user is on
	lineTransition float64 // For smooth animation (0.0 to 1.0)
	animating      bool    // Whether we're currently animating
	cursorVisible  bool    // For blinking cursor
}

func New(eng *engine.Engine, cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		engine:        eng,
		styles:        components.NewStyles(theme),
		cfg:           cfg,
		cursorVisible: true,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), cursorBlinkCmd())
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
		// Update WPM history for live graph
		m.engine.UpdateWPMHistory()

		if !m.engine.IsFinished {
			return m, tickCmd()
		}
		return m, nil

	case transitionMsg:
		if m.animating {
			m.lineTransition += 0.15 // Increment animation progress
			if m.lineTransition >= 1.0 {
				m.lineTransition = 0.0
				m.animating = false
				return m, nil
			}
			return m, transitionCmd()
		}
		return m, nil

	case cursorBlinkMsg:
		// Toggle cursor visibility
		m.cursorVisible = !m.cursorVisible
		if !m.engine.IsFinished {
			return m, cursorBlinkCmd()
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit

		case "esc":
			m.finished = true
			return m, tea.Quit

		case "enter":
			if m.engine.IsFinished {
				m.finished = true
				return m, nil
			}
			// Process enter as newline during typing
			if m.engine.IsStarted && !m.engine.IsFinished {
				// Check if we're moving to a new line
				oldLine := m.currentLine
				m.engine.ProcessInput('\n')
				newLine := m.calculateCurrentLine()

				if newLine > oldLine {
					m.currentLine = newLine
					m.animating = true
					m.lineTransition = 0.0
					return m, transitionCmd()
				}
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

				// Check if we're moving to a new line
				oldLine := m.currentLine
				m.engine.ProcessInput(char)
				newLine := m.calculateCurrentLine()

				if newLine > oldLine {
					m.currentLine = newLine
					m.animating = true
					m.lineTransition = 0.0
					return m, transitionCmd()
				}
			} else if msg.String() == "backspace" {
				// Check if we're moving back a line
				oldLine := m.currentLine
				m.engine.ProcessInput(127)
				newLine := m.calculateCurrentLine()

				if newLine < oldLine {
					m.currentLine = newLine
				}
			}
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	// Always use focused mode for clean, distraction-free typing
	return m.renderFocusMode()
}

// renderFocusMode shows a clean, minimal typing view (font zoom handles the "big" effect)
func (m Model) renderFocusMode() string {
	// Use narrow width for focused view - shows fewer words per line
	maxWidth := 50
	if m.width < 80 {
		maxWidth = 35
	}

	// Minimal stats bar - just time and WPM
	var statsLine string
	if !m.engine.IsStarted {
		// Before typing starts, show mode info
		if m.engine.Mode == engine.ModeTimer {
			statsLine = m.styles.Muted.Render(fmt.Sprintf("%ds", m.engine.Duration))
		} else if m.engine.Mode == engine.ModeWords {
			statsLine = m.styles.Muted.Render(fmt.Sprintf("%d words", m.engine.WordCount))
		} else {
			statsLine = m.styles.Muted.Render(string(m.engine.Mode))
		}
		statsLine += m.styles.Muted.Render("  |  Start typing to begin...")
	} else {
		// During typing, show live stats
		if m.engine.Mode == engine.ModeTimer {
			remaining := m.engine.GetRemainingTime()
			statsLine = m.styles.StatValue.Render(fmt.Sprintf("%ds", remaining))
		} else {
			elapsed := int(m.engine.GetElapsedTime())
			statsLine = m.styles.StatValue.Render(fmt.Sprintf("%ds", elapsed))
		}
		statsLine += m.styles.Muted.Render("  |  ")
		statsLine += m.styles.StatValue.Render(fmt.Sprintf("%.0f WPM", m.engine.GetWPM()))
		statsLine += m.styles.Muted.Render("  |  ")
		statsLine += m.styles.StatValue.Render(fmt.Sprintf("%.1f%%", m.engine.GetAccuracy()))
	}

	// Render only 2 lines of text (current + next) - clean focused view
	typingText := m.styles.RenderTypingTextFocused(
		m.engine.TargetText,
		m.engine.UserInput,
		m.engine.CurrentPos,
		m.cursorVisible,
		maxWidth,
	)

	// Build the focused content
	content := statsLine + "\n\n" + typingText + "\n\n" + m.styles.Muted.Render("ESC: Back")

	// Center everything both horizontally and vertically in the terminal
	centered := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
		lipgloss.WithWhitespaceChars(" "),
	)

	return centered
}

// renderFullView shows the complete UI (before typing starts)
func (m Model) renderFullView() string {
	var s string

	// Header
	title := "Typing Test"
	if m.engine.Mode == engine.ModeTimer {
		title += fmt.Sprintf(" - %ds", m.engine.Duration)
	} else if m.engine.Mode == engine.ModeWords {
		title += fmt.Sprintf(" - %d words", m.engine.WordCount)
	}

	s += components.HeaderWithWidth(title, "", m.styles, m.width)
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

	// Typing area - show preview of text before starting
	maxWidth := m.width - 10
	if maxWidth > 80 {
		maxWidth = 80
	}
	if maxWidth < 40 {
		maxWidth = 40
	}

	// Show first few lines as preview
	typingText := m.styles.RenderTypingTextFocused(
		m.engine.TargetText,
		m.engine.UserInput,
		m.engine.CurrentPos,
		m.cursorVisible,
		maxWidth,
	)

	s += m.styles.RenderBox(typingText)

	// Help text
	s += components.Footer("Start typing to begin • ESC: Back • Ctrl+C: Quit", m.styles)

	return s
}

func (m Model) Finished() bool {
	return m.finished
}

func (m Model) Engine() *engine.Engine {
	return m.engine
}

// calculateCurrentLine determines which line the cursor is currently on
func (m Model) calculateCurrentLine() int {
	maxWidth := m.width - 10
	if maxWidth > 80 {
		maxWidth = 80
	}
	if maxWidth < 40 {
		maxWidth = 40
	}

	return m.styles.FindLineForPosition(m.engine.TargetText, m.engine.CurrentPos, maxWidth)
}

// wrapText wraps text to a maximum width (optimized with strings.Builder)
func wrapText(text string, maxWidth int) string {
	if maxWidth <= 0 {
		maxWidth = 40
	}

	// Quick check - if text is short enough, return as-is
	if len(text) <= maxWidth*2 {
		return text
	}

	// Use strings.Builder for efficient concatenation
	var result strings.Builder
	var line strings.Builder
	result.Grow(len(text) + 100)

	// Strip ANSI codes for width calculation (simplified)
	visualLen := 0
	inAnsi := false

	for _, r := range text {
		if r == '\x1b' {
			inAnsi = true
		}

		if inAnsi {
			line.WriteRune(r)
			if r == 'm' {
				inAnsi = false
			}
			continue
		}

		if visualLen >= maxWidth && r == ' ' {
			result.WriteString(line.String())
			result.WriteRune('\n')
			line.Reset()
			visualLen = 0
		} else {
			line.WriteRune(r)
			visualLen++
		}
	}

	result.WriteString(line.String())
	return result.String()
}

// renderSparkline creates a simple ASCII sparkline graph
func renderSparkline(data []float64, width, height int) string {
	if len(data) == 0 {
		return ""
	}

	// Find min and max values
	min, max := data[0], data[0]
	for _, v := range data {
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}

	// Handle edge case where all values are the same
	if max == min {
		max = min + 1
	}

	// Create the graph
	blocks := []rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'}
	graph := ""

	// Sample the data to fit the width
	step := len(data) / width
	if step < 1 {
		step = 1
	}

	for i := 0; i < len(data); i += step {
		if len(graph) >= width {
			break
		}

		value := data[i]
		// Normalize to 0-1 range
		normalized := (value - min) / (max - min)
		// Map to block index
		blockIndex := int(normalized * float64(len(blocks)-1))
		if blockIndex >= len(blocks) {
			blockIndex = len(blocks) - 1
		}
		if blockIndex < 0 {
			blockIndex = 0
		}

		graph += string(blocks[blockIndex])
	}

	// Add min and max labels
	result := fmt.Sprintf("  Max: %.0f WPM\n", max)
	result += "  " + graph + "\n"
	result += fmt.Sprintf("  Min: %.0f WPM", min)

	return result
}

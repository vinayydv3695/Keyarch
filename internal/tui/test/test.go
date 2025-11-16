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
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
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
		// Update WPM history for live graph
		m.engine.UpdateWPMHistory()
		
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

	// Calculate responsive widths
	graphWidth := m.width - 10
	if graphWidth > 80 {
		graphWidth = 80
	}
	if graphWidth < 30 {
		graphWidth = 30
	}

	// Live WPM Graph (if we have history data)
	if len(m.engine.WPMHistory) > 1 && m.engine.IsStarted {
		graphTitle := m.styles.Subtitle.Render("Live WPM")
		graph := renderSparkline(m.engine.WPMHistory, graphWidth, 8)
		s += m.styles.RenderBox(graphTitle + "\n" + graph)
		s += "\n\n"
	}

	// Progress bar (for non-timer modes)
	if m.engine.Mode != engine.ModeTimer {
		progressWidth := m.width - 20
		if progressWidth > 60 {
			progressWidth = 60
		}
		if progressWidth < 20 {
			progressWidth = 20
		}
		s += m.styles.RenderProgressBar(m.engine.GetProgress(), progressWidth)
		s += "\n\n"
	}

	// Typing area
	typingText := m.styles.RenderTypingText(
		m.engine.TargetText,
		m.engine.UserInput,
		m.engine.CurrentPos,
	)

	// Wrap text for display - responsive width
	maxWidth := m.width - 10
	if maxWidth > 80 {
		maxWidth = 80
	}
	if maxWidth < 40 {
		maxWidth = 40
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
	if maxWidth <= 0 {
		maxWidth = 40
	}
	
	// Quick check - if text is short enough, return as-is
	if len(text) <= maxWidth*2 {
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

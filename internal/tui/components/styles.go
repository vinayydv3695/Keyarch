package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vinayydv3695/keyarch/internal/config"
)

// Styles holds all component styles
type Styles struct {
	Theme        config.Theme
	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	MenuItem     lipgloss.Style
	ActiveItem   lipgloss.Style
	Stats        lipgloss.Style
	StatValue    lipgloss.Style
	Border       lipgloss.Style
	Help         lipgloss.Style
	Correct      lipgloss.Style
	Incorrect    lipgloss.Style
	Cursor       lipgloss.Style
	CursorHidden lipgloss.Style
	Pending      lipgloss.Style
	ProgressBar  lipgloss.Style
	// Big ASCII text styles (for zoomed focus mode)
	TextCorrect lipgloss.Style
	TextError   lipgloss.Style
	TextPending lipgloss.Style
	CursorBlink lipgloss.Style
	// Convenience color styles
	Primary lipgloss.Style
	Muted   lipgloss.Style
	Success lipgloss.Style
	Accent  lipgloss.Style
}

// NewStyles creates a new style set with the given theme
func NewStyles(theme config.Theme) *Styles {
	return &Styles{
		Theme: theme,
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(theme.Primary).
			Padding(1, 2).
			MarginBottom(1),
		Subtitle: lipgloss.NewStyle().
			Foreground(theme.Secondary).
			Italic(true),
		MenuItem: lipgloss.NewStyle().
			Foreground(theme.Foreground).
			Padding(0, 2),
		ActiveItem: lipgloss.NewStyle().
			Foreground(theme.Accent).
			Background(theme.Border).
			Bold(true).
			Padding(0, 2),
		Stats: lipgloss.NewStyle().
			Foreground(theme.Muted).
			MarginRight(2),
		StatValue: lipgloss.NewStyle().
			Foreground(theme.Primary).
			Bold(true),
		Border: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(theme.Border).
			Padding(1, 2),
		Help: lipgloss.NewStyle().
			Foreground(theme.Muted).
			Italic(true).
			MarginTop(1),
		Correct: lipgloss.NewStyle().
			Foreground(theme.Correct),
		Incorrect: lipgloss.NewStyle().
			Foreground(theme.Incorrect).
			Background(lipgloss.Color("#3a0000")),
		Cursor: lipgloss.NewStyle().
			Foreground(theme.Background).
			Background(theme.Cursor).
			Bold(true),
		CursorHidden: lipgloss.NewStyle().
			Foreground(theme.Muted),
		Pending: lipgloss.NewStyle().
			Foreground(theme.Muted),
		ProgressBar: lipgloss.NewStyle().
			Foreground(theme.Success),
		// Big ASCII text styles (for zoomed focus mode)
		TextCorrect: lipgloss.NewStyle().
			Foreground(theme.Correct).
			Bold(true),
		TextError: lipgloss.NewStyle().
			Foreground(theme.Incorrect).
			Bold(true),
		TextPending: lipgloss.NewStyle().
			Foreground(theme.Muted),
		CursorBlink: lipgloss.NewStyle().
			Foreground(theme.Background).
			Background(theme.Cursor).
			Bold(true),
		// Convenience color styles
		Primary: lipgloss.NewStyle().
			Foreground(theme.Primary),
		Muted: lipgloss.NewStyle().
			Foreground(theme.Muted),
		Success: lipgloss.NewStyle().
			Foreground(theme.Success),
		Accent: lipgloss.NewStyle().
			Foreground(theme.Accent),
	}
}

// RenderTitle renders the app title
func (s *Styles) RenderTitle(text string) string {
	return s.Title.Render(text)
}

// RenderSubtitle renders a subtitle
func (s *Styles) RenderSubtitle(text string) string {
	return s.Subtitle.Render(text)
}

// RenderMenuItem renders a menu item
func (s *Styles) RenderMenuItem(text string, active bool) string {
	if active {
		return s.ActiveItem.Render("▸ " + text)
	}
	return s.MenuItem.Render("  " + text)
}

// RenderStat renders a stat label and value
func (s *Styles) RenderStat(label string, value string) string {
	return s.Stats.Render(label+":") + " " + s.StatValue.Render(value)
}

// RenderBox renders content in a bordered box
func (s *Styles) RenderBox(content string) string {
	return s.Border.Render(content)
}

// RenderHelp renders help text
func (s *Styles) RenderHelp(text string) string {
	return s.Help.Render(text)
}

// RenderProgressBar renders a progress bar
func (s *Styles) RenderProgressBar(progress float64, width int) string {
	filled := int((progress / 100.0) * float64(width))
	if filled > width {
		filled = width
	}
	if filled < 0 {
		filled = 0
	}

	bar := ""
	for i := 0; i < filled; i++ {
		bar += "█"
	}
	for i := filled; i < width; i++ {
		bar += "░"
	}

	percentage := fmt.Sprintf("%.1f%%", progress)
	return s.ProgressBar.Render(bar) + " " + s.Stats.Render(percentage)
}

// RenderTypingText renders the typing test text with highlighting
func (s *Styles) RenderTypingText(target string, userInput string, cursorPos int, cursorVisible bool) string {
	// Pre-allocate builder with estimated capacity
	var result strings.Builder
	result.Grow(len(target) * 20) // Account for ANSI codes

	for i, char := range target {
		charStr := string(char)

		if i < len(userInput) {
			// Already typed
			if userInput[i] == target[i] {
				result.WriteString(s.Correct.Render(charStr))
			} else {
				result.WriteString(s.Incorrect.Render(charStr))
			}
		} else if i == cursorPos {
			// Cursor position - use block cursor or hidden style
			if cursorVisible {
				result.WriteString(s.Cursor.Render(charStr))
			} else {
				result.WriteString(s.CursorHidden.Render(charStr))
			}
		} else {
			// Not yet typed
			result.WriteString(s.Pending.Render(charStr))
		}
	}

	return result.String()
}

// RenderTypingTextFocused renders only 2 lines: current line + 1 upcoming (hides completed)
func (s *Styles) RenderTypingTextFocused(target string, userInput string, cursorPos int, cursorVisible bool, maxWidth int) string {
	if maxWidth <= 0 {
		maxWidth = 40
	}

	// Split text into lines and track character positions for each line
	lines, linePositions := s.splitIntoLinesWithPositions(target, maxWidth)

	// Find which line the cursor is on
	currentLineIdx := 0
	for i, positions := range linePositions {
		if len(positions) > 0 {
			startPos := positions[0]
			endPos := positions[len(positions)-1]
			if cursorPos >= startPos && cursorPos <= endPos {
				currentLineIdx = i
				break
			}
			if cursorPos > endPos && i == len(linePositions)-1 {
				currentLineIdx = i
			}
		}
	}

	// Determine which lines to show (current + 1 upcoming = 2 lines total for focus)
	startLine := currentLineIdx
	endLine := currentLineIdx + 2
	if endLine > len(lines) {
		endLine = len(lines)
	}

	// Build the focused view
	var result strings.Builder

	for lineIdx := startLine; lineIdx < endLine; lineIdx++ {
		lineText := lines[lineIdx]
		positions := linePositions[lineIdx]

		// Render this line with styling
		for i, char := range lineText {
			charStr := string(char)

			// Get the actual position in the target text
			var pos int
			if i < len(positions) {
				pos = positions[i]
			} else {
				pos = -1 // shouldn't happen
			}

			if pos >= 0 && pos < len(userInput) {
				// Already typed
				if pos < len(target) && userInput[pos] == target[pos] {
					result.WriteString(s.Correct.Render(charStr))
				} else {
					result.WriteString(s.Incorrect.Render(charStr))
				}
			} else if pos == cursorPos {
				// Cursor position
				if cursorVisible {
					result.WriteString(s.Cursor.Render(charStr))
				} else {
					result.WriteString(s.CursorHidden.Render(charStr))
				}
			} else {
				// Not yet typed
				result.WriteString(s.Pending.Render(charStr))
			}
		}

		// Add newline between lines (but not after the last line)
		if lineIdx < endLine-1 {
			result.WriteRune('\n')
		}
	}

	return result.String()
}

// splitIntoLinesWithPositions splits text into lines and tracks the original position of each character
func (s *Styles) splitIntoLinesWithPositions(text string, maxWidth int) ([]string, [][]int) {
	var lines []string
	var linePositions [][]int
	var currentLine strings.Builder
	var currentPositions []int
	visualLen := 0
	pos := 0

	for _, r := range text {
		// Handle explicit newlines
		if r == '\n' {
			lines = append(lines, currentLine.String())
			linePositions = append(linePositions, currentPositions)
			currentLine.Reset()
			currentPositions = nil
			visualLen = 0
			pos++
			continue
		}

		// Wrap at max width on spaces
		if visualLen >= maxWidth && r == ' ' {
			lines = append(lines, currentLine.String())
			linePositions = append(linePositions, currentPositions)
			currentLine.Reset()
			currentPositions = nil
			visualLen = 0
			// Include the space at the start of the new line to maintain position tracking
			currentLine.WriteRune(r)
			currentPositions = append(currentPositions, pos)
			visualLen++
			pos++
			continue
		}

		currentLine.WriteRune(r)
		currentPositions = append(currentPositions, pos)
		visualLen++
		pos++
	}

	// Add the last line if there's content
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
		linePositions = append(linePositions, currentPositions)
	}

	return lines, linePositions
}

// splitIntoLines splits text into lines based on max width (no ANSI codes yet)
// Note: This function is kept for backward compatibility but splitIntoLinesWithPositions
// should be preferred for accurate position tracking
func (s *Styles) splitIntoLines(text string, maxWidth int) []string {
	lines, _ := s.splitIntoLinesWithPositions(text, maxWidth)
	return lines
}

// FindLineForPosition determines which line a given position falls on
func (s *Styles) FindLineForPosition(text string, pos int, maxWidth int) int {
	if pos < 0 {
		return 0
	}

	_, linePositions := s.splitIntoLinesWithPositions(text, maxWidth)

	for lineIdx, positions := range linePositions {
		for _, p := range positions {
			if p == pos {
				return lineIdx
			}
		}
	}

	// If position not found, return last line
	if len(linePositions) > 0 {
		return len(linePositions) - 1
	}
	return 0
}

// Header renders a consistent header across screens
func Header(title, subtitle string, styles *Styles) string {
	return HeaderWithWidth(title, subtitle, styles, 120)
}

// HeaderWithWidth renders a responsive header based on terminal width
func HeaderWithWidth(title, subtitle string, styles *Styles, termWidth int) string {
	var header string

	// Show logo only on wider terminals (> 80 cols)
	if termWidth > 80 {
		logo := `
 ██ ▄█▀▓█████▓██   ██▓ ▄▄▄       ██▀███   ▄████▄   ██░ ██ 
 ██▄█▒ ▓█   ▀ ▒██  ██▒▒████▄    ▓██ ▒ ██▒▒██▀ ▀█  ▓██░ ██▒
▓███▄░ ▒███    ▒██ ██░▒██  ▀█▄  ▓██ ░▄█ ▒▒▓█    ▄ ▒██▀▀██░
▓██ █▄ ▒▓█  ▄  ░ ▐██▓░░██▄▄▄▄██ ▒██▀▀█▄  ▒▓▓▄ ▄██▒░▓█ ░██ 
▒██▒ █▄░▒████▒ ░ ██▒▓░ ▓█   ▓██▒░██▓ ▒██▒▒ ▓███▀ ░░▓█▒░██▓
▒ ▒▒ ▓▒░░ ▒░ ░  ██▒▒▒  ▒▒   ▓▒█░░ ▒▓ ░▒▓░░ ░▒ ▒  ░ ▒ ░░▒░▒
░ ░▒ ▒░ ░ ░  ░▓██ ░▒░   ▒   ▒▒ ░  ░▒ ░ ▒░  ░  ▒    ▒ ░▒░ ░
░ ░░ ░    ░   ▒ ▒ ░░    ░   ▒     ░░   ░ ░         ░  ░░ ░
░  ░      ░  ░░ ░           ░  ░   ░     ░ ░       ░  ░  ░
              ░ ░                        ░                `
		header = styles.Title.Foreground(styles.Theme.Primary).Render(logo) + "\n\n"
	}

	header += styles.RenderTitle(title) + "\n"
	if subtitle != "" {
		header += styles.RenderSubtitle(subtitle) + "\n"
	}

	return header
}

// Footer renders help text at the bottom
func Footer(helpText string, styles *Styles) string {
	return "\n" + styles.RenderHelp(helpText)
}

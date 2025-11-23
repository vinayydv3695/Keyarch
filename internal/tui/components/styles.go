package components

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vinayydv3695/keyarch/internal/config"
)

// Styles holds all component styles
type Styles struct {
	Theme       config.Theme
	Title       lipgloss.Style
	Subtitle    lipgloss.Style
	MenuItem    lipgloss.Style
	ActiveItem  lipgloss.Style
	Stats       lipgloss.Style
	StatValue   lipgloss.Style
	Border      lipgloss.Style
	Help        lipgloss.Style
	Correct     lipgloss.Style
	Incorrect   lipgloss.Style
	Cursor      lipgloss.Style
	Pending     lipgloss.Style
	ProgressBar lipgloss.Style
	// Convenience color styles
	Primary     lipgloss.Style
	Muted       lipgloss.Style
	Success     lipgloss.Style
	Accent      lipgloss.Style
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
			Foreground(theme.Cursor).
			Background(theme.Cursor).
			Underline(true),
		Pending: lipgloss.NewStyle().
			Foreground(theme.Muted),
		ProgressBar: lipgloss.NewStyle().
			Foreground(theme.Success),
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
func (s *Styles) RenderTypingText(target string, userInput string, cursorPos int) string {
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
			// Cursor position
			result.WriteString(s.Cursor.Render(charStr))
		} else {
			// Not yet typed
			result.WriteString(s.Pending.Render(charStr))
		}
	}

	return result.String()
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

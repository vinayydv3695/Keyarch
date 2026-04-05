package settings

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type settingType int

const (
	settingTheme settingType = iota
	settingSound
	settingSoundProfile
	settingBlindMode
	settingBack
)

type Model struct {
	styles       *components.Styles
	cfg          *config.Config
	cursor       int
	width        int
	height       int
	done         bool
	themes       []config.Theme
	themeIndex   int
	expandThemes bool // Whether theme selector is expanded
}

func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	themes := config.AllThemes()

	// Find current theme index
	themeIndex := 0
	for i, t := range themes {
		if t.Name == cfg.Theme {
			themeIndex = i
			break
		}
	}

	return Model{
		styles:     components.NewStyles(theme),
		cfg:        cfg,
		themes:     themes,
		themeIndex: themeIndex,
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
			if m.expandThemes {
				// Navigate themes
				if m.themeIndex > 0 {
					m.themeIndex--
					m.applyTheme()
				}
			} else {
				if m.cursor > 0 {
					m.cursor--
				}
			}

		case "down", "j":
			if m.expandThemes {
				// Navigate themes
				if m.themeIndex < len(m.themes)-1 {
					m.themeIndex++
					m.applyTheme()
				}
			} else {
				if m.cursor < int(settingBack) {
					m.cursor++
				}
			}

		case "left", "h":
			if m.expandThemes {
				return m, nil
			}
			m.handleLeftRight(-1)

		case "right", "l":
			if m.expandThemes {
				return m, nil
			}
			m.handleLeftRight(1)

		case "enter", " ":
			if m.expandThemes {
				// Confirm theme selection
				m.expandThemes = false
				m.cfg.Theme = m.themes[m.themeIndex].Name
				m.cfg.Save()
				return m, nil
			}

			switch settingType(m.cursor) {
			case settingTheme:
				m.expandThemes = true
			case settingSound:
				m.cfg.Sound = !m.cfg.Sound
				m.cfg.Save()
			case settingSoundProfile:
				m.cycleSoundProfile()
				m.cfg.Save()
			case settingBlindMode:
				m.cfg.BlindMode = !m.cfg.BlindMode
				m.cfg.Save()
			case settingBack:
				m.done = true
				return m, tea.Quit
			}

		case "esc":
			if m.expandThemes {
				m.expandThemes = false
				// Revert to saved theme
				for i, t := range m.themes {
					if t.Name == m.cfg.Theme {
						m.themeIndex = i
						break
					}
				}
				m.applyTheme()
				return m, nil
			}
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *Model) handleLeftRight(dir int) {
	switch settingType(m.cursor) {
	case settingTheme:
		newIndex := m.themeIndex + dir
		if newIndex >= 0 && newIndex < len(m.themes) {
			m.themeIndex = newIndex
			m.applyTheme()
			m.cfg.Theme = m.themes[m.themeIndex].Name
			m.cfg.Save()
		}
	case settingSound:
		m.cfg.Sound = !m.cfg.Sound
		m.cfg.Save()
	case settingSoundProfile:
		if dir > 0 {
			m.cycleSoundProfile()
		} else {
			m.cycleSoundProfileReverse()
		}
		m.cfg.Save()
	case settingBlindMode:
		m.cfg.BlindMode = !m.cfg.BlindMode
		m.cfg.Save()
	}
}

func (m *Model) applyTheme() {
	theme := m.themes[m.themeIndex]
	m.styles = components.NewStyles(theme)
}

func (m *Model) cycleSoundProfile() {
	switch m.cfg.SoundProfile {
	case config.SoundProfileOff:
		m.cfg.SoundProfile = config.SoundProfileSubtle
	case config.SoundProfileSubtle:
		m.cfg.SoundProfile = config.SoundProfileMechanical
	case config.SoundProfileMechanical:
		m.cfg.SoundProfile = config.SoundProfileOff
	default:
		m.cfg.SoundProfile = config.SoundProfileSubtle
	}
}

func (m *Model) cycleSoundProfileReverse() {
	switch m.cfg.SoundProfile {
	case config.SoundProfileOff:
		m.cfg.SoundProfile = config.SoundProfileMechanical
	case config.SoundProfileSubtle:
		m.cfg.SoundProfile = config.SoundProfileOff
	case config.SoundProfileMechanical:
		m.cfg.SoundProfile = config.SoundProfileSubtle
	default:
		m.cfg.SoundProfile = config.SoundProfileSubtle
	}
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.HeaderWithWidth("Settings", "Configure your preferences", m.styles, m.width)
	s += "\n\n"

	if m.expandThemes {
		s += m.renderThemeSelector()
	} else {
		s += m.renderSettings()
	}

	return s
}

func (m Model) renderSettings() string {
	var s string

	// Theme setting with inline preview
	themeActive := m.cursor == int(settingTheme)
	currentTheme := m.themes[m.themeIndex]
	themeValue := fmt.Sprintf("< %s >", currentTheme.Name)

	s += m.renderSettingRow("Theme", themeValue, "Color scheme for the app", themeActive)

	// Show color preview dots
	if themeActive {
		preview := "    "
		preview += lipgloss.NewStyle().Foreground(currentTheme.Primary).Render("[primary]") + " "
		preview += lipgloss.NewStyle().Foreground(currentTheme.Correct).Render("[correct]") + " "
		preview += lipgloss.NewStyle().Foreground(currentTheme.Incorrect).Render("[error]") + " "
		preview += lipgloss.NewStyle().Foreground(currentTheme.Accent).Render("[accent]")
		s += m.styles.Muted.Render(preview) + "\n"
	}
	s += "\n"

	// Sound setting
	soundActive := m.cursor == int(settingSound)
	soundValue := "< OFF >"
	if m.cfg.Sound {
		soundValue = "< ON >"
	}
	s += m.renderSettingRow("Sound Effects", soundValue, "Enable typing sounds", soundActive)
	s += "\n"

	// Sound profile setting
	profileActive := m.cursor == int(settingSoundProfile)
	profileValue := fmt.Sprintf("< %s >", m.cfg.SoundProfile)
	if m.cfg.SoundProfile == "" {
		profileValue = "< subtle >"
	}
	s += m.renderSettingRow("Sound Profile", profileValue, "off / subtle / mechanical", profileActive)
	s += "\n"

	// Blind mode setting
	blindActive := m.cursor == int(settingBlindMode)
	blindValue := "< OFF >"
	if m.cfg.BlindMode {
		blindValue = "< ON >"
	}
	s += m.renderSettingRow("Blind Mode", blindValue, "Hide stats while typing", blindActive)
	s += "\n"

	// Back option
	backActive := m.cursor == int(settingBack)
	if backActive {
		s += m.styles.ActiveItem.Render("> Back to Menu") + "\n"
	} else {
		s += m.styles.MenuItem.Render("  Back to Menu") + "\n"
	}

	s += "\n"
	s += components.Footer("Up/Down: Navigate | Left/Right: Change | Enter: Select | ESC: Back", m.styles)

	return s
}

func (m Model) renderSettingRow(label, value, desc string, active bool) string {
	var s string

	// Calculate widths for alignment
	labelWidth := 16
	valueWidth := 20

	labelStr := fmt.Sprintf("%-*s", labelWidth, label)
	valueStr := fmt.Sprintf("%-*s", valueWidth, value)

	if active {
		s += m.styles.ActiveItem.Render("> "+labelStr) + " "
		s += m.styles.Primary.Bold(true).Render(valueStr) + "\n"
		s += "    " + m.styles.Muted.Render(desc) + "\n"
	} else {
		s += m.styles.MenuItem.Render("  "+labelStr) + " "
		s += m.styles.Muted.Render(valueStr) + "\n"
	}

	return s
}

func (m Model) renderThemeSelector() string {
	var s string

	s += m.styles.Subtitle.Render("Select a theme (changes apply in real-time):") + "\n\n"

	// Show themes in a nice grid/list
	for i, theme := range m.themes {
		active := i == m.themeIndex

		// Theme name with selection indicator
		var line string
		if active {
			line = m.styles.ActiveItem.Render("> " + theme.Name)
		} else {
			line = m.styles.MenuItem.Render("  " + theme.Name)
		}

		// Color preview dots
		preview := "  "
		preview += lipgloss.NewStyle().Foreground(theme.Primary).Render("●") + " "
		preview += lipgloss.NewStyle().Foreground(theme.Correct).Render("●") + " "
		preview += lipgloss.NewStyle().Foreground(theme.Incorrect).Render("●") + " "
		preview += lipgloss.NewStyle().Foreground(theme.Accent).Render("●") + " "
		preview += lipgloss.NewStyle().Foreground(theme.Muted).Render("●")

		s += line + preview + "\n"
	}

	s += "\n"
	s += m.styles.RenderBox("Preview: The colors you see now are from the selected theme")

	s += "\n"
	s += components.Footer("Up/Down: Browse | Enter: Confirm | ESC: Cancel", m.styles)

	return s
}

func (m Model) Done() bool {
	return m.done
}

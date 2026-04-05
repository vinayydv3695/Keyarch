package settings

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type Model struct {
	styles *components.Styles
	cfg    *config.Config
	cursor int
	width  int
	height int
	done   bool
}

func New(cfg *config.Config) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
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

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < 3 { // 4 settings (0-3)
				m.cursor++
			}

		case "enter", " ":
			// Toggle the selected setting
			switch m.cursor {
			case 0:
				m.cfg.Sound = !m.cfg.Sound
				m.cfg.Save()
			case 1:
				// Cycle through sound profiles
				m.cycleSoundProfile()
				m.cfg.Save()
			case 2:
				m.cfg.BlindMode = !m.cfg.BlindMode
				m.cfg.Save()
			case 3:
				m.done = true
				return m, tea.Quit
			}

		case "esc":
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
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

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.HeaderWithWidth("Settings", "Configure your preferences", m.styles, m.width)
	s += "\n\n"

	// Sound setting
	soundValue := "OFF"
	if m.cfg.Sound {
		soundValue = "ON"
	}
	soundDesc := "Enable sound effects during typing"
	if m.cursor == 0 {
		s += m.styles.ActiveItem.Render("> Sound Effects: "+soundValue) + "\n"
		s += "  " + m.styles.Muted.Render(soundDesc) + "\n\n"
	} else {
		s += m.styles.MenuItem.Render("  Sound Effects: "+soundValue) + "\n"
		s += "  " + m.styles.Muted.Render(soundDesc) + "\n\n"
	}

	// Sound profile setting
	profileValue := string(m.cfg.SoundProfile)
	if profileValue == "" {
		profileValue = "subtle"
	}
	profileDesc := "Sound style: off, subtle, or mechanical"
	if m.cursor == 1 {
		s += m.styles.ActiveItem.Render("> Sound Profile: "+profileValue) + "\n"
		s += "  " + m.styles.Muted.Render(profileDesc) + "\n\n"
	} else {
		s += m.styles.MenuItem.Render("  Sound Profile: "+profileValue) + "\n"
		s += "  " + m.styles.Muted.Render(profileDesc) + "\n\n"
	}

	// Blind mode setting
	blindValue := "OFF"
	if m.cfg.BlindMode {
		blindValue = "ON"
	}
	blindDesc := "Hide WPM and stats during typing (hardcore mode)"
	if m.cursor == 2 {
		s += m.styles.ActiveItem.Render("> Blind Mode: "+blindValue) + "\n"
		s += "  " + m.styles.Muted.Render(blindDesc) + "\n\n"
	} else {
		s += m.styles.MenuItem.Render("  Blind Mode: "+blindValue) + "\n"
		s += "  " + m.styles.Muted.Render(blindDesc) + "\n\n"
	}

	// Back option
	if m.cursor == 3 {
		s += m.styles.ActiveItem.Render("> Back to Menu") + "\n"
	} else {
		s += m.styles.MenuItem.Render("  Back to Menu") + "\n"
	}

	s += "\n"
	s += components.Footer("Up/Down: Navigate | Enter/Space: Toggle | ESC: Back", m.styles)

	return s
}

func (m Model) Done() bool {
	return m.done
}

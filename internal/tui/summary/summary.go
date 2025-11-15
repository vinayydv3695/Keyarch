package summary

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/engine"
	"github.com/vinayydv3695/keyarch/internal/storage"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type Model struct {
	engine   *engine.Engine
	styles   *components.Styles
	cfg      *config.Config
	db       *storage.DB
	width    int
	height   int
	saved    bool
	done     bool
}

func New(eng *engine.Engine, cfg *config.Config, db *storage.DB) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		engine: eng,
		styles: components.NewStyles(theme),
		cfg:    cfg,
		db:     db,
	}
}

func (m Model) Init() tea.Cmd {
	return m.saveResult
}

func (m Model) saveResult() tea.Msg {
	if m.db == nil || m.saved {
		return nil
	}

	result := &storage.TestResult{
		WPM:         m.engine.GetWPM(),
		CPM:         m.engine.GetCPM(),
		Accuracy:    m.engine.GetAccuracy(),
		Mistakes:    m.engine.Mistakes,
		TestType:    string(m.engine.Mode),
		TestMode:    m.getTestMode(),
		Duration:    int(m.engine.GetElapsedTime()),
		CreatedAt:   time.Now(),
		TotalChars:  len(m.engine.UserInput),
		CorrectChar: m.countCorrect(),
	}

	m.db.SaveResult(result)
	m.saved = true
	return nil
}

func (m Model) getTestMode() string {
	switch m.engine.Mode {
	case engine.ModeTimer:
		return fmt.Sprintf("%ds", m.engine.Duration)
	case engine.ModeWords:
		return fmt.Sprintf("%dw", m.engine.WordCount)
	default:
		return "custom"
	}
}

func (m Model) countCorrect() int {
	correct := 0
	minLen := len(m.engine.UserInput)
	if len(m.engine.TargetText) < minLen {
		minLen = len(m.engine.TargetText)
	}

	for i := 0; i < minLen; i++ {
		if m.engine.UserInput[i] == m.engine.TargetText[i] {
			correct++
		}
	}
	return correct
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
			m.done = true
			return m, tea.Quit

		case "enter", "esc", " ":
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	s := components.Header("Test Complete!", "Here are your results", m.styles)
	s += "\n\n"

	// Main stats
	mainStats := ""
	mainStats += m.styles.Title.Render(fmt.Sprintf("%.0f WPM", m.engine.GetWPM())) + "\n"
	mainStats += m.styles.Subtitle.Render(fmt.Sprintf("%.0f CPM", m.engine.GetCPM())) + "\n"
	mainStats += m.styles.Subtitle.Render(fmt.Sprintf("%.1f%% Accuracy", m.engine.GetAccuracy())) + "\n\n"

	// Detailed stats
	elapsed := int(m.engine.GetElapsedTime())
	mainStats += m.styles.RenderStat("Time", fmt.Sprintf("%ds", elapsed)) + "\n"
	mainStats += m.styles.RenderStat("Characters", fmt.Sprintf("%d", len(m.engine.UserInput))) + "\n"
	mainStats += m.styles.RenderStat("Mistakes", fmt.Sprintf("%d", m.engine.Mistakes)) + "\n"
	mainStats += m.styles.RenderStat("Correct", fmt.Sprintf("%d", m.countCorrect())) + "\n"

	s += m.styles.RenderBox(mainStats)
	s += "\n"

	// Weak keys
	weakKeys := m.engine.GetWeakKeys(5)
	if len(weakKeys) > 0 {
		keysInfo := m.styles.Title.Render("Weak Keys") + "\n\n"
		for i, k := range weakKeys {
			if i >= 3 {
				break
			}
			keysInfo += fmt.Sprintf("  '%c' - %.0f%% accuracy\n", k.Key, k.Accuracy)
		}
		s += m.styles.RenderBox(keysInfo)
		s += "\n"
	}

	// Performance message
	wpm := m.engine.GetWPM()
	var msg string
	switch {
	case wpm >= 100:
		msg = "🔥 Incredible! You're a typing master!"
	case wpm >= 80:
		msg = "⚡ Excellent! Very fast typing!"
	case wpm >= 60:
		msg = "👍 Great job! Above average!"
	case wpm >= 40:
		msg = "✓ Good work! Keep practicing!"
	default:
		msg = "💪 Keep practicing to improve!"
	}

	s += "\n" + m.styles.Subtitle.Render(msg) + "\n"

	s += components.Footer("Enter: Return to menu • Ctrl+C: Quit", m.styles)

	return s
}

func (m Model) Done() bool {
	return m.done
}

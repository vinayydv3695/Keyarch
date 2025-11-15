package stats

import (
	"fmt"
	"math"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/storage"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type Model struct {
	styles  *components.Styles
	cfg     *config.Config
	db      *storage.DB
	stats   *storage.UserStats
	recent  []storage.TestResult
	graph   map[string]float64
	width   int
	height  int
	done    bool
	loading bool
	err     error
}

func New(cfg *config.Config, db *storage.DB) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		styles:  components.NewStyles(theme),
		cfg:     cfg,
		db:      db,
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadStats
}

func (m Model) loadStats() tea.Msg {
	stats, err := m.db.GetStats()
	if err != nil {
		return err
	}

	recent, err := m.db.GetRecentResults(10)
	if err != nil {
		return err
	}

	graph, err := m.db.GetLast7DaysWPM()
	if err != nil {
		return err
	}

	return statsLoadedMsg{
		stats:  stats,
		recent: recent,
		graph:  graph,
	}
}

type statsLoadedMsg struct {
	stats  *storage.UserStats
	recent []storage.TestResult
	graph  map[string]float64
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case statsLoadedMsg:
		m.stats = msg.stats
		m.recent = msg.recent
		m.graph = msg.graph
		m.loading = false
		return m, nil

	case error:
		m.err = msg
		m.loading = false
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit
			
		case "q":
			m.done = true
			return m, tea.Quit

		case "esc", "enter", " ":
			m.done = true
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.loading {
		return "Loading statistics..."
	}

	if m.err != nil {
		return fmt.Sprintf("Error loading stats: %v\n\nPress ESC to go back", m.err)
	}

	s := components.Header("Statistics", "Your typing journey", m.styles)
	s += "\n\n"

	if m.stats.TotalTests == 0 {
		s += m.styles.RenderBox("No tests completed yet.\nStart typing to see your stats!")
		s += components.Footer("ESC: Back • Ctrl+C: Quit", m.styles)
		return s
	}

	// Overview
	overview := m.styles.Title.Render("Overview") + "\n\n"
	overview += m.styles.RenderStat("Total Tests", fmt.Sprintf("%d", m.stats.TotalTests)) + "\n"
	overview += m.styles.RenderStat("Best WPM", fmt.Sprintf("%.0f", m.stats.BestWPM)) + "\n"
	overview += m.styles.RenderStat("Average WPM", fmt.Sprintf("%.0f", m.stats.AverageWPM)) + "\n"
	overview += m.styles.RenderStat("Best Accuracy", fmt.Sprintf("%.1f%%", m.stats.BestAccuracy)) + "\n"
	overview += m.styles.RenderStat("Avg Accuracy", fmt.Sprintf("%.1f%%", m.stats.AverageAccuracy)) + "\n"
	overview += m.styles.RenderStat("Total Time", fmt.Sprintf("%dm", m.stats.TotalTime/60)) + "\n"
	overview += m.styles.RenderStat("Current Streak", fmt.Sprintf("%d days", m.stats.CurrentStreak)) + "\n"

	s += m.styles.RenderBox(overview)
	s += "\n"

	// Last 7 days graph
	if len(m.graph) > 0 {
		graphView := m.styles.Title.Render("Last 7 Days WPM") + "\n\n"
		graphView += m.renderGraph(m.graph)
		s += m.styles.RenderBox(graphView)
		s += "\n"
	}

	// Recent tests
	if len(m.recent) > 0 {
		recentView := m.styles.Title.Render("Recent Tests") + "\n\n"
		for i, r := range m.recent {
			if i >= 5 {
				break
			}
			recentView += fmt.Sprintf("%s - %.0f WPM - %.1f%% - %s\n",
				r.CreatedAt.Format("Jan 02"),
				r.WPM,
				r.Accuracy,
				r.TestMode,
			)
		}
		s += m.styles.RenderBox(recentView)
	}

	s += components.Footer("ESC: Back • Ctrl+C: Quit", m.styles)

	return s
}

func (m Model) renderGraph(data map[string]float64) string {
	if len(data) == 0 {
		return "No data available"
	}

	// Find max value for scaling
	maxWPM := 0.0
	for _, wpm := range data {
		if wpm > maxWPM {
			maxWPM = wpm
		}
	}

	if maxWPM == 0 {
		return "No data available"
	}

	// Create simple bar graph
	graph := ""
	maxHeight := 10
	
	// Get sorted dates (simplified - just use map iteration)
	for date, wpm := range data {
		height := int(math.Round((wpm / maxWPM) * float64(maxHeight)))
		bar := ""
		for i := 0; i < height; i++ {
			bar += "█"
		}
		graph += fmt.Sprintf("%s: %s %.0f\n", date[5:], bar, wpm)
	}

	return graph
}

func (m Model) Done() bool {
	return m.done
}

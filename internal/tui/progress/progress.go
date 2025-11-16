package progress

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/achievements"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/goals"
	"github.com/vinayydv3695/keyarch/internal/storage"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

type Model struct {
	cfg          *config.Config
	db           *storage.DB
	styles       *components.Styles
	width        int
	height       int
	achievements []achievements.Achievement
	dailyGoals   []goals.Goal
	weeklyGoals  []goals.Goal
	selected     string
	tab          int // 0=achievements, 1=daily goals, 2=weekly goals
}

func New(cfg *config.Config, db *storage.DB) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		cfg:    cfg,
		db:     db,
		styles: components.NewStyles(theme),
		tab:    0,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadData()
}

func (m Model) loadData() tea.Cmd {
	return func() tea.Msg {
		return dataLoadedMsg{}
	}
}

type dataLoadedMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case dataLoadedMsg:
		// Load achievements
		m.achievements = m.loadAchievements()
		
		// Load goals
		m.dailyGoals, m.weeklyGoals = m.loadGoals()
		
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit

		case "q", "esc":
			m.selected = "back"
			return m, tea.Quit

		case "tab", "right", "l":
			m.tab = (m.tab + 1) % 3

		case "shift+tab", "left", "h":
			m.tab = (m.tab - 1 + 3) % 3

		case "enter", " ":
			m.selected = "back"
			return m, tea.Quit
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
	s += components.HeaderWithWidth("PROGRESS & ACHIEVEMENTS", "Track your typing journey", m.styles, m.width)
	s += "\n\n"

	// Tabs
	tabs := []string{"Achievements", "Daily Goals", "Weekly Goals"}
	tabStr := ""
	for i, tab := range tabs {
		if i == m.tab {
			tabStr += m.styles.Primary.Render("▶ " + tab + " ◀") + "  "
		} else {
			tabStr += m.styles.Muted.Render(tab) + "  "
		}
	}
	s += tabStr + "\n\n"

	// Content based on tab
	switch m.tab {
	case 0:
		s += m.renderAchievements()
	case 1:
		s += m.renderGoals(m.dailyGoals, "Daily")
	case 2:
		s += m.renderGoals(m.weeklyGoals, "Weekly")
	}

	s += "\n"
	s += components.Footer("Tab: Switch • Enter: Back • q/Ctrl+C: Quit", m.styles)

	return s
}

func (m Model) renderAchievements() string {
	var s string
	
	unlocked := achievements.GetUnlockedCount(m.achievements)
	total := len(m.achievements)
	
	s += m.styles.Subtitle.Render(fmt.Sprintf("Unlocked: %d/%d\n\n", unlocked, total))

	// Group by category
	categories := []string{"speed", "accuracy", "dedication", "consistency"}
	categoryNames := map[string]string{
		"speed":       "⚡ Speed",
		"accuracy":    "🎯 Accuracy",
		"dedication":  "📚 Dedication",
		"consistency": "🔥 Consistency",
	}

	for _, cat := range categories {
		catAchievements := achievements.GetByCategory(m.achievements, cat)
		if len(catAchievements) == 0 {
			continue
		}

		s += m.styles.Primary.Bold(true).Render(categoryNames[cat]) + "\n"
		
		for _, achievement := range catAchievements {
			if achievement.Hidden && !achievement.Unlocked {
				continue
			}

			line := achievement.Icon + " " + achievement.Name
			if achievement.Unlocked {
				line = m.styles.Success.Render(line + " ✓")
			} else {
				line = m.styles.Muted.Render(line)
				line += m.styles.Muted.Render(fmt.Sprintf(" (%d/%d)", achievement.Progress, achievement.Target))
			}
			s += "  " + line + "\n"
			s += "     " + m.styles.Muted.Render(achievement.Description) + "\n"
		}
		s += "\n"
	}

	return s
}

func (m Model) renderGoals(goalsList []goals.Goal, period string) string {
	var s string
	
	completed := goals.GetCompletedCount(goalsList)
	total := len(goalsList)
	
	s += m.styles.Subtitle.Render(fmt.Sprintf("%s Goals: %d/%d completed\n\n", period, completed, total))

	for _, goal := range goalsList {
		progress := goals.GetProgress(goal)
		
		// Goal name and icon
		line := goal.Icon + " " + goal.Name
		if goal.Completed {
			line = m.styles.Success.Render(line + " ✓")
		} else {
			line = m.styles.Primary.Render(line)
		}
		s += line + "\n"

		// Description
		s += "  " + m.styles.Muted.Render(goal.Description) + "\n"

		// Progress bar
		barWidth := 30
		filledWidth := (progress * barWidth) / 100
		bar := strings.Repeat("█", filledWidth) + strings.Repeat("░", barWidth-filledWidth)
		
		if goal.Completed {
			s += "  " + m.styles.Success.Render(bar) + " " + m.styles.Success.Render(fmt.Sprintf("%d%%", progress)) + "\n"
		} else {
			s += "  " + m.styles.Accent.Render(bar) + " " + m.styles.Primary.Render(fmt.Sprintf("%d%%", progress)) + "\n"
		}

		// Current vs target
		unitStr := goal.Unit
		if goal.Unit == "seconds" {
			unitStr = "minutes"
			goal.Current = goal.Current / 60
			goal.Target = goal.Target / 60
		}
		if goal.Unit == "percent" {
			unitStr = "%"
		}
		
		s += "  " + m.styles.Muted.Render(fmt.Sprintf("%d/%d %s", goal.Current, goal.Target, unitStr)) + "\n\n"
	}

	return s
}

func (m Model) loadAchievements() []achievements.Achievement {
	allAchievements := achievements.AllAchievements()
	
	// Get unlocked achievement IDs from database
	unlockedIDs, err := m.db.GetUnlockedAchievements()
	if err != nil {
		return allAchievements
	}

	unlockedMap := make(map[string]bool)
	for _, id := range unlockedIDs {
		unlockedMap[id] = true
	}

	// Get stats for progress calculation
	stats, err := m.db.GetStats()
	if err != nil {
		return allAchievements
	}

	modeStats, _ := m.db.GetModeStats()

	// Update achievement status
	for i := range allAchievements {
		if unlockedMap[allAchievements[i].ID] {
			allAchievements[i].Unlocked = true
		}

		// Calculate progress
		statsMap := map[string]interface{}{
			"best_wpm":         stats.BestWPM,
			"best_accuracy":    stats.BestAccuracy,
			"total_tests":      stats.TotalTests,
			"current_streak":   stats.CurrentStreak,
			"total_time":       stats.TotalTime,
			"modes_completed":  len(modeStats),
			"code_tests":       modeStats["code"],
			"quote_tests":      modeStats["quote"],
		}

		unlocked, progress := achievements.CheckAchievement(allAchievements[i], statsMap)
		allAchievements[i].Progress = progress

		// Auto-unlock if criteria met
		if unlocked && !allAchievements[i].Unlocked {
			allAchievements[i].Unlocked = true
			m.db.SaveAchievement(allAchievements[i].ID)
		}
	}

	return allAchievements
}

func (m Model) loadGoals() ([]goals.Goal, []goals.Goal) {
	// Get saved goals from database
	savedGoals, err := m.db.GetGoals()
	if err != nil {
		// Return defaults if no saved goals
		return goals.DefaultDailyGoals(), goals.DefaultWeeklyGoals()
	}

	// Convert to goal objects
	var daily, weekly []goals.Goal
	
	if len(savedGoals) == 0 {
		return goals.DefaultDailyGoals(), goals.DefaultWeeklyGoals()
	}

	for _, record := range savedGoals {
		goal := goals.Goal{
			ID:        record.ID,
			Type:      record.Type,
			Target:    record.Target,
			Current:   record.Current,
			Completed: record.Completed,
			StartDate: record.StartDate,
			EndDate:   record.EndDate,
			Active:    true,
		}

		// Set name, description, unit, icon based on ID
		defaultDaily := goals.DefaultDailyGoals()
		defaultWeekly := goals.DefaultWeeklyGoals()
		
		for _, d := range append(defaultDaily, defaultWeekly...) {
			if d.ID == goal.ID {
				goal.Name = d.Name
				goal.Description = d.Description
				goal.Unit = d.Unit
				goal.Icon = d.Icon
				break
			}
		}

		if record.Type == "daily" {
			daily = append(daily, goal)
		} else if record.Type == "weekly" {
			weekly = append(weekly, goal)
		}
	}

	// Reset expired goals
	daily = goals.ResetDailyGoals(daily)
	weekly = goals.ResetWeeklyGoals(weekly)

	// Save reset goals
	for _, goal := range append(daily, weekly...) {
		m.db.SaveGoal(storage.GoalRecord{
			ID:        goal.ID,
			Type:      goal.Type,
			Target:    goal.Target,
			Current:   goal.Current,
			Completed: goal.Completed,
			StartDate: goal.StartDate,
			EndDate:   goal.EndDate,
		})
	}

	return daily, weekly
}

func (m Model) Selected() string {
	return m.selected
}

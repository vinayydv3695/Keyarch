package summary

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/achievements"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/engine"
	"github.com/vinayydv3695/keyarch/internal/goals"
	"github.com/vinayydv3695/keyarch/internal/storage"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
	"github.com/vinayydv3695/keyarch/internal/tui/heatmap"
)

type Model struct {
	engine              *engine.Engine
	styles              *components.Styles
	cfg                 *config.Config
	db                  *storage.DB
	width               int
	height              int
	saved               bool
	done                bool
	newAchievements     []achievements.Achievement
	completedGoals      []goals.Goal
	achievementsChecked bool
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

	// Update goals after saving result
	m.updateGoals()

	// Check for new achievements
	m.checkAchievements()

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

func (m *Model) updateGoals() {
	if m.db == nil {
		return
	}

	// Load current goal records from database
	currentRecords, err := m.db.GetGoals()
	if err != nil {
		return
	}

	// Get all goal templates
	allDailyGoals := goals.DefaultDailyGoals()
	allWeeklyGoals := goals.DefaultWeeklyGoals()
	allGoalTemplates := append(allDailyGoals, allWeeklyGoals...)

	// Initialize with defaults if no goals exist
	if len(currentRecords) == 0 {
		for _, g := range allDailyGoals {
			record := goals.ToRecord(g)
			m.db.SaveGoal(record)
		}

		for _, g := range allWeeklyGoals {
			record := goals.ToRecord(g)
			m.db.SaveGoal(record)
		}

		// Reload after initialization
		currentRecords, _ = m.db.GetGoals()
	}

	// Convert records to goals and update progress
	for i := range currentRecords {
		record := &currentRecords[i]
		goal := goals.FromRecord(*record, allGoalTemplates)

		// Check if goal needs to be reset
		if goals.ShouldResetGoal(&goal) {
			if goal.Type == "daily" {
				goal = goals.ResetDailyGoal(&goal)
			} else if goal.Type == "weekly" {
				goal = goals.ResetWeeklyGoal(&goal)
			}
		}

		// Update progress based on this test
		updated := false
		switch {
		case goal.ID == "daily_tests_5" || goal.ID == "weekly_tests_25":
			goal.Current++
			updated = true

		case goal.ID == "daily_wpm_50":
			wpm := m.engine.GetWPM()
			if wpm > goal.Current {
				goal.Current = wpm
				updated = true
			}

		case goal.ID == "daily_accuracy_95":
			accuracy := m.engine.GetAccuracy()
			if accuracy > goal.Current {
				goal.Current = accuracy
				updated = true
			}

		case goal.ID == "weekly_time_30":
			duration := int(m.engine.GetElapsedTime())
			goal.Current += float64(duration) / 60.0 // Add minutes
			updated = true

		case goal.ID == "weekly_streak_7":
			// Streak is calculated separately based on dates
			updated = true
		}

		// Check if goal is completed
		if goal.Current >= goal.Target && !goal.Completed {
			goal.Completed = true
			m.completedGoals = append(m.completedGoals, goal)
		}

		// Save updated goal
		if updated {
			updatedRecord := goals.ToRecord(goal)
			m.db.SaveGoal(updatedRecord)
		}
	}
}

func (m *Model) checkAchievements() {
	if m.db == nil || m.achievementsChecked {
		return
	}
	m.achievementsChecked = true

	// Get all achievements
	allAchievements := achievements.AllAchievements()

	// Get already unlocked achievements
	unlocked, err := m.db.GetUnlockedAchievements()
	if err != nil {
		return
	}

	// Create a map of unlocked achievement IDs
	unlockedMap := make(map[string]bool)
	for _, achievementID := range unlocked {
		unlockedMap[achievementID] = true
	}

	// Get current stats for checking
	stats, err := m.db.GetStats()
	if err != nil {
		return
	}

	// Get mode stats
	modeStats, err := m.db.GetModeStats()
	if err != nil {
		modeStats = make(map[string]int)
	}

	// Convert stats to map for achievement checking
	statsMap := make(map[string]interface{})
	statsMap["best_wpm"] = stats.BestWPM
	statsMap["average_wpm"] = stats.AverageWPM
	statsMap["best_accuracy"] = stats.BestAccuracy
	statsMap["average_accuracy"] = stats.AverageAccuracy
	statsMap["total_tests"] = stats.TotalTests
	statsMap["total_time"] = stats.TotalTime
	statsMap["current_streak"] = stats.CurrentStreak

	// Add mode-specific stats
	statsMap["code_tests"] = modeStats["code"]
	statsMap["quote_tests"] = modeStats["quote"]
	
	// Count unique modes completed
	modesCompleted := 0
	for _, count := range modeStats {
		if count > 0 {
			modesCompleted++
		}
	}
	statsMap["modes_completed"] = modesCompleted

	// Check each achievement
	for _, ach := range allAchievements {
		// Skip if already unlocked
		if unlockedMap[ach.ID] {
			continue
		}

		// Check if this achievement should be unlocked
		shouldUnlock, _ := achievements.CheckAchievement(ach, statsMap)
		if shouldUnlock {
			// Unlock the achievement
			m.db.SaveAchievement(ach.ID)

			// Add to new achievements list
			m.newAchievements = append(m.newAchievements, ach)
		}
	}
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

	s := components.HeaderWithWidth("Test Complete!", "Here are your results", m.styles, m.width)
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

	// New achievements unlocked
	if len(m.newAchievements) > 0 {
		achievementInfo := m.styles.Title.Render("🎉 New Achievements Unlocked!") + "\n\n"
		for _, ach := range m.newAchievements {
			achievementInfo += fmt.Sprintf("  %s %s\n", ach.Icon, m.styles.Accent.Render(ach.Name))
			achievementInfo += fmt.Sprintf("     %s\n", m.styles.Muted.Render(ach.Description))
		}
		s += m.styles.RenderBox(achievementInfo)
		s += "\n"
	}

	// Completed goals
	if len(m.completedGoals) > 0 {
		goalInfo := m.styles.Title.Render("✅ Goals Completed!") + "\n\n"
		for _, goal := range m.completedGoals {
			goalInfo += fmt.Sprintf("  %s %s\n", "🎯", m.styles.Accent.Render(goal.Name))
		}
		s += m.styles.RenderBox(goalInfo)
		s += "\n"
	}

	// Weak keys (only show on wider terminals)
	if m.width > 60 {
		weakKeys := m.engine.GetWeakKeys(5)
		if len(weakKeys) > 0 {
			keysInfo := m.styles.Title.Render("Weak Keys") + "\n\n"
			maxKeys := 3
			if m.width < 80 {
				maxKeys = 2
			}
			for i, k := range weakKeys {
				if i >= maxKeys {
					break
				}
				keysInfo += fmt.Sprintf("  '%c' - %.0f%% accuracy\n", k.Key, k.Accuracy)
			}
			s += m.styles.RenderBox(keysInfo)
			s += "\n"
		}
	}

	// Typing Heatmap (only show on wider terminals)
	if m.width > 80 && len(m.engine.KeyStrokes) > 0 {
		heatmapView := heatmap.RenderHeatmap(m.engine.KeyStrokes)
		s += m.styles.RenderBox(heatmapView)
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

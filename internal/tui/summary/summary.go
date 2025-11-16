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

	s := components.Header("✅ Test Complete!", "Excellent work! Here are your results", m.styles)
	s += "\n"

	// Performance message first
	wpm := m.engine.GetWPM()
	accuracy := m.engine.GetAccuracy()
	var msg string
	var emoji string
	switch {
	case wpm >= 100:
		msg = "Incredible! You're a typing master!"
		emoji = "🔥"
	case wpm >= 80:
		msg = "Excellent! Very fast typing!"
		emoji = "⚡"
	case wpm >= 60:
		msg = "Great job! Above average!"
		emoji = "👍"
	case wpm >= 40:
		msg = "Good work! Keep practicing!"
		emoji = "✓"
	default:
		msg = "Keep practicing to improve!"
		emoji = "💪"
	}
	
	msgBox := emoji + "  " + m.styles.Accent.Bold(true).Render(msg)
	s += m.styles.Border.Render(msgBox) + "\n\n"

	// Main stats in a prominent card
	mainStats := m.styles.Title.Render("📊 PERFORMANCE") + "\n\n"
	
	// Big WPM display
	mainStats += "  " + m.styles.Primary.Bold(true).Render(fmt.Sprintf("%.0f", wpm)) + 
		m.styles.Muted.Render(" WPM") + "\n"
	
	// Accuracy with color
	accuracyColor := m.styles.Success
	if accuracy < 95 {
		accuracyColor = m.styles.Accent
	}
	mainStats += "  " + accuracyColor.Bold(true).Render(fmt.Sprintf("%.1f%%", accuracy)) + 
		m.styles.Muted.Render(" Accuracy") + "\n\n"

	// Detailed stats in columns
	elapsed := int(m.engine.GetElapsedTime())
	mainStats += "  " + m.styles.Muted.Render("⏱️  Time:       ") + m.styles.Primary.Render(fmt.Sprintf("%ds", elapsed)) + "\n"
	mainStats += "  " + m.styles.Muted.Render("📝 Characters: ") + m.styles.Primary.Render(fmt.Sprintf("%d", len(m.engine.UserInput))) + "\n"
	mainStats += "  " + m.styles.Muted.Render("✓  Correct:    ") + m.styles.Success.Render(fmt.Sprintf("%d", m.countCorrect())) + "\n"
	
	if m.engine.Mistakes > 0 {
		mainStats += "  " + m.styles.Muted.Render("✗  Mistakes:   ") + m.styles.Accent.Render(fmt.Sprintf("%d", m.engine.Mistakes)) + "\n"
	} else {
		mainStats += "  " + m.styles.Success.Render("✓  Perfect! No mistakes! 🎯") + "\n"
	}

	s += m.styles.Border.Render(mainStats) + "\n"

	// New achievements unlocked with celebration
	if len(m.newAchievements) > 0 {
		achievementInfo := m.styles.Title.Render("🎉 NEW ACHIEVEMENTS") + "\n\n"
		for _, ach := range m.newAchievements {
			badgeBox := m.styles.Accent.Bold(true).Render(ach.Icon+" "+ach.Name) + "\n" +
				"  " + m.styles.Muted.Render(ach.Description)
			achievementInfo += m.styles.Border.Copy().
				BorderForeground(m.styles.Theme.Success).
				Render(badgeBox) + "\n"
		}
		s += achievementInfo + "\n"
	}

	// Completed goals
	if len(m.completedGoals) > 0 {
		goalInfo := m.styles.Title.Render("✅ GOALS COMPLETED") + "\n\n"
		for _, goal := range m.completedGoals {
			goalBox := "🎯 " + m.styles.Success.Bold(true).Render(goal.Name)
			goalInfo += m.styles.Border.Copy().
				BorderForeground(m.styles.Theme.Success).
				Render(goalBox) + "\n"
		}
		s += goalInfo + "\n"
	}

	// Weak keys analysis
	weakKeys := m.engine.GetWeakKeys(5)
	if len(weakKeys) > 0 {
		keysInfo := m.styles.Title.Render("🎯 FOCUS AREAS") + "\n\n"
		keysInfo += m.styles.Muted.Render("  Practice these keys to improve:\n\n")
		for i, k := range weakKeys {
			if i >= 3 {
				break
			}
			keyDisplay := fmt.Sprintf("'%c'", k.Key)
			accuracy := fmt.Sprintf("%.0f%%", k.Accuracy)
			
			barLength := int(k.Accuracy / 10)
			bar := ""
			for j := 0; j < 10; j++ {
				if j < barLength {
					bar += "█"
				} else {
					bar += "░"
				}
			}
			
			keysInfo += fmt.Sprintf("  %s  %s %s\n", 
				m.styles.Primary.Bold(true).Render(keyDisplay),
				bar,
				m.styles.Muted.Render(accuracy))
		}
		s += m.styles.Border.Render(keysInfo) + "\n"
	}

	// Typing Heatmap
	if len(m.engine.KeyStrokes) > 0 {
		heatmapBox := m.styles.Title.Render("🔥 TYPING HEATMAP") + "\n"
		heatmapBox += heatmap.RenderHeatmap(m.engine.KeyStrokes)
		s += m.styles.Border.Render(heatmapBox) + "\n"
	}

	s += "\n" + components.Footer("Press Enter to return to menu", m.styles)

	return s
}

func (m Model) Done() bool {
	return m.done
}

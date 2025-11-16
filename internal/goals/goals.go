package goals

import (
	"time"
	"github.com/vinayydv3695/keyarch/internal/storage"
)

// Goal represents a daily/weekly goal
type Goal struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Type        string    `json:"type"`        // daily, weekly, custom
	Target      float64   `json:"target"`      // target value
	Current     float64   `json:"current"`     // current progress
	Unit        string    `json:"unit"`        // tests, wpm, minutes
	Icon        string    `json:"icon"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
	Completed   bool      `json:"completed"`
	Active      bool      `json:"active"`
}

// DefaultDailyGoals returns default daily goals
func DefaultDailyGoals() []Goal {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	return []Goal{
		{
			ID:          "daily_tests_5",
			Name:        "Daily Practice",
			Description: "Complete 5 tests today",
			Type:        "daily",
			Target:      5,
			Current:     0,
			Unit:        "tests",
			Icon:        "🎯",
			StartDate:   today,
			EndDate:     tomorrow,
			Active:      true,
		},
		{
			ID:          "daily_wpm_50",
			Name:        "Speed Goal",
			Description: "Reach 50 WPM in any test",
			Type:        "daily",
			Target:      50,
			Current:     0,
			Unit:        "wpm",
			Icon:        "⚡",
			StartDate:   today,
			EndDate:     tomorrow,
			Active:      true,
		},
		{
			ID:          "daily_accuracy_95",
			Name:        "Accuracy Goal",
			Description: "Achieve 95% accuracy",
			Type:        "daily",
			Target:      95,
			Current:     0,
			Unit:        "percent",
			Icon:        "🎯",
			StartDate:   today,
			EndDate:     tomorrow,
			Active:      true,
		},
	}
}

// DefaultWeeklyGoals returns default weekly goals
func DefaultWeeklyGoals() []Goal {
	now := time.Now()
	// Start of week (Monday)
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	return []Goal{
		{
			ID:          "weekly_tests_25",
			Name:        "Weekly Commitment",
			Description: "Complete 25 tests this week",
			Type:        "weekly",
			Target:      25,
			Current:     0,
			Unit:        "tests",
			Icon:        "📚",
			StartDate:   startOfWeek,
			EndDate:     endOfWeek,
			Active:      true,
		},
		{
			ID:          "weekly_time_30min",
			Name:        "Time Investment",
			Description: "Practice for 30 minutes this week",
			Type:        "weekly",
			Target:      1800, // 30 minutes in seconds
			Current:     0,
			Unit:        "seconds",
			Icon:        "⏱️",
			StartDate:   startOfWeek,
			EndDate:     endOfWeek,
			Active:      true,
		},
		{
			ID:          "weekly_streak_7",
			Name:        "Week Warrior",
			Description: "Maintain 7-day streak",
			Type:        "weekly",
			Target:      7,
			Current:     0,
			Unit:        "days",
			Icon:        "🔥",
			StartDate:   startOfWeek,
			EndDate:     endOfWeek,
			Active:      true,
		},
	}
}

// UpdateGoalProgress updates goal progress based on test result
func UpdateGoalProgress(goal *Goal, value float64) {
	if !goal.Active || goal.Completed {
		return
	}

	// For WPM and accuracy goals, update if value is higher
	if goal.Unit == "wpm" || goal.Unit == "percent" {
		if value > goal.Current {
			goal.Current = value
		}
	} else {
		// For cumulative goals (tests, time), add to current
		goal.Current += value
	}

	// Check if goal is completed
	if goal.Current >= goal.Target {
		goal.Completed = true
	}
}

// IsExpired checks if a goal has expired
func IsExpired(goal Goal) bool {
	return time.Now().After(goal.EndDate)
}

// ResetDailyGoals resets goals if they're from a previous day
func ResetDailyGoals(goals []Goal) []Goal {
	today := time.Now().Truncate(24 * time.Hour)
	var activeGoals []Goal

	for _, goal := range goals {
		if goal.Type == "daily" {
			if goal.StartDate.Before(today) {
				// Reset this goal
				newGoals := DefaultDailyGoals()
				for _, newGoal := range newGoals {
					if newGoal.ID == goal.ID {
						activeGoals = append(activeGoals, newGoal)
						break
					}
				}
			} else {
				activeGoals = append(activeGoals, goal)
			}
		} else {
			activeGoals = append(activeGoals, goal)
		}
	}

	return activeGoals
}

// ResetWeeklyGoals resets goals if they're from a previous week
func ResetWeeklyGoals(goals []Goal) []Goal {
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	
	var activeGoals []Goal

	for _, goal := range goals {
		if goal.Type == "weekly" {
			if goal.StartDate.Before(startOfWeek) {
				// Reset this goal
				newGoals := DefaultWeeklyGoals()
				for _, newGoal := range newGoals {
					if newGoal.ID == goal.ID {
						activeGoals = append(activeGoals, newGoal)
						break
					}
				}
			} else {
				activeGoals = append(activeGoals, goal)
			}
		} else {
			activeGoals = append(activeGoals, goal)
		}
	}

	return activeGoals
}

// GetProgress returns progress percentage
func GetProgress(goal Goal) int {
	if goal.Target == 0 {
		return 0
	}
	progress := (goal.Current * 100) / goal.Target
	if progress > 100 {
		progress = 100
	}
	return int(progress)
}

// GetCompletedCount returns number of completed goals
func GetCompletedCount(goals []Goal) int {
	count := 0
	for _, goal := range goals {
		if goal.Completed {
			count++
		}
	}
	return count
}

// ShouldResetGoal checks if a goal should be reset based on its end date
func ShouldResetGoal(goal *Goal) bool {
	if !goal.Active {
		return false
	}
	return time.Now().After(goal.EndDate)
}

// ResetDailyGoal resets a daily goal for the next day
func ResetDailyGoal(goal *Goal) Goal {
	today := time.Now().Truncate(24 * time.Hour)
	tomorrow := today.Add(24 * time.Hour)

	return Goal{
		ID:          goal.ID,
		Name:        goal.Name,
		Description: goal.Description,
		Type:        goal.Type,
		Target:      goal.Target,
		Current:     0,
		Unit:        goal.Unit,
		Icon:        goal.Icon,
		StartDate:   today,
		EndDate:     tomorrow,
		Completed:   false,
		Active:      true,
	}
}

// ResetWeeklyGoal resets a weekly goal for the next week
func ResetWeeklyGoal(goal *Goal) Goal {
	now := time.Now()
	startOfWeek := now.Truncate(24 * time.Hour)
	// Calculate days until next Monday
	daysUntilMonday := (8 - int(startOfWeek.Weekday())) % 7
	if daysUntilMonday == 0 {
		daysUntilMonday = 7
	}
	nextMonday := startOfWeek.Add(time.Duration(daysUntilMonday) * 24 * time.Hour)
	endOfWeek := nextMonday.Add(7 * 24 * time.Hour)

	return Goal{
		ID:          goal.ID,
		Name:        goal.Name,
		Description: goal.Description,
		Type:        goal.Type,
		Target:      goal.Target,
		Current:     0,
		Unit:        goal.Unit,
		Icon:        goal.Icon,
		StartDate:   nextMonday,
		EndDate:     endOfWeek,
		Completed:   false,
		Active:      true,
	}
}

// ToRecord converts a Goal to a GoalRecord for database storage
func ToRecord(goal Goal) storage.GoalRecord {
	return storage.GoalRecord{
		ID:        goal.ID,
		Type:      goal.Type,
		Target:    goal.Target,
		Current:   goal.Current,
		Completed: goal.Completed,
		StartDate: goal.StartDate,
		EndDate:   goal.EndDate,
	}
}

// FromRecord converts a GoalRecord to a Goal
func FromRecord(record storage.GoalRecord, goals []Goal) Goal {
	// Find the matching goal template
	for _, g := range goals {
		if g.ID == record.ID {
			return Goal{
				ID:          record.ID,
				Name:        g.Name,
				Description: g.Description,
				Type:        record.Type,
				Target:      record.Target,
				Current:     record.Current,
				Unit:        g.Unit,
				Icon:        g.Icon,
				StartDate:   record.StartDate,
				EndDate:     record.EndDate,
				Completed:   record.Completed,
				Active:      true,
			}
		}
	}
	// Return a basic goal if not found
	return Goal{
		ID:        record.ID,
		Type:      record.Type,
		Target:    record.Target,
		Current:   record.Current,
		StartDate: record.StartDate,
		EndDate:   record.EndDate,
		Completed: record.Completed,
		Active:    true,
	}
}

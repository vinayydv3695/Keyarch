package achievements

import (
	"time"
)

// Achievement represents a typing achievement/badge
type Achievement struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Icon        string    `json:"icon"`
	Category    string    `json:"category"` // speed, accuracy, consistency, dedication
	Unlocked    bool      `json:"unlocked"`
	UnlockedAt  time.Time `json:"unlocked_at,omitempty"`
	Progress    int       `json:"progress"` // current progress
	Target      int       `json:"target"`   // target to unlock
	Hidden      bool      `json:"hidden"`   // hide until unlocked
}

// AllAchievements returns all available achievements
func AllAchievements() []Achievement {
	return []Achievement{
		// Speed Achievements
		{
			ID:          "speed_rookie",
			Name:        "Speed Rookie",
			Description: "Reach 40 WPM",
			Icon:        "*",
			Category:    "speed",
			Target:      40,
		},
		{
			ID:          "speed_intermediate",
			Name:        "Speed Demon",
			Description: "Reach 60 WPM",
			Icon:        "**",
			Category:    "speed",
			Target:      60,
		},
		{
			ID:          "speed_advanced",
			Name:        "Lightning Fingers",
			Description: "Reach 80 WPM",
			Icon:        "***",
			Category:    "speed",
			Target:      80,
		},
		{
			ID:          "speed_expert",
			Name:        "Typing Master",
			Description: "Reach 100 WPM",
			Icon:        "****",
			Category:    "speed",
			Target:      100,
		},
		{
			ID:          "speed_godlike",
			Name:        "Keyboard God",
			Description: "Reach 120 WPM",
			Icon:        "*****",
			Category:    "speed",
			Target:      120,
			Hidden:      true,
		},

		// Accuracy Achievements
		{
			ID:          "accuracy_good",
			Name:        "Accuracy Focused",
			Description: "Achieve 95% accuracy",
			Icon:        "+",
			Category:    "accuracy",
			Target:      95,
		},
		{
			ID:          "accuracy_great",
			Name:        "Precision Typist",
			Description: "Achieve 98% accuracy",
			Icon:        "++",
			Category:    "accuracy",
			Target:      98,
		},
		{
			ID:          "accuracy_perfect",
			Name:        "Flawless",
			Description: "Achieve 100% accuracy",
			Icon:        "+++",
			Category:    "accuracy",
			Target:      100,
		},

		// Dedication Achievements
		{
			ID:          "tests_10",
			Name:        "Getting Started",
			Description: "Complete 10 tests",
			Icon:        "#",
			Category:    "dedication",
			Target:      10,
		},
		{
			ID:          "tests_50",
			Name:        "Committed",
			Description: "Complete 50 tests",
			Icon:        "##",
			Category:    "dedication",
			Target:      50,
		},
		{
			ID:          "tests_100",
			Name:        "Century Club",
			Description: "Complete 100 tests",
			Icon:        "###",
			Category:    "dedication",
			Target:      100,
		},
		{
			ID:          "tests_500",
			Name:        "Dedicated Typist",
			Description: "Complete 500 tests",
			Icon:        "####",
			Category:    "dedication",
			Target:      500,
			Hidden:      true,
		},
		{
			ID:          "tests_1000",
			Name:        "Legend",
			Description: "Complete 1000 tests",
			Icon:        "#####",
			Category:    "dedication",
			Target:      1000,
			Hidden:      true,
		},

		// Streak Achievements
		{
			ID:          "streak_3",
			Name:        "On a Roll",
			Description: "3-day streak",
			Icon:        "~",
			Category:    "consistency",
			Target:      3,
		},
		{
			ID:          "streak_7",
			Name:        "Week Warrior",
			Description: "7-day streak",
			Icon:        "~~",
			Category:    "consistency",
			Target:      7,
		},
		{
			ID:          "streak_30",
			Name:        "Monthly Master",
			Description: "30-day streak",
			Icon:        "~~~",
			Category:    "consistency",
			Target:      30,
		},
		{
			ID:          "streak_100",
			Name:        "Unstoppable",
			Description: "100-day streak",
			Icon:        "~~~~",
			Category:    "consistency",
			Target:      100,
			Hidden:      true,
		},

		// Mode Achievements
		{
			ID:          "mode_all",
			Name:        "Jack of All Modes",
			Description: "Complete tests in all modes",
			Icon:        "@",
			Category:    "dedication",
			Target:      4, // timer, words, quote, code
		},
		{
			ID:          "code_master",
			Name:        "Code Master",
			Description: "Complete 20 code tests",
			Icon:        "</>",
			Category:    "dedication",
			Target:      20,
		},
		{
			ID:          "quote_lover",
			Name:        "Quote Enthusiast",
			Description: "Complete 20 quote tests",
			Icon:        "\"",
			Category:    "dedication",
			Target:      20,
		},

		// Time Achievements
		{
			ID:          "time_1hour",
			Name:        "Hour Hand",
			Description: "Type for 1 hour total",
			Icon:        "[1h]",
			Category:    "dedication",
			Target:      3600, // seconds
		},
		{
			ID:          "time_10hours",
			Name:        "Practice Makes Perfect",
			Description: "Type for 10 hours total",
			Icon:        "[10h]",
			Category:    "dedication",
			Target:      36000,
			Hidden:      true,
		},
	}
}

// CheckAchievement checks if an achievement should be unlocked based on stats
func CheckAchievement(achievement Achievement, stats map[string]interface{}) (bool, int) {
	switch achievement.Category {
	case "speed":
		if bestWPM, ok := stats["best_wpm"].(float64); ok {
			progress := int(bestWPM)
			return bestWPM >= float64(achievement.Target), progress
		}

	case "accuracy":
		if bestAccuracy, ok := stats["best_accuracy"].(float64); ok {
			progress := int(bestAccuracy)
			return bestAccuracy >= float64(achievement.Target), progress
		}

	case "dedication":
		// Check test count
		if achievement.ID == "tests_10" || achievement.ID == "tests_50" ||
			achievement.ID == "tests_100" || achievement.ID == "tests_500" ||
			achievement.ID == "tests_1000" {
			if totalTests, ok := stats["total_tests"].(int); ok {
				return totalTests >= achievement.Target, totalTests
			}
		}

		// Check mode diversity
		if achievement.ID == "mode_all" {
			if modesCompleted, ok := stats["modes_completed"].(int); ok {
				return modesCompleted >= achievement.Target, modesCompleted
			}
		}

		// Check code tests
		if achievement.ID == "code_master" {
			if codeTests, ok := stats["code_tests"].(int); ok {
				return codeTests >= achievement.Target, codeTests
			}
		}

		// Check quote tests
		if achievement.ID == "quote_lover" {
			if quoteTests, ok := stats["quote_tests"].(int); ok {
				return quoteTests >= achievement.Target, quoteTests
			}
		}

		// Check time
		if achievement.ID == "time_1hour" || achievement.ID == "time_10hours" {
			if totalTime, ok := stats["total_time"].(int); ok {
				return totalTime >= achievement.Target, totalTime
			}
		}

	case "consistency":
		if currentStreak, ok := stats["current_streak"].(int); ok {
			return currentStreak >= achievement.Target, currentStreak
		}
	}

	return false, 0
}

// GetUnlockedCount returns the number of unlocked achievements
func GetUnlockedCount(achievements []Achievement) int {
	count := 0
	for _, a := range achievements {
		if a.Unlocked {
			count++
		}
	}
	return count
}

// GetByCategory returns achievements filtered by category
func GetByCategory(achievements []Achievement, category string) []Achievement {
	var filtered []Achievement
	for _, a := range achievements {
		if a.Category == category {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

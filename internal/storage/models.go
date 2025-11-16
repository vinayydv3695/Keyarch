package storage

import "time"

// TestResult represents a completed typing test
type TestResult struct {
	ID          int64     `json:"id"`
	WPM         float64   `json:"wpm"`
	CPM         float64   `json:"cpm"`
	Accuracy    float64   `json:"accuracy"`
	Mistakes    int       `json:"mistakes"`
	TestType    string    `json:"test_type"` // timer, words, code, quote
	TestMode    string    `json:"test_mode"` // 15s, 30s, 60s, 25w, 50w, go, js, etc.
	Duration    int       `json:"duration"`  // in seconds
	CreatedAt   time.Time `json:"created_at"`
	TotalChars  int       `json:"total_chars"`
	CorrectChar int       `json:"correct_chars"`
}

// KeyStats represents statistics for individual keys
type KeyStats struct {
	Key        string `json:"key"`
	Correct    int    `json:"correct"`
	Incorrect  int    `json:"incorrect"`
	Accuracy   float64 `json:"accuracy"`
}

// UserStats represents aggregated user statistics
type UserStats struct {
	TotalTests      int       `json:"total_tests"`
	BestWPM         float64   `json:"best_wpm"`
	AverageWPM      float64   `json:"average_wpm"`
	BestAccuracy    float64   `json:"best_accuracy"`
	AverageAccuracy float64   `json:"average_accuracy"`
	TotalTime       int       `json:"total_time"` // in seconds
	CurrentStreak   int       `json:"current_streak"`
	LastTestDate    time.Time `json:"last_test_date"`
}

// AchievementRecord represents an unlocked achievement
type AchievementRecord struct {
	ID         string    `json:"id"`
	UnlockedAt time.Time `json:"unlocked_at"`
}

// GoalRecord represents a goal and its progress
type GoalRecord struct {
	ID        string    `json:"id"`
	Type      string    `json:"type"`
	Target    float64   `json:"target"`
	Current   float64   `json:"current"`
	Completed bool      `json:"completed"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

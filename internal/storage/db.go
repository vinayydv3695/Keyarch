package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn *sql.DB
}

// New creates a new database connection
func New() (*DB, error) {
	// Create data directory in user's home
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, ".keyarch")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dbPath := filepath.Join(dataDir, "keyarch.db")
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.initialize(); err != nil {
		conn.Close()
		return nil, err
	}

	return db, nil
}

// initialize creates the necessary tables
func (db *DB) initialize() error {
	schema := `
	CREATE TABLE IF NOT EXISTS test_results (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		wpm REAL NOT NULL,
		cpm REAL NOT NULL,
		accuracy REAL NOT NULL,
		mistakes INTEGER NOT NULL,
		test_type TEXT NOT NULL,
		test_mode TEXT NOT NULL,
		duration INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		total_chars INTEGER NOT NULL,
		correct_chars INTEGER NOT NULL
	);

	CREATE TABLE IF NOT EXISTS key_stats (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		test_id INTEGER NOT NULL,
		key TEXT NOT NULL,
		correct INTEGER DEFAULT 0,
		incorrect INTEGER DEFAULT 0,
		FOREIGN KEY (test_id) REFERENCES test_results(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS achievements (
		id TEXT PRIMARY KEY,
		unlocked_at DATETIME NOT NULL
	);

	CREATE TABLE IF NOT EXISTS goals (
		id TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		target INTEGER NOT NULL,
		current INTEGER NOT NULL,
		completed BOOLEAN DEFAULT 0,
		start_date DATETIME NOT NULL,
		end_date DATETIME NOT NULL
	);

	CREATE INDEX IF NOT EXISTS idx_created_at ON test_results(created_at);
	CREATE INDEX IF NOT EXISTS idx_test_type ON test_results(test_type);
	`

	_, err := db.conn.Exec(schema)
	return err
}

// SaveResult saves a test result to the database
func (db *DB) SaveResult(result *TestResult) error {
	query := `
	INSERT INTO test_results 
	(wpm, cpm, accuracy, mistakes, test_type, test_mode, duration, total_chars, correct_chars, created_at)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	_, err := db.conn.Exec(query,
		result.WPM,
		result.CPM,
		result.Accuracy,
		result.Mistakes,
		result.TestType,
		result.TestMode,
		result.Duration,
		result.TotalChars,
		result.CorrectChar,
		result.CreatedAt,
	)

	return err
}

// GetStats returns aggregated user statistics
func (db *DB) GetStats() (*UserStats, error) {
	stats := &UserStats{}

	// Get basic stats
	query := `
	SELECT 
		COUNT(*) as total_tests,
		COALESCE(MAX(wpm), 0) as best_wpm,
		COALESCE(AVG(wpm), 0) as average_wpm,
		COALESCE(MAX(accuracy), 0) as best_accuracy,
		COALESCE(AVG(accuracy), 0) as average_accuracy,
		COALESCE(SUM(duration), 0) as total_time,
		MAX(created_at) as last_test_date
	FROM test_results
	`

	var lastTestDateStr sql.NullString
	err := db.conn.QueryRow(query).Scan(
		&stats.TotalTests,
		&stats.BestWPM,
		&stats.AverageWPM,
		&stats.BestAccuracy,
		&stats.AverageAccuracy,
		&stats.TotalTime,
		&lastTestDateStr,
	)

	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	// Parse the date string if it's valid
	if lastTestDateStr.Valid && lastTestDateStr.String != "" {
		// Try parsing with different formats
		for _, layout := range []string{
			"2006-01-02 15:04:05",
			"2006-01-02T15:04:05Z",
			time.RFC3339,
		} {
			if parsedTime, err := time.Parse(layout, lastTestDateStr.String); err == nil {
				stats.LastTestDate = parsedTime
				stats.CurrentStreak = db.calculateStreak()
				break
			}
		}
	}

	return stats, nil
}

// GetRecentResults returns the last N test results
func (db *DB) GetRecentResults(limit int) ([]TestResult, error) {
	query := `
	SELECT id, wpm, cpm, accuracy, mistakes, test_type, test_mode, 
	       duration, created_at, total_chars, correct_chars
	FROM test_results
	ORDER BY created_at DESC
	LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []TestResult
	for rows.Next() {
		var r TestResult
		err := rows.Scan(
			&r.ID,
			&r.WPM,
			&r.CPM,
			&r.Accuracy,
			&r.Mistakes,
			&r.TestType,
			&r.TestMode,
			&r.Duration,
			&r.CreatedAt,
			&r.TotalChars,
			&r.CorrectChar,
		)
		if err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

// GetLast7DaysWPM returns WPM data for the last 7 days
func (db *DB) GetLast7DaysWPM() (map[string]float64, error) {
	query := `
	SELECT DATE(created_at) as date, AVG(wpm) as avg_wpm
	FROM test_results
	WHERE created_at >= datetime('now', '-7 days')
	GROUP BY DATE(created_at)
	ORDER BY date
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	data := make(map[string]float64)
	for rows.Next() {
		var date sql.NullString
		var avgWPM float64
		if err := rows.Scan(&date, &avgWPM); err != nil {
			return nil, err
		}
		if date.Valid {
			data[date.String] = avgWPM
		}
	}

	return data, nil
}

// calculateStreak calculates the current daily streak
func (db *DB) calculateStreak() int {
	query := `
	SELECT DATE(created_at) as test_date
	FROM test_results
	ORDER BY created_at DESC
	LIMIT 365
	`

	rows, err := db.conn.Query(query)
	if err != nil {
		return 0
	}
	defer rows.Close()

	var dates []time.Time
	for rows.Next() {
		var dateStr sql.NullString
		if err := rows.Scan(&dateStr); err != nil {
			return 0
		}
		if !dateStr.Valid {
			continue
		}
		date, err := time.Parse("2006-01-02", dateStr.String)
		if err != nil {
			return 0
		}
		dates = append(dates, date)
	}

	if len(dates) == 0 {
		return 0
	}

	streak := 1
	today := time.Now().Truncate(24 * time.Hour)
	lastDate := dates[0].Truncate(24 * time.Hour)

	// Check if the last test was today or yesterday
	daysSince := int(today.Sub(lastDate).Hours() / 24)
	if daysSince > 1 {
		return 0
	}

	for i := 1; i < len(dates); i++ {
		currentDate := dates[i].Truncate(24 * time.Hour)
		diff := int(lastDate.Sub(currentDate).Hours() / 24)

		if diff == 1 {
			streak++
			lastDate = currentDate
		} else if diff > 1 {
			break
		}
	}

	return streak
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.conn.Close()
}

// SaveAchievement saves an unlocked achievement
func (db *DB) SaveAchievement(achievementID string) error {
	query := `INSERT OR IGNORE INTO achievements (id, unlocked_at) VALUES (?, ?)`
	_, err := db.conn.Exec(query, achievementID, time.Now())
	return err
}

// GetUnlockedAchievements returns all unlocked achievement IDs
func (db *DB) GetUnlockedAchievements() ([]string, error) {
	query := `SELECT id FROM achievements ORDER BY unlocked_at DESC`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var achievements []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		achievements = append(achievements, id)
	}

	return achievements, nil
}

// SaveGoal saves or updates a goal
func (db *DB) SaveGoal(goal GoalRecord) error {
	query := `
	INSERT OR REPLACE INTO goals (id, type, target, current, completed, start_date, end_date)
	VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := db.conn.Exec(query, goal.ID, goal.Type, goal.Target, goal.Current,
		goal.Completed, goal.StartDate, goal.EndDate)
	return err
}

// GetGoals returns all active goals
func (db *DB) GetGoals() ([]GoalRecord, error) {
	query := `SELECT id, type, target, current, completed, start_date, end_date FROM goals`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goals []GoalRecord
	for rows.Next() {
		var goal GoalRecord
		var startDateStr, endDateStr string
		if err := rows.Scan(&goal.ID, &goal.Type, &goal.Target, &goal.Current,
			&goal.Completed, &startDateStr, &endDateStr); err != nil {
			return nil, err
		}

		// Parse dates
		for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", time.RFC3339} {
			if t, err := time.Parse(layout, startDateStr); err == nil {
				goal.StartDate = t
				break
			}
		}
		for _, layout := range []string{"2006-01-02 15:04:05", "2006-01-02T15:04:05Z", time.RFC3339} {
			if t, err := time.Parse(layout, endDateStr); err == nil {
				goal.EndDate = t
				break
			}
		}

		goals = append(goals, goal)
	}

	return goals, nil
}

// GetModeStats returns test counts by mode type
func (db *DB) GetModeStats() (map[string]int, error) {
	query := `
	SELECT test_type, COUNT(*) as count
	FROM test_results
	GROUP BY test_type
	`
	rows, err := db.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]int)
	for rows.Next() {
		var testType string
		var count int
		if err := rows.Scan(&testType, &count); err != nil {
			return nil, err
		}
		stats[testType] = count
	}

	return stats, nil
}

// SaveKeyStats saves key statistics from a test
func (db *DB) SaveKeyStats(testID int64, keyStats map[rune]struct{ Correct, Incorrect int }) error {
	query := `INSERT INTO key_stats (test_id, key, correct, incorrect) VALUES (?, ?, ?, ?)`

	for key, stat := range keyStats {
		_, err := db.conn.Exec(query, testID, string(key), stat.Correct, stat.Incorrect)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetWeakKeys returns the keys with lowest accuracy across all tests
func (db *DB) GetWeakKeys(limit int) ([]WeakKeyRecord, error) {
	query := `
	SELECT 
		key,
		SUM(correct) as total_correct,
		SUM(incorrect) as total_incorrect,
		SUM(correct) + SUM(incorrect) as total,
		CAST(SUM(correct) AS REAL) / (SUM(correct) + SUM(incorrect)) * 100 as accuracy
	FROM key_stats
	GROUP BY key
	HAVING total >= 10
	ORDER BY accuracy ASC
	LIMIT ?
	`

	rows, err := db.conn.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var keys []WeakKeyRecord
	for rows.Next() {
		var k WeakKeyRecord
		if err := rows.Scan(&k.Key, &k.Correct, &k.Incorrect, &k.Total, &k.Accuracy); err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}

	return keys, nil
}

// GetProgressHistory returns daily progress history for the last N days
func (db *DB) GetProgressHistory(days int) ([]ProgressHistory, error) {
	query := `
	SELECT 
		DATE(created_at) as date,
		AVG(wpm) as avg_wpm,
		AVG(accuracy) as avg_accuracy,
		COUNT(*) as test_count
	FROM test_results
	WHERE created_at >= datetime('now', '-' || ? || ' days')
	GROUP BY DATE(created_at)
	ORDER BY date ASC
	`

	rows, err := db.conn.Query(query, days)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var history []ProgressHistory
	for rows.Next() {
		var h ProgressHistory
		var dateStr string
		if err := rows.Scan(&dateStr, &h.AvgWPM, &h.AvgAccuracy, &h.TestCount); err != nil {
			return nil, err
		}
		// Parse date
		for _, layout := range []string{"2006-01-02", "2006-01-02 15:04:05"} {
			if t, err := time.Parse(layout, dateStr); err == nil {
				h.Date = t
				break
			}
		}
		history = append(history, h)
	}

	return history, nil
}

// GetLastTestID returns the ID of the most recently inserted test
func (db *DB) GetLastTestID() (int64, error) {
	var id int64
	err := db.conn.QueryRow("SELECT last_insert_rowid()").Scan(&id)
	return id, err
}

// SaveResultWithKeys saves a test result and its key statistics
func (db *DB) SaveResultWithKeys(result *TestResult, keyStats map[rune]struct{ Correct, Incorrect int }) error {
	// Save the test result
	if err := db.SaveResult(result); err != nil {
		return err
	}

	// Get the last inserted ID
	testID, err := db.GetLastTestID()
	if err != nil {
		return err
	}

	// Save key stats
	return db.SaveKeyStats(testID, keyStats)
}

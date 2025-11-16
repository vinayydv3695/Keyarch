package engine

import (
	"time"
)

// TestMode represents the type of typing test
type TestMode string

const (
	ModeTimer   TestMode = "timer"
	ModeWords   TestMode = "words"
	ModeQuote   TestMode = "quote"
	ModeCode    TestMode = "code"
)

// Engine manages the typing test logic
type Engine struct {
	TargetText     string
	UserInput      string
	StartTime      time.Time
	EndTime        time.Time
	IsStarted      bool
	IsFinished     bool
	Mode           TestMode
	Duration       int // for timer mode (seconds)
	WordCount      int // for word mode
	CurrentPos     int
	Mistakes       int
	MistakeMap     map[int]bool // tracks position of mistakes
	KeyStrokes     map[rune]KeyStat
	WPMHistory     []float64 // WPM samples over time
	LastWPMUpdate  time.Time
}

// KeyStat tracks statistics for individual keys
type KeyStat struct {
	Correct   int
	Incorrect int
}

// WeakKey represents a key with its accuracy
type WeakKey struct {
	Key      rune
	Accuracy float64
}

// New creates a new typing engine
func New(targetText string, mode TestMode, duration int, wordCount int) *Engine {
	return &Engine{
		TargetText:    targetText,
		Mode:          mode,
		Duration:      duration,
		WordCount:     wordCount,
		MistakeMap:    make(map[int]bool),
		KeyStrokes:    make(map[rune]KeyStat),
		WPMHistory:    make([]float64, 0),
		LastWPMUpdate: time.Now(),
	}
}

// Start begins the typing test
func (e *Engine) Start() {
	e.IsStarted = true
	e.StartTime = time.Now()
}

// ProcessInput processes a character input
func (e *Engine) ProcessInput(char rune) {
	if !e.IsStarted {
		e.Start()
	}

	if e.IsFinished {
		return
	}

	// Handle backspace
	if char == 127 || char == 8 {
		if len(e.UserInput) > 0 {
			e.UserInput = e.UserInput[:len(e.UserInput)-1]
			e.CurrentPos = len(e.UserInput)
		}
		return
	}

	// Add character
	e.UserInput += string(char)
	
	// Track if it's correct or incorrect
	if e.CurrentPos < len(e.TargetText) {
		expected := rune(e.TargetText[e.CurrentPos])
		stat := e.KeyStrokes[expected]
		
		if char == expected {
			stat.Correct++
		} else {
			stat.Incorrect++
			e.Mistakes++
			e.MistakeMap[e.CurrentPos] = true
		}
		
		e.KeyStrokes[expected] = stat
	}
	
	e.CurrentPos = len(e.UserInput)

	// Check if test is complete
	if e.Mode == ModeWords && e.CountWords() >= e.WordCount {
		e.Finish()
	} else if len(e.UserInput) >= len(e.TargetText) {
		e.Finish()
	}
}

// CountWords counts words in user input
func (e *Engine) CountWords() int {
	if len(e.UserInput) == 0 {
		return 0
	}
	count := 1
	for _, char := range e.UserInput {
		if char == ' ' {
			count++
		}
	}
	return count
}

// Finish ends the typing test
func (e *Engine) Finish() {
	if e.IsFinished {
		return
	}
	e.IsFinished = true
	e.EndTime = time.Now()
}

// GetElapsedTime returns the elapsed time in seconds
func (e *Engine) GetElapsedTime() float64 {
	if !e.IsStarted {
		return 0
	}
	if e.IsFinished {
		return e.EndTime.Sub(e.StartTime).Seconds()
	}
	return time.Since(e.StartTime).Seconds()
}

// GetWPM calculates words per minute
func (e *Engine) GetWPM() float64 {
	elapsed := e.GetElapsedTime()
	if elapsed == 0 {
		return 0
	}

	// Standard: 5 characters = 1 word
	words := float64(len(e.UserInput)) / 5.0
	minutes := elapsed / 60.0

	if minutes == 0 {
		return 0
	}

	return words / minutes
}

// GetCPM calculates characters per minute
func (e *Engine) GetCPM() float64 {
	elapsed := e.GetElapsedTime()
	if elapsed == 0 {
		return 0
	}

	minutes := elapsed / 60.0
	if minutes == 0 {
		return 0
	}

	return float64(len(e.UserInput)) / minutes
}

// GetAccuracy calculates typing accuracy as a percentage
func (e *Engine) GetAccuracy() float64 {
	if len(e.UserInput) == 0 {
		return 100.0
	}

	correct := 0
	minLen := len(e.UserInput)
	if len(e.TargetText) < minLen {
		minLen = len(e.TargetText)
	}

	for i := 0; i < minLen; i++ {
		if e.UserInput[i] == e.TargetText[i] {
			correct++
		}
	}

	return (float64(correct) / float64(len(e.UserInput))) * 100.0
}

// GetProgress returns the progress as a percentage
func (e *Engine) GetProgress() float64 {
	if len(e.TargetText) == 0 {
		return 0
	}

	switch e.Mode {
	case ModeTimer:
		elapsed := e.GetElapsedTime()
		if elapsed >= float64(e.Duration) {
			return 100.0
		}
		return (elapsed / float64(e.Duration)) * 100.0
	case ModeWords:
		words := float64(e.CountWords())
		return (words / float64(e.WordCount)) * 100.0
	default:
		return (float64(len(e.UserInput)) / float64(len(e.TargetText))) * 100.0
	}
}

// IsCorrectChar checks if character at position is correct
func (e *Engine) IsCorrectChar(pos int) bool {
	if pos >= len(e.UserInput) {
		return true // not typed yet
	}
	if pos >= len(e.TargetText) {
		return false // extra characters
	}
	return e.UserInput[pos] == e.TargetText[pos]
}

// ShouldFinish checks if the test should end (for timer mode)
func (e *Engine) ShouldFinish() bool {
	if e.Mode == ModeTimer && e.IsStarted && !e.IsFinished {
		return e.GetElapsedTime() >= float64(e.Duration)
	}
	return false
}

// GetRemainingTime returns remaining time for timer mode
func (e *Engine) GetRemainingTime() int {
	if e.Mode != ModeTimer {
		return 0
	}

	elapsed := int(e.GetElapsedTime())
	remaining := e.Duration - elapsed

	if remaining < 0 {
		return 0
	}

	return remaining
}

// GetWeakKeys returns keys with lowest accuracy
func (e *Engine) GetWeakKeys(limit int) []WeakKey {
	var keys []WeakKey
	for key, stat := range e.KeyStrokes {
		total := stat.Correct + stat.Incorrect
		if total > 0 {
			accuracy := (float64(stat.Correct) / float64(total)) * 100.0
			keys = append(keys, WeakKey{Key: key, Accuracy: accuracy})
		}
	}

	// Sort by accuracy (bubble sort for simplicity)
	for i := 0; i < len(keys); i++ {
		for j := i + 1; j < len(keys); j++ {
			if keys[i].Accuracy > keys[j].Accuracy {
				keys[i], keys[j] = keys[j], keys[i]
			}
		}
	}

	if len(keys) > limit {
		keys = keys[:limit]
	}

	return keys
}

// UpdateWPMHistory samples current WPM for the live graph
func (e *Engine) UpdateWPMHistory() {
	if !e.IsStarted || e.IsFinished {
		return
	}

	// Update every 2 seconds
	if time.Since(e.LastWPMUpdate) < 2*time.Second {
		return
	}

	currentWPM := e.GetWPM()
	e.WPMHistory = append(e.WPMHistory, currentWPM)
	e.LastWPMUpdate = time.Now()

	// Keep only last 30 samples (60 seconds of data)
	if len(e.WPMHistory) > 30 {
		e.WPMHistory = e.WPMHistory[1:]
	}
}

// GetWPMHistory returns the WPM history for graphing
func (e *Engine) GetWPMHistory() []float64 {
	return e.WPMHistory
}

package heatmap

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/vinayydv3695/keyarch/internal/engine"
)

// KeyboardLayout represents the visual keyboard layout
var KeyboardLayout = [][]string{
	{"`", "1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "-", "="},
	{"Q", "W", "E", "R", "T", "Y", "U", "I", "O", "P", "[", "]", "\\"},
	{"A", "S", "D", "F", "G", "H", "J", "K", "L", ";", "'"},
	{"Z", "X", "C", "V", "B", "N", "M", ",", ".", "/"},
	{"SPACE"},
}

// GetHeatmapColor returns a color based on accuracy percentage
func GetHeatmapColor(accuracy float64) lipgloss.Color {
	switch {
	case accuracy >= 95:
		return lipgloss.Color("#a6e3a1") // Green
	case accuracy >= 85:
		return lipgloss.Color("#94e2d5") // Teal
	case accuracy >= 75:
		return lipgloss.Color("#f9e2af") // Yellow
	case accuracy >= 60:
		return lipgloss.Color("#fab387") // Peach
	case accuracy > 0:
		return lipgloss.Color("#f38ba8") // Red
	default:
		return lipgloss.Color("#6c7086") // Gray (not used)
	}
}

// RenderHeatmap creates a visual keyboard heatmap
func RenderHeatmap(keystats map[rune]engine.KeyStat) string {
	var result strings.Builder

	// Calculate accuracy for each key
	keyAccuracy := make(map[rune]float64)
	for key, stats := range keystats {
		total := stats.Correct + stats.Incorrect
		if total > 0 {
			keyAccuracy[key] = (float64(stats.Correct) / float64(total)) * 100
		}
	}

	// Render keyboard layout
	result.WriteString("\n")
	result.WriteString("  Typing Heatmap - Color coded by accuracy:\n\n")

	for _, row := range KeyboardLayout {
		result.WriteString("  ")
		for _, key := range row {
			keyRune := rune(strings.ToLower(key)[0])
			
			// Special handling for space
			if key == "SPACE" {
				spaceRune := ' '
				accuracy := keyAccuracy[spaceRune]
				color := GetHeatmapColor(accuracy)
				
				keyStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#1e1e2e")).
					Background(color).
					Width(20).
					Align(lipgloss.Center).
					Bold(true)
				
				if accuracy > 0 {
					result.WriteString(keyStyle.Render(fmt.Sprintf("SPACE %.0f%%", accuracy)))
				} else {
					result.WriteString(keyStyle.Render("SPACE"))
				}
				continue
			}

			accuracy := keyAccuracy[keyRune]
			color := GetHeatmapColor(accuracy)

			keyStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("#1e1e2e")).
				Background(color).
				Width(4).
				Align(lipgloss.Center).
				Bold(true)

			if accuracy > 0 {
				result.WriteString(keyStyle.Render(fmt.Sprintf("%s", key)))
			} else {
				// Key not used
				grayStyle := lipgloss.NewStyle().
					Foreground(lipgloss.Color("#6c7086")).
					Background(lipgloss.Color("#313244")).
					Width(4).
					Align(lipgloss.Center)
				result.WriteString(grayStyle.Render(key))
			}
			result.WriteString(" ")
		}
		result.WriteString("\n")
	}

	// Legend
	result.WriteString("\n  Legend:\n")
	
	legendItems := []struct {
		label    string
		color    lipgloss.Color
		accuracy string
	}{
		{"Excellent", lipgloss.Color("#a6e3a1"), "95-100%"},
		{"Good", lipgloss.Color("#94e2d5"), "85-94%"},
		{"Average", lipgloss.Color("#f9e2af"), "75-84%"},
		{"Below Avg", lipgloss.Color("#fab387"), "60-74%"},
		{"Weak", lipgloss.Color("#f38ba8"), "< 60%"},
		{"Unused", lipgloss.Color("#6c7086"), "Not typed"},
	}

	result.WriteString("  ")
	for _, item := range legendItems {
		boxStyle := lipgloss.NewStyle().
			Background(item.color).
			Foreground(lipgloss.Color("#1e1e2e")).
			Width(3).
			Align(lipgloss.Center).
			Bold(true)
		
		if item.label == "Unused" {
			boxStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#313244")).
				Foreground(item.color).
				Width(3).
				Align(lipgloss.Center)
		}

		result.WriteString(boxStyle.Render("  "))
		result.WriteString(" " + item.label + " (" + item.accuracy + ")  ")
	}
	result.WriteString("\n")

	return result.String()
}

// GetWeakKeysFromHeatmap returns the weakest keys based on accuracy
func GetWeakKeysFromHeatmap(keystats map[rune]engine.KeyStat, limit int) []engine.WeakKey {
	var keys []engine.WeakKey

	for key, stats := range keystats {
		total := stats.Correct + stats.Incorrect
		if total > 0 {
			accuracy := (float64(stats.Correct) / float64(total)) * 100
			keys = append(keys, engine.WeakKey{
				Key:      key,
				Accuracy: accuracy,
			})
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

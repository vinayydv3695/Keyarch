package training

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/storage"
	"github.com/vinayydv3695/keyarch/internal/tui/components"
)

// Lesson types
const (
	LessonHomeRow   = "home_row"
	LessonTopRow    = "top_row"
	LessonBottomRow = "bottom_row"
	LessonNumbers   = "numbers"
	LessonSymbols   = "symbols"
	LessonWeakKeys  = "weak_keys"
	LessonCustom    = "custom"
)

// Lesson represents a typing lesson
type Lesson struct {
	ID          string
	Name        string
	Description string
	Keys        string // The keys to practice
	Words       []string
}

// GetLessons returns all available lessons
func GetLessons() []Lesson {
	return []Lesson{
		{
			ID:          LessonHomeRow,
			Name:        "Home Row",
			Description: "Master the home row keys: asdfjkl;",
			Keys:        "asdfghjkl;",
			Words: []string{
				"sad", "dad", "add", "fad", "had", "lad", "ash", "ask",
				"fall", "hall", "shall", "lass", "lass", "glass", "flask",
				"salad", "shall", "flash", "slash", "dash", "hash", "gash",
			},
		},
		{
			ID:          LessonTopRow,
			Name:        "Top Row",
			Description: "Practice the top row: qwertyuiop",
			Keys:        "qwertyuiop",
			Words: []string{
				"quit", "quote", "quiet", "quite", "queue",
				"write", "writer", "wrote", "writ", "wry",
				"type", "typo", "trip", "trap", "true",
				"pier", "pier", "pure", "pour", "your",
			},
		},
		{
			ID:          LessonBottomRow,
			Name:        "Bottom Row",
			Description: "Practice the bottom row: zxcvbnm",
			Keys:        "zxcvbnm,./",
			Words: []string{
				"zero", "zone", "zoom", "zinc",
				"exam", "next", "text", "vex",
				"come", "some", "home", "dome",
				"name", "same", "game", "fame",
				"cab", "cub", "club", "crab",
			},
		},
		{
			ID:          LessonNumbers,
			Name:        "Numbers",
			Description: "Practice number keys: 1234567890",
			Keys:        "1234567890",
			Words: []string{
				"123", "456", "789", "100", "200", "300",
				"2024", "1999", "2000", "1234", "5678", "9012",
				"42", "007", "404", "500", "101", "303",
			},
		},
		{
			ID:          LessonSymbols,
			Name:        "Symbols",
			Description: "Practice common symbols: !@#$%^&*()",
			Keys:        "!@#$%^&*()_+-=[]{}|;':\",./<>?",
			Words: []string{
				"@email", "#hashtag", "$100", "100%", "a&b",
				"func()", "arr[0]", "{}", "->", "=>", "!=",
				"path/to/file", "key:value", "yes/no", "a+b",
			},
		},
		{
			ID:          LessonWeakKeys,
			Name:        "Weak Keys",
			Description: "Practice your weakest keys based on history",
			Keys:        "", // Will be filled dynamically
			Words:       nil,
		},
	}
}

// Model represents the training screen
type Model struct {
	styles   *components.Styles
	cfg      *config.Config
	db       *storage.DB
	lessons  []Lesson
	weakKeys []storage.WeakKeyRecord
	cursor   int
	width    int
	height   int
	selected string
	loading  bool
}

// New creates a new training model
func New(cfg *config.Config, db *storage.DB) Model {
	theme := config.GetThemeByName(cfg.Theme)
	return Model{
		styles:  components.NewStyles(theme),
		cfg:     cfg,
		db:      db,
		lessons: GetLessons(),
		loading: true,
	}
}

func (m Model) Init() tea.Cmd {
	return m.loadWeakKeys
}

func (m Model) loadWeakKeys() tea.Msg {
	if m.db == nil {
		return weakKeysLoadedMsg{keys: nil}
	}

	keys, err := m.db.GetWeakKeys(10)
	if err != nil {
		return weakKeysLoadedMsg{keys: nil}
	}
	return weakKeysLoadedMsg{keys: keys}
}

type weakKeysLoadedMsg struct {
	keys []storage.WeakKeyRecord
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case weakKeysLoadedMsg:
		m.weakKeys = msg.keys
		m.loading = false

		// Update weak keys lesson if we have data
		if len(m.weakKeys) > 0 {
			for i, lesson := range m.lessons {
				if lesson.ID == LessonWeakKeys {
					m.lessons[i].Words = generateWeakKeyWords(m.weakKeys)
					break
				}
			}
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			os.Exit(0)
			return m, tea.Quit

		case "q", "esc":
			m.selected = "back"
			return m, tea.Quit

		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.cursor < len(m.lessons) {
				m.cursor++
			}

		case "enter", " ":
			if m.cursor == len(m.lessons) {
				m.selected = "back"
			} else {
				m.selected = m.lessons[m.cursor].ID
			}
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m Model) View() string {
	if m.width == 0 || m.loading {
		return "Loading training lessons..."
	}

	s := components.HeaderWithWidth("TRAINING MODE", "Practice specific keys and improve weak areas", m.styles, m.width)
	s += "\n\n"

	// Lessons list
	for i, lesson := range m.lessons {
		name := lesson.Name
		desc := lesson.Description

		// Special handling for weak keys
		if lesson.ID == LessonWeakKeys {
			if len(m.weakKeys) > 0 {
				weakKeyChars := ""
				for j, wk := range m.weakKeys {
					if j > 5 {
						break
					}
					if j > 0 {
						weakKeyChars += ", "
					}
					weakKeyChars += fmt.Sprintf("'%s' (%.0f%%)", wk.Key, wk.Accuracy)
				}
				desc = "Focus on: " + weakKeyChars
			} else {
				desc = "Complete more tests to identify weak keys"
			}
		}

		if i == m.cursor {
			s += m.styles.ActiveItem.Render(fmt.Sprintf("> %s", name)) + "\n"
			s += "  " + m.styles.Muted.Render(desc) + "\n\n"
		} else {
			s += m.styles.MenuItem.Render(fmt.Sprintf("  %s", name)) + "\n"
			s += "  " + m.styles.Muted.Render(desc) + "\n\n"
		}
	}

	// Back option
	if m.cursor == len(m.lessons) {
		s += m.styles.ActiveItem.Render("> Back to Menu") + "\n"
	} else {
		s += m.styles.MenuItem.Render("  Back to Menu") + "\n"
	}

	s += "\n"
	s += components.Footer("Up/Down: Navigate | Enter: Select | ESC/q: Back", m.styles)

	return s
}

// Selected returns the selected lesson ID or "back"
func (m Model) Selected() string {
	return m.selected
}

// GetSelectedLesson returns the selected lesson
func (m Model) GetSelectedLesson() *Lesson {
	for _, lesson := range m.lessons {
		if lesson.ID == m.selected {
			return &lesson
		}
	}
	return nil
}

// GenerateLessonText generates practice text for a lesson
func GenerateLessonText(lesson *Lesson, wordCount int) string {
	if len(lesson.Words) == 0 {
		return generateKeyPractice(lesson.Keys, wordCount)
	}

	var result []string
	for i := 0; i < wordCount; i++ {
		word := lesson.Words[rand.Intn(len(lesson.Words))]
		result = append(result, word)
	}

	return strings.Join(result, " ")
}

// generateKeyPractice generates practice text from specific keys
func generateKeyPractice(keys string, wordCount int) string {
	if len(keys) == 0 {
		return ""
	}

	runes := []rune(keys)
	var result []string

	for i := 0; i < wordCount; i++ {
		// Generate a "word" of 3-6 characters from the key set
		wordLen := 3 + rand.Intn(4)
		word := ""
		for j := 0; j < wordLen; j++ {
			word += string(runes[rand.Intn(len(runes))])
		}
		result = append(result, word)
	}

	return strings.Join(result, " ")
}

// generateWeakKeyWords generates practice words focusing on weak keys
func generateWeakKeyWords(weakKeys []storage.WeakKeyRecord) []string {
	if len(weakKeys) == 0 {
		return nil
	}

	// Common words containing each weak key
	keyWordMap := map[rune][]string{
		'q': {"quick", "quite", "quote", "queen", "quest", "equal"},
		'w': {"with", "what", "when", "where", "which", "would", "world", "work", "write"},
		'e': {"the", "be", "then", "them", "here", "there", "every", "even", "never"},
		'r': {"are", "for", "from", "or", "more", "your", "their", "first", "great"},
		't': {"that", "this", "to", "not", "but", "at", "it", "out", "just", "about"},
		'y': {"you", "your", "they", "by", "my", "any", "very", "only", "year"},
		'u': {"you", "use", "out", "up", "but", "just", "us", "must", "could"},
		'i': {"in", "is", "it", "if", "with", "will", "this", "which", "think"},
		'o': {"of", "on", "or", "to", "not", "from", "do", "so", "no", "who"},
		'p': {"up", "people", "part", "place", "point", "help", "keep", "stop"},
		'a': {"a", "and", "as", "at", "all", "are", "an", "have", "had", "that"},
		's': {"so", "she", "say", "see", "some", "same", "also", "just", "us", "as"},
		'd': {"do", "did", "had", "would", "could", "should", "and", "old", "day"},
		'f': {"for", "from", "if", "of", "off", "after", "first", "find", "few"},
		'g': {"go", "get", "got", "give", "good", "great", "going", "again"},
		'h': {"he", "his", "have", "had", "has", "her", "how", "here", "him"},
		'j': {"just", "job", "join", "jump", "judge", "major", "enjoy"},
		'k': {"know", "like", "look", "make", "take", "think", "work", "back"},
		'l': {"like", "all", "will", "well", "also", "only", "last", "still"},
		';': {"don;t", "can;t", "it;s", "that;s", "what;s", "there;s"},
		'z': {"zero", "zone", "size", "realize", "organize", "amazing"},
		'x': {"next", "example", "exactly", "except", "extra", "six"},
		'c': {"can", "come", "could", "case", "each", "such", "much", "back"},
		'v': {"very", "have", "over", "even", "every", "never", "five"},
		'b': {"be", "but", "by", "been", "both", "about", "before", "between"},
		'n': {"no", "not", "now", "on", "in", "and", "one", "new", "than"},
		'm': {"my", "me", "may", "make", "made", "more", "most", "many", "time"},
	}

	var words []string
	for _, wk := range weakKeys {
		if len(wk.Key) > 0 {
			r := rune(wk.Key[0])
			if keyWords, ok := keyWordMap[r]; ok {
				words = append(words, keyWords...)
			}
		}
	}

	if len(words) == 0 {
		return []string{"the", "and", "for", "are", "but", "not", "you", "all"}
	}

	return words
}

package data

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// CustomWordList represents a user-defined word list
type CustomWordList struct {
	Name  string
	Words []string
	Path  string
}

// GetWordListsDir returns the directory where custom word lists are stored
func GetWordListsDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(homeDir, ".keyarch", "wordlists"), nil
}

// EnsureWordListsDir creates the wordlists directory if it doesn't exist
func EnsureWordListsDir() error {
	dir, err := GetWordListsDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(dir, 0755)
}

// ListCustomWordLists returns all available custom word lists
func ListCustomWordLists() ([]CustomWordList, error) {
	dir, err := GetWordListsDir()
	if err != nil {
		return nil, err
	}

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var lists []CustomWordList
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		// Accept .txt and .words files
		if !strings.HasSuffix(name, ".txt") && !strings.HasSuffix(name, ".words") {
			continue
		}

		listPath := filepath.Join(dir, name)
		words, err := loadWordsFromFile(listPath)
		if err != nil {
			continue // Skip invalid files
		}

		// Get name without extension
		baseName := strings.TrimSuffix(name, filepath.Ext(name))

		lists = append(lists, CustomWordList{
			Name:  baseName,
			Words: words,
			Path:  listPath,
		})
	}

	return lists, nil
}

// LoadCustomWordList loads a custom word list by name
func LoadCustomWordList(name string) (*CustomWordList, error) {
	dir, err := GetWordListsDir()
	if err != nil {
		return nil, err
	}

	// Try .txt first, then .words
	for _, ext := range []string{".txt", ".words"} {
		path := filepath.Join(dir, name+ext)
		if _, err := os.Stat(path); err == nil {
			words, err := loadWordsFromFile(path)
			if err != nil {
				return nil, err
			}
			return &CustomWordList{
				Name:  name,
				Words: words,
				Path:  path,
			}, nil
		}
	}

	// Try the name as-is (full path)
	if _, err := os.Stat(name); err == nil {
		words, err := loadWordsFromFile(name)
		if err != nil {
			return nil, err
		}
		baseName := strings.TrimSuffix(filepath.Base(name), filepath.Ext(name))
		return &CustomWordList{
			Name:  baseName,
			Words: words,
			Path:  name,
		}, nil
	}

	return nil, os.ErrNotExist
}

// loadWordsFromFile reads words from a file
// Supports multiple formats:
// - One word per line
// - Space-separated words
// - Comma-separated words
func loadWordsFromFile(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var words []string
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue // Skip empty lines and comments
		}

		// Check if line contains commas (CSV format)
		if strings.Contains(line, ",") {
			parts := strings.Split(line, ",")
			for _, part := range parts {
				word := strings.TrimSpace(part)
				if word != "" {
					words = append(words, word)
				}
			}
		} else if strings.Contains(line, " ") {
			// Space-separated words on a single line
			parts := strings.Fields(line)
			words = append(words, parts...)
		} else {
			// Single word per line
			words = append(words, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return words, nil
}

// CreateSampleWordList creates a sample word list file for users
func CreateSampleWordList() error {
	dir, err := GetWordListsDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	samplePath := filepath.Join(dir, "sample.txt")
	if _, err := os.Stat(samplePath); err == nil {
		return nil // Already exists
	}

	content := `# Sample Custom Word List for Keyarch
# Lines starting with # are comments
# You can use one word per line:
hello
world
typing
practice

# Or space-separated words on a single line:
keyboard keys fingers hands

# Or comma-separated:
fast, accurate, consistent, smooth

# Create your own word lists by adding .txt or .words files
# to this directory: ~/.keyarch/wordlists/
`

	return os.WriteFile(samplePath, []byte(content), 0644)
}

// GeneratorWithCustom extends Generator to support custom word lists
type GeneratorWithCustom struct {
	*Generator
	customWords []string
	useCustom   bool
}

// NewGeneratorWithCustom creates a generator that can use custom word lists
func NewGeneratorWithCustom(customListName string) *GeneratorWithCustom {
	gen := &GeneratorWithCustom{
		Generator: NewGenerator(),
	}

	if customListName != "" {
		list, err := LoadCustomWordList(customListName)
		if err == nil && len(list.Words) > 0 {
			gen.customWords = list.Words
			gen.useCustom = true
		}
	}

	return gen
}

// GenerateWords generates words using custom list if available
func (g *GeneratorWithCustom) GenerateWords(count int, difficulty string) string {
	if g.useCustom && len(g.customWords) > 0 {
		return g.generateFromCustom(count)
	}
	return g.Generator.GenerateWords(count, difficulty)
}

// generateFromCustom generates words from the custom list
func (g *GeneratorWithCustom) generateFromCustom(count int) string {
	if len(g.customWords) == 0 {
		return g.Generator.GenerateWords(count, "medium")
	}

	var result []string
	for i := 0; i < count; i++ {
		word := g.customWords[randInt(len(g.customWords))]
		result = append(result, word)
	}

	return strings.Join(result, " ")
}

// GenerateByTime generates text for timed tests using custom words if available
func (g *GeneratorWithCustom) GenerateByTime(seconds int, difficulty string) string {
	// Generate enough words for the time
	estimatedWords := ((seconds * 120) / 60) + ((seconds * 120) / 60 / 2)
	if estimatedWords < 100 {
		estimatedWords = 100
	}

	return g.GenerateWords(estimatedWords, difficulty)
}

// randInt returns a random integer in [0, n)
func randInt(n int) int {
	if n <= 0 {
		return 0
	}
	// Use a simple approach for randomness
	return int(uint64(randSeed()) % uint64(n))
}

var seed uint64 = uint64(os.Getpid())

func randSeed() uint64 {
	seed = seed*1103515245 + 12345
	return seed
}

// HasCustomList returns true if custom words are being used
func (g *GeneratorWithCustom) HasCustomList() bool {
	return g.useCustom
}

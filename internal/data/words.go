package data

import (
	"math/rand"
	"strings"
)

// WordList contains word lists of different difficulties
type WordList struct {
	Easy   []string
	Medium []string
	Hard   []string
}

// DefaultWords provides default word lists
var DefaultWords = WordList{
	Easy: []string{
		"the", "be", "to", "of", "and", "a", "in", "that", "have", "I",
		"it", "for", "not", "on", "with", "he", "as", "you", "do", "at",
		"this", "but", "his", "by", "from", "they", "we", "say", "her", "she",
		"or", "an", "will", "my", "one", "all", "would", "there", "their", "what",
		"so", "up", "out", "if", "about", "who", "get", "which", "go", "me",
		"when", "make", "can", "like", "time", "no", "just", "him", "know", "take",
		"people", "into", "year", "your", "good", "some", "could", "them", "see", "other",
		"than", "then", "now", "look", "only", "come", "its", "over", "think", "also",
		"back", "after", "use", "two", "how", "our", "work", "first", "well", "way",
		"even", "new", "want", "because", "any", "these", "give", "day", "most", "us",
	},
	Medium: []string{
		"system", "program", "question", "work", "government", "number", "night", "point",
		"world", "school", "state", "company", "problem", "service", "place", "hand",
		"party", "group", "money", "story", "fact", "month", "different", "study",
		"book", "issue", "side", "business", "information", "power", "room", "important",
		"health", "person", "level", "office", "door", "student", "education", "music",
		"family", "national", "research", "history", "market", "political", "result", "community",
		"interest", "development", "change", "reason", "action", "experience", "process", "management",
		"practice", "position", "quality", "evidence", "performance", "model", "product", "technology",
		"structure", "support", "treatment", "analysis", "activity", "benefit", "patient", "design",
		"material", "relationship", "control", "individual", "building", "culture", "behavior", "series",
	},
	Hard: []string{
		"administration", "implementation", "infrastructure", "acknowledgment", "characteristic", "recommendation",
		"responsibility", "consciousness", "undergraduate", "ということ", "sophisticated", "differentiation",
		"entrepreneurship", "pharmaceutical", "semiconductor", "archaeological", "constitutional", "refrigerator",
		"miscellaneous", "correspondence", "comprehensive", "vulnerability", "authentication", "extraordinary",
		"simultaneously", "telecommunications", "representative", "interconnected", "mediterranean", "neighborhood",
		"substantially", "investigation", "transformation", "psychological", "environmental", "establishment",
		"consideration", "participation", "revolutionary", "sophisticated", "appropriately", "complementary",
		"distinguished", "occasionally", "specification", "configuration", "functionality", "approximately",
		"documentation", "unprecedented", "controversial", "substantially", "considerably", "significantly",
		"predominantly", "fundamental", "experimental", "contemporary", "conventional", "instrumental",
	},
}

// Generator generates text for typing tests
type Generator struct {
	words WordList
}

// NewGenerator creates a new text generator
func NewGenerator() *Generator {
	return &Generator{
		words: DefaultWords,
	}
}

// GenerateWords generates a random word sequence
func (g *Generator) GenerateWords(count int, difficulty string) string {
	var wordList []string

	switch difficulty {
	case "easy":
		wordList = g.words.Easy
	case "hard":
		wordList = g.words.Hard
	default:
		wordList = g.words.Medium
	}

	if len(wordList) == 0 {
		wordList = g.words.Easy
	}

	var result []string
	for i := 0; i < count; i++ {
		word := wordList[rand.Intn(len(wordList))]
		result = append(result, word)
	}

	return strings.Join(result, " ")
}

// GenerateByTime generates text that should take approximately the specified duration
func (g *Generator) GenerateByTime(seconds int, difficulty string) string {
	// Generate enough words so it's nearly impossible to finish in the given time
	// Assume very fast typing: 120 WPM (professional level)
	// 120 WPM = 600 characters per minute = 10 characters per second
	// Add 50% buffer to ensure there's always more text than needed
	estimatedWords := ((seconds * 120) / 60) + ((seconds * 120) / 60 / 2)
	
	// Minimum 100 words for any timed test
	if estimatedWords < 100 {
		estimatedWords = 100
	}

	return g.GenerateWords(estimatedWords, difficulty)
}

// Quotes contains various quotes for typing tests
var Quotes = []Quote{
	{Text: "The only way to do great work is to love what you do.", Author: "Steve Jobs", Category: "motivation"},
	{Text: "Life is what happens when you're busy making other plans.", Author: "John Lennon", Category: "philosophy"},
	{Text: "The future belongs to those who believe in the beauty of their dreams.", Author: "Eleanor Roosevelt", Category: "motivation"},
	{Text: "It is during our darkest moments that we must focus to see the light.", Author: "Aristotle", Category: "philosophy"},
	{Text: "The purpose of our lives is to be happy.", Author: "Dalai Lama", Category: "philosophy"},
	{Text: "In three words I can sum up everything I've learned about life: it goes on.", Author: "Robert Frost", Category: "philosophy"},
	{Text: "The way to get started is to quit talking and begin doing.", Author: "Walt Disney", Category: "motivation"},
	{Text: "Don't let yesterday take up too much of today.", Author: "Will Rogers", Category: "motivation"},
	{Text: "You learn more from failure than from success. Don't let it stop you. Failure builds character.", Author: "Unknown", Category: "motivation"},
	{Text: "If you are working on something that you really care about, you don't have to be pushed. The vision pulls you.", Author: "Steve Jobs", Category: "motivation"},
	{Text: "People who are crazy enough to think they can change the world, are the ones who do.", Author: "Rob Siltanen", Category: "motivation"},
	{Text: "Knowing is not enough; we must apply. Willing is not enough; we must do.", Author: "Johann Wolfgang Von Goethe", Category: "philosophy"},
	{Text: "Whether you think you can or you think you can't, you're right.", Author: "Henry Ford", Category: "motivation"},
	{Text: "Perfection is not attainable, but if we chase perfection we can catch excellence.", Author: "Vince Lombardi", Category: "motivation"},
	{Text: "Life is 10% what happens to me and 90% of how I react to it.", Author: "Charles Swindoll", Category: "philosophy"},
	{Text: "The mind is everything. What you think you become.", Author: "Buddha", Category: "philosophy"},
	{Text: "The best time to plant a tree was 20 years ago. The second best time is now.", Author: "Chinese Proverb", Category: "philosophy"},
	{Text: "An unexamined life is not worth living.", Author: "Socrates", Category: "philosophy"},
	{Text: "Your time is limited, so don't waste it living someone else's life.", Author: "Steve Jobs", Category: "motivation"},
	{Text: "The only impossible journey is the one you never begin.", Author: "Tony Robbins", Category: "motivation"},
}

// AnimeQuotes contains anime-related quotes
var AnimeQuotes = []Quote{
	{Text: "People's lives don't end when they die. It ends when they lose faith.", Author: "Itachi Uchiha", Category: "anime"},
	{Text: "If you don't take risks, you can't create a future.", Author: "Monkey D. Luffy", Category: "anime"},
	{Text: "Hard work is worthless for those that don't believe in themselves.", Author: "Naruto Uzumaki", Category: "anime"},
	{Text: "It's just pathetic to give up on something before you even give it a shot.", Author: "Reiko Mikami", Category: "anime"},
	{Text: "The world isn't perfect, but it's there for us trying the best it can. That's what makes it so beautiful.", Author: "Roy Mustang", Category: "anime"},
	{Text: "Power comes in response to a need, not a desire.", Author: "Goku", Category: "anime"},
	{Text: "Whatever you lose, you'll find it again. But what you throw away you'll never get back.", Author: "Kenshin Himura", Category: "anime"},
	{Text: "We are all like fireworks: we climb, we shine and always go our separate ways and become further apart.", Author: "Hitsugaya Toushiro", Category: "anime"},
	{Text: "Don't be so quick to throw away your life. No matter how disgraceful or embarrassing it may be, you need to keep struggling to find your way out until the very end.", Author: "Clare", Category: "anime"},
	{Text: "A lesson without pain is meaningless. For you cannot gain anything without sacrificing something else in return.", Author: "Edward Elric", Category: "anime"},
}

// Quote represents a quote for typing
type Quote struct {
	Text     string
	Author   string
	Category string
}

// GetRandomQuote returns a random quote
func GetRandomQuote() Quote {
	allQuotes := append(Quotes, AnimeQuotes...)
	return allQuotes[rand.Intn(len(allQuotes))]
}

// GetQuoteByCategory returns a random quote from a specific category
func GetQuoteByCategory(category string) Quote {
	var filtered []Quote
	allQuotes := append(Quotes, AnimeQuotes...)

	for _, q := range allQuotes {
		if q.Category == category {
			filtered = append(filtered, q)
		}
	}

	if len(filtered) == 0 {
		return GetRandomQuote()
	}

	return filtered[rand.Intn(len(filtered))]
}

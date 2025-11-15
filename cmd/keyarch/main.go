package main

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/vinayydv3695/keyarch/internal/config"
	"github.com/vinayydv3695/keyarch/internal/data"
	"github.com/vinayydv3695/keyarch/internal/engine"
	"github.com/vinayydv3695/keyarch/internal/storage"
	"github.com/vinayydv3695/keyarch/internal/tui/duration"
	"github.com/vinayydv3695/keyarch/internal/tui/home"
	"github.com/vinayydv3695/keyarch/internal/tui/language"
	"github.com/vinayydv3695/keyarch/internal/tui/stats"
	"github.com/vinayydv3695/keyarch/internal/tui/summary"
	"github.com/vinayydv3695/keyarch/internal/tui/test"
	"github.com/vinayydv3695/keyarch/internal/tui/theme"
	"github.com/vinayydv3695/keyarch/internal/tui/wordcount"
)

var (
	themeFlag    string
	durationFlag int
	wordsFlag    int
	modeFlag     string
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "keyarch",
		Short: "Keyarch - A minimal typing test for your terminal",
		Long:  `Keyarch is a fast, beautiful, feature-rich typing tester built with Go, BubbleTea, and Lipgloss.`,
		Run:   run,
	}

	rootCmd.Flags().StringVar(&themeFlag, "theme", "", "Theme to use (catppuccin-mocha, nord, dracula, gruvbox-dark)")
	rootCmd.Flags().IntVar(&durationFlag, "duration", 15, "Test duration in seconds (for quick test)")
	rootCmd.Flags().IntVar(&wordsFlag, "words", 25, "Number of words (for word mode)")
	rootCmd.Flags().StringVar(&modeFlag, "mode", "", "Direct mode: quick, timed, words, quote, code")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Apply theme flag if provided
	if themeFlag != "" {
		cfg.Theme = themeFlag
	}

	// Initialize database
	db, err := storage.New()
	if err != nil {
		log.Printf("Warning: Could not initialize database: %v", err)
	}
	if db != nil {
		defer db.Close()
	}

	// Handle direct mode flag
	if modeFlag != "" {
		runDirectMode(cfg, db)
		return
	}

	// Run TUI application
	app := NewApp(cfg, db)
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}

func runDirectMode(cfg *config.Config, db *storage.DB) {
	gen := data.NewGenerator()
	var eng *engine.Engine

	switch modeFlag {
	case "quick":
		text := gen.GenerateByTime(durationFlag, "medium")
		eng = engine.New(text, engine.ModeTimer, durationFlag, 0)

	case "words":
		text := gen.GenerateWords(wordsFlag, "medium")
		eng = engine.New(text, engine.ModeWords, 0, wordsFlag)

	case "quote":
		quote := data.GetRandomQuote()
		eng = engine.New(quote.Text, engine.ModeQuote, 0, 0)

	case "code":
		snippet := data.GetCodeSnippet("go")
		eng = engine.New(snippet.Code, engine.ModeCode, 0, 0)

	default:
		fmt.Println("Invalid mode. Use: quick, words, quote, or code")
		return
	}

	// Run test
	m := test.New(eng, cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		log.Fatal(err)
	}

	testModel := finalModel.(test.Model)
	if testModel.Finished() && testModel.Engine().IsFinished {
		// Show summary
		s := summary.New(testModel.Engine(), cfg, db)
		p2 := tea.NewProgram(s, tea.WithAltScreen())
		if _, err := p2.Run(); err != nil {
			log.Fatal(err)
		}
	}
}

// App represents the main application
type App struct {
	cfg   *config.Config
	db    *storage.DB
	state string
}

// NewApp creates a new application
func NewApp(cfg *config.Config, db *storage.DB) *App {
	return &App{
		cfg:   cfg,
		db:    db,
		state: "home",
	}
}

// Run starts the application
func (a *App) Run() error {
	for {
		switch a.state {
		case "home":
			if err := a.runHome(); err != nil {
				return err
			}

		case "quick":
			if err := a.runQuickTest(); err != nil {
				return err
			}

		case "timed":
			if err := a.runTimedTest(); err != nil {
				return err
			}

		case "words":
			if err := a.runWordTest(); err != nil {
				return err
			}

		case "quote":
			if err := a.runQuoteMode(); err != nil {
				return err
			}

		case "code":
			if err := a.runCodeMode(); err != nil {
				return err
			}

		case "stats":
			if err := a.runStats(); err != nil {
				return err
			}

		case "themes":
			if err := a.runThemes(); err != nil {
				return err
			}

		case "quit":
			return nil

		default:
			a.state = "home"
		}
	}
}

func (a *App) runHome() error {
	m := home.New(a.cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	homeModel := finalModel.(home.Model)
	selected := homeModel.Selected()

	switch selected {
	case "Quick Test":
		a.state = "quick"
	case "Timed Test":
		a.state = "timed"
	case "Word Test":
		a.state = "words"
	case "Quote Mode":
		a.state = "quote"
	case "Code Mode":
		a.state = "code"
	case "Statistics":
		a.state = "stats"
	case "Themes":
		a.state = "themes"
	case "Quit", "quit":
		a.state = "quit"
	default:
		a.state = "home"
	}

	return nil
}

func (a *App) runQuickTest() error {
	gen := data.NewGenerator()
	text := gen.GenerateByTime(15, "medium")
	eng := engine.New(text, engine.ModeTimer, 15, 0)

	return a.runTest(eng)
}

func (a *App) runTimedTest() error {
	// Show duration selector
	m := duration.New(a.cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	durationModel := finalModel.(duration.Model)
	selectedDuration := durationModel.Selected()

	// -1 means back button
	if selectedDuration == -1 {
		a.state = "home"
		return nil
	}

	gen := data.NewGenerator()
	text := gen.GenerateByTime(selectedDuration, "medium")
	eng := engine.New(text, engine.ModeTimer, selectedDuration, 0)

	return a.runTest(eng)
}

func (a *App) runWordTest() error {
	// Show word count selector
	m := wordcount.New(a.cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	wordcountModel := finalModel.(wordcount.Model)
	selectedCount := wordcountModel.Selected()

	// -1 means back button
	if selectedCount == -1 {
		a.state = "home"
		return nil
	}

	gen := data.NewGenerator()
	text := gen.GenerateWords(selectedCount, "medium")
	eng := engine.New(text, engine.ModeWords, 0, selectedCount)

	return a.runTest(eng)
}

func (a *App) runQuoteMode() error {
	quote := data.GetRandomQuote()
	eng := engine.New(quote.Text, engine.ModeQuote, 0, 0)

	return a.runTest(eng)
}

func (a *App) runCodeMode() error {
	// Show language selector
	m := language.New(a.cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	languageModel := finalModel.(language.Model)
	selectedLang := languageModel.Selected()

	// "back" means back button
	if selectedLang == "back" {
		a.state = "home"
		return nil
	}

	snippet := data.GetCodeSnippet(selectedLang)
	eng := engine.New(snippet.Code, engine.ModeCode, 0, 0)

	return a.runTest(eng)
}

func (a *App) runTest(eng *engine.Engine) error {
	// Run the test
	m := test.New(eng, a.cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	testModel := finalModel.(test.Model)
	
	// If test was finished, show summary
	if testModel.Finished() && testModel.Engine().IsFinished {
		s := summary.New(testModel.Engine(), a.cfg, a.db)
		p2 := tea.NewProgram(s, tea.WithAltScreen())
		summaryFinal, err := p2.Run()
		if err != nil {
			return err
		}

		summaryModel := summaryFinal.(summary.Model)
		if summaryModel.Done() {
			a.state = "home"
		}
	} else {
		a.state = "home"
	}

	return nil
}

func (a *App) runStats() error {
	m := stats.New(a.cfg, a.db)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	statsModel := finalModel.(stats.Model)
	if statsModel.Done() {
		a.state = "home"
	}

	return nil
}

func (a *App) runThemes() error {
	m := theme.New(a.cfg)
	p := tea.NewProgram(m, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return err
	}

	themeModel := finalModel.(theme.Model)
	if themeModel.Selected() {
		// Reload config to get the updated theme
		newCfg, _ := config.Load()
		if newCfg != nil {
			a.cfg = newCfg
		}
		a.state = "home"
	}

	return nil
}

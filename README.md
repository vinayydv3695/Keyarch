# Keyarch ⌨️

<div align="center">

<pre>
 ██ ▄█▀▓█████▓██   ██▓ ▄▄▄       ██▀███   ▄████▄   ██░ ██ 
 ██▄█▒ ▓█   ▀ ▒██  ██▒▒████▄    ▓██ ▒ ██▒▒██▀ ▀█  ▓██░ ██▒
▓███▄░ ▒███    ▒██ ██░▒██  ▀█▄  ▓██ ░▄█ ▒▒▓█    ▄ ▒██▀▀██░
▓██ █▄ ▒▓█  ▄  ░ ▐██▓░░██▄▄▄▄██ ▒██▀▀█▄  ▒▓▓▄ ▄██▒░▓█ ░██ 
▒██▒ █▄░▒████▒ ░ ██▒▓░ ▓█   ▓██▒░██▓ ▒██▒▒ ▓███▀ ░░▓█▒░██▓
▒ ▒▒ ▓▒░░ ▒░ ░  ██▒▒▒  ▒▒   ▓▒█░░ ▒▓ ░▒▓░░ ░▒ ▒  ░ ▒ ░░▒░▒
</pre>

**A minimal, beautiful typing test for your terminal**

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

---

## ✨ Features

### 🎯 Core Features
- **Multiple Test Modes**: Timer (15s-120s), Word (25-200), Quote, Code (10 languages)
- **Real-Time Stats**: WPM, CPM, Accuracy, Mistakes
- **10 Beautiful Themes**: Catppuccin, Nord, Dracula, Gruvbox, Rose Pine, Tokyo Night, One Dark, Monokai, Solarized, Material
- **Performance Tracking**: History, streaks, weak keys analysis
- **Local Storage**: All data stored in `~/.keyarch/`

### 🎮 Gamification
- **23 Achievements**: Speed, Accuracy, Dedication, Consistency badges
- **Daily & Weekly Goals**: Auto-reset goals with progress tracking
- **Progress Dashboard**: 3-tab interface (Achievements, Daily, Weekly)
- **Hidden Achievements**: Unlock surprises at expert milestones

### 🚀 Advanced Features
- **📊 Live WPM Graph**: Real-time sparkline showing speed trends during test
- **🎯 Difficulty Levels**: Easy/Medium/Hard word sets (300+ words)
- **🔊 Sound Effects**: Keystroke sounds, achievement/goal notifications
- **🔥 Typing Heatmap**: Color-coded keyboard showing key accuracy

---

## 🚀 Installation

### Prerequisites
- Go 1.22 or higher

### From Source

```bash
git clone https://github.com/vinayydv3695/keyarch.git
cd keyarch
go mod tidy
go build -o keyarch ./cmd/keyarch
sudo mv keyarch /usr/local/bin/  # Optional
```

### Quick Install

```bash
go install github.com/vinayydv3695/keyarch/cmd/keyarch@latest
```

---

## 🎮 Usage

```bash
# Launch TUI
keyarch

# Direct modes
keyarch --mode quick --duration 30
keyarch --mode words --words 50
keyarch --mode quote
keyarch --mode code

# With theme
keyarch --theme dracula
```

---

## ⌨️ Keyboard Shortcuts

- `↑/↓` or `j/k` - Navigate
- `Enter` - Select
- `ESC` - Back
- `Ctrl+C` - Quit
- `Backspace` - Delete character

---

## 🎨 Themes

Switch themes from the menu or use `--theme` flag:

### 🌙 Available Themes
1. **Catppuccin Mocha** (default) - Soft pastels
2. **Nord** - Arctic blues
3. **Dracula** - Vibrant dark
4. **Gruvbox Dark** - Warm retro
5. **Rose Pine** - Elegant cozy
6. **Tokyo Night** - Modern clean
7. **One Dark** - Classic Atom
8. **Monokai** - Sublime colors
9. **Solarized Dark** - Precision colors
10. **Material Dark** - Google design

---

## 💻 Code Mode Languages

Practice typing in 10 programming languages:
- Go, JavaScript, TypeScript, Python, Rust
- C++, Java, C#, Ruby, PHP

---

## 📊 Statistics

Track your progress with:
- Test history with timestamps
- Best and average WPM
- Accuracy trends
- Daily streaks
- Weak keys identification

Data stored in: `~/.keyarch/keyarch.db`

---

## 🛠️ Development

```bash
# Run
go run ./cmd/keyarch

# Build
go build -ldflags="-s -w" -o keyarch ./cmd/keyarch

# Test
go test ./...
```

---

## 📝 License

MIT License - see [LICENSE](LICENSE) file

---

## 🙏 Credits

Built with:
- [BubbleTea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling
- [Cobra](https://github.com/spf13/cobra) - CLI
- [SQLite](https://gitlab.com/cznic/sqlite) - Database

---

<div align="center">

**Made with ❤️ and ⌨️**

GitHub: [@vinayydv3695](https://github.com/vinayydv3695)

</div>
````

<div align="center">

**Made with ❤️ and ⌨️**

GitHub: [@vinayydv3695](https://github.com/vinayydv3695)

</div>

---

## 🚀 Installation

### Prerequisites
- Go 1.22 or higher
- A terminal with Unicode support

### From Source

```bash
# Clone the repository
git clone https://github.com/vinayydv3695/keyarch.git
cd keyarch

# Install dependencies
go mod tidy

# Build the application
go build -o keyarch ./cmd/keyarch

# Move to PATH (optional)
sudo mv keyarch /usr/local/bin/
```

### Quick Install

```bash
go install github.com/vinayydv3695/keyarch/cmd/keyarch@latest
```

---

## 🎮 Usage

### Launch the TUI

```bash
keyarch
```

### Direct Mode (Skip Menu)

```bash
# Quick 15-second test
keyarch --mode quick

# 30-second timed test
keyarch --mode quick --duration 30

# 50-word test
keyarch --mode words --words 50

# Random quote
keyarch --mode quote

# Code snippet
keyarch --mode code
```

### Theme Selection

```bash
# Set theme via CLI
keyarch --theme dracula
keyarch --theme nord
keyarch --theme gruvbox-dark
keyarch --theme catppuccin-mocha
```

---

## ⌨️ Keyboard Shortcuts

### Navigation
- `↑/↓` or `j/k` - Navigate menus
- `Enter` or `Space` - Select/Confirm
- `ESC` - Go back
- `Ctrl+C` or `q` - Quit

### During Test
- `Backspace` - Delete character
- `ESC` - Exit to menu
- `Ctrl+C` - Quit application

---

## 📁 Project Structure

```
keyarch/
├── cmd/
│   └── keyarch/
│       └── main.go              # Application entry point
│
├── internal/
│   ├── engine/                  # Typing test logic
│   │   └── engine.go            # WPM, accuracy, input handling
│   │
│   ├── tui/                     # BubbleTea UI
│   │   ├── home/                # Home menu screen
│   │   ├── test/                # Active test screen
│   │   ├── summary/             # Results summary screen
│   │   ├── stats/               # Statistics dashboard
│   │   ├── theme/               # Theme selector
│   │   └── components/          # Reusable UI components
│   │
│   ├── storage/                 # SQLite database
│   │   ├── db.go                # Database operations
│   │   └── models.go            # Data models
│   │
│   ├── config/                  # Configuration
│   │   ├── config.go            # App settings
│   │   └── theme.go             # Theme definitions
│   │
│   └── data/                    # Text generation
│       ├── words.go             # Word lists & quotes
│       └── snippets.go          # Code snippets
│
├── assets/
│   ├── words/                   # Word lists
│   │   ├── easy.txt
│   │   ├── medium.txt
│   │   └── hard.txt
│   └── snippets/                # Code snippets
│       ├── go/
│       ├── js/
│       ├── py/
│       └── rs/
│
├── go.mod
├── go.sum
└── README.md
```

---

## 🎨 Themes

### Catppuccin Mocha (Default)
Soft, pastel colors with excellent contrast. Perfect for long typing sessions.

### Nord
Cool, arctic-inspired palette. Clean and minimal.

### Dracula
Dark theme with vibrant accent colors. High contrast for visibility.

### Gruvbox Dark
Warm, retro-inspired colors. Easy on the eyes.

**Switch themes**: Navigate to "Themes" in the main menu or use `--theme` flag.

---

## 📊 Statistics

Keyarch tracks your progress locally:

- **Test History**: All your typing tests with timestamps
- **Best Scores**: Track your personal records
- **Accuracy Trends**: See how your accuracy improves
- **Weak Keys**: Identify which keys you mistype most
- **Daily Streaks**: Stay motivated with daily practice

All data is stored in `~/.keyarch/keyarch.db`

---

## 🛠️ Development

### Run in Development

```bash
go run ./cmd/keyarch
```

### Run Tests

```bash
go test ./...
```

### Build for Production

```bash
# Current platform
go build -ldflags="-s -w" -o keyarch ./cmd/keyarch

# Cross-compile for Linux
GOOS=linux GOARCH=amd64 go build -o keyarch-linux ./cmd/keyarch

# Cross-compile for macOS
GOOS=darwin GOARCH=arm64 go build -o keyarch-macos ./cmd/keyarch

# Cross-compile for Windows
GOOS=windows GOARCH=amd64 go build -o keyarch.exe ./cmd/keyarch
```

---

## 📚 Documentation

- **[GAMIFICATION.md](GAMIFICATION.md)** - Complete guide to achievements and goals system
- **[ADVANCED_FEATURES.md](ADVANCED_FEATURES.md)** - Detailed docs for WPM Graph, Difficulty Levels, Sound Effects, and Heatmap
- **[CODE_SNIPPETS.md](CODE_SNIPPETS.md)** - Available programming languages and snippet examples

---

## 🤝 Contributing

Contributions are welcome! Here are some ideas:

- 🌍 Add more language support
- 🎨 Create new themes
- 📝 Add more quotes and code snippets
- 🐛 Fix bugs and improve performance
- 📚 Improve documentation

### Steps to Contribute

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 🎯 Roadmap

### ✅ Completed
- [x] 23 Achievement system with badges
- [x] Daily and weekly goals
- [x] Live WPM graph during tests
- [x] Difficulty levels (Easy/Medium/Hard)
- [x] Sound effects system
- [x] Typing heatmap visualization
- [x] 10 beautiful themes
- [x] 10 programming languages

### 🚧 In Progress
- [ ] Online quote API integration
- [ ] Custom word lists
- [ ] Multiplayer mode (local network)

### 📋 Planned
- [ ] GitHub Gist leaderboard sync
- [ ] Custom themes (JSON config)
- [ ] Export statistics as JSON/CSV
- [ ] Mobile-friendly version
- [ ] Plugin system

---

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

Built with amazing open-source libraries:

- [BubbleTea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [Bubbles](https://github.com/charmbracelet/bubbles) - TUI components
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [modernc.org/sqlite](https://gitlab.com/cznic/sqlite) - Pure Go SQLite

Inspired by:
- [Monkeytype](https://monkeytype.com/)
- [Typeracer](https://play.typeracer.com/)
- The Arch Linux philosophy of simplicity

---

## 📬 Contact

**Author**: Vinay Yadav  
**GitHub**: [@vinayydv3695](https://github.com/vinayydv3695)  
**Project**: [Keyarch](https://github.com/vinayydv3695/keyarch)

---

<div align="center">

**Made with ❤️ and ⌨️**

*Happy typing!*

</div>

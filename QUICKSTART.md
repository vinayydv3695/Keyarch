# Quick Start Guide for Keyarch

## 🚀 Getting Started

### Step 1: Ensure Go is Installed

Check if Go is installed:

```bash
go version
```

If not installed, visit: https://golang.org/dl/

### Step 2: Navigate to Project

```bash
cd "/home/zura/Personal/coding cuff/Keyarch"
```

### Step 3: Install Dependencies

```bash
go mod tidy
```

This will download:
- BubbleTea (TUI framework)
- Lipgloss (styling)
- Bubbles (components)
- Cobra (CLI)
- SQLite (database)

### Step 4: Build the Application

```bash
go build -o keyarch ./cmd/keyarch
```

Or use the Makefile:

```bash
make build
```

### Step 5: Run Keyarch!

```bash
./keyarch
```

Or directly:

```bash
go run ./cmd/keyarch
```

## 🎯 First Test

1. Launch: `./keyarch`
2. Select "Quick Test" with Enter
3. Start typing when text appears
4. View your results after 15 seconds
5. Press Enter to return to menu

## ⌨️ Navigation Cheat Sheet

**Main Menu:**
- `↑/↓` - Navigate options
- `Enter` - Select
- `q` or `Ctrl+C` - Quit

**During Test:**
- Type normally
- `Backspace` - Correct mistakes
- `ESC` - Back to menu

**Themes:**
- Navigate to "Themes" in main menu
- Use `↑/↓` to preview
- Press `Enter` to save

## 📊 View Your Stats

1. Select "Statistics" from main menu
2. See your progress, best WPM, and more
3. Data stored in `~/.keyarch/keyarch.db`

## 🎨 Try Different Modes

- **Quick Test** - 15 seconds, perfect for warming up
- **Timed Test** - Longer focused sessions
- **Word Test** - Type exactly N words
- **Quote Mode** - Inspirational quotes
- **Code Mode** - Practice programming syntax

## 🐛 Troubleshooting

**Problem: "go: command not found"**
- Install Go from https://golang.org/dl/

**Problem: Colors not showing**
- Use a modern terminal (Alacritty, Kitty, iTerm2, Windows Terminal)

**Problem: Unicode issues**
- Ensure UTF-8 encoding in your terminal
- Install a Nerd Font or similar

## 📈 Tips for Better Scores

1. **Warm up** - Do a quick test before longer sessions
2. **Posture** - Sit comfortably with wrists relaxed
3. **Look at screen** - Not at keyboard
4. **Consistency** - Daily practice improves speed
5. **Accuracy first** - Speed comes with accuracy

## 🎯 Challenge Yourself

- Try all 4 themes to find your favorite
- Beat your best WPM
- Aim for 95%+ accuracy
- Build a daily streak
- Try code mode for programming practice

## 💡 CLI Power User Tips

```bash
# Direct quick test
./keyarch --mode quick

# 60-second test
./keyarch --mode quick --duration 60

# 100-word test
./keyarch --mode words --words 100

# With theme
./keyarch --theme dracula --mode quick
```

## 🔧 Advanced

**Custom installation:**
```bash
sudo make install
# Now run from anywhere: keyarch
```

**Build for distribution:**
```bash
make build-all
# Creates binaries for Linux, macOS, Windows
```

**Reset everything:**
```bash
rm -rf ~/.keyarch
# Fresh start on next run
```

## 📚 Learn More

- Full documentation: [README.md](README.md)
- Build instructions: [BUILD.md](BUILD.md)
- GitHub: https://github.com/vinayydv3695/keyarch

---

**Happy Typing! ⌨️**

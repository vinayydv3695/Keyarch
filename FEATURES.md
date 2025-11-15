# Keyarch Features Documentation

## 📋 Complete Feature List

### 🎮 Test Modes

#### 1. Timer Mode
Type as much as you can within a time limit.

**Available Durations:**
- 15 seconds (Quick Test)
- 30 seconds
- 60 seconds
- 120 seconds
- Custom duration via CLI

**Features:**
- Real-time countdown
- Automatic completion
- WPM/CPM tracking
- Accuracy monitoring

**Usage:**
```bash
keyarch --mode quick --duration 30
```

#### 2. Word Mode
Type a specific number of words.

**Available Counts:**
- 25 words
- 50 words
- 100 words
- Custom count via CLI

**Features:**
- Progress bar
- Word counter
- No time pressure
- Focus on accuracy

**Usage:**
```bash
keyarch --mode words --words 50
```

#### 3. Quote Mode
Type inspirational and motivational quotes.

**Quote Categories:**
- Motivational
- Philosophical
- Anime quotes
- Random selection

**Features:**
- Meaningful text
- Variable length
- Author attribution
- Category variety

**Usage:**
```bash
keyarch --mode quote
```

#### 4. Code Mode
Practice typing programming syntax.

**Supported Languages:**
- Go
- JavaScript
- Python
- Rust

**Features:**
- Real code snippets
- Syntax-focused
- Multiple snippets per language
- Programming practice

**Usage:**
```bash
keyarch --mode code
```

---

## 📊 Real-Time Statistics

### During Test

**WPM (Words Per Minute)**
- Standard calculation: 5 characters = 1 word
- Updates in real-time
- Shown prominently during test

**CPM (Characters Per Minute)**
- Raw character count
- More granular than WPM
- Real-time updates

**Accuracy Percentage**
- Correct characters / total characters
- Displayed as percentage
- Updated live

**Mistake Counter**
- Tracks every error
- Increments on incorrect key
- Helps identify weak areas

**Progress Tracking**
- Visual progress bar (Word/Quote modes)
- Time remaining (Timer mode)
- Completion percentage

---

## 🎨 Themes

### Built-in Themes

Keyarch includes **10 beautiful themes**:

#### 1. Catppuccin Mocha (Default)
- **Style**: Modern, soft pastels
- **Best for**: Long typing sessions
- **Colors**: Purple, blue, pink accents
- **Contrast**: Medium-high

#### 2. Nord
- **Style**: Arctic, cool tones
- **Best for**: Clean, minimal aesthetic
- **Colors**: Blues, grays, teals
- **Contrast**: Medium

#### 3. Dracula
- **Style**: Dark with vibrant accents
- **Best for**: High visibility
- **Colors**: Purple, pink, cyan
- **Contrast**: High

#### 4. Gruvbox Dark
- **Style**: Warm, retro-inspired
- **Best for**: Easy on eyes
- **Colors**: Orange, yellow, green
- **Contrast**: Medium

#### 5. Rose Pine
- **Style**: All natural pine, faux fur and a bit of soho vibes
- **Best for**: Elegant, cozy aesthetic
- **Colors**: Rose, pine, foam, gold
- **Contrast**: Medium

#### 6. Tokyo Night
- **Style**: Clean, dark Tokyo-inspired
- **Best for**: Modern developers
- **Colors**: Blue, purple, orange
- **Contrast**: High

#### 7. One Dark
- **Style**: Iconic Atom editor theme
- **Best for**: Familiar, professional look
- **Colors**: Blue, purple, red, green
- **Contrast**: Medium-high

#### 8. Monokai
- **Style**: Classic Sublime Text colors
- **Best for**: Vibrant coding sessions
- **Colors**: Pink, cyan, orange, green
- **Contrast**: High

#### 9. Solarized Dark
- **Style**: Precision colors for machines and people
- **Best for**: Scientifically designed readability
- **Colors**: Blue, cyan, magenta, yellow
- **Contrast**: Medium

#### 10. Material Dark
- **Style**: Google's Material Design
- **Best for**: Clean, modern interface
- **Colors**: Blue, purple, pink, yellow
- **Contrast**: Medium-high

### Theme Features

- **Real-time preview** - See changes instantly
- **Persistent settings** - Saved to config
- **Consistent styling** - Applied across all screens
- **Custom colors** for:
  - Correct characters (green)
  - Incorrect characters (red)
  - Cursor position (yellow)
  - Pending text (muted)
  - UI elements (themed)

### Switching Themes

**Via TUI:**
1. Select "Themes" from main menu
2. Navigate with ↑/↓
3. Preview in real-time
4. Press Enter to save

**Via CLI:**
```bash
keyarch --theme nord
keyarch --theme dracula
keyarch --theme gruvbox-dark
keyarch --theme catppuccin-mocha
keyarch --theme rose-pine
keyarch --theme tokyo-night
keyarch --theme one-dark
keyarch --theme monokai
keyarch --theme solarized-dark
keyarch --theme material-dark
```

---

## 📈 Statistics Dashboard

### Overview Stats

- **Total Tests**: Lifetime test count
- **Best WPM**: Personal record
- **Average WPM**: Mean across all tests
- **Best Accuracy**: Highest accuracy achieved
- **Average Accuracy**: Mean accuracy
- **Total Time**: Cumulative typing time
- **Current Streak**: Consecutive days with tests

### Performance Graph

- **7-Day WPM Graph**: Visual bar chart
- **Date-based tracking**: Daily averages
- **Trend analysis**: See improvement over time

### Recent Tests

- **Last 10 tests**: Quick history view
- **Date stamps**: When tests were taken
- **Key metrics**: WPM, accuracy, mode
- **Test types**: Distinguish between modes

### Weak Keys Analysis

- **Character-level accuracy**: Per-key statistics
- **Top mistakes**: Most frequently missed keys
- **Improvement hints**: Focus areas
- **Accuracy percentages**: Per character

---

## 🎯 Text Highlighting

### Visual Feedback

**Correct Characters**
- Color: Green
- Indicates: Typed correctly
- Real-time: Immediate feedback

**Incorrect Characters**
- Color: Red
- Background: Dark red
- Indicates: Mistake made
- Visibility: High contrast

**Cursor Position**
- Color: Yellow/Gold
- Style: Underlined
- Indicates: Current position
- Always visible

**Pending Characters**
- Color: Muted gray
- Indicates: Not yet typed
- Context: Shows what's ahead

### Benefits

- **Instant feedback**: Know immediately if correct
- **Error awareness**: See mistakes as they happen
- **Progress tracking**: Visual completion indicator
- **Focus aid**: Cursor shows where to look

---

## 💾 Data Persistence

### Local Storage

**Location**: `~/.keyarch/`

**Files:**
- `keyarch.db` - SQLite database
- `config.json` - User preferences

### Database Schema

**test_results table:**
- Test ID
- WPM
- CPM
- Accuracy
- Mistakes count
- Test type
- Test mode
- Duration
- Timestamp
- Character counts

**key_stats table:**
- Per-test key statistics
- Correct/incorrect counts
- Character-level data

### Data Privacy

- **Fully local**: No cloud sync
- **No telemetry**: No data sent
- **User owned**: Full control
- **Portable**: Easy to backup/restore

### Managing Data

**View location:**
```bash
ls ~/.keyarch/
```

**Backup:**
```bash
cp -r ~/.keyarch/ ~/keyarch-backup/
```

**Reset:**
```bash
rm -rf ~/.keyarch/
```

---

## ⌨️ Keyboard Controls

### Navigation

| Key | Action |
|-----|--------|
| `↑` or `k` | Move up |
| `↓` or `j` | Move down |
| `Enter` | Select/Confirm |
| `Space` | Alternative select |
| `ESC` | Back/Cancel |
| `q` | Quit (menus) |
| `Ctrl+C` | Force quit |

### During Test

| Key | Action |
|-----|--------|
| Any letter | Type character |
| `Space` | Space character |
| `Backspace` | Delete character |
| `ESC` | Exit test |
| `Ctrl+C` | Force quit |

### Vim-style Navigation

- `j` = Down
- `k` = Up
- Works in all menus

---

## 🚀 CLI Flags & Options

### Global Flags

```bash
--theme string      # Set theme
--duration int      # Test duration (seconds)
--words int         # Word count for word mode
--mode string       # Direct mode selection
```

### Mode Values

- `quick` - Quick timer test
- `timed` - Timed test
- `words` - Word count test
- `quote` - Quote mode
- `code` - Code snippet mode

### Examples

```bash
# Quick 30-second test with Nord theme
keyarch --theme nord --mode quick --duration 30

# 100-word test
keyarch --mode words --words 100

# Random quote with Dracula theme
keyarch --theme dracula --mode quote

# Code practice
keyarch --mode code
```

---

## 🎓 Advanced Features

### Accuracy Algorithm

- **Character-by-character**: Each character scored
- **Position-aware**: Order matters
- **Real-time calculation**: Updates instantly
- **Percentage-based**: 0-100% scale

### WPM Calculation

**Standard Formula:**
```
WPM = (Total Characters / 5) / (Time in Minutes)
```

- 5 characters = 1 word (standard)
- Includes spaces and punctuation
- Real-time updates

### Streak Calculation

- **Daily basis**: One test per day counts
- **Consecutive days**: No gaps allowed
- **Timezone aware**: Uses local date
- **Resets on gap**: Miss a day, start over

### Progress Tracking

**Timer Mode:**
- Based on elapsed time
- Percentage of duration complete

**Word Mode:**
- Based on words typed
- Percentage of target reached

**Quote/Code Mode:**
- Based on characters typed
- Percentage of text complete

---

## 🎨 UI/UX Features

### Layout

- **Centered design**: Focus on content
- **Generous padding**: Comfortable spacing
- **Clear hierarchy**: Important info prominent
- **Consistent spacing**: Aligned elements

### Typography

- **ASCII art logo**: Distinctive header
- **Bold titles**: Clear sections
- **Muted labels**: Subtle UI text
- **Monospace text**: Typing area

### Visual Effects

- **Rounded borders**: Soft, modern look
- **Color gradients**: Theme-appropriate
- **Smooth transitions**: Screen changes
- **Hover states**: Active selections

### Responsive Design

- **Terminal size aware**: Adapts to window
- **Text wrapping**: Handles long lines
- **Scalable stats**: Adjusts to space
- **Mobile-friendly**: Works on small screens

---

## 📦 File Organization

### Project Structure

```
keyarch/
├── cmd/keyarch/        # Entry point
├── internal/
│   ├── engine/         # Core logic
│   ├── tui/            # UI components
│   ├── storage/        # Database
│   ├── config/         # Configuration
│   └── data/           # Text generation
└── assets/             # Static data
```

### Code Organization

- **Package-based**: Logical separation
- **Internal packages**: Not importable
- **Clear responsibilities**: Single purpose
- **Minimal dependencies**: Clean imports

---

## 🔧 Configuration

### Config File

**Location**: `~/.keyarch/config.json`

**Structure:**
```json
{
  "theme": "catppuccin-mocha",
  "sound": false,
  "blind_mode": false
}
```

### Settings

- **theme**: Active theme name
- **sound**: Sound effects (future)
- **blind_mode**: Hide text (future)

### Modification

**Via TUI**: Theme selector menu
**Manually**: Edit config.json
**CLI**: Use --theme flag

---

## 🎯 Performance Tips

### For Users

1. **Practice daily**: Build muscle memory
2. **Start slow**: Accuracy before speed
3. **Use proper posture**: Ergonomics matter
4. **Touch type**: Don't look at keyboard
5. **Warm up**: Quick test first

### For Developers

1. **Efficient rendering**: Minimal redraws
2. **SQLite**: Fast local storage
3. **Go concurrency**: Responsive UI
4. **Minimal dependencies**: Small binary
5. **Optimized builds**: Stripped symbols

---

## 🐛 Known Limitations

1. **Terminal-only**: No GUI version
2. **Local data**: No cloud sync
3. **Single player**: No multiplayer (yet)
4. **English only**: No internationalization (yet)
5. **Basic stats**: Limited analytics

---

## 🔮 Future Features

See [README.md](README.md) Roadmap section for:
- Online quote API
- Custom word lists
- Multiplayer mode
- GitHub Gist leaderboard
- More programming languages
- Sound effects
- Blind mode
- Export statistics

---

## 📚 Resources

- **Documentation**: [README.md](README.md)
- **Build Guide**: [BUILD.md](BUILD.md)
- **Quick Start**: [QUICKSTART.md](QUICKSTART.md)
- **GitHub**: https://github.com/vinayydv3695/keyarch

---

**Last Updated**: November 2025

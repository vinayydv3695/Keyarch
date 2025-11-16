# Advanced Features Guide

Keyarch includes several advanced features to enhance your typing practice experience and help you improve faster.

## 📊 Live WPM Graph

### What is it?
A real-time sparkline graph that displays your typing speed (Words Per Minute) as you type during a test.

### How it works:
- **Sampling**: Updates every 2 seconds during active typing
- **Display**: Shows last 60 seconds of WPM data (30 samples)
- **Visualization**: Uses ASCII blocks (▁▂▃▄▅▆▇█) to create a sparkline
- **Stats**: Displays minimum and maximum WPM values

### Benefits:
- See your speed trends in real-time
- Identify when you slow down or speed up
- Visual feedback helps maintain consistent pace
- Track performance within a single test

### Example:
```
┌─ Live WPM ─────────────────────────────────────────┐
│  Max: 85 WPM                                        │
│  ▃▄▅▆▇█▇▆▅▄▅▆▇▇▆▅▄▃▄▅▆▇█▇▆▅▆▇█▇▆▅▄▃              │
│  Min: 45 WPM                                        │
└─────────────────────────────────────────────────────┘
```

---

## 🎯 Difficulty Levels

### Overview
Choose from three difficulty levels to match your skill level and challenge yourself.

### Difficulty Options:

#### 🟢 Easy Mode
- **Word Length**: 2-4 letters
- **Vocabulary**: Common everyday words
- **Best For**: Beginners, warming up, building confidence
- **Examples**: "the", "and", "can", "run", "dog", "cat"

#### 🟡 Medium Mode (Default)
- **Word Length**: 5-8 letters
- **Vocabulary**: Balanced mix of common and moderate words
- **Best For**: Regular practice, balanced improvement
- **Examples**: "about", "system", "program", "question"

#### 🔴 Hard Mode
- **Word Length**: 8+ letters
- **Vocabulary**: Complex, rare, and challenging words
- **Best For**: Advanced typists, vocabulary expansion
- **Examples**: "aberration", "comprehensive", "philosophical"

### How to Use:
1. Select "Word Test" from main menu
2. Choose your difficulty level
3. Select word count
4. Start typing!

### Tips:
- Start with Easy to build muscle memory
- Progress to Medium when comfortable with accuracy
- Challenge yourself with Hard for vocabulary growth
- Mix difficulties based on your goals

---

## 🔊 Sound Effects

### Features
Immersive audio feedback to make typing more engaging and provide instant feedback.

### Sound Types:

#### 1. Keystroke Sounds
- **When**: Every key press
- **Sound**: Subtle single beep
- **Purpose**: Auditory confirmation of input
- **Note**: Can be disabled in settings

#### 2. Success Sounds
- **When**: Test completion with good results
- **Sound**: Pleasant ascending beeps (3 tones)
- **Purpose**: Positive reinforcement

#### 3. Error Sounds
- **When**: Mistakes during typing
- **Sound**: Two quick beeps
- **Purpose**: Alert you to errors

#### 4. Achievement Unlock
- **When**: New achievement earned
- **Sound**: Triumphant sequence (5 tones)
- **Purpose**: Celebrate milestones

#### 5. Goal Completion
- **When**: Daily/weekly goal completed
- **Sound**: Success fanfare (4 tones)
- **Purpose**: Motivate consistency

### Platform Support:
- **Linux**: Uses paplay or aplay for system sounds
- **macOS**: Uses afplay with system sounds
- **Windows**: Uses terminal bell
- **Fallback**: Terminal bell if no audio system available

### How to Enable/Disable:
1. Go to main menu
2. Select "Settings"
3. Toggle "Sound Effects: ON/OFF"
4. Changes save automatically

### Blind Mode:
- Hides WPM and statistics during typing
- Hardcore mode for experienced typists
- Forces you to type without visual feedback
- Toggle in Settings menu

---

## 🔥 Typing Heatmap

### What is it?
A visual representation of your keyboard showing accuracy for each key you typed, displayed as a color-coded QWERTY layout.

### Color Coding:

| Color | Range | Meaning |
|-------|-------|---------|
| 🟢 Green | 95-100% | Excellent - Master these keys! |
| 🔷 Teal | 85-94% | Good - Consistent performance |
| 🟡 Yellow | 75-84% | Average - Room for improvement |
| 🟠 Peach | 60-74% | Below Average - Needs practice |
| 🔴 Red | <60% | Weak - Focus here! |
| ⬜ Gray | Unused | Not typed in this test |

### Layout:
```
  ` 1 2 3 4 5 6 7 8 9 0 - =
  Q W E R T Y U I O P [ ] \
  A S D F G H J K L ; '
  Z X C V B N M , . /
  SPACE
```

### Benefits:
1. **Instant Feedback**: See which keys you struggle with
2. **Targeted Practice**: Focus on weak keys (red/orange)
3. **Progress Tracking**: Watch colors improve over time
4. **Visual Learning**: Spatial memory of problem areas

### How to Read:
- **Bright Green Keys**: Your strong suit - maintain accuracy!
- **Red/Orange Keys**: Practice these in isolation
- **Gray Keys**: Consider including in future tests
- **Space Bar**: Often overlooked but crucial for accuracy

### Tips for Improvement:
1. **Focus on Red Keys**: Create custom tests with these letters
2. **Finger Placement**: Check if you're using correct fingers
3. **Slow Down**: Accuracy > speed for weak keys
4. **Repeat Tests**: Watch colors shift from red → orange → yellow → green

---

**Happy typing! Use these features to accelerate your improvement!** 🚀

package sound

import (
	"fmt"
	"os/exec"
	"runtime"
)

// SoundType represents different sound effects
type SoundType string

const (
	SoundKeystroke    SoundType = "keystroke"
	SoundSuccess      SoundType = "success"
	SoundError        SoundType = "error"
	SoundAchievement  SoundType = "achievement"
	SoundGoalComplete SoundType = "goal_complete"
)

// Player handles sound effects
type Player struct {
	Enabled bool
}

// NewPlayer creates a new sound player
func NewPlayer(enabled bool) *Player {
	return &Player{
		Enabled: enabled,
	}
}

// Play plays a sound effect
func (p *Player) Play(soundType SoundType) {
	if !p.Enabled {
		return
	}

	go p.playSound(soundType)
}

// playSound plays the actual sound (in a goroutine to avoid blocking)
func (p *Player) playSound(soundType SoundType) {
	switch soundType {
	case SoundKeystroke:
		// Subtle single beep for keystroke
		p.beep(1, 10)
	case SoundSuccess:
		// Pleasant ascending beeps
		p.beep(3, 30)
	case SoundError:
		// Two quick beeps
		p.beep(2, 15)
	case SoundAchievement:
		// Triumphant sequence
		p.beep(5, 40)
	case SoundGoalComplete:
		// Success fanfare
		p.beep(4, 35)
	}
}

// beep creates a beep sound
func (p *Player) beep(count int, duration int) {
	for i := 0; i < count; i++ {
		p.systemBeep()
	}
}

// systemBeep uses system-specific beep command
func (p *Player) systemBeep() {
	switch runtime.GOOS {
	case "linux":
		// Try different methods for Linux
		if err := exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/message.oga").Run(); err != nil {
			if err := exec.Command("aplay", "/usr/share/sounds/alsa/Front_Center.wav").Run(); err != nil {
				// Fallback to terminal bell
				fmt.Print("\a")
			}
		}
	case "darwin": // macOS
		exec.Command("afplay", "/System/Library/Sounds/Tink.aiff").Run()
	case "windows":
		// Windows has built-in beep
		fmt.Print("\a")
	default:
		// Fallback to terminal bell
		fmt.Print("\a")
	}
}

// PlayKeystroke plays a subtle keystroke sound
func (p *Player) PlayKeystroke() {
	if !p.Enabled {
		return
	}
	// Very subtle - just terminal bell
	fmt.Print("\a")
}

// PlaySuccess plays a success sound
func (p *Player) PlaySuccess() {
	p.Play(SoundSuccess)
}

// PlayError plays an error sound
func (p *Player) PlayError() {
	p.Play(SoundError)
}

// PlayAchievement plays an achievement unlock sound
func (p *Player) PlayAchievement() {
	p.Play(SoundAchievement)
}

// PlayGoalComplete plays a goal completion sound
func (p *Player) PlayGoalComplete() {
	p.Play(SoundGoalComplete)
}

// Toggle enables/disables sound
func (p *Player) Toggle() {
	p.Enabled = !p.Enabled
}

// SetEnabled sets the enabled state
func (p *Player) SetEnabled(enabled bool) {
	p.Enabled = enabled
}

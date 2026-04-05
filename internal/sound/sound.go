package sound

import (
	"fmt"
	"os/exec"
	"runtime"
	"sync"
	"time"
)

// SoundType represents different sound effects
type SoundType string

const (
	SoundKeystroke    SoundType = "keystroke"
	SoundSuccess      SoundType = "success"
	SoundError        SoundType = "error"
	SoundAchievement  SoundType = "achievement"
	SoundGoalComplete SoundType = "goal_complete"
	SoundMistake      SoundType = "mistake"
)

// Profile represents sound intensity/style
type Profile string

const (
	ProfileOff        Profile = "off"
	ProfileSubtle     Profile = "subtle"
	ProfileMechanical Profile = "mechanical"
)

// Player handles sound effects
type Player struct {
	Enabled      bool
	Profile      Profile
	mu           sync.Mutex
	lastKeySound time.Time
}

// NewPlayer creates a new sound player
func NewPlayer(enabled bool) *Player {
	return &Player{
		Enabled: enabled,
		Profile: ProfileSubtle,
	}
}

// NewPlayerWithProfile creates a new sound player with specific profile
func NewPlayerWithProfile(enabled bool, profile Profile) *Player {
	return &Player{
		Enabled: enabled,
		Profile: profile,
	}
}

// Play plays a sound effect
func (p *Player) Play(soundType SoundType) {
	if !p.Enabled || p.Profile == ProfileOff {
		return
	}

	go p.playSound(soundType)
}

// playSound plays the actual sound (in a goroutine to avoid blocking)
func (p *Player) playSound(soundType SoundType) {
	switch p.Profile {
	case ProfileSubtle:
		p.playSubtle(soundType)
	case ProfileMechanical:
		p.playMechanical(soundType)
	}
}

// playSubtle plays subtle/quiet sounds
func (p *Player) playSubtle(soundType SoundType) {
	switch soundType {
	case SoundKeystroke:
		// Very subtle - just terminal bell occasionally
		p.mu.Lock()
		if time.Since(p.lastKeySound) > 50*time.Millisecond {
			fmt.Print("\a")
			p.lastKeySound = time.Now()
		}
		p.mu.Unlock()
	case SoundSuccess:
		p.systemSound("message")
	case SoundError, SoundMistake:
		fmt.Print("\a")
	case SoundAchievement:
		p.systemSound("complete")
	case SoundGoalComplete:
		p.systemSound("complete")
	}
}

// playMechanical plays louder mechanical keyboard-like sounds
func (p *Player) playMechanical(soundType SoundType) {
	switch soundType {
	case SoundKeystroke:
		p.mu.Lock()
		if time.Since(p.lastKeySound) > 30*time.Millisecond {
			fmt.Print("\a")
			p.lastKeySound = time.Now()
		}
		p.mu.Unlock()
	case SoundSuccess:
		p.beepSequence(3, 100*time.Millisecond)
	case SoundError, SoundMistake:
		p.beepSequence(2, 50*time.Millisecond)
	case SoundAchievement:
		p.beepSequence(5, 80*time.Millisecond)
	case SoundGoalComplete:
		p.beepSequence(4, 100*time.Millisecond)
	}
}

// beepSequence plays multiple beeps
func (p *Player) beepSequence(count int, delay time.Duration) {
	for i := 0; i < count; i++ {
		fmt.Print("\a")
		time.Sleep(delay)
	}
}

// systemSound plays a system sound by name
func (p *Player) systemSound(name string) {
	switch runtime.GOOS {
	case "linux":
		// Try freedesktop sounds first
		sounds := map[string][]string{
			"message":  {"/usr/share/sounds/freedesktop/stereo/message.oga", "/usr/share/sounds/alsa/Front_Center.wav"},
			"complete": {"/usr/share/sounds/freedesktop/stereo/complete.oga", "/usr/share/sounds/alsa/Front_Center.wav"},
			"error":    {"/usr/share/sounds/freedesktop/stereo/dialog-error.oga"},
		}

		if paths, ok := sounds[name]; ok {
			for _, path := range paths {
				if err := exec.Command("paplay", path).Run(); err == nil {
					return
				}
				if err := exec.Command("aplay", "-q", path).Run(); err == nil {
					return
				}
			}
		}
		// Fallback to terminal bell
		fmt.Print("\a")

	case "darwin": // macOS
		sounds := map[string]string{
			"message":  "/System/Library/Sounds/Tink.aiff",
			"complete": "/System/Library/Sounds/Glass.aiff",
			"error":    "/System/Library/Sounds/Basso.aiff",
		}
		if path, ok := sounds[name]; ok {
			exec.Command("afplay", path).Run()
		}

	case "windows":
		// Windows uses system beep
		fmt.Print("\a")

	default:
		fmt.Print("\a")
	}
}

// PlayKeystroke plays a subtle keystroke sound
func (p *Player) PlayKeystroke() {
	p.Play(SoundKeystroke)
}

// PlaySuccess plays a success sound
func (p *Player) PlaySuccess() {
	p.Play(SoundSuccess)
}

// PlayError plays an error sound
func (p *Player) PlayError() {
	p.Play(SoundError)
}

// PlayMistake plays a mistake sound (different from error)
func (p *Player) PlayMistake() {
	p.Play(SoundMistake)
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

// SetProfile sets the sound profile
func (p *Player) SetProfile(profile Profile) {
	p.Profile = profile
}

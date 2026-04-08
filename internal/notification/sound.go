package notification

import (
	"fmt"
	"os/exec"
	"runtime"
)

// SoundPlayer plays notification sounds.
type SoundPlayer struct {
	enabled bool
}

// NewSoundPlayer creates a new sound player.
func NewSoundPlayer(enabled bool) *SoundPlayer {
	return &SoundPlayer{enabled: enabled}
}

// Play plays the notification sound.
func (s *SoundPlayer) Play() {
	if !s.enabled {
		return
	}
	go s.playSound()
}

func (s *SoundPlayer) playSound() {
	switch runtime.GOOS {
	case "linux":
		// Try paplay with system sound.
		cmd := exec.Command("paplay", "/usr/share/sounds/freedesktop/stereo/message.oga")
		if err := cmd.Run(); err != nil {
			// Fallback: canberra-gtk-play.
			cmd = exec.Command("canberra-gtk-play", "-i", "message-new-instant")
			if err := cmd.Run(); err != nil {
				// Last resort: terminal bell.
				fmt.Print("\a")
			}
		}
	case "darwin":
		cmd := exec.Command("afplay", "/System/Library/Sounds/Ping.aiff")
		cmd.Run()
	default:
		fmt.Print("\a")
	}
}

package media

import (
	"fmt"
	"os/exec"
	"sync"
)

// PlaybackState represents the state of audio/video playback.
type PlaybackState int

const (
	PlaybackStopped PlaybackState = iota
	PlaybackPlaying
	PlaybackPaused
)

// VoicePlayer handles voice message playback.
type VoicePlayer struct {
	mu      sync.Mutex
	cmd     *exec.Cmd
	state   PlaybackState
	player  string
}

// NewVoicePlayer creates a new voice player using the specified player binary.
func NewVoicePlayer(player string) *VoicePlayer {
	return &VoicePlayer{
		player: player,
	}
}

// Play starts playing a voice note file.
func (v *VoicePlayer) Play(filePath string) error {
	v.mu.Lock()
	defer v.mu.Unlock()

	// Stop any current playback.
	v.stopLocked()

	var cmd *exec.Cmd
	switch v.player {
	case "mpv":
		cmd = exec.Command("mpv", "--no-video", "--really-quiet", filePath)
	case "ffplay":
		cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", filePath)
	case "paplay":
		cmd = exec.Command("paplay", filePath)
	default:
		cmd = exec.Command("mpv", "--no-video", "--really-quiet", filePath)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("starting voice player: %w", err)
	}

	v.cmd = cmd
	v.state = PlaybackPlaying

	// Wait for completion in background.
	go func() {
		cmd.Wait()
		v.mu.Lock()
		defer v.mu.Unlock()
		if v.cmd == cmd {
			v.state = PlaybackStopped
			v.cmd = nil
		}
	}()

	return nil
}

// Stop stops the current playback.
func (v *VoicePlayer) Stop() {
	v.mu.Lock()
	defer v.mu.Unlock()
	v.stopLocked()
}

func (v *VoicePlayer) stopLocked() {
	if v.cmd != nil && v.cmd.Process != nil {
		v.cmd.Process.Kill()
		v.cmd = nil
		v.state = PlaybackStopped
	}
}

// State returns the current playback state.
func (v *VoicePlayer) State() PlaybackState {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.state
}

// IsPlaying returns true if audio is currently playing.
func (v *VoicePlayer) IsPlaying() bool {
	return v.State() == PlaybackPlaying
}

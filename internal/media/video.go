package media

import (
	"fmt"
	"os/exec"
	"runtime"
)

// VideoPlayer launches videos in an external player.
type VideoPlayer struct {
	player string
}

// NewVideoPlayer creates a new video player launcher.
func NewVideoPlayer(player string) *VideoPlayer {
	return &VideoPlayer{
		player: player,
	}
}

// Play launches a video file in the external player.
func (v *VideoPlayer) Play(filePath string) error {
	var cmd *exec.Cmd

	switch v.player {
	case "mpv":
		cmd = exec.Command("mpv", filePath)
	case "vlc":
		cmd = exec.Command("vlc", filePath)
	case "xdg-open":
		cmd = exec.Command("xdg-open", filePath)
	default:
		cmd = exec.Command(v.defaultPlayer(), filePath)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("launching video player: %w", err)
	}

	// Don't wait — let it run independently.
	go cmd.Wait()

	return nil
}

func (v *VideoPlayer) defaultPlayer() string {
	switch runtime.GOOS {
	case "darwin":
		return "open"
	case "windows":
		return "start"
	default:
		return "xdg-open"
	}
}

package notification

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notifier sends desktop notifications.
type Notifier struct {
	enabled     bool
	showPreview bool
}

// NewNotifier creates a new notification dispatcher.
func NewNotifier(enabled, showPreview bool) *Notifier {
	return &Notifier{
		enabled:     enabled,
		showPreview: showPreview,
	}
}

// Notify sends a desktop notification.
func (n *Notifier) Notify(title, body string) {
	if !n.enabled {
		return
	}

	if !n.showPreview {
		body = "New message"
	}

	go n.send(title, body)
}

func (n *Notifier) send(title, body string) {
	switch runtime.GOOS {
	case "linux":
		n.sendLinux(title, body)
	case "darwin":
		n.sendMacOS(title, body)
	default:
		// Unsupported platform.
	}
}

func (n *Notifier) sendLinux(title, body string) {
	// Try notify-send first.
	cmd := exec.Command("notify-send",
		"--app-name=Tele-TUI",
		"--icon=telegram",
		"--urgency=normal",
		title,
		body,
	)
	if err := cmd.Run(); err != nil {
		// Fallback: terminal bell.
		fmt.Print("\a")
	}
}

func (n *Notifier) sendMacOS(title, body string) {
	script := fmt.Sprintf(
		`display notification %q with title %q`,
		body, title,
	)
	cmd := exec.Command("osascript", "-e", script)
	cmd.Run()
}

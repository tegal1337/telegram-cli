package widgets

import (
	"time"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// SpinnerTickMsg triggers the next spinner frame.
type SpinnerTickMsg struct{}

// Spinner is a loading spinner widget.
type Spinner struct {
	Frame  int
	Label  string
	Style  lipgloss.Style
	Active bool
}

// NewSpinner creates a new spinner.
func NewSpinner(label string) Spinner {
	return Spinner{
		Label:  label,
		Active: true,
	}
}

// Tick returns a command that triggers the next spinner frame.
func (s Spinner) Tick() tea.Cmd {
	if !s.Active {
		return nil
	}
	return tea.Tick(80*time.Millisecond, func(t time.Time) tea.Msg {
		return SpinnerTickMsg{}
	})
}

// Update advances the spinner frame.
func (s *Spinner) Update(msg tea.Msg) tea.Cmd {
	if _, ok := msg.(SpinnerTickMsg); ok && s.Active {
		s.Frame = (s.Frame + 1) % len(spinnerFrames)
		return s.Tick()
	}
	return nil
}

// View renders the spinner.
func (s *Spinner) View() string {
	if !s.Active {
		return ""
	}
	frame := s.Style.Render(spinnerFrames[s.Frame])
	if s.Label != "" {
		return frame + " " + s.Label
	}
	return frame
}

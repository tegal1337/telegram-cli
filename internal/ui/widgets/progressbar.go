package widgets

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ProgressBar renders a download/upload progress bar.
type ProgressBar struct {
	Progress float64 // 0.0 to 1.0
	Width    int
	Label    string

	StyleBar  lipgloss.Style
	StyleFill lipgloss.Style
}

// NewProgressBar creates a new progress bar.
func NewProgressBar(width int) ProgressBar {
	return ProgressBar{
		Width: width,
	}
}

// View renders the progress bar.
func (p *ProgressBar) View() string {
	barWidth := p.Width - 8 // space for percentage
	if barWidth < 5 {
		barWidth = 5
	}

	filled := int(p.Progress * float64(barWidth))
	if filled > barWidth {
		filled = barWidth
	}
	empty := barWidth - filled

	bar := p.StyleFill.Render(strings.Repeat("█", filled)) +
		p.StyleBar.Render(strings.Repeat("░", empty))

	pct := fmt.Sprintf(" %3d%%", int(p.Progress*100))

	result := bar + pct
	if p.Label != "" {
		result = p.Label + "\n" + result
	}
	return result
}

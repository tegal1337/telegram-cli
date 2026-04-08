package widgets

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

// Tabs is a tab bar widget for switching between views.
type Tabs struct {
	Labels   []string
	Active   int
	Width    int

	StyleTab       lipgloss.Style
	StyleTabActive lipgloss.Style
}

// NewTabs creates a new tab bar.
func NewTabs(labels []string) Tabs {
	return Tabs{
		Labels: labels,
	}
}

// Update handles tab switching via left/right keys or number keys.
func (t *Tabs) Update(msg tea.Msg) bool {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "left", "h":
			if t.Active > 0 {
				t.Active--
				return true
			}
		case "right", "l":
			if t.Active < len(t.Labels)-1 {
				t.Active++
				return true
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			idx := int(msg.String()[0]-'0') - 1
			if idx >= 0 && idx < len(t.Labels) {
				t.Active = idx
				return true
			}
		case "tab":
			t.Active = (t.Active + 1) % len(t.Labels)
			return true
		}
	}
	return false
}

// View renders the tab bar.
func (t *Tabs) View() string {
	var tabs []string
	for i, label := range t.Labels {
		if i == t.Active {
			tabs = append(tabs, t.StyleTabActive.Render(label))
		} else {
			tabs = append(tabs, t.StyleTab.Render(label))
		}
	}
	return lipgloss.NewStyle().Width(t.Width).Render(strings.Join(tabs, " "))
}

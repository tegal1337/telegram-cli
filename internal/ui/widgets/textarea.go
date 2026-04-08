package widgets

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
)

// TextArea is a simple multi-line text input widget.
type TextArea struct {
	Value       string
	Cursor      int
	Width       int
	Height      int
	Focused     bool
	Placeholder string
	Style       lipgloss.Style
}

// NewTextArea creates a new text area widget.
func NewTextArea() TextArea {
	return TextArea{
		Height: 3,
	}
}

// Update handles key events for text input.
func (t *TextArea) Update(msg tea.Msg) (submitted bool) {
	if !t.Focused {
		return false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "enter":
			return true
		case "backspace":
			if t.Cursor > 0 {
				runes := []rune(t.Value)
				t.Value = string(runes[:t.Cursor-1]) + string(runes[t.Cursor:])
				t.Cursor--
			}
		case "delete":
			runes := []rune(t.Value)
			if t.Cursor < len(runes) {
				t.Value = string(runes[:t.Cursor]) + string(runes[t.Cursor+1:])
			}
		case "left":
			if t.Cursor > 0 {
				t.Cursor--
			}
		case "right":
			if t.Cursor < len([]rune(t.Value)) {
				t.Cursor++
			}
		case "home", "ctrl+a":
			t.Cursor = 0
		case "end", "ctrl+e":
			t.Cursor = len([]rune(t.Value))
		case "ctrl+u":
			t.Value = string([]rune(t.Value)[t.Cursor:])
			t.Cursor = 0
		case "ctrl+k":
			t.Value = string([]rune(t.Value)[:t.Cursor])
		default:
			// Insert character.
			if len(msg.String()) == 1 || msg.String() == " " {
				runes := []rune(t.Value)
				char := []rune(msg.String())
				newRunes := make([]rune, 0, len(runes)+len(char))
				newRunes = append(newRunes, runes[:t.Cursor]...)
				newRunes = append(newRunes, char...)
				newRunes = append(newRunes, runes[t.Cursor:]...)
				t.Value = string(newRunes)
				t.Cursor += len(char)
			}
		}
	}
	return false
}

// Reset clears the text area.
func (t *TextArea) Reset() {
	t.Value = ""
	t.Cursor = 0
}

// View renders the text area.
func (t *TextArea) View() string {
	content := t.Value
	if content == "" && t.Placeholder != "" && !t.Focused {
		content = lipgloss.NewStyle().Foreground(lipgloss.Color("#565F89")).Render(t.Placeholder)
	}

	if t.Focused && content == t.Value {
		// Show cursor.
		runes := []rune(t.Value)
		before := string(runes[:t.Cursor])
		cursor := "▏"
		after := ""
		if t.Cursor < len(runes) {
			after = string(runes[t.Cursor:])
		}
		content = before + cursor + after
	}

	lines := strings.Split(content, "\n")
	if len(lines) > t.Height {
		lines = lines[len(lines)-t.Height:]
	}

	return t.Style.Width(t.Width).Height(t.Height).Render(strings.Join(lines, "\n"))
}

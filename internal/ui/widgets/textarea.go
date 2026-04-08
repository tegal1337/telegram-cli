package widgets

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
)

type TextArea struct {
	Value       string
	Cursor      int
	Width       int
	Height      int
	Focused     bool
	Placeholder string
	Style       lipgloss.Style
}

func NewTextArea() TextArea {
	return TextArea{
		Height: 1,
	}
}

func (t *TextArea) Update(msg tea.Msg) (submitted bool) {
	if !t.Focused {
		return false
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		key := msg.String()
		switch key {
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
		case "ctrl+w":
			// Delete word backwards
			runes := []rune(t.Value)
			i := t.Cursor - 1
			for i >= 0 && runes[i] == ' ' {
				i--
			}
			for i >= 0 && runes[i] != ' ' {
				i--
			}
			t.Value = string(runes[:i+1]) + string(runes[t.Cursor:])
			t.Cursor = i + 1
		default:
			// Map bubbletea v2 key names to actual characters
			text := key
			switch key {
			case "space":
				text = " "
			case "tab", "escape", "up", "down", "pgup", "pgdown",
				"shift+tab", "f1", "f2", "f3", "f4", "f5",
				"f6", "f7", "f8", "f9", "f10", "f11", "f12":
				return false
			}
			if strings.HasPrefix(text, "ctrl+") || strings.HasPrefix(text, "alt+") || strings.HasPrefix(text, "shift+") {
				return false
			}
			if len(text) >= 1 {
				runes := []rune(t.Value)
				insert := []rune(text)
				newRunes := make([]rune, 0, len(runes)+len(insert))
				newRunes = append(newRunes, runes[:t.Cursor]...)
				newRunes = append(newRunes, insert...)
				newRunes = append(newRunes, runes[t.Cursor:]...)
				t.Value = string(newRunes)
				t.Cursor += len(insert)
			}
		}

	case tea.PasteMsg:
		// Handle bracketed paste
		text := msg.Content
		runes := []rune(t.Value)
		insert := []rune(text)
		newRunes := make([]rune, 0, len(runes)+len(insert))
		newRunes = append(newRunes, runes[:t.Cursor]...)
		newRunes = append(newRunes, insert...)
		newRunes = append(newRunes, runes[t.Cursor:]...)
		t.Value = string(newRunes)
		t.Cursor += len(insert)
	}
	return false
}

func (t *TextArea) Reset() {
	t.Value = ""
	t.Cursor = 0
}

func (t *TextArea) View() string {
	content := t.Value
	if content == "" && t.Placeholder != "" && !t.Focused {
		content = lipgloss.NewStyle().Foreground(lipgloss.Color("#565F89")).Render(t.Placeholder)
	}

	if t.Focused && content == t.Value {
		runes := []rune(t.Value)
		before := string(runes[:t.Cursor])
		after := ""
		if t.Cursor < len(runes) {
			after = string(runes[t.Cursor:])
		}
		content = before + "█" + after
	}

	return t.Style.Width(t.Width).Render(content)
}

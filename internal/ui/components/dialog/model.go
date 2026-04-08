package dialog

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
)

// DialogResultMsg is emitted when a dialog is closed.
type DialogResultMsg struct {
	ID        string
	Confirmed bool
	Input     string
}

// Kind represents the type of dialog.
type Kind int

const (
	KindConfirm Kind = iota
	KindPrompt
	KindAlert
)

// Model is a modal dialog component.
type Model struct {
	theme       *theme.Theme
	visible     bool
	kind        Kind
	id          string
	title       string
	message     string
	input       string
	cursor      int
	buttonIdx   int
	buttons     []string
	width       int
	height      int
}

// NewConfirm creates a confirmation dialog.
func NewConfirm(th *theme.Theme, id, title, message string) Model {
	return Model{
		theme:   th,
		visible: true,
		kind:    KindConfirm,
		id:      id,
		title:   title,
		message: message,
		buttons: []string{"Cancel", "Confirm"},
	}
}

// NewAlert creates an alert dialog.
func NewAlert(th *theme.Theme, id, title, message string) Model {
	return Model{
		theme:   th,
		visible: true,
		kind:    KindAlert,
		id:      id,
		title:   title,
		message: message,
		buttons: []string{"OK"},
	}
}

// NewPrompt creates a prompt dialog with text input.
func NewPrompt(th *theme.Theme, id, title, message string) Model {
	return Model{
		theme:   th,
		visible: true,
		kind:    KindPrompt,
		id:      id,
		title:   title,
		message: message,
		buttons: []string{"Cancel", "OK"},
	}
}

// IsVisible returns whether the dialog is visible.
func (m Model) IsVisible() bool {
	return m.visible
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "escape":
			m.visible = false
			return m, func() tea.Msg {
				return DialogResultMsg{ID: m.id, Confirmed: false}
			}

		case "tab", "left", "right":
			m.buttonIdx = (m.buttonIdx + 1) % len(m.buttons)

		case "enter":
			m.visible = false
			confirmed := m.buttonIdx == len(m.buttons)-1
			return m, func() tea.Msg {
				return DialogResultMsg{
					ID:        m.id,
					Confirmed: confirmed,
					Input:     m.input,
				}
			}

		case "backspace":
			if m.kind == KindPrompt && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}

		default:
			if m.kind == KindPrompt && len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}
	}

	return m, nil
}

// View renders the dialog.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	title := m.theme.DialogTitle.Render(m.title)
	message := lipgloss.NewStyle().Foreground(m.theme.Text).Render(m.message)

	var content string
	if m.kind == KindPrompt {
		inputStyle := m.theme.AuthInput.Width(30)
		inputText := m.input + "▏"
		content = fmt.Sprintf("%s\n\n%s\n\n%s", title, message, inputStyle.Render(inputText))
	} else {
		content = fmt.Sprintf("%s\n\n%s", title, message)
	}

	// Buttons
	var buttonRow string
	for i, label := range m.buttons {
		style := m.theme.DialogButton
		if i == m.buttonIdx {
			style = m.theme.DialogButtonActive
		}
		if i > 0 {
			buttonRow += "  "
		}
		buttonRow += style.Render(label)
	}

	content += "\n\n" + lipgloss.NewStyle().Align(lipgloss.Center).Render(buttonRow)

	return m.theme.DialogBox.Render(content)
}

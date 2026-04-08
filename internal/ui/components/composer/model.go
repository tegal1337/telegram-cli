package composer

import (
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/tegal1337/telegram-cli/internal/ui/widgets"
)

// Mode represents the composer's current mode.
type Mode int

const (
	ModeNormal Mode = iota
	ModeReply
	ModeEdit
)

// Model is the message composer component.
type Model struct {
	textarea   widgets.TextArea
	theme      *theme.Theme
	width      int
	height     int
	focused    bool
	mode       Mode
	chatID     int64
	replyToID  int64
	editMsgID  int64
	replyText  string
	attachment string
}

// New creates a new composer model.
func New(th *theme.Theme) Model {
	ta := widgets.NewTextArea()
	ta.Placeholder = "Type a message..."
	ta.Style = th.ComposerInput

	return Model{
		textarea: ta,
		theme:    th,
		height:   3,
	}
}

// SetSize sets the component dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.textarea.Width = width - 2
	m.textarea.Height = height - 1
}

// SetFocused sets focus state.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
	m.textarea.Focused = focused
}

// SetChatId sets the active chat for the composer.
func (m *Model) SetChatId(chatID int64) {
	m.chatID = chatID
	m.Reset()
}

// EnterReplyMode starts replying to a message.
func (m *Model) EnterReplyMode(messageID int64, previewText string) {
	m.mode = ModeReply
	m.replyToID = messageID
	m.replyText = previewText
}

// EnterEditMode starts editing a message.
func (m *Model) EnterEditMode(messageID int64, currentText string) {
	m.mode = ModeEdit
	m.editMsgID = messageID
	m.textarea.Value = currentText
	m.textarea.Cursor = len([]rune(currentText))
}

// Reset clears the composer state.
func (m *Model) Reset() {
	m.textarea.Reset()
	m.mode = ModeNormal
	m.replyToID = 0
	m.editMsgID = 0
	m.replyText = ""
	m.attachment = ""
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.focused {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "escape":
			if m.mode != ModeNormal {
				m.Reset()
				return m, nil
			}
		case "enter":
			if m.textarea.Value != "" {
				text := m.textarea.Value
				submitted := MessageSubmittedMsg{
					ChatId: m.chatID,
					Text:   text,
				}

				switch m.mode {
				case ModeReply:
					submitted.ReplyToId = m.replyToID
				case ModeEdit:
					submitted.EditMessageId = m.editMsgID
				}

				m.Reset()
				return m, func() tea.Msg { return submitted }
			}
		default:
			m.textarea.Update(msg)
		}
	}

	return m, nil
}

// View renders the composer.
func (m Model) View() string {
	var parts []string

	// Reply/edit bar.
	switch m.mode {
	case ModeReply:
		replyBar := m.theme.ComposerReplyBar.
			Width(m.width).
			Render(fmt.Sprintf("↩ Reply: %s", truncate(m.replyText, m.width-12)))
		parts = append(parts, replyBar)
	case ModeEdit:
		editBar := m.theme.ComposerReplyBar.
			Width(m.width).
			Render("✏ Editing message")
		parts = append(parts, editBar)
	}

	// Attachment indicator.
	if m.attachment != "" {
		attBar := m.theme.ComposerHint.Render(fmt.Sprintf("📎 %s", m.attachment))
		parts = append(parts, attBar)
	}

	// Input area.
	input := m.textarea.View()
	parts = append(parts, input)

	// Hint.
	hint := m.theme.ComposerHint.Render("Enter: send | Esc: cancel | Ctrl+A: attach")
	parts = append(parts, hint)

	content := lipgloss.JoinVertical(lipgloss.Left, parts...)
	return m.theme.ComposerPane.Width(m.width).Render(content)
}

func truncate(s string, maxWidth int) string {
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	if maxWidth <= 3 {
		return string(runes[:maxWidth])
	}
	return string(runes[:maxWidth-3]) + "..."
}

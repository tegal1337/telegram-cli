package chatview

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/tegal1337/telegram-cli/internal/render"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/zelenin/go-tdlib/client"
)

// Model is the chat view component showing messages.
type Model struct {
	store       *store.Store
	tg          *telegram.Client
	theme       *theme.Theme
	renderer    *render.MessageRenderer
	width       int
	height      int
	headerH     int
	focused     bool
	chatID      int64
	chatTitle   string
	scrollOffset int
	cursor      int
	loading     bool
	myUserID    int64
}

// New creates a new chat view model.
func New(s *store.Store, tg *telegram.Client, th *theme.Theme) Model {
	return Model{
		store:    s,
		tg:       tg,
		theme:    th,
		renderer: render.NewMessageRenderer(th),
		headerH:  1,
	}
}

// SetSize sets the component dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
}

// SetFocused sets focus state.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
}

// SetMyUserID sets the current user's ID for message alignment.
func (m *Model) SetMyUserID(id int64) {
	m.myUserID = id
}

// OpenChat switches to a new chat.
func (m *Model) OpenChat(chatID int64, title string) tea.Cmd {
	m.chatID = chatID
	m.chatTitle = title
	m.scrollOffset = 0
	m.cursor = 0
	m.loading = true

	return tea.Batch(
		m.loadHistoryCmd(chatID, 0),
		m.openChatCmd(chatID),
	)
}

func (m *Model) openChatCmd(chatID int64) tea.Cmd {
	return func() tea.Msg {
		m.tg.OpenChat(context.Background(), chatID)
		return nil
	}
}

type historyLoadedMsg struct {
	ChatID   int64
	Messages []*client.Message
	Err      error
}

func (m *Model) loadHistoryCmd(chatID int64, fromMessageID int64) tea.Cmd {
	return func() tea.Msg {
		msgs, err := m.tg.GetChatHistory(context.Background(), chatID, fromMessageID, 0, 30)
		if err != nil {
			return historyLoadedMsg{ChatID: chatID, Err: err}
		}
		return historyLoadedMsg{ChatID: chatID, Messages: msgs.Messages}
	}
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case historyLoadedMsg:
		if msg.ChatID != m.chatID {
			return m, nil
		}
		m.loading = false
		if msg.Err == nil && len(msg.Messages) > 0 {
			// Messages come newest-first from TDLib; reverse for display.
			reversed := make([]*client.Message, len(msg.Messages))
			for i, m := range msg.Messages {
				reversed[len(msg.Messages)-1-i] = m
			}
			m.store.Messages.Prepend(m.chatID, reversed)
		}

	case telegram.NewMessageMsg:
		if msg.Message.ChatID == m.chatID {
			m.store.Messages.Append(m.chatID, msg.Message)
			// Auto-scroll to bottom if already at bottom.
			if m.scrollOffset == 0 {
				// Already at bottom, stay there.
			}
			// Mark as read.
			return m, m.viewMessagesCmd(m.chatID, []int64{msg.Message.ID})
		}

	case telegram.MessageEditedMsg:
		if msg.ChatID == m.chatID {
			// Re-fetch the edited message.
			return m, m.fetchMessageCmd(msg.ChatID, msg.MessageID)
		}

	case telegram.MessageDeletedMsg:
		if msg.ChatID == m.chatID {
			m.store.Messages.Delete(m.chatID, msg.MessageIDs)
		}

	case telegram.MessageSendSucceededMsg:
		if msg.Message.ChatID == m.chatID {
			m.store.Messages.ReplaceMessageID(m.chatID, msg.OldMessageID, msg.Message)
		}

	case messageFetchedMsg:
		if msg.chatID == m.chatID && msg.message != nil {
			m.store.Messages.UpdateMessage(m.chatID, msg.message.ID, msg.message)
		}

	case tea.KeyPressMsg:
		if m.focused {
			return m.handleKey(msg)
		}
	}

	return m, nil
}

type messageFetchedMsg struct {
	chatID  int64
	message *client.Message
}

func (m *Model) fetchMessageCmd(chatID, messageID int64) tea.Cmd {
	return func() tea.Msg {
		msg, err := m.tg.GetMessage(context.Background(), chatID, messageID)
		if err != nil {
			return nil
		}
		return messageFetchedMsg{chatID: chatID, message: msg}
	}
}

func (m *Model) viewMessagesCmd(chatID int64, messageIDs []int64) tea.Cmd {
	return func() tea.Msg {
		m.tg.ViewMessages(context.Background(), chatID, messageIDs)
		return nil
	}
}

func (m Model) handleKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	messages := m.store.Messages.Get(m.chatID)

	switch msg.String() {
	case "up", "k":
		if m.cursor > 0 {
			m.cursor--
		} else if len(messages) > 0 {
			// Load more history.
			oldestID := m.store.Messages.OldestMessageID(m.chatID)
			return m, m.loadHistoryCmd(m.chatID, oldestID)
		}
	case "down", "j":
		if m.cursor < len(messages)-1 {
			m.cursor++
		}
	case "G":
		m.cursor = max(0, len(messages)-1)
		m.scrollOffset = 0
	case "g":
		m.cursor = 0
	case "r":
		if m.cursor >= 0 && m.cursor < len(messages) {
			return m, func() tea.Msg {
				return MessageActionMsg{
					Action:    "reply",
					ChatID:    m.chatID,
					MessageID: messages[m.cursor].ID,
				}
			}
		}
	case "e":
		if m.cursor >= 0 && m.cursor < len(messages) {
			msg := messages[m.cursor]
			if isOwnMessage(msg, m.myUserID) {
				return m, func() tea.Msg {
					return MessageActionMsg{
						Action:    "edit",
						ChatID:    m.chatID,
						MessageID: msg.ID,
					}
				}
			}
		}
	case "d":
		if m.cursor >= 0 && m.cursor < len(messages) {
			msg := messages[m.cursor]
			return m, func() tea.Msg {
				return MessageActionMsg{
					Action:    "delete",
					ChatID:    m.chatID,
					MessageID: msg.ID,
				}
			}
		}
	case "f":
		if m.cursor >= 0 && m.cursor < len(messages) {
			return m, func() tea.Msg {
				return MessageActionMsg{
					Action:    "forward",
					ChatID:    m.chatID,
					MessageID: messages[m.cursor].ID,
				}
			}
		}
	case "ctrl+u":
		m.scrollOffset += m.height / 2
	case "ctrl+d":
		m.scrollOffset -= m.height / 2
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	}

	return m, nil
}

func isOwnMessage(msg *client.Message, myUserID int64) bool {
	if sender, ok := msg.SenderId.(*client.MessageSenderUser); ok {
		return sender.UserID == myUserID
	}
	return false
}

// View renders the chat view.
func (m Model) View() string {
	if m.chatID == 0 {
		return m.renderEmpty()
	}

	header := m.renderHeader()
	messages := m.renderMessages()

	return lipgloss.JoinVertical(lipgloss.Left, header, messages)
}

func (m Model) renderEmpty() string {
	return m.theme.ChatViewPane.
		Width(m.width).
		Height(m.height + m.headerH).
		Align(lipgloss.Center, lipgloss.Center).
		Render("Select a chat to start messaging")
}

func (m Model) renderHeader() string {
	title := m.chatTitle
	if m.loading {
		title += " (loading...)"
	}

	// Show typing indicators if any.
	return m.theme.ChatViewHeader.
		Width(m.width).
		Render(title)
}

func (m Model) renderMessages() string {
	messages := m.store.Messages.Get(m.chatID)

	if len(messages) == 0 {
		return m.theme.ChatViewPane.
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("No messages yet")
	}

	var lines []string
	for i, msg := range messages {
		isOwn := isOwnMessage(msg, m.myUserID)
		isSelected := i == m.cursor && m.focused
		rendered := m.renderer.RenderMessage(msg, m.store, isOwn, isSelected, m.width-4)
		lines = append(lines, rendered)
	}

	content := strings.Join(lines, "\n")

	// Simple scroll: show the last N lines that fit.
	contentLines := strings.Split(content, "\n")
	if len(contentLines) > m.height {
		start := len(contentLines) - m.height - m.scrollOffset
		if start < 0 {
			start = 0
		}
		end := start + m.height
		if end > len(contentLines) {
			end = len(contentLines)
		}
		contentLines = contentLines[start:end]
	}

	return m.theme.ChatViewPane.
		Width(m.width).
		Height(m.height).
		Render(strings.Join(contentLines, "\n"))
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

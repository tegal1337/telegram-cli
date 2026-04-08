package statusbar

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/zelenin/go-tdlib/client"
)

// Model is the status bar component.
type Model struct {
	store          *store.Store
	theme          *theme.Theme
	width          int
	connected      bool
	userName       string
	typing         map[int64][]int64 // chatID -> userIDs typing
	unreadCount    int32
	activeChatID   int64
}

// New creates a new status bar model.
func New(s *store.Store, th *theme.Theme) Model {
	return Model{
		store:     s,
		theme:     th,
		connected: false,
		typing:    make(map[int64][]int64),
	}
}

// SetSize sets the component width.
func (m *Model) SetSize(width int) {
	m.width = width
}

// SetUserName sets the current user's display name.
func (m *Model) SetUserName(name string) {
	m.userName = name
}

// SetActiveChatID sets the currently viewed chat.
func (m *Model) SetActiveChatID(chatID int64) {
	m.activeChatID = chatID
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case telegram.ConnectionStateMsg:
		switch msg.State.(type) {
		case *client.ConnectionStateReady:
			m.connected = true
		default:
			m.connected = false
		}

	case telegram.ChatActionMsg:
		if msg.UserID != 0 {
			users := m.typing[msg.ChatID]
			// Add user if typing, remove if stopped.
			switch msg.Action.(type) {
			case *client.ChatActionTyping:
				found := false
				for _, uid := range users {
					if uid == msg.UserID {
						found = true
						break
					}
				}
				if !found {
					m.typing[msg.ChatID] = append(users, msg.UserID)
				}
			case *client.ChatActionCancel:
				filtered := users[:0]
				for _, uid := range users {
					if uid != msg.UserID {
						filtered = append(filtered, uid)
					}
				}
				m.typing[msg.ChatID] = filtered
			}
		}

	case telegram.UnreadCountMsg:
		m.unreadCount = msg.UnreadCount
	}

	return m, nil
}

// TypingIndicator returns a typing indicator string for the active chat.
func (m Model) TypingIndicator() string {
	users, ok := m.typing[m.activeChatID]
	if !ok || len(users) == 0 {
		return ""
	}

	var names []string
	for _, uid := range users {
		names = append(names, m.store.Users.DisplayName(uid))
	}

	if len(names) == 1 {
		return fmt.Sprintf("%s is typing...", names[0])
	}
	return fmt.Sprintf("%s are typing...", strings.Join(names, ", "))
}

// View renders the status bar.
func (m Model) View() string {
	// Connection status
	connStatus := m.theme.StatusBarConnected.Render("● Connected")
	if !m.connected {
		connStatus = m.theme.StatusBar.Foreground(m.theme.Error).Render("● Disconnected")
	}

	// Typing indicator
	typingText := ""
	if indicator := m.TypingIndicator(); indicator != "" {
		typingText = m.theme.StatusBarTyping.Render(indicator)
	}

	// User name
	userName := m.theme.StatusBar.Render(m.userName)

	// Unread count
	unread := ""
	if m.unreadCount > 0 {
		unread = m.theme.StatusBar.Foreground(m.theme.Primary).
			Render(fmt.Sprintf(" [%d unread]", m.unreadCount))
	}

	// Keybind hints
	hints := m.theme.StatusBar.Foreground(m.theme.TextMuted).
		Render("Ctrl+1:Chats  Ctrl+2:View  Ctrl+3:Compose  /:Search  Ctrl+K:Contacts")

	left := fmt.Sprintf("%s  %s%s", connStatus, userName, unread)
	center := typingText
	right := hints

	// Calculate padding
	leftW := lipgloss.Width(left)
	centerW := lipgloss.Width(center)
	rightW := lipgloss.Width(right)
	padding := m.width - leftW - centerW - rightW
	if padding < 0 {
		padding = 0
	}

	pad1 := padding / 2
	pad2 := padding - pad1

	return m.theme.StatusBar.
		Width(m.width).
		Render(left + strings.Repeat(" ", pad1) + center + strings.Repeat(" ", pad2) + right)
}

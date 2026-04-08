package chatlist

import (
	"context"
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/tegal1337/telegram-cli/internal/ui/widgets"
	"github.com/zelenin/go-tdlib/client"
)

// Model is the chat list component.
type Model struct {
	list     widgets.List
	store    *store.Store
	tg       *telegram.Client
	theme    *theme.Theme
	width    int
	height   int
	focused  bool
	filter   string
	loading  bool
	spinner  widgets.Spinner
	activeChatID int64
}

// New creates a new chat list model.
func New(s *store.Store, tg *telegram.Client, th *theme.Theme) Model {
	l := widgets.NewList()
	l.StyleNormal = th.ChatListItem
	l.StyleActive = th.ChatListItemActive
	l.StyleTitle = th.ChatListTitle
	l.StyleSub = th.ChatListPreview
	l.StyleMeta = th.ChatListTime
	l.StyleBadge = th.ChatListUnread
	l.StyleOnline = th.ChatListOnline

	sp := widgets.NewSpinner("Loading chats...")
	sp.Style = th.Spinner

	return Model{
		list:    l,
		store:   s,
		tg:      tg,
		theme:   th,
		loading: true,
		spinner: sp,
	}
}

// Init loads the initial chat list.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.loadChatsCmd(),
		m.spinner.Tick(),
	)
}

func (m Model) loadChatsCmd() tea.Cmd {
	return func() tea.Msg {
		err := m.tg.LoadChats(context.Background(), &client.ChatListMain{}, 50)
		if err != nil {
			return chatsLoadedMsg{err: err}
		}
		return chatsLoadedMsg{}
	}
}

type chatsLoadedMsg struct {
	err error
}

// SetSize sets the component dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.list.Width = width
	m.list.Height = height
}

// SetFocused sets focus state.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
	m.list.Focused = focused
}

// ActiveChatID returns the currently selected chat ID.
func (m *Model) ActiveChatID() int64 {
	return m.activeChatID
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case chatsLoadedMsg:
		m.loading = false
		m.spinner.Active = false
		m.refreshList()

	case telegram.ChatLastMessageMsg:
		m.store.Chats.UpdateLastMessage(msg.ChatID, msg.LastMessage, msg.Positions)
		m.refreshList()

	case telegram.ChatPositionMsg:
		m.store.Chats.UpdatePosition(msg.ChatID, msg.Positions)
		m.refreshList()

	case telegram.ChatReadInboxMsg:
		m.store.Chats.UpdateReadInbox(msg.ChatID, msg.UnreadCount)
		m.refreshList()

	case telegram.ChatUpdateMsg:
		if msg.Chat != nil {
			m.store.Chats.Set(msg.Chat)
			m.refreshList()
		}

	case telegram.NewMessageMsg:
		m.refreshList()

	case widgets.SpinnerTickMsg:
		cmd := m.spinner.Update(msg)
		cmds = append(cmds, cmd)

	case tea.KeyPressMsg:
		if m.focused {
			if selected := m.list.Update(msg); selected {
				item := m.list.SelectedItem()
				if item != nil {
					var chatID int64
					fmt.Sscanf(item.ID, "%d", &chatID)
					m.activeChatID = chatID
					return m, func() tea.Msg {
						return ChatSelectedMsg{ChatID: chatID}
					}
				}
			}
		}
	}

	return m, tea.Batch(cmds...)
}

func (m *Model) refreshList() {
	chats := m.store.Chats.OrderedChats()
	items := make([]widgets.ListItem, 0, len(chats))

	for _, entry := range chats {
		if entry.Chat == nil {
			continue
		}

		if m.filter != "" {
			if !strings.Contains(strings.ToLower(entry.Chat.Title), strings.ToLower(m.filter)) {
				continue
			}
		}

		preview := ""
		meta := ""
		if entry.LastMessage != nil {
			preview = messagePreview(entry.LastMessage)
			meta = formatTime(entry.LastMessage.Date)
		}

		badge := ""
		if entry.UnreadCount > 0 {
			badge = fmt.Sprintf("%d", entry.UnreadCount)
		}

		online := false
		if entry.Chat.Type != nil {
			if pt, ok := entry.Chat.Type.(*client.ChatTypePrivate); ok {
				online = m.store.Users.IsOnline(pt.UserID)
			}
		}

		items = append(items, widgets.ListItem{
			ID:       fmt.Sprintf("%d", entry.Chat.ID),
			Title:    chatIcon(entry.Chat) + " " + entry.Chat.Title,
			Subtitle: preview,
			Badge:    badge,
			Meta:     meta,
			Online:   online,
		})
	}

	m.list.SetItems(items)
}

func chatIcon(chat *client.Chat) string {
	switch chat.Type.(type) {
	case *client.ChatTypePrivate:
		return "👤"
	case *client.ChatTypeBasicGroup:
		return "👥"
	case *client.ChatTypeSupergroup:
		sg := chat.Type.(*client.ChatTypeSupergroup)
		if sg.IsChannel {
			return "📢"
		}
		return "👥"
	case *client.ChatTypeSecret:
		return "🔒"
	default:
		return "💬"
	}
}

func messagePreview(msg *client.Message) string {
	if msg == nil || msg.Content == nil {
		return ""
	}

	switch c := msg.Content.(type) {
	case *client.MessageText:
		text := c.Text.Text
		if len(text) > 50 {
			text = text[:50] + "..."
		}
		return text
	case *client.MessagePhoto:
		return "📷 Photo"
	case *client.MessageVideo:
		return "🎥 Video"
	case *client.MessageDocument:
		return "📎 " + c.Document.FileName
	case *client.MessageVoiceNote:
		return "🎤 Voice message"
	case *client.MessageVideoNote:
		return "📹 Video message"
	case *client.MessageSticker:
		return "🏷 " + c.Sticker.Emoji + " Sticker"
	case *client.MessageAnimation:
		return "🎬 GIF"
	case *client.MessageAudio:
		return "🎵 Audio"
	case *client.MessageLocation:
		return "📍 Location"
	case *client.MessageContact:
		return "👤 Contact"
	case *client.MessagePoll:
		return "📊 Poll"
	default:
		return "💬 Message"
	}
}

func formatTime(timestamp int32) string {
	// Simplified; the render package handles full formatting.
	if timestamp == 0 {
		return ""
	}
	return fmt.Sprintf("%02d:%02d", (timestamp/3600)%24, (timestamp/60)%60)
}

// View renders the chat list.
func (m Model) View() string {
	if m.loading {
		return m.theme.ChatListPane.
			Width(m.width).
			Height(m.height).
			Align(1, 1). // center
			Render(m.spinner.View())
	}

	return m.theme.ChatListPane.
		Width(m.width).
		Height(m.height).
		Render(m.list.View())
}

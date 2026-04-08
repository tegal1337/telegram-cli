package search

import (
	
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/tegal1337/telegram-cli/internal/ui/widgets"
	"github.com/zelenin/go-tdlib/client"
)

// SearchResultMsg is emitted when a search result is selected.
type SearchResultMsg struct {
	ChatId    int64
	MessageId int64
}

// Tab represents a search tab.
type Tab int

const (
	TabChats Tab = iota
	TabMessages
	TabGlobal
)

// Model is the search overlay component.
type Model struct {
	input   widgets.TextArea
	tabs    widgets.Tabs
	list    widgets.List
	store   *store.Store
	tg      *telegram.Client
	theme   *theme.Theme
	width   int
	height  int
	visible bool
	focused bool
	query   string
}

// New creates a new search model.
func New(s *store.Store, tg *telegram.Client, th *theme.Theme) Model {
	ta := widgets.NewTextArea()
	ta.Placeholder = "Search..."
	ta.Style = th.SearchInput
	ta.Focused = true

	tabs := widgets.NewTabs([]string{"Chats", "Messages", "Global"})
	tabs.StyleTab = th.Tab
	tabs.StyleTabActive = th.TabActive

	l := widgets.NewList()
	l.StyleNormal = th.SearchResult
	l.StyleActive = th.SearchResultActive
	l.StyleTitle = th.ChatListTitle
	l.StyleSub = th.ChatListPreview

	return Model{
		input: ta,
		tabs:  tabs,
		list:  l,
		store: s,
		tg:    tg,
		theme: th,
	}
}

// SetSize sets the component dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.input.Width = width - 4
	m.tabs.Width = width
	m.list.Width = width - 2
	m.list.Height = height - 6
}

// SetVisible shows or hides the search overlay.
func (m *Model) SetVisible(visible bool) {
	m.visible = visible
	if visible {
		m.input.Focused = true
		m.input.Reset()
		m.list.SetItems(nil)
	}
}

// IsVisible returns whether the search is visible.
func (m Model) IsVisible() bool {
	return m.visible
}

type searchResultsMsg struct {
	tab   Tab
	items []widgets.ListItem
}

func (m *Model) searchCmd() tea.Cmd {
	query := m.query
	tab := Tab(m.tabs.Active)

	return func() tea.Msg {
		var items []widgets.ListItem

		switch tab {
		case TabChats:
			chats, err := m.tg.SearchChats(query, 20)
			if err == nil {
				for _, chatID := range chats.ChatIds {
					chat, err := m.tg.GetChat(chatID)
					if err == nil {
						items = append(items, widgets.ListItem{
							ID:       fmt.Sprintf("%d", chat.Id),
							Title:    chat.Title,
							Subtitle: chatTypeLabel(chat),
						})
					}
				}
			}

		case TabMessages:
			found, err := m.tg.SearchMessages(query, 20)
			if err == nil {
				for _, msg := range found.Messages {
					chat, _ := m.tg.GetChat(msg.ChatId)
					title := fmt.Sprintf("%d", msg.ChatId)
					if chat != nil {
						title = chat.Title
					}
					items = append(items, widgets.ListItem{
						ID:       fmt.Sprintf("%d:%d", msg.ChatId, msg.Id),
						Title:    title,
						Subtitle: messagePreview(msg),
					})
				}
			}

		case TabGlobal:
			chats, err := m.tg.SearchChats(query, 20)
			if err == nil {
				for _, chatID := range chats.ChatIds {
					chat, err := m.tg.GetChat(chatID)
					if err == nil {
						items = append(items, widgets.ListItem{
							ID:       fmt.Sprintf("%d", chat.Id),
							Title:    chat.Title,
							Subtitle: chatTypeLabel(chat),
						})
					}
				}
			}
		}

		return searchResultsMsg{tab: tab, items: items}
	}
}

func chatTypeLabel(chat *client.Chat) string {
	switch chat.Type.(type) {
	case *client.ChatTypePrivate:
		return "Private chat"
	case *client.ChatTypeBasicGroup:
		return "Group"
	case *client.ChatTypeSupergroup:
		sg := chat.Type.(*client.ChatTypeSupergroup)
		if sg.IsChannel {
			return "Channel"
		}
		return "Supergroup"
	case *client.ChatTypeSecret:
		return "Secret chat"
	default:
		return "Chat"
	}
}

func messagePreview(msg *client.Message) string {
	if msg == nil || msg.Content == nil {
		return ""
	}
	if text, ok := msg.Content.(*client.MessageText); ok {
		t := text.Text.Text
		if len(t) > 60 {
			t = t[:60] + "..."
		}
		return t
	}
	return "Media message"
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case searchResultsMsg:
		m.list.SetItems(msg.items)
		m.list.Focused = true
		m.input.Focused = false

	case tea.KeyPressMsg:
		switch msg.String() {
		case "escape":
			m.visible = false
			return m, nil
		case "tab":
			m.tabs.Update(msg)
			if m.query != "" {
				return m, m.searchCmd()
			}
		case "enter":
			if m.input.Focused {
				m.query = m.input.Value
				return m, m.searchCmd()
			}
			if m.list.Focused {
				item := m.list.SelectedItem()
				if item != nil {
					var chatID, messageID int64
					n, _ := fmt.Sscanf(item.ID, "%d:%d", &chatID, &messageID)
					if n == 1 {
						fmt.Sscanf(item.ID, "%d", &chatID)
					}
					m.visible = false
					return m, func() tea.Msg {
						return SearchResultMsg{ChatId: chatID, MessageId: messageID}
					}
				}
			}
		default:
			if m.input.Focused {
				m.input.Update(msg)
			} else {
				m.list.Update(msg)
			}
		}
	}

	return m, nil
}

// View renders the search overlay.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	title := m.theme.AuthTitle.Render("Search")
	input := m.input.View()
	tabs := m.tabs.View()
	results := m.list.View()

	content := lipgloss.JoinVertical(lipgloss.Left,
		title, input, tabs, results,
	)

	return m.theme.DialogBox.
		Width(m.width - 4).
		Height(m.height - 4).
		Render(content)
}

package contacts

import (
	"context"
	"fmt"
	"sort"

	tea "charm.land/bubbletea/v2"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/tegal1337/telegram-cli/internal/ui/widgets"
	"github.com/zelenin/go-tdlib/client"
)

// ContactSelectedMsg is emitted when a contact is selected.
type ContactSelectedMsg struct {
	UserID int64
}

// Model is the contact list component.
type Model struct {
	list    widgets.List
	store   *store.Store
	tg      *telegram.Client
	theme   *theme.Theme
	width   int
	height  int
	focused bool
	visible bool
	loaded  bool
}

// New creates a new contacts model.
func New(s *store.Store, tg *telegram.Client, th *theme.Theme) Model {
	l := widgets.NewList()
	l.StyleNormal = th.ChatListItem
	l.StyleActive = th.ChatListItemActive
	l.StyleTitle = th.ChatListTitle
	l.StyleSub = th.ChatListPreview
	l.StyleOnline = th.ChatListOnline

	return Model{
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
	m.list.Width = width
	m.list.Height = height - 2
}

// SetFocused sets focus state.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
	m.list.Focused = focused
}

// SetVisible shows or hides the contacts panel.
func (m *Model) SetVisible(visible bool) {
	m.visible = visible
}

// IsVisible returns whether the contacts panel is visible.
func (m Model) IsVisible() bool {
	return m.visible
}

type contactsLoadedMsg struct {
	users []*client.User
	err   error
}

// LoadContacts fetches the contact list.
func (m *Model) LoadContacts() tea.Cmd {
	return func() tea.Msg {
		contacts, err := m.tg.GetContacts(context.Background())
		if err != nil {
			return contactsLoadedMsg{err: err}
		}

		var users []*client.User
		for _, userID := range contacts.UserIDs {
			user, err := m.tg.GetUser(context.Background(), userID)
			if err == nil {
				users = append(users, user)
				m.store.Users.Set(user)
			}
		}
		return contactsLoadedMsg{users: users}
	}
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case contactsLoadedMsg:
		m.loaded = true
		if msg.err == nil {
			m.refreshList(msg.users)
		}

	case tea.KeyPressMsg:
		if m.focused {
			if selected := m.list.Update(msg); selected {
				item := m.list.SelectedItem()
				if item != nil {
					var userID int64
					fmt.Sscanf(item.ID, "%d", &userID)
					return m, func() tea.Msg {
						return ContactSelectedMsg{UserID: userID}
					}
				}
			}

			if msg.String() == "escape" {
				m.visible = false
			}
		}
	}

	return m, nil
}

func (m *Model) refreshList(users []*client.User) {
	sort.Slice(users, func(i, j int) bool {
		return users[i].FirstName < users[j].FirstName
	})

	items := make([]widgets.ListItem, 0, len(users))
	for _, user := range users {
		name := user.FirstName
		if user.LastName != "" {
			name += " " + user.LastName
		}

		subtitle := ""
		if user.Usernames != nil && len(user.Usernames.ActiveUsernames) > 0 {
			subtitle = "@" + user.Usernames.ActiveUsernames[0]
		}

		_, online := user.Status.(*client.UserStatusOnline)

		items = append(items, widgets.ListItem{
			ID:       fmt.Sprintf("%d", user.ID),
			Title:    name,
			Subtitle: subtitle,
			Online:   online,
		})
	}

	m.list.SetItems(items)
}

// View renders the contacts list.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	title := m.theme.AuthTitle.Width(m.width).Render("Contacts")
	content := m.list.View()

	return m.theme.ChatListPane.
		Width(m.width).
		Height(m.height).
		Render(title + "\n" + content)
}

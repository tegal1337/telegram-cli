package groupinfo

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/tegal1337/telegram-cli/internal/ui/widgets"
	"github.com/zelenin/go-tdlib/client"
)

// Model is the group/channel info panel component.
type Model struct {
	store       *store.Store
	tg          *telegram.Client
	theme       *theme.Theme
	memberList  widgets.List
	width       int
	height      int
	visible     bool
	focused     bool
	chatID      int64
	title       string
	description string
	memberCount int32
	isChannel   bool
}

// New creates a new group info model.
func New(s *store.Store, tg *telegram.Client, th *theme.Theme) Model {
	l := widgets.NewList()
	l.StyleNormal = th.ChatListItem
	l.StyleActive = th.ChatListItemActive
	l.StyleTitle = th.ChatListTitle
	l.StyleSub = th.ChatListPreview

	return Model{
		store:      s,
		tg:         tg,
		theme:      th,
		memberList: l,
	}
}

// SetSize sets the component dimensions.
func (m *Model) SetSize(width, height int) {
	m.width = width
	m.height = height
	m.memberList.Width = width - 2
	m.memberList.Height = height - 8
}

// SetVisible shows or hides the panel.
func (m *Model) SetVisible(visible bool) {
	m.visible = visible
}

// IsVisible returns whether the panel is visible.
func (m Model) IsVisible() bool {
	return m.visible
}

// SetFocused sets focus state.
func (m *Model) SetFocused(focused bool) {
	m.focused = focused
	m.memberList.Focused = focused
}

type groupInfoLoadedMsg struct {
	title       string
	description string
	memberCount int32
	isChannel   bool
	members     []*client.ChatMember
}

// OpenGroupInfo loads group/channel info for a chat.
func (m *Model) OpenGroupInfo(chatID int64) tea.Cmd {
	m.chatID = chatID
	m.visible = true

	return func() tea.Msg {
		chat, err := m.tg.GetChat(context.Background(), chatID)
		if err != nil {
			return nil
		}

		result := groupInfoLoadedMsg{
			title: chat.Title,
		}

		switch t := chat.Type.(type) {
		case *client.ChatTypeSupergroup:
			result.isChannel = t.IsChannel
			info, err := m.tg.GetSupergroupFullInfo(context.Background(), t.SupergroupID)
			if err == nil {
				result.description = info.Description
				result.memberCount = info.MemberCount
			}
			members, err := m.tg.GetSupergroupMembers(context.Background(), t.SupergroupID, 0, 50)
			if err == nil {
				result.members = members.Members
			}

		case *client.ChatTypeBasicGroup:
			info, err := m.tg.GetBasicGroupFullInfo(context.Background(), t.BasicGroupID)
			if err == nil {
				result.description = info.Description
				result.members = info.Members
			}
		}

		return result
	}
}

// Update handles messages.
func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	if !m.visible {
		return m, nil
	}

	switch msg := msg.(type) {
	case groupInfoLoadedMsg:
		m.title = msg.title
		m.description = msg.description
		m.memberCount = msg.memberCount
		m.isChannel = msg.isChannel
		m.refreshMembers(msg.members)

	case tea.KeyPressMsg:
		if m.focused {
			switch msg.String() {
			case "escape":
				m.visible = false
			default:
				m.memberList.Update(msg)
			}
		}
	}

	return m, nil
}

func (m *Model) refreshMembers(members []*client.ChatMember) {
	items := make([]widgets.ListItem, 0, len(members))
	for _, member := range members {
		if sender, ok := member.MemberID.(*client.MessageSenderUser); ok {
			name := m.store.Users.DisplayName(sender.UserID)
			role := memberRole(member.Status)
			online := m.store.Users.IsOnline(sender.UserID)

			items = append(items, widgets.ListItem{
				ID:       fmt.Sprintf("%d", sender.UserID),
				Title:    name,
				Subtitle: role,
				Online:   online,
			})
		}
	}
	m.memberList.SetItems(items)
}

func memberRole(status client.ChatMemberStatus) string {
	switch status.(type) {
	case *client.ChatMemberStatusCreator:
		return "Owner"
	case *client.ChatMemberStatusAdministrator:
		return "Admin"
	case *client.ChatMemberStatusMember:
		return "Member"
	case *client.ChatMemberStatusRestricted:
		return "Restricted"
	case *client.ChatMemberStatusBanned:
		return "Banned"
	case *client.ChatMemberStatusLeft:
		return "Left"
	default:
		return ""
	}
}

// View renders the group info panel.
func (m Model) View() string {
	if !m.visible {
		return ""
	}

	typeLabel := "Group"
	if m.isChannel {
		typeLabel = "Channel"
	}

	header := lipgloss.JoinVertical(lipgloss.Left,
		m.theme.AuthTitle.Render(m.title),
		m.theme.ChatListPreview.Render(typeLabel),
		"",
	)

	if m.description != "" {
		header += m.theme.ChatListPreview.Render(m.description) + "\n\n"
	}

	membersTitle := m.theme.ChatListTitle.Render(
		fmt.Sprintf("Members (%d)", m.memberCount),
	)

	content := lipgloss.JoinVertical(lipgloss.Left,
		header,
		membersTitle,
		m.memberList.View(),
	)

	return m.theme.ChatListPane.
		Width(m.width).
		Height(m.height).
		Render(content)
}

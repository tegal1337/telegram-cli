package app

import (
	"context"
	"fmt"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/tegal1337/telegram-cli/internal/config"
	"github.com/tegal1337/telegram-cli/internal/notification"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/components/auth"
	"github.com/tegal1337/telegram-cli/internal/ui/components/chatlist"
	"github.com/tegal1337/telegram-cli/internal/ui/components/chatview"
	"github.com/tegal1337/telegram-cli/internal/ui/components/composer"
	"github.com/tegal1337/telegram-cli/internal/ui/components/contacts"
	"github.com/tegal1337/telegram-cli/internal/ui/components/dialog"
	"github.com/tegal1337/telegram-cli/internal/ui/components/groupinfo"
	"github.com/tegal1337/telegram-cli/internal/ui/components/search"
	"github.com/tegal1337/telegram-cli/internal/ui/components/statusbar"
	"github.com/tegal1337/telegram-cli/internal/ui/layout"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/zelenin/go-tdlib/client"
)

// Model is the root bubbletea model that composes all sub-components.
type Model struct {
	// Sub-models
	auth      auth.Model
	chatList  chatlist.Model
	chatView  chatview.Model
	composer  composer.Model
	contacts  contacts.Model
	search    search.Model
	groupInfo groupinfo.Model
	statusBar statusbar.Model
	dialog    *dialog.Model

	// State
	screen    ScreenState
	focus     FocusPanel
	layout    layout.Layout

	// Dependencies
	tg        *telegram.Client
	store     *store.Store
	config    *config.Config
	theme     *theme.Theme
	notifier  *notification.Notifier
	sound     *notification.SoundPlayer
	authorizer *telegram.TUIAuthorizer

	// Dimensions
	width  int
	height int

	// User info
	myUserID int64
}

// New creates the root application model.
func New(
	cfg *config.Config,
	tg *telegram.Client,
	s *store.Store,
	authorizer *telegram.TUIAuthorizer,
) Model {
	th := theme.ForName(cfg.UI.Theme)
	notifier := notification.NewNotifier(cfg.Notifications.Enabled, cfg.Notifications.ShowPreview)
	sound := notification.NewSoundPlayer(cfg.Notifications.Sound)

	return Model{
		auth:       auth.New(th, authorizer),
		chatList:   chatlist.New(s, tg, th),
		chatView:   chatview.New(s, tg, th),
		composer:   composer.New(th),
		contacts:   contacts.New(s, tg, th),
		search:     search.New(s, tg, th),
		groupInfo:  groupinfo.New(s, tg, th),
		statusBar:  statusbar.New(s, th),

		screen:     ScreenAuth,
		focus:      PanelChatList,

		tg:         tg,
		store:      s,
		config:     cfg,
		theme:      th,
		notifier:   notifier,
		sound:      sound,
		authorizer: authorizer,
	}
}

// Init returns the initial command.
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handles all incoming messages.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		return m, nil

	case tea.KeyPressMsg:
		// Global keybindings.
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "ctrl+1":
			m.setFocus(PanelChatList)
			return m, nil
		case "ctrl+2":
			m.setFocus(PanelChatView)
			return m, nil
		case "ctrl+3":
			m.setFocus(PanelComposer)
			return m, nil
		case "/":
			if m.screen == ScreenMain && m.focus != PanelComposer {
				m.search.SetVisible(true)
				m.setFocus(PanelSearch)
				return m, nil
			}
		case "ctrl+k":
			if m.screen == ScreenMain {
				m.contacts.SetVisible(!m.contacts.IsVisible())
				if m.contacts.IsVisible() {
					m.setFocus(PanelContacts)
					cmd := m.contacts.LoadContacts()
					cmds = append(cmds, cmd)
				} else {
					m.setFocus(PanelChatList)
				}
				return m, tea.Batch(cmds...)
			}
		}

	// Auth state from TDLib.
	case telegram.AuthStateMsg:
		return m.handleAuthState(msg)

	// Authentication complete.
	case AuthenticatedMsg:
		m.screen = ScreenMain
		m.myUserID = msg.UserID
		m.chatView.SetMyUserID(msg.UserID)
		m.statusBar.SetUserName(fmt.Sprintf("%s %s", msg.FirstName, msg.LastName))
		m.setFocus(PanelChatList)
		m.updateLayout()
		cmd := m.chatList.Init()
		return m, cmd

	// New message notification.
	case telegram.NewMessageMsg:
		if msg.Message.ChatID != m.chatList.ActiveChatID() {
			// Notify for messages in non-active chats.
			entry, ok := m.store.Chats.Get(msg.Message.ChatID)
			title := "New Message"
			if ok && entry.Chat != nil {
				title = entry.Chat.Title
			}
			body := "New message received"
			if text, ok := msg.Message.Content.(*client.MessageText); ok {
				body = text.Text.Text
			}
			m.notifier.Notify(title, body)
			m.sound.Play()
		}

	// Chat selected from chat list.
	case chatlist.ChatSelectedMsg:
		entry, ok := m.store.Chats.Get(msg.ChatID)
		title := ""
		if ok && entry.Chat != nil {
			title = entry.Chat.Title
		}
		cmd := m.chatView.OpenChat(msg.ChatID, title)
		m.composer.SetChatID(msg.ChatID)
		m.statusBar.SetActiveChatID(msg.ChatID)
		m.setFocus(PanelChatView)
		cmds = append(cmds, cmd)

	// Contact selected.
	case contacts.ContactSelectedMsg:
		m.contacts.SetVisible(false)
		cmd := m.openPrivateChat(msg.UserID)
		cmds = append(cmds, cmd)

	// Search result selected.
	case search.SearchResultMsg:
		entry, ok := m.store.Chats.Get(msg.ChatID)
		title := ""
		if ok && entry.Chat != nil {
			title = entry.Chat.Title
		}
		cmd := m.chatView.OpenChat(msg.ChatID, title)
		m.composer.SetChatID(msg.ChatID)
		m.setFocus(PanelChatView)
		cmds = append(cmds, cmd)

	// Message submitted from composer.
	case composer.MessageSubmittedMsg:
		cmd := m.handleMessageSubmit(msg)
		cmds = append(cmds, cmd)

	// Message actions (reply, edit, delete, forward).
	case chatview.MessageActionMsg:
		return m.handleMessageAction(msg)

	// Dialog results.
	case dialog.DialogResultMsg:
		m.dialog = nil
		if msg.ID == "delete" && msg.Confirmed {
			// Delete confirmed — handled by the stored command.
		}
	}

	// Dispatch to all relevant sub-models.
	if m.screen == ScreenAuth {
		var cmd tea.Cmd
		m.auth, cmd = m.auth.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		// Dispatch to sub-models.
		var cmd tea.Cmd

		m.chatList, cmd = m.chatList.Update(msg)
		cmds = append(cmds, cmd)

		m.chatView, cmd = m.chatView.Update(msg)
		cmds = append(cmds, cmd)

		m.composer, cmd = m.composer.Update(msg)
		cmds = append(cmds, cmd)

		m.contacts, cmd = m.contacts.Update(msg)
		cmds = append(cmds, cmd)

		m.search, cmd = m.search.Update(msg)
		cmds = append(cmds, cmd)

		m.groupInfo, cmd = m.groupInfo.Update(msg)
		cmds = append(cmds, cmd)

		m.statusBar, cmd = m.statusBar.Update(msg)
		cmds = append(cmds, cmd)

		if m.dialog != nil {
			var d dialog.Model
			d, cmd = m.dialog.Update(msg)
			m.dialog = &d
			cmds = append(cmds, cmd)
		}
	}

	return m, tea.Batch(cmds...)
}

func (m Model) handleAuthState(msg telegram.AuthStateMsg) (tea.Model, tea.Cmd) {
	switch msg.State.(type) {
	case *client.AuthorizationStateWaitPhoneNumber:
		m.auth.SetStep(auth.StepPhone)
	case *client.AuthorizationStateWaitCode:
		m.auth.SetStep(auth.StepCode)
	case *client.AuthorizationStateWaitPassword:
		m.auth.SetStep(auth.StepPassword)
	case *client.AuthorizationStateReady:
		m.auth.SetStep(auth.StepDone)
		return m, m.fetchMeCmd()
	}
	return m, nil
}

func (m *Model) fetchMeCmd() tea.Cmd {
	return func() tea.Msg {
		me, err := m.tg.GetMe(context.Background())
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return AuthenticatedMsg{
			UserID:    me.ID,
			FirstName: me.FirstName,
			LastName:  me.LastName,
		}
	}
}

func (m *Model) openPrivateChat(userID int64) tea.Cmd {
	return func() tea.Msg {
		chat, err := m.tg.CreatePrivateChat(context.Background(), userID)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return chatlist.ChatSelectedMsg{ChatID: chat.ID}
	}
}

func (m Model) handleMessageSubmit(msg composer.MessageSubmittedMsg) tea.Cmd {
	if msg.EditMessageID != 0 {
		return func() tea.Msg {
			_, err := m.tg.EditTextMessage(context.Background(), msg.ChatID, msg.EditMessageID, msg.Text)
			if err != nil {
				return ErrorMsg{Err: err}
			}
			return nil
		}
	}

	return func() tea.Msg {
		_, err := m.tg.SendTextMessage(context.Background(), msg.ChatID, msg.Text, msg.ReplyToID)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return nil
	}
}

func (m Model) handleMessageAction(msg chatview.MessageActionMsg) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case "reply":
		// Get message preview for reply bar.
		msgs := m.store.Messages.Get(msg.ChatID)
		preview := ""
		for _, message := range msgs {
			if message.ID == msg.MessageID {
				if text, ok := message.Content.(*client.MessageText); ok {
					preview = text.Text.Text
				} else {
					preview = "[Media]"
				}
				break
			}
		}
		m.composer.EnterReplyMode(msg.MessageID, preview)
		m.setFocus(PanelComposer)

	case "edit":
		msgs := m.store.Messages.Get(msg.ChatID)
		for _, message := range msgs {
			if message.ID == msg.MessageID {
				if text, ok := message.Content.(*client.MessageText); ok {
					m.composer.EnterEditMode(msg.MessageID, text.Text.Text)
					m.setFocus(PanelComposer)
				}
				break
			}
		}

	case "delete":
		d := dialog.NewConfirm(m.theme, "delete", "Delete Message", "Are you sure you want to delete this message?")
		m.dialog = &d

	case "forward":
		// TODO: show chat picker for forwarding.
	}

	return m, nil
}

func (m *Model) setFocus(panel FocusPanel) {
	m.focus = panel
	m.chatList.SetFocused(panel == PanelChatList)
	m.chatView.SetFocused(panel == PanelChatView)
	m.composer.SetFocused(panel == PanelComposer)
	m.contacts.SetFocused(panel == PanelContacts)
	m.groupInfo.SetFocused(panel == PanelGroupInfo)
}

func (m *Model) updateLayout() {
	l := layout.Compute(m.width, m.height, m.config.UI.ChatListWidth)
	m.layout = l

	m.auth.SetSize(m.width, m.height)
	m.chatList.SetSize(l.ChatListWidth, l.ChatListHeight)
	m.chatView.SetSize(l.ChatViewWidth, l.ChatViewHeight)
	m.composer.SetSize(l.ComposerWidth, l.ComposerHeight)
	m.contacts.SetSize(l.ChatListWidth, l.ChatListHeight)
	m.search.SetSize(m.width/2, m.height/2)
	m.groupInfo.SetSize(l.ChatListWidth, l.ChatListHeight)
	m.statusBar.SetSize(l.StatusBarWidth)
}

// View renders the entire UI.
func (m Model) View() tea.View {
	var content string

	switch m.screen {
	case ScreenAuth:
		content = m.auth.View()

	case ScreenLoading:
		content = lipgloss.NewStyle().
			Width(m.width).
			Height(m.height).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Loading...")

	case ScreenMain:
		content = m.renderMainScreen()
	}

	// Overlay dialog if present.
	if m.dialog != nil && m.dialog.IsVisible() {
		dialogView := m.dialog.View()
		content = overlayCenter(content, dialogView, m.width, m.height)
	}

	// Overlay search if present.
	if m.search.IsVisible() {
		searchView := m.search.View()
		content = overlayCenter(content, searchView, m.width, m.height)
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) renderMainScreen() string {
	// Left panel: chat list or contacts.
	var leftPanel string
	if m.contacts.IsVisible() {
		leftPanel = m.contacts.View()
	} else {
		leftPanel = m.chatList.View()
	}

	// Right panel: chat view + composer.
	chatHeader := m.chatView.View()
	composerView := m.composer.View()

	rightPanel := lipgloss.JoinVertical(lipgloss.Left,
		chatHeader,
		composerView,
	)

	// Group info panel (if visible).
	if m.groupInfo.IsVisible() {
		rightPanel = lipgloss.JoinHorizontal(lipgloss.Top,
			rightPanel,
			m.groupInfo.View(),
		)
	}

	// Main area: left + right.
	var mainArea string
	if m.layout.SinglePanel {
		switch m.focus {
		case PanelChatList, PanelContacts:
			mainArea = leftPanel
		default:
			mainArea = rightPanel
		}
	} else {
		mainArea = lipgloss.JoinHorizontal(lipgloss.Top,
			leftPanel,
			rightPanel,
		)
	}

	// Status bar at the bottom.
	statusBar := m.statusBar.View()

	return lipgloss.JoinVertical(lipgloss.Left,
		mainArea,
		statusBar,
	)
}

// overlayCenter places an overlay in the center of the base content.
func overlayCenter(base, overlay string, width, height int) string {
	// Simple overlay: just replace center lines.
	overlayLines := lipgloss.Height(overlay)
	overlayWidth := lipgloss.Width(overlay)

	baseLines := make([]string, height)
	for i := range baseLines {
		baseLines[i] = ""
	}

	startY := (height - overlayLines) / 2
	startX := (width - overlayWidth) / 2
	if startX < 0 {
		startX = 0
	}

	_ = startY
	_ = startX

	// For simplicity, just render overlay on top.
	return lipgloss.Place(width, height,
		lipgloss.Center, lipgloss.Center,
		overlay,
		lipgloss.WithWhitespaceBackground(lipgloss.Color("#00000088")),
	)
}

package app

import (
	"fmt"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
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

type Model struct {
	auth      auth.Model
	chatList  chatlist.Model
	chatView  chatview.Model
	composer  composer.Model
	contacts  contacts.Model
	search    search.Model
	groupInfo groupinfo.Model
	statusBar statusbar.Model
	dialog    *dialog.Model

	screen     ScreenState
	focus      FocusPanel
	layout     layout.Layout
	tg         *telegram.Client
	store      *store.Store
	config     *config.Config
	theme      *theme.Theme
	notifier   *notification.Notifier
	sound      *notification.SoundPlayer
	authorizer *telegram.TUIAuthorizer
	width      int
	height     int
	myUserId   int64
}

func New(cfg *config.Config, tg *telegram.Client, s *store.Store, authorizer *telegram.TUIAuthorizer) Model {
	th := theme.ForName(cfg.UI.Theme)
	return Model{
		auth:       auth.New(th, authorizer),
		chatList:   chatlist.New(s, tg, th),
		chatView:   chatview.New(s, tg, th),
		composer:   composer.New(th),
		contacts:   contacts.New(s, tg, th),
		search:     search.New(s, tg, th),
		groupInfo:  groupinfo.New(s, tg, th),
		statusBar:  statusbar.New(s, th),
		screen:     ScreenLoading,
		focus:      PanelChatList,
		tg:         tg,
		store:      s,
		config:     cfg,
		theme:      th,
		notifier:   notification.NewNotifier(cfg.Notifications.Enabled, cfg.Notifications.ShowPreview),
		sound:      notification.NewSoundPlayer(cfg.Notifications.Sound),
		authorizer: authorizer,
	}
}

func (m Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.updateLayout()
		return m, nil

	case tea.KeyPressMsg:
		key := msg.String()

		// Quit
		if key == "ctrl+c" || key == "ctrl+q" {
			return m, tea.Quit
		}

		if m.screen == ScreenMain {
			// Tab / Shift+Tab cycle panels
			if key == "tab" && m.focus != PanelSearch && m.focus != PanelComposer {
				switch m.focus {
				case PanelChatList:
					m.setFocus(PanelChatView)
				case PanelChatView:
					m.setFocus(PanelComposer)
				default:
					m.setFocus(PanelChatList)
				}
				return m, nil
			}
			if key == "shift+tab" && m.focus != PanelSearch {
				switch m.focus {
				case PanelComposer:
					m.setFocus(PanelChatView)
				case PanelChatView:
					m.setFocus(PanelChatList)
				default:
					m.setFocus(PanelComposer)
				}
				return m, nil
			}

			// Escape: close overlay or go back
			if key == "escape" {
				if m.search.IsVisible() {
					m.search.SetVisible(false)
					m.setFocus(PanelChatList)
					return m, nil
				}
				if m.contacts.IsVisible() {
					m.contacts.SetVisible(false)
					m.setFocus(PanelChatList)
					return m, nil
				}
				if m.focus == PanelComposer {
					m.setFocus(PanelChatView)
					return m, nil
				}
				if m.focus != PanelChatList {
					m.setFocus(PanelChatList)
					return m, nil
				}
			}

			// Alt+1/2/3 for panel focus (works in all terminals)
			if key == "alt+1" || key == "F1" {
				m.setFocus(PanelChatList)
				return m, nil
			}
			if key == "alt+2" || key == "F2" {
				m.setFocus(PanelChatView)
				return m, nil
			}
			if key == "alt+3" || key == "F3" {
				m.setFocus(PanelComposer)
				return m, nil
			}

			// Alt+j / Alt+k: next/prev chat (works when not typing)
			if key == "alt+j" && m.focus != PanelComposer {
				// next chat handled by chatlist
			}
			if key == "alt+k" && m.focus != PanelComposer {
				// prev chat handled by chatlist
			}

			// Search (not when typing in composer)
			if key == "/" && m.focus != PanelComposer {
				m.search.SetVisible(true)
				m.setFocus(PanelSearch)
				return m, nil
			}

			// Contacts toggle
			if key == "alt+c" {
				m.contacts.SetVisible(!m.contacts.IsVisible())
				if m.contacts.IsVisible() {
					m.setFocus(PanelContacts)
					return m, m.contacts.LoadContacts()
				}
				m.setFocus(PanelChatList)
				return m, nil
			}

			// Quick compose: just start typing from chatview
			if key == "i" && m.focus == PanelChatView {
				m.setFocus(PanelComposer)
				return m, nil
			}
		}

	case AuthStateChangedMsg:
		return m.handleAuthStateChanged(msg)

	case AuthenticatedMsg:
		m.screen = ScreenMain
		m.myUserId = msg.UserId
		m.chatView.SetMyUserId(msg.UserId)
		m.statusBar.SetUserName(fmt.Sprintf("%s %s", msg.FirstName, msg.LastName))
		m.setFocus(PanelChatList)
		m.updateLayout()
		return m, m.chatList.Init()

	case telegram.NewMessageMsg:
		if msg.Message.ChatId != m.chatList.ActiveChatId() {
			entry, ok := m.store.Chats.Get(msg.Message.ChatId)
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

	case chatlist.ChatSelectedMsg:
		entry, ok := m.store.Chats.Get(msg.ChatId)
		title := ""
		if ok && entry.Chat != nil {
			title = entry.Chat.Title
		}
		cmd := m.chatView.OpenChat(msg.ChatId, title)
		m.composer.SetChatId(msg.ChatId)
		m.statusBar.SetActiveChatId(msg.ChatId)
		m.setFocus(PanelChatView)
		cmds = append(cmds, cmd)

	case contacts.ContactSelectedMsg:
		m.contacts.SetVisible(false)
		cmds = append(cmds, m.openPrivateChat(msg.UserId))

	case search.SearchResultMsg:
		entry, ok := m.store.Chats.Get(msg.ChatId)
		title := ""
		if ok && entry.Chat != nil {
			title = entry.Chat.Title
		}
		cmd := m.chatView.OpenChat(msg.ChatId, title)
		m.composer.SetChatId(msg.ChatId)
		m.setFocus(PanelChatView)
		cmds = append(cmds, cmd)

	case composer.MessageSubmittedMsg:
		cmds = append(cmds, m.handleMessageSubmit(msg))

	case chatview.MessageActionMsg:
		return m.handleMessageAction(msg)

	case dialog.DialogResultMsg:
		m.dialog = nil
	}

	// Dispatch to sub-models
	if m.screen == ScreenAuth {
		var cmd tea.Cmd
		m.auth, cmd = m.auth.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		var cmd tea.Cmd

		// Key events only go to the focused panel.
		// Non-key events (telegram updates, spinner ticks, etc.) go to all.
		_, isKey := msg.(tea.KeyPressMsg)
		_, isPaste := msg.(tea.PasteMsg)
		isInputEvent := isKey || isPaste

		if !isInputEvent || m.focus == PanelChatList {
			m.chatList, cmd = m.chatList.Update(msg)
			cmds = append(cmds, cmd)
		}
		if !isInputEvent || m.focus == PanelChatView {
			m.chatView, cmd = m.chatView.Update(msg)
			cmds = append(cmds, cmd)
		}
		if !isInputEvent || m.focus == PanelComposer {
			m.composer, cmd = m.composer.Update(msg)
			cmds = append(cmds, cmd)
		}
		if !isInputEvent || m.focus == PanelContacts {
			m.contacts, cmd = m.contacts.Update(msg)
			cmds = append(cmds, cmd)
		}
		if !isInputEvent || m.focus == PanelSearch {
			m.search, cmd = m.search.Update(msg)
			cmds = append(cmds, cmd)
		}
		if !isInputEvent || m.focus == PanelGroupInfo {
			m.groupInfo, cmd = m.groupInfo.Update(msg)
			cmds = append(cmds, cmd)
		}

		// Status bar always gets all events (non-interactive).
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

func (m Model) handleAuthStateChanged(msg AuthStateChangedMsg) (tea.Model, tea.Cmd) {
	switch telegram.AuthState(msg.State) {
	case telegram.AuthStateWaitPhone:
		m.screen = ScreenAuth
		m.auth.SetStep(auth.StepPhone)
	case telegram.AuthStateWaitCode:
		m.screen = ScreenAuth
		m.auth.SetStep(auth.StepCode)
	case telegram.AuthStateWaitPassword:
		m.screen = ScreenAuth
		m.auth.SetStep(auth.StepPassword)
	case telegram.AuthStateReady:
		m.screen = ScreenAuth
		m.auth.SetStep(auth.StepDone)
	}
	return m, nil
}

func (m *Model) openPrivateChat(userID int64) tea.Cmd {
	return func() tea.Msg {
		chat, err := m.tg.CreatePrivateChat(userID)
		if err != nil {
			return ErrorMsg{Err: err}
		}
		return chatlist.ChatSelectedMsg{ChatId: chat.Id}
	}
}

func (m Model) handleMessageSubmit(msg composer.MessageSubmittedMsg) tea.Cmd {
	if msg.EditMessageId != 0 {
		return func() tea.Msg {
			m.tg.EditTextMessage(msg.ChatId, msg.EditMessageId, msg.Text)
			return nil
		}
	}
	return func() tea.Msg {
		m.tg.SendTextMessage(msg.ChatId, msg.Text, msg.ReplyToId)
		return nil
	}
}

func (m Model) handleMessageAction(msg chatview.MessageActionMsg) (tea.Model, tea.Cmd) {
	switch msg.Action {
	case "reply":
		msgs := m.store.Messages.Get(msg.ChatId)
		preview := ""
		for _, message := range msgs {
			if message.Id == msg.MessageId {
				if text, ok := message.Content.(*client.MessageText); ok {
					preview = text.Text.Text
				} else {
					preview = "[Media]"
				}
				break
			}
		}
		m.composer.EnterReplyMode(msg.MessageId, preview)
		m.setFocus(PanelComposer)
	case "edit":
		msgs := m.store.Messages.Get(msg.ChatId)
		for _, message := range msgs {
			if message.Id == msg.MessageId {
				if text, ok := message.Content.(*client.MessageText); ok {
					m.composer.EnterEditMode(msg.MessageId, text.Text.Text)
					m.setFocus(PanelComposer)
				}
				break
			}
		}
	case "delete":
		d := dialog.NewConfirm(m.theme, "delete", "Delete Message", "Are you sure?")
		m.dialog = &d
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
	// Inner dimensions (subtract 2 for border)
	m.chatList.SetSize(l.ChatListWidth-2, l.ChatListHeight-2)
	m.chatView.SetSize(l.ChatViewWidth-2, l.ChatViewHeight-2)
	m.composer.SetSize(l.ComposerWidth-2, l.ComposerHeight-2)
	m.contacts.SetSize(l.ChatListWidth-2, l.ChatListHeight-2)
	m.search.SetSize(m.width/2, m.height/2)
	m.groupInfo.SetSize(l.ChatListWidth-2, l.ChatListHeight-2)
	m.statusBar.SetSize(l.StatusBarWidth)
}

func (m Model) View() tea.View {
	var content string

	switch m.screen {
	case ScreenAuth:
		content = m.auth.View()
	case ScreenLoading:
		blue := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Bold(true)
		cyan := lipgloss.NewStyle().Foreground(lipgloss.Color("51"))
		dim := lipgloss.NewStyle().Foreground(lipgloss.Color("244"))

		tgLogo := cyan.Render(
			"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣀⣤⣴⣾⣿⣿⣿⡄\n" +
				"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣠⣴⣶⣿⣿⡿⠿⠛⢙⣿⣿⠃\n" +
				"⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣤⣶⣾⣿⣿⠿⠛⠋⠁⠀⠀⠀⣸⣿⣿⠀\n" +
				"⠀⠀⠀⠀⣀⣤⣴⣾⣿⣿⡿⠟⠛⠉⠀⠀⣠⣤⠞⠁⠀⠀⣿⣿⡇⠀\n" +
				"⠀⣴⣾⣿⣿⡿⠿⠛⠉⠀⠀⠀⢀⣠⣶⣿⠟⠁⠀⠀⠀⢸⣿⣿⠀⠀\n" +
				"⠸⣿⣿⣿⣧⣄⣀⠀⠀⣀⣴⣾⣿⣿⠟⠁⠀⠀⠀⠀⠀⣼⣿⡿⠀⠀\n" +
				"⠀⠈⠙⠻⠿⣿⣿⣿⣿⣿⣿⣿⠟⠁⠀⠀⠀⠀⠀⠀⢠⣿⣿⠇⠀⠀\n" +
				"⠀⠀⠀⠀⠀⠀⠘⣿⣿⣿⣿⡇⠀⣀⣄⡀⠀⠀⠀⠀⢸⣿⣿⠀⠀⠀\n" +
				"⠀⠀⠀⠀⠀⠀⠀⠸⣿⣿⣿⣠⣾⣿⣿⣿⣦⡀⠀⠀⣿⣿⡏⠀⠀⠀\n" +
				"⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⡿⠋⠈⠻⣿⣿⣦⣸⣿⣿⠁⠀⠀⠀\n" +
				"⠀⠀⠀⠀⠀⠀⠀⠀⠀⠙⠛⠁⠀⠀⠀⠀⠈⠻⣿⣿⣿⠏⠀⠀⠀⠀")

		title := blue.Render("  T E L E G R A M   C L I")
		sub := dim.Render("  Terminal Client for Telegram")
		author := dim.Render("  github.com/tegal1337")

		box := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(1, 4).
			Render(tgLogo + "\n\n" + title + "\n\n" + sub + "\n" + author)

		spinner := blue.Render("  ⣾ Connecting to Telegram...")
		content = lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, box+"\n\n"+spinner)
	case ScreenMain:
		content = m.renderMainScreen()
	}

	if m.dialog != nil && m.dialog.IsVisible() {
		content = lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			m.dialog.View())
	}

	if m.search.IsVisible() {
		content = lipgloss.Place(m.width, m.height,
			lipgloss.Center, lipgloss.Center,
			m.search.View())
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.MouseMode = tea.MouseModeCellMotion
	return v
}

func (m Model) renderMainScreen() string {
	// Build left panel with rounded border
	var leftContent string
	if m.contacts.IsVisible() {
		leftContent = m.contacts.View()
	} else {
		leftContent = m.chatList.View()
	}

	leftStyle := m.theme.PanelNormal
	if m.focus == PanelChatList || m.focus == PanelContacts {
		leftStyle = m.theme.PanelFocused
	}
	leftPanel := leftStyle.
		Width(m.layout.ChatListWidth - 2).
		Height(m.layout.ChatListHeight - 2).
		Render(leftContent)

	// Build chat view with rounded border
	chatViewStyle := m.theme.PanelNormal
	if m.focus == PanelChatView {
		chatViewStyle = m.theme.PanelFocused
	}
	chatPanel := chatViewStyle.
		Width(m.layout.ChatViewWidth - 2).
		Height(m.layout.ChatViewHeight - 2).
		Render(m.chatView.View())

	// Build composer with rounded border
	composerStyle := m.theme.PanelNormal
	if m.focus == PanelComposer {
		composerStyle = m.theme.PanelFocused
	}
	composerPanel := composerStyle.
		Width(m.layout.ComposerWidth - 2).
		Height(m.layout.ComposerHeight - 2).
		Render(m.composer.View())

	// Right side = chat + composer stacked
	rightPanel := lipgloss.JoinVertical(lipgloss.Left, chatPanel, composerPanel)

	// Main area = left + right
	var mainArea string
	if m.layout.SinglePanel {
		switch m.focus {
		case PanelChatList, PanelContacts:
			mainArea = leftPanel
		default:
			mainArea = lipgloss.JoinVertical(lipgloss.Left, chatPanel, composerPanel)
		}
	} else {
		mainArea = lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
	}

	// Status bar
	statusBar := m.statusBar.View()

	// Keybind help line
	helpStyle := lipgloss.NewStyle().Foreground(m.theme.TextMuted)
	focusName := [...]string{"CHATS", "MESSAGES", "COMPOSE", "SEARCH", "CONTACTS", "INFO"}
	fi := int(m.focus)
	if fi >= len(focusName) {
		fi = 0
	}
	help := helpStyle.Render(fmt.Sprintf(
		" Tab:switch │ Esc:back │ /:search │ Alt+C:contacts │ F1/F2/F3:panels │ i:compose │ %s",
		focusName[fi],
	))

	// Pad help to full width
	helpW := lipgloss.Width(help)
	if helpW < m.width {
		help += strings.Repeat(" ", m.width-helpW)
	}

	return lipgloss.JoinVertical(lipgloss.Left, mainArea, statusBar, help)
}

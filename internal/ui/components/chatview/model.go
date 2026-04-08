package chatview

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/charmbracelet/lipgloss"
	"github.com/tegal1337/telegram-cli/internal/render"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/telegram"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/zelenin/go-tdlib/client"
)

type Model struct {
	store        *store.Store
	tg           *telegram.Client
	theme        *theme.Theme
	renderer     *render.MessageRenderer
	width        int
	height       int
	focused      bool
	chatID       int64
	chatTitle    string
	scrollOffset int
	loading      bool
	loadStatus   string // "Loading messages...", "Fetching users...", etc.
	loadProgress int    // 0-100
	myUserId     int64
	mediaStatus  string
}

func New(s *store.Store, tg *telegram.Client, th *theme.Theme) Model {
	return Model{
		store:    s,
		tg:       tg,
		theme:    th,
		renderer: render.NewMessageRenderer(th),
	}
}

func (m *Model) SetSize(w, h int)       { m.width = w; m.height = h }
func (m *Model) SetFocused(focused bool) { m.focused = focused }
func (m *Model) SetMyUserId(id int64)    { m.myUserId = id }

func (m *Model) OpenChat(chatID int64, title string) tea.Cmd {
	m.chatID = chatID
	m.chatTitle = title
	m.scrollOffset = 0
	m.loading = true
	m.loadStatus = "⟳ Loading messages..."
	m.loadProgress = 10
	m.mediaStatus = ""
	return tea.Batch(
		m.loadHistoryCmd(chatID, 0),
		func() tea.Msg { m.tg.OpenChat(chatID); return nil },
	)
}

type historyLoadedMsg struct {
	chatID   int64
	messages []*client.Message
	err      error
}

func (m *Model) loadHistoryCmd(chatID int64, fromMsgId int64) tea.Cmd {
	return func() tea.Msg {
		msgs, err := m.tg.GetChatHistory(chatID, fromMsgId, 0, 50)
		if err != nil {
			return historyLoadedMsg{chatID: chatID, err: err}
		}
		return historyLoadedMsg{chatID: chatID, messages: msgs.Messages}
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case historyLoadedMsg:
		if msg.chatID != m.chatID {
			return m, nil
		}
		if msg.err == nil && len(msg.messages) > 0 {
			reversed := make([]*client.Message, len(msg.messages))
			for i, v := range msg.messages {
				reversed[len(msg.messages)-1-i] = v
			}
			m.store.Messages.Prepend(m.chatID, reversed)
			m.loadStatus = "⟳ Fetching user info..."
			m.loadProgress = 40
			return m, m.fetchMessageMeta(reversed)
		}
		m.loading = false
		m.loadProgress = 100

	case telegram.NewMessageMsg:
		if msg.Message.ChatId == m.chatID {
			m.store.Messages.Append(m.chatID, msg.Message)
			return m, func() tea.Msg {
				m.tg.ViewMessages(m.chatID, []int64{msg.Message.Id})
				return nil
			}
		}

	case telegram.MessageEditedMsg:
		if msg.ChatId == m.chatID {
			return m, func() tea.Msg {
				fetched, _ := m.tg.GetMessage(msg.ChatId, msg.MessageId)
				if fetched != nil {
					return messageFetchedMsg{chatID: msg.ChatId, message: fetched}
				}
				return nil
			}
		}

	case telegram.MessageDeletedMsg:
		if msg.ChatId == m.chatID {
			m.store.Messages.Delete(m.chatID, msg.MessageIds)
		}

	case telegram.MessageSendSucceededMsg:
		if msg.Message.ChatId == m.chatID {
			m.store.Messages.ReplaceMessageId(m.chatID, msg.OldMessageId, msg.Message)
		}

	case messageFetchedMsg:
		if msg.chatID == m.chatID && msg.message != nil {
			m.store.Messages.UpdateMessage(m.chatID, msg.message.Id, msg.message)
		}

	case telegram.FileUpdateMsg:
		if msg.File != nil {
			m.store.Files.Update(msg.File)
		}

	case metaFetchedMsg:
		m.loading = false
		m.loadProgress = 100
		m.loadStatus = ""

	case MediaPlayMsg:
		m.mediaStatus = msg.Info

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

type metaFetchedMsg struct{}

type loadProgressMsg struct {
	status   string
	progress int
}

func (m Model) fetchMessageMeta(msgs []*client.Message) tea.Cmd {
	return func() tea.Msg {
		total := len(msgs)
		seen := make(map[int64]bool)

		// Phase 1: fetch user info (40% → 70%)
		for i, msg := range msgs {
			if sender, ok := msg.SenderId.(*client.MessageSenderUser); ok {
				if !seen[sender.UserId] {
					seen[sender.UserId] = true
					if _, exists := m.store.Users.Get(sender.UserId); !exists {
						user, err := m.tg.GetUser(sender.UserId)
						if err == nil {
							m.store.Users.Set(user)
						}
					}
				}
			}
			_ = i
		}

		// Phase 2: download photos (70% → 100%)
		photoCount := 0
		for _, msg := range msgs {
			if _, ok := msg.Content.(*client.MessagePhoto); ok {
				photoCount++
			}
		}

		downloaded := 0
		for _, msg := range msgs {
			if photo, ok := msg.Content.(*client.MessagePhoto); ok {
				if photo.Photo != nil && len(photo.Photo.Sizes) > 0 {
					target := photo.Photo.Sizes[0]
					for _, sz := range photo.Photo.Sizes {
						if sz.Width <= 320 && sz.Width > target.Width {
							target = sz
						}
					}
					if target.Photo != nil {
						if target.Photo.Local == nil || !target.Photo.Local.IsDownloadingCompleted {
							file, err := m.tg.DownloadFileSync(target.Photo.Id)
							if err == nil && file != nil {
								m.store.Files.Update(file)
							}
						}
					}
				}
				downloaded++
				_ = total
			}
		}

		return metaFetchedMsg{}
	}
}

func (m Model) handleKey(msg tea.KeyPressMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "up", "k":
		m.scrollOffset += 3
		msgs := m.store.Messages.Get(m.chatID)
		if m.scrollOffset > len(msgs)*4 {
			m.scrollOffset = len(msgs) * 4
			oldest := m.store.Messages.OldestMessageId(m.chatID)
			if oldest != 0 {
				return m, m.loadHistoryCmd(m.chatID, oldest)
			}
		}
	case "down", "j":
		m.scrollOffset -= 3
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	case "G", "end":
		m.scrollOffset = 0
	case "g", "home":
		msgs := m.store.Messages.Get(m.chatID)
		m.scrollOffset = len(msgs) * 4
	case "ctrl+u":
		m.scrollOffset += m.height
	case "ctrl+d":
		m.scrollOffset -= m.height
		if m.scrollOffset < 0 {
			m.scrollOffset = 0
		}
	case "r":
		return m, m.messageAction("reply")
	case "e":
		return m, m.messageAction("edit")
	case "d":
		return m, m.messageAction("delete")

	// Enter: play/open media of the bottom-visible message
	case "enter":
		return m, m.playMedia()

	// 'o' also opens media
	case "o":
		return m, m.playMedia()

	// 's' saves/downloads file
	case "s":
		return m, m.downloadFile()
	}
	return m, nil
}

func (m Model) getTargetMessage() *client.Message {
	msgs := m.store.Messages.Get(m.chatID)
	if len(msgs) == 0 {
		return nil
	}
	idx := len(msgs) - 1
	if m.scrollOffset > 0 {
		offset := m.scrollOffset / 4
		idx = len(msgs) - 1 - offset
		if idx < 0 {
			idx = 0
		}
	}
	return msgs[idx]
}

func (m Model) messageAction(action string) tea.Cmd {
	msg := m.getTargetMessage()
	if msg == nil {
		return nil
	}
	if action == "edit" && !isOwnMessage(msg, m.myUserId) {
		return nil
	}
	return func() tea.Msg {
		return MessageActionMsg{Action: action, ChatId: m.chatID, MessageId: msg.Id}
	}
}

// playMedia downloads and plays the media in the target message.
func (m Model) playMedia() tea.Cmd {
	msg := m.getTargetMessage()
	if msg == nil || msg.Content == nil {
		return nil
	}

	switch c := msg.Content.(type) {
	case *client.MessageVoiceNote:
		return m.downloadAndPlay(c.VoiceNote.Voice.Id, "voice", "🎤 Playing voice...")

	case *client.MessageAudio:
		return m.downloadAndPlay(c.Audio.Audio.Id, "audio", fmt.Sprintf("🎵 Playing %s...", c.Audio.Title))

	case *client.MessageVideoNote:
		return m.downloadAndPlay(c.VideoNote.Video.Id, "video", "📹 Playing video note...")

	case *client.MessageVideo:
		return m.downloadAndPlay(c.Video.Video.Id, "video", "🎥 Opening video...")

	case *client.MessageAnimation:
		return m.downloadAndPlay(c.Animation.Animation.Id, "video", "🎬 Opening GIF...")

	case *client.MessageDocument:
		return m.downloadAndOpen(c.Document.Document.Id, fmt.Sprintf("📎 Opening %s...", c.Document.FileName))

	case *client.MessagePhoto:
		if c.Photo != nil && len(c.Photo.Sizes) > 0 {
			best := c.Photo.Sizes[len(c.Photo.Sizes)-1]
			return m.downloadAndOpen(best.Photo.Id, "🖼 Opening photo...")
		}

	case *client.MessageSticker:
		return m.downloadAndOpen(c.Sticker.Sticker.Id, "Opening sticker...")
	}

	return nil
}

// downloadFile saves the file from the target message.
func (m Model) downloadFile() tea.Cmd {
	msg := m.getTargetMessage()
	if msg == nil || msg.Content == nil {
		return nil
	}

	var fileId int32
	var name string

	switch c := msg.Content.(type) {
	case *client.MessageDocument:
		fileId = c.Document.Document.Id
		name = c.Document.FileName
	case *client.MessagePhoto:
		if c.Photo != nil && len(c.Photo.Sizes) > 0 {
			best := c.Photo.Sizes[len(c.Photo.Sizes)-1]
			fileId = best.Photo.Id
			name = "photo"
		}
	case *client.MessageVideo:
		fileId = c.Video.Video.Id
		name = c.Video.FileName
	case *client.MessageAudio:
		fileId = c.Audio.Audio.Id
		name = c.Audio.FileName
	case *client.MessageVoiceNote:
		fileId = c.VoiceNote.Voice.Id
		name = "voice"
	default:
		return nil
	}

	return func() tea.Msg {
		file, err := m.tg.DownloadFileSync(fileId)
		if err != nil {
			return MediaPlayMsg{Status: "error", Info: fmt.Sprintf("Download failed: %v", err)}
		}
		return MediaPlayMsg{Status: "downloaded", Info: fmt.Sprintf("💾 Saved %s → %s", name, file.Local.Path)}
	}
}

func (m Model) downloadAndPlay(fileId int32, mediaType string, statusMsg string) tea.Cmd {
	return func() tea.Msg {
		// Download
		file, err := m.tg.DownloadFileSync(fileId)
		if err != nil {
			return MediaPlayMsg{Status: "error", Info: fmt.Sprintf("Download error: %v", err)}
		}

		path := file.Local.Path

		// Play based on type
		var cmd *exec.Cmd
		switch mediaType {
		case "voice", "audio":
			// Try mpv first, then ffplay
			if _, err := exec.LookPath("mpv"); err == nil {
				cmd = exec.Command("mpv", "--no-video", "--really-quiet", path)
			} else if _, err := exec.LookPath("ffplay"); err == nil {
				cmd = exec.Command("ffplay", "-nodisp", "-autoexit", "-loglevel", "quiet", path)
			} else {
				// Open with default app
				cmd = defaultOpenCmd(path)
			}
		case "video":
			if _, err := exec.LookPath("mpv"); err == nil {
				cmd = exec.Command("mpv", path)
			} else {
				cmd = defaultOpenCmd(path)
			}
		}

		if cmd != nil {
			cmd.Start()
			go cmd.Wait()
		}

		return MediaPlayMsg{Status: "playing", Info: statusMsg}
	}
}

func (m Model) downloadAndOpen(fileId int32, statusMsg string) tea.Cmd {
	return func() tea.Msg {
		file, err := m.tg.DownloadFileSync(fileId)
		if err != nil {
			return MediaPlayMsg{Status: "error", Info: fmt.Sprintf("Download error: %v", err)}
		}

		cmd := defaultOpenCmd(file.Local.Path)
		if cmd != nil {
			cmd.Start()
			go cmd.Wait()
		}

		return MediaPlayMsg{Status: "opened", Info: statusMsg}
	}
}

func defaultOpenCmd(path string) *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("open", path)
	case "windows":
		return exec.Command("cmd", "/c", "start", path)
	default:
		return exec.Command("xdg-open", path)
	}
}

func isOwnMessage(msg *client.Message, myUserId int64) bool {
	if s, ok := msg.SenderId.(*client.MessageSenderUser); ok {
		return s.UserId == myUserId
	}
	return false
}

func (m Model) View() string {
	if m.chatID == 0 {
		return lipgloss.NewStyle().
			Width(m.width).Height(m.height).
			Foreground(lipgloss.Color("244")).
			Align(lipgloss.Center, lipgloss.Center).
			Render("Select a chat\n\nTab to switch panels")
	}

	// Header
	headerText := "  " + m.chatTitle
	if m.mediaStatus != "" {
		headerText += "  │  " + m.mediaStatus
	}
	header := m.theme.ChatViewHeader.Width(m.width).Render(headerText)

	// Progress bar under header during loading
	progressBar := ""
	if m.loading {
		progressBar = renderProgressBar(m.width, m.loadProgress, m.loadStatus)
	}

	bodyH := m.height - 1
	if m.loading {
		bodyH -= 2 // progress bar takes 2 lines
	}
	if bodyH < 1 {
		bodyH = 1
	}

	messages := m.store.Messages.Get(m.chatID)

	if len(messages) == 0 {
		label := "No messages"
		if m.loading {
			label = m.loadStatus
		}
		body := lipgloss.NewStyle().
			Width(m.width).Height(bodyH).
			Foreground(lipgloss.Color("244")).
			Align(lipgloss.Center, lipgloss.Center).
			Render(label)
		return header + "\n" + body
	}

	var bubbles []string
	for _, msg := range messages {
		isOwn := isOwnMessage(msg, m.myUserId)
		bubble := m.renderer.RenderMessage(msg, m.store, isOwn, false, m.width)
		bubbles = append(bubbles, bubble)
	}

	allContent := strings.Join(bubbles, "\n")
	lines := strings.Split(allContent, "\n")

	total := len(lines)
	end := total - m.scrollOffset
	if end > total {
		end = total
	}
	if end < 0 {
		end = 0
	}
	start := end - bodyH
	if start < 0 {
		start = 0
	}

	visible := lines[start:end]

	for len(visible) < bodyH {
		visible = append([]string{""}, visible...)
	}
	if len(visible) > bodyH {
		visible = visible[len(visible)-bodyH:]
	}

	body := strings.Join(visible, "\n")
	if progressBar != "" {
		return header + "\n" + progressBar + "\n" + body
	}
	return header + "\n" + body
}

func renderProgressBar(width, percent int, label string) string {
	if percent > 100 {
		percent = 100
	}
	if percent < 0 {
		percent = 0
	}

	barW := width - 10
	if barW < 10 {
		barW = 10
	}
	filled := barW * percent / 100
	empty := barW - filled

	bar := lipgloss.NewStyle().Foreground(lipgloss.Color("39")).Render(strings.Repeat("━", filled))
	bar += lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", empty))
	pct := fmt.Sprintf(" %3d%%", percent)

	line1 := lipgloss.NewStyle().Foreground(lipgloss.Color("244")).PaddingLeft(1).Render(label)
	line2 := "  " + bar + pct

	return line1 + "\n" + line2
}

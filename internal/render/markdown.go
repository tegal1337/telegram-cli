package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/tegal1337/telegram-cli/internal/media"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/zelenin/go-tdlib/client"
)

type MessageRenderer struct {
	theme    *theme.Theme
	glamour  *glamour.TermRenderer
	imgCache *media.Cache
	imgRend  *media.ImageRenderer
}

func NewMessageRenderer(th *theme.Theme) *MessageRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)
	protocol := media.DetectProtocol()
	return &MessageRenderer{
		theme:    th,
		glamour:  r,
		imgCache: media.NewCache(50),
		imgRend:  media.NewImageRenderer(protocol, 50, 25),
	}
}

func (r *MessageRenderer) RenderMessage(msg *client.Message, s *store.Store, isOwn, isSelected bool, maxWidth int) string {
	if msg == nil {
		return ""
	}

	// Sender name
	senderName := r.getSenderName(msg, s)

	// Content
	content := r.renderContent(msg.Content, s, maxWidth-8)
	if content == "" {
		content = "[empty]"
	}

	timeStr := FormatTimestamp(msg.Date)

	var lines []string

	if msg.ForwardInfo != nil {
		lines = append(lines, lipgloss.NewStyle().Foreground(r.theme.TextMuted).Italic(true).Render("↪ Forwarded"))
	}

	if msg.ReplyTo != nil {
		if rt, ok := msg.ReplyTo.(*client.MessageReplyToMessage); ok {
			lines = append(lines, lipgloss.NewStyle().Foreground(r.theme.Primary).Italic(true).Render(fmt.Sprintf("┃ reply #%d", rt.MessageId)))
		}
	}

	if !isOwn && senderName != "" {
		lines = append(lines, lipgloss.NewStyle().Foreground(r.theme.Accent).Bold(true).Render(senderName))
	}

	lines = append(lines, content)

	footer := lipgloss.NewStyle().Foreground(r.theme.TextMuted).Render(timeStr)
	if isOwn {
		if msg.Id < 0 {
			footer += " " + lipgloss.NewStyle().Foreground(r.theme.Warning).Render("⏳")
		} else {
			footer += " " + lipgloss.NewStyle().Foreground(r.theme.Success).Render("✓✓")
		}
	}
	lines = append(lines, footer)

	inner := strings.Join(lines, "\n")

	bubbleW := maxWidth * 65 / 100
	if bubbleW < 15 {
		bubbleW = 15
	}

	var bubble string
	if isOwn {
		style := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Foreground(lipgloss.Color("252")).
			PaddingLeft(1).PaddingRight(1).
			MaxWidth(bubbleW)
		if isSelected {
			style = style.Copy().BorderForeground(lipgloss.Color("214"))
		}
		bubble = style.Render(inner)

		w := lipgloss.Width(bubble)
		pad := maxWidth - w
		if pad > 0 {
			var padded []string
			for _, line := range strings.Split(bubble, "\n") {
				padded = append(padded, strings.Repeat(" ", pad)+line)
			}
			bubble = strings.Join(padded, "\n")
		}
	} else {
		style := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("245")).
			Foreground(lipgloss.Color("252")).
			PaddingLeft(1).PaddingRight(1).
			MaxWidth(bubbleW)
		if isSelected {
			style = style.Copy().BorderForeground(lipgloss.Color("214"))
		}
		bubble = style.Render(inner)
	}

	return bubble
}

func (r *MessageRenderer) getSenderName(msg *client.Message, s *store.Store) string {
	switch sender := msg.SenderId.(type) {
	case *client.MessageSenderUser:
		name := s.Users.DisplayName(sender.UserId)
		if name != "Unknown" {
			return name
		}
		// Fallback: use first/last from message if available
		return fmt.Sprintf("User#%d", sender.UserId)
	case *client.MessageSenderChat:
		if entry, ok := s.Chats.Get(sender.ChatId); ok && entry.Chat != nil {
			return entry.Chat.Title
		}
		return fmt.Sprintf("Chat#%d", sender.ChatId)
	}
	return ""
}

func (r *MessageRenderer) renderContent(content client.MessageContent, s *store.Store, maxWidth int) string {
	if content == nil {
		return "[unsupported]"
	}

	switch c := content.(type) {
	case *client.MessageText:
		if c.Text == nil || c.Text.Text == "" {
			return "[empty]"
		}
		md := EntitiesToMarkdown(c.Text)
		if r.glamour != nil && strings.Contains(md, "```") {
			rendered, err := r.glamour.Render(md)
			if err == nil {
				return strings.TrimSpace(rendered)
			}
		}
		return EntitiesToANSI(c.Text)

	case *client.MessagePhoto:
		imgStr := r.renderPhoto(c.Photo, s)
		caption := ""
		if c.Caption != nil && c.Caption.Text != "" {
			caption = "\n" + c.Caption.Text
		}
		return imgStr + caption

	case *client.MessageVideo:
		s := fmt.Sprintf("🎥 Video [%s]", fmtDur(c.Video.Duration))
		if c.Caption != nil && c.Caption.Text != "" {
			s += "\n" + c.Caption.Text
		}
		return s

	case *client.MessageDocument:
		s := fmt.Sprintf("📎 %s (%s)", c.Document.FileName, fmtSize(c.Document.Document.ExpectedSize))
		if c.Caption != nil && c.Caption.Text != "" {
			s += "\n" + c.Caption.Text
		}
		return s

	case *client.MessageVoiceNote:
		s := fmt.Sprintf("🎤 Voice [%s]", fmtDur(c.VoiceNote.Duration))
		if c.Caption != nil && c.Caption.Text != "" {
			s += "\n" + c.Caption.Text
		}
		return s

	case *client.MessageVideoNote:
		return fmt.Sprintf("📹 Video msg [%s]", fmtDur(c.VideoNote.Duration))

	case *client.MessageSticker:
		return c.Sticker.Emoji + " Sticker"

	case *client.MessageAnimation:
		return "🎬 GIF"

	case *client.MessageAudio:
		title := c.Audio.Title
		if title == "" {
			title = c.Audio.FileName
		}
		return fmt.Sprintf("🎵 %s [%s]", title, fmtDur(c.Audio.Duration))

	case *client.MessageLocation:
		return fmt.Sprintf("📍 %.4f, %.4f", c.Location.Latitude, c.Location.Longitude)

	case *client.MessageContact:
		return fmt.Sprintf("👤 %s %s", c.Contact.FirstName, c.Contact.LastName)

	case *client.MessagePoll:
		return fmt.Sprintf("📊 %s", c.Poll.Question.Text)

	case *client.MessagePinMessage:
		return "📌 Pinned"
	case *client.MessageChatAddMembers:
		return "➕ Members added"
	case *client.MessageChatDeleteMember:
		return "➖ Member left"
	case *client.MessageChatChangeTitle:
		return "✏ " + c.Title
	case *client.MessageChatChangePhoto:
		return "🖼 Photo changed"
	case *client.MessageChatJoinByLink:
		return "🔗 Joined via link"
	default:
		return fmt.Sprintf("[%s]", content.MessageContentType())
	}
}

func (r *MessageRenderer) renderPhoto(photo *client.Photo, s *store.Store) string {
	if photo == nil || len(photo.Sizes) == 0 {
		return "🖼  [Photo]"
	}

	// Try to find a downloaded file in the photo sizes
	// Check smallest first (thumbnail), then larger
	for _, size := range photo.Sizes {
		if size.Photo == nil {
			continue
		}
		fileState, ok := s.Files.Get(size.Photo.Id)
		if ok && fileState.IsComplete && fileState.LocalPath != "" {
			// Check cache first
			cacheKey := fmt.Sprintf("img:%d", size.Photo.Id)
			if cached, ok := r.imgCache.Get(cacheKey); ok {
				return cached
			}
			// Render image from local file
			rendered, err := r.imgRend.RenderFile(fileState.LocalPath)
			if err == nil && rendered != "" {
				r.imgCache.Set(cacheKey, rendered)
				return rendered
			}
		}

		// Check if local path is available directly from TDLib file info
		if size.Photo.Local != nil && size.Photo.Local.IsDownloadingCompleted && size.Photo.Local.Path != "" {
			cacheKey := fmt.Sprintf("img:%d", size.Photo.Id)
			if cached, ok := r.imgCache.Get(cacheKey); ok {
				return cached
			}
			rendered, err := r.imgRend.RenderFile(size.Photo.Local.Path)
			if err == nil && rendered != "" {
				r.imgCache.Set(cacheKey, rendered)
				return rendered
			}
		}
	}

	// Not downloaded yet — show placeholder with size info
	best := photo.Sizes[len(photo.Sizes)-1]
	return fmt.Sprintf("🖼  [Photo %dx%d] ⬇ downloading...", best.Width, best.Height)
}

func fmtDur(s int32) string {
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func fmtSize(b int64) string {
	switch {
	case b >= 1<<30:
		return fmt.Sprintf("%.1fGB", float64(b)/(1<<30))
	case b >= 1<<20:
		return fmt.Sprintf("%.1fMB", float64(b)/(1<<20))
	case b >= 1<<10:
		return fmt.Sprintf("%.1fKB", float64(b)/(1<<10))
	default:
		return fmt.Sprintf("%dB", b)
	}
}

package render

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss/v2"
	"github.com/tegal1337/telegram-cli/internal/store"
	"github.com/tegal1337/telegram-cli/internal/ui/theme"
	"github.com/zelenin/go-tdlib/client"
)

// MessageRenderer renders Telegram messages for terminal display.
type MessageRenderer struct {
	theme    *theme.Theme
	glamour  *glamour.TermRenderer
}

// NewMessageRenderer creates a new message renderer.
func NewMessageRenderer(th *theme.Theme) *MessageRenderer {
	r, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(80),
	)

	return &MessageRenderer{
		theme:   th,
		glamour: r,
	}
}

// RenderMessage renders a single message for terminal display.
func (r *MessageRenderer) RenderMessage(msg *client.Message, s *store.Store, isOwn, isSelected bool, maxWidth int) string {
	if msg == nil {
		return ""
	}

	// Get sender name.
	senderName := ""
	if sender, ok := msg.SenderId.(*client.MessageSenderUser); ok {
		senderName = s.Users.DisplayName(sender.UserID)
	} else if sender, ok := msg.SenderId.(*client.MessageSenderChat); ok {
		if entry, ok := s.Chats.Get(sender.ChatID); ok && entry.Chat != nil {
			senderName = entry.Chat.Title
		}
	}

	// Render message content.
	content := r.renderContent(msg.Content, maxWidth-4)

	// Format timestamp.
	time := FormatTimestamp(msg.Date)

	// Message status (for own messages).
	status := ""
	if isOwn {
		if msg.ID < 0 {
			status = "⏳" // sending
		} else {
			status = "✓" // sent
		}
	}

	// Build the message bubble.
	bubbleStyle := r.theme.MessageBubbleOther
	if isOwn {
		bubbleStyle = r.theme.MessageBubbleOwn
	}

	if isSelected {
		bubbleStyle = bubbleStyle.Copy().
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(r.theme.Primary)
	}

	// Header line: sender + time.
	header := ""
	if !isOwn && senderName != "" {
		header = r.theme.MessageSender.Render(senderName) + "  " +
			r.theme.MessageTime.Render(time)
	} else {
		header = r.theme.MessageTime.Render(time)
		if status != "" {
			header += " " + r.theme.MessageStatus.Render(status)
		}
	}

	// Reply context.
	replyLine := ""
	if msg.ReplyTo != nil {
		if replyTo, ok := msg.ReplyTo.(*client.MessageReplyToMessage); ok {
			replyLine = r.theme.MessageReply.Render(
				fmt.Sprintf("Reply to message %d", replyTo.MessageID),
			) + "\n"
		}
	}

	// Forward info.
	forwardLine := ""
	if msg.ForwardInfo != nil {
		forwardLine = r.theme.MessageTime.Render("↪ Forwarded") + "\n"
	}

	bubble := fmt.Sprintf("%s%s%s\n%s", forwardLine, replyLine, header, content)

	// Align own messages to the right.
	if isOwn {
		return bubbleStyle.
			Width(maxWidth).
			Align(lipgloss.Right).
			Render(bubble)
	}

	return bubbleStyle.
		Width(maxWidth).
		Render(bubble)
}

func (r *MessageRenderer) renderContent(content client.MessageContent, maxWidth int) string {
	switch c := content.(type) {
	case *client.MessageText:
		md := EntitiesToMarkdown(c.Text)
		// Try glamour rendering for rich text.
		if r.glamour != nil && strings.Contains(md, "```") {
			rendered, err := r.glamour.Render(md)
			if err == nil {
				return strings.TrimSpace(rendered)
			}
		}
		return EntitiesToANSI(c.Text)

	case *client.MessagePhoto:
		caption := ""
		if c.Caption != nil && c.Caption.Text != "" {
			caption = "\n" + EntitiesToANSI(c.Caption)
		}
		size := ""
		if len(c.Photo.Sizes) > 0 {
			best := c.Photo.Sizes[len(c.Photo.Sizes)-1]
			size = fmt.Sprintf(" %dx%d", best.Width, best.Height)
		}
		return fmt.Sprintf("📷 [Photo%s]%s", size, caption)

	case *client.MessageVideo:
		caption := ""
		if c.Caption != nil && c.Caption.Text != "" {
			caption = "\n" + EntitiesToANSI(c.Caption)
		}
		duration := formatDuration(c.Video.Duration)
		return fmt.Sprintf("🎥 [Video %s]%s", duration, caption)

	case *client.MessageDocument:
		caption := ""
		if c.Caption != nil && c.Caption.Text != "" {
			caption = "\n" + EntitiesToANSI(c.Caption)
		}
		size := formatFileSize(c.Document.Document.ExpectedSize)
		return fmt.Sprintf("📎 %s (%s)%s", c.Document.FileName, size, caption)

	case *client.MessageVoiceNote:
		duration := formatDuration(c.VoiceNote.Duration)
		caption := ""
		if c.Caption != nil && c.Caption.Text != "" {
			caption = "\n" + EntitiesToANSI(c.Caption)
		}
		return fmt.Sprintf("🎤 [Voice %s] ▶ Press Enter to play%s", duration, caption)

	case *client.MessageVideoNote:
		duration := formatDuration(c.VideoNote.Duration)
		return fmt.Sprintf("📹 [Video message %s]", duration)

	case *client.MessageSticker:
		return fmt.Sprintf("%s [Sticker]", c.Sticker.Emoji)

	case *client.MessageAnimation:
		caption := ""
		if c.Caption != nil && c.Caption.Text != "" {
			caption = "\n" + EntitiesToANSI(c.Caption)
		}
		return fmt.Sprintf("🎬 [GIF]%s", caption)

	case *client.MessageAudio:
		caption := ""
		if c.Caption != nil && c.Caption.Text != "" {
			caption = "\n" + EntitiesToANSI(c.Caption)
		}
		title := c.Audio.Title
		if title == "" {
			title = c.Audio.FileName
		}
		duration := formatDuration(c.Audio.Duration)
		return fmt.Sprintf("🎵 %s (%s)%s", title, duration, caption)

	case *client.MessageLocation:
		return fmt.Sprintf("📍 Location (%.4f, %.4f)", c.Location.Latitude, c.Location.Longitude)

	case *client.MessageContact:
		return fmt.Sprintf("👤 Contact: %s %s (%s)", c.Contact.FirstName, c.Contact.LastName, c.Contact.PhoneNumber)

	case *client.MessagePoll:
		var options []string
		for _, opt := range c.Poll.Options {
			bar := strings.Repeat("█", int(opt.VotePercentage)/5)
			options = append(options, fmt.Sprintf("  %s %d%% %s", opt.Text.Text, opt.VotePercentage, bar))
		}
		return fmt.Sprintf("📊 %s\n%s", c.Poll.Question.Text, strings.Join(options, "\n"))

	case *client.MessagePinMessage:
		return "📌 Message pinned"

	case *client.MessageChatAddMembers:
		return "➕ Members added to the group"

	case *client.MessageChatDeleteMember:
		return "➖ Member removed from the group"

	case *client.MessageChatChangeTitle:
		return fmt.Sprintf("✏ Chat title changed to: %s", c.Title)

	case *client.MessageChatChangePhoto:
		return "🖼 Chat photo changed"

	case *client.MessageChatJoinByLink:
		return "🔗 Joined via invite link"

	default:
		return "[Unsupported message type]"
	}
}

func formatDuration(seconds int32) string {
	m := seconds / 60
	s := seconds % 60
	return fmt.Sprintf("%d:%02d", m, s)
}

func formatFileSize(bytes int64) string {
	switch {
	case bytes >= 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

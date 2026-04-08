package theme

import (
	"github.com/charmbracelet/lipgloss/v2"
)

// Theme holds all the styles used across the TUI.
type Theme struct {
	// Base colors
	Primary     lipgloss.Color
	Secondary   lipgloss.Color
	Accent      lipgloss.Color
	Background  lipgloss.Color
	Surface     lipgloss.Color
	Text        lipgloss.Color
	TextMuted   lipgloss.Color
	Error       lipgloss.Color
	Success     lipgloss.Color
	Warning     lipgloss.Color

	// Chat list styles
	ChatListPane        lipgloss.Style
	ChatListItem        lipgloss.Style
	ChatListItemActive  lipgloss.Style
	ChatListTitle       lipgloss.Style
	ChatListPreview     lipgloss.Style
	ChatListTime        lipgloss.Style
	ChatListUnread      lipgloss.Style
	ChatListOnline      lipgloss.Style

	// Chat view styles
	ChatViewPane        lipgloss.Style
	ChatViewHeader      lipgloss.Style
	MessageBubbleOwn    lipgloss.Style
	MessageBubbleOther  lipgloss.Style
	MessageSender       lipgloss.Style
	MessageTime         lipgloss.Style
	MessageStatus       lipgloss.Style
	MessageReply        lipgloss.Style
	MessageSystem       lipgloss.Style

	// Composer styles
	ComposerPane        lipgloss.Style
	ComposerInput       lipgloss.Style
	ComposerHint        lipgloss.Style
	ComposerReplyBar    lipgloss.Style

	// Status bar styles
	StatusBar           lipgloss.Style
	StatusBarConnected  lipgloss.Style
	StatusBarTyping     lipgloss.Style

	// Dialog styles
	DialogOverlay       lipgloss.Style
	DialogBox           lipgloss.Style
	DialogTitle         lipgloss.Style
	DialogButton        lipgloss.Style
	DialogButtonActive  lipgloss.Style

	// Auth screen styles
	AuthPane            lipgloss.Style
	AuthTitle           lipgloss.Style
	AuthInput           lipgloss.Style
	AuthLabel           lipgloss.Style

	// Search styles
	SearchInput         lipgloss.Style
	SearchResult        lipgloss.Style
	SearchResultActive  lipgloss.Style

	// Generic styles
	Border              lipgloss.Style
	FocusedBorder       lipgloss.Style
	Separator           lipgloss.Style
	Badge               lipgloss.Style
	Spinner             lipgloss.Style
	ProgressBar         lipgloss.Style
	ProgressBarFill     lipgloss.Style
	Tab                 lipgloss.Style
	TabActive           lipgloss.Style
}

// DarkTheme returns the default dark theme.
func DarkTheme() *Theme {
	primary := lipgloss.Color("#7AA2F7")
	secondary := lipgloss.Color("#9ECE6A")
	accent := lipgloss.Color("#BB9AF7")
	bg := lipgloss.Color("#1A1B26")
	surface := lipgloss.Color("#24283B")
	text := lipgloss.Color("#C0CAF5")
	textMuted := lipgloss.Color("#565F89")
	errColor := lipgloss.Color("#F7768E")
	success := lipgloss.Color("#9ECE6A")
	warning := lipgloss.Color("#E0AF68")

	return &Theme{
		Primary:    primary,
		Secondary:  secondary,
		Accent:     accent,
		Background: bg,
		Surface:    surface,
		Text:       text,
		TextMuted:  textMuted,
		Error:      errColor,
		Success:    success,
		Warning:    warning,

		ChatListPane: lipgloss.NewStyle().
			Background(bg).
			Foreground(text).
			BorderRight(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(textMuted),

		ChatListItem: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Foreground(text),

		ChatListItemActive: lipgloss.NewStyle().
			PaddingLeft(1).
			PaddingRight(1).
			Background(surface).
			Foreground(primary).
			Bold(true),

		ChatListTitle: lipgloss.NewStyle().
			Foreground(text).
			Bold(true),

		ChatListPreview: lipgloss.NewStyle().
			Foreground(textMuted),

		ChatListTime: lipgloss.NewStyle().
			Foreground(textMuted).
			Align(lipgloss.Right),

		ChatListUnread: lipgloss.NewStyle().
			Background(primary).
			Foreground(bg).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1),

		ChatListOnline: lipgloss.NewStyle().
			Foreground(success),

		ChatViewPane: lipgloss.NewStyle().
			Background(bg).
			Foreground(text),

		ChatViewHeader: lipgloss.NewStyle().
			Background(surface).
			Foreground(text).
			Bold(true).
			PaddingLeft(2).
			PaddingRight(2).
			Height(1),

		MessageBubbleOwn: lipgloss.NewStyle().
			Background(lipgloss.Color("#2E3A59")).
			Foreground(text).
			PaddingLeft(1).
			PaddingRight(1).
			MarginLeft(4),

		MessageBubbleOther: lipgloss.NewStyle().
			Background(surface).
			Foreground(text).
			PaddingLeft(1).
			PaddingRight(1).
			MarginRight(4),

		MessageSender: lipgloss.NewStyle().
			Foreground(accent).
			Bold(true),

		MessageTime: lipgloss.NewStyle().
			Foreground(textMuted),

		MessageStatus: lipgloss.NewStyle().
			Foreground(secondary),

		MessageReply: lipgloss.NewStyle().
			Foreground(primary).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(primary).
			PaddingLeft(1),

		MessageSystem: lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true).
			Align(lipgloss.Center),

		ComposerPane: lipgloss.NewStyle().
			Background(surface).
			BorderTop(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(textMuted),

		ComposerInput: lipgloss.NewStyle().
			Background(surface).
			Foreground(text).
			PaddingLeft(1),

		ComposerHint: lipgloss.NewStyle().
			Foreground(textMuted).
			Italic(true),

		ComposerReplyBar: lipgloss.NewStyle().
			Foreground(primary).
			Background(surface).
			PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(primary),

		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("#16161E")).
			Foreground(textMuted).
			Height(1).
			PaddingLeft(1).
			PaddingRight(1),

		StatusBarConnected: lipgloss.NewStyle().
			Foreground(success),

		StatusBarTyping: lipgloss.NewStyle().
			Foreground(warning).
			Italic(true),

		DialogOverlay: lipgloss.NewStyle().
			Background(lipgloss.Color("#00000088")),

		DialogBox: lipgloss.NewStyle().
			Background(surface).
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			Padding(1, 2),

		DialogTitle: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			MarginBottom(1),

		DialogButton: lipgloss.NewStyle().
			Background(surface).
			Foreground(text).
			Padding(0, 2),

		DialogButtonActive: lipgloss.NewStyle().
			Background(primary).
			Foreground(bg).
			Bold(true).
			Padding(0, 2),

		AuthPane: lipgloss.NewStyle().
			Background(bg).
			Foreground(text).
			Align(lipgloss.Center),

		AuthTitle: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			MarginBottom(2),

		AuthInput: lipgloss.NewStyle().
			Background(surface).
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			Padding(0, 1).
			Width(40),

		AuthLabel: lipgloss.NewStyle().
			Foreground(text).
			MarginBottom(1),

		SearchInput: lipgloss.NewStyle().
			Background(surface).
			Foreground(text).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			PaddingLeft(1),

		SearchResult: lipgloss.NewStyle().
			Foreground(text).
			PaddingLeft(2),

		SearchResultActive: lipgloss.NewStyle().
			Foreground(primary).
			Background(surface).
			Bold(true).
			PaddingLeft(2),

		Border: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(textMuted),

		FocusedBorder: lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(primary),

		Separator: lipgloss.NewStyle().
			Foreground(textMuted),

		Badge: lipgloss.NewStyle().
			Background(primary).
			Foreground(bg).
			Bold(true).
			PaddingLeft(1).
			PaddingRight(1),

		Spinner: lipgloss.NewStyle().
			Foreground(primary),

		ProgressBar: lipgloss.NewStyle().
			Foreground(textMuted),

		ProgressBarFill: lipgloss.NewStyle().
			Foreground(primary),

		Tab: lipgloss.NewStyle().
			Foreground(textMuted).
			Padding(0, 2),

		TabActive: lipgloss.NewStyle().
			Foreground(primary).
			Bold(true).
			Padding(0, 2).
			BorderBottom(true).
			BorderStyle(lipgloss.NormalBorder()).
			BorderForeground(primary),
	}
}

// LightTheme returns a light color theme.
func LightTheme() *Theme {
	primary := lipgloss.Color("#2E59A8")
	secondary := lipgloss.Color("#587B2E")
	accent := lipgloss.Color("#7B4BAA")
	bg := lipgloss.Color("#FFFFFF")
	surface := lipgloss.Color("#F0F0F0")
	text := lipgloss.Color("#1A1A1A")
	textMuted := lipgloss.Color("#888888")
	errColor := lipgloss.Color("#CC3333")
	success := lipgloss.Color("#338833")
	warning := lipgloss.Color("#CC8833")

	t := DarkTheme()
	t.Primary = primary
	t.Secondary = secondary
	t.Accent = accent
	t.Background = bg
	t.Surface = surface
	t.Text = text
	t.TextMuted = textMuted
	t.Error = errColor
	t.Success = success
	t.Warning = warning

	return t
}

// ForName returns the appropriate theme for the given name.
func ForName(name string) *Theme {
	switch name {
	case "light":
		return LightTheme()
	default:
		return DarkTheme()
	}
}

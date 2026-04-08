package theme

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
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

	// Panel styles
	PanelNormal  lipgloss.Style
	PanelFocused lipgloss.Style

	// Chat list
	ChatListItem       lipgloss.Style
	ChatListItemActive lipgloss.Style
	ChatListTitle      lipgloss.Style
	ChatListPreview    lipgloss.Style
	ChatListTime       lipgloss.Style
	ChatListUnread     lipgloss.Style
	ChatListOnline     lipgloss.Style

	// Chat view
	ChatViewHeader     lipgloss.Style
	MessageBubbleOwn   lipgloss.Style
	MessageBubbleOther lipgloss.Style
	MessageSender      lipgloss.Style
	MessageTime        lipgloss.Style
	MessageStatus      lipgloss.Style
	MessageReply       lipgloss.Style
	MessageSystem      lipgloss.Style

	// Composer
	ComposerPane     lipgloss.Style
	ComposerInput    lipgloss.Style
	ComposerHint     lipgloss.Style
	ComposerReplyBar lipgloss.Style

	// Status bar
	StatusBar          lipgloss.Style
	StatusBarConnected lipgloss.Style
	StatusBarTyping    lipgloss.Style

	// Dialog
	DialogBox          lipgloss.Style
	DialogTitle        lipgloss.Style
	DialogButton       lipgloss.Style
	DialogButtonActive lipgloss.Style

	// Auth
	AuthPane  lipgloss.Style
	AuthTitle lipgloss.Style
	AuthInput lipgloss.Style
	AuthLabel lipgloss.Style

	// Search
	SearchInput        lipgloss.Style
	SearchResult       lipgloss.Style
	SearchResultActive lipgloss.Style

	// Generic
	Spinner        lipgloss.Style
	Badge          lipgloss.Style
	Tab            lipgloss.Style
	TabActive      lipgloss.Style
	ProgressBar    lipgloss.Style
	ProgressBarFill lipgloss.Style
}

func DarkTheme() *Theme {
	primary := lipgloss.Color("39")     // bright blue
	secondary := lipgloss.Color("42")   // green
	accent := lipgloss.Color("177")     // purple
	bg := lipgloss.Color("234")         // dark gray
	surface := lipgloss.Color("236")    // slightly lighter
	text := lipgloss.Color("252")       // light gray
	textMuted := lipgloss.Color("244")  // mid gray
	errColor := lipgloss.Color("196")   // red
	success := lipgloss.Color("42")     // green
	warning := lipgloss.Color("214")    // orange

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

		PanelNormal: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")),

		PanelFocused: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary),

		ChatListItem: lipgloss.NewStyle().
			PaddingLeft(1).PaddingRight(1).
			Foreground(text),

		ChatListItemActive: lipgloss.NewStyle().
			PaddingLeft(1).PaddingRight(1).
			Background(lipgloss.Color("237")).
			Foreground(primary).Bold(true),

		ChatListTitle: lipgloss.NewStyle().
			Foreground(text).Bold(true),

		ChatListPreview: lipgloss.NewStyle().
			Foreground(textMuted),

		ChatListTime: lipgloss.NewStyle().
			Foreground(textMuted),

		ChatListUnread: lipgloss.NewStyle().
			Background(primary).Foreground(lipgloss.Color("232")).
			Bold(true).Padding(0, 1),

		ChatListOnline: lipgloss.NewStyle().
			Foreground(success),

		ChatViewHeader: lipgloss.NewStyle().
			Foreground(text).Bold(true).
			PaddingLeft(1).PaddingRight(1).
			Background(lipgloss.Color("236")),

		MessageBubbleOwn: lipgloss.NewStyle().
			Foreground(text).
			Background(lipgloss.Color("24")).
			Padding(0, 1),

		MessageBubbleOther: lipgloss.NewStyle().
			Foreground(text).
			Background(lipgloss.Color("237")).
			Padding(0, 1),

		MessageSender: lipgloss.NewStyle().
			Foreground(accent).Bold(true),

		MessageTime: lipgloss.NewStyle().
			Foreground(textMuted),

		MessageStatus: lipgloss.NewStyle().
			Foreground(secondary),

		MessageReply: lipgloss.NewStyle().
			Foreground(primary).Italic(true).
			PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(primary),

		MessageSystem: lipgloss.NewStyle().
			Foreground(textMuted).Italic(true),

		ComposerPane: lipgloss.NewStyle(),

		ComposerInput: lipgloss.NewStyle().
			Foreground(text).PaddingLeft(1),

		ComposerHint: lipgloss.NewStyle().
			Foreground(textMuted).Italic(true).PaddingLeft(1),

		ComposerReplyBar: lipgloss.NewStyle().
			Foreground(primary).PaddingLeft(1).
			BorderLeft(true).
			BorderStyle(lipgloss.ThickBorder()).
			BorderForeground(primary),

		StatusBar: lipgloss.NewStyle().
			Background(lipgloss.Color("235")).
			Foreground(textMuted).
			PaddingLeft(1).PaddingRight(1),

		StatusBarConnected: lipgloss.NewStyle().
			Foreground(success).Bold(true),

		StatusBarTyping: lipgloss.NewStyle().
			Foreground(warning).Italic(true),

		DialogBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			Padding(1, 2).
			Background(lipgloss.Color("236")),

		DialogTitle: lipgloss.NewStyle().
			Foreground(primary).Bold(true),

		DialogButton: lipgloss.NewStyle().
			Foreground(text).Padding(0, 2),

		DialogButtonActive: lipgloss.NewStyle().
			Background(primary).
			Foreground(lipgloss.Color("232")).
			Bold(true).Padding(0, 2),

		AuthPane: lipgloss.NewStyle().
			Foreground(text),

		AuthTitle: lipgloss.NewStyle().
			Foreground(primary).Bold(true),

		AuthInput: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			Foreground(text).
			Padding(0, 1).Width(40),

		AuthLabel: lipgloss.NewStyle().
			Foreground(text),

		SearchInput: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primary).
			Foreground(text).PaddingLeft(1),

		SearchResult: lipgloss.NewStyle().
			Foreground(text).PaddingLeft(2),

		SearchResultActive: lipgloss.NewStyle().
			Foreground(primary).Bold(true).
			Background(lipgloss.Color("237")).
			PaddingLeft(2),

		Spinner: lipgloss.NewStyle().Foreground(primary),

		Badge: lipgloss.NewStyle().
			Background(primary).Foreground(lipgloss.Color("232")).
			Bold(true).Padding(0, 1),

		Tab: lipgloss.NewStyle().
			Foreground(textMuted).Padding(0, 2),

		TabActive: lipgloss.NewStyle().
			Foreground(primary).Bold(true).Padding(0, 2).
			Underline(true),

		ProgressBar: lipgloss.NewStyle().Foreground(textMuted),

		ProgressBarFill: lipgloss.NewStyle().Foreground(primary),
	}
}

func LightTheme() *Theme {
	t := DarkTheme()
	t.Primary = lipgloss.Color("33")
	t.Background = lipgloss.Color("231")
	t.Surface = lipgloss.Color("254")
	t.Text = lipgloss.Color("234")
	t.TextMuted = lipgloss.Color("245")
	return t
}

func ForName(name string) *Theme {
	if name == "light" {
		return LightTheme()
	}
	return DarkTheme()
}

package layout

// Layout holds computed panel dimensions based on terminal size.
type Layout struct {
	Width  int
	Height int

	// Chat list panel
	ChatListWidth  int
	ChatListHeight int

	// Chat view panel
	ChatViewWidth  int
	ChatViewHeight int
	ChatViewHeaderHeight int

	// Composer panel
	ComposerWidth  int
	ComposerHeight int

	// Status bar
	StatusBarWidth  int
	StatusBarHeight int

	// Whether to use single-panel mode (narrow terminals)
	SinglePanel bool
}

// Compute calculates layout dimensions from terminal size and config.
func Compute(width, height, chatListWidthPercent int) Layout {
	l := Layout{
		Width:  width,
		Height: height,
	}

	if chatListWidthPercent <= 0 {
		chatListWidthPercent = 30
	}

	// Status bar is always 1 row at the bottom.
	l.StatusBarWidth = width
	l.StatusBarHeight = 1

	// Composer is 3 rows (border + input + padding).
	l.ComposerHeight = 3

	// Available height for chat list and chat view.
	availableHeight := height - l.StatusBarHeight

	// Single-panel mode for narrow terminals.
	if width < 80 {
		l.SinglePanel = true
		l.ChatListWidth = width
		l.ChatListHeight = availableHeight

		l.ChatViewWidth = width
		l.ChatViewHeaderHeight = 1
		l.ChatViewHeight = availableHeight - l.ComposerHeight - l.ChatViewHeaderHeight

		l.ComposerWidth = width
		return l
	}

	// Dual-panel mode.
	l.ChatListWidth = width * chatListWidthPercent / 100
	if l.ChatListWidth < 25 {
		l.ChatListWidth = 25
	}
	if l.ChatListWidth > 60 {
		l.ChatListWidth = 60
	}

	l.ChatListHeight = availableHeight

	l.ChatViewWidth = width - l.ChatListWidth
	l.ChatViewHeaderHeight = 1
	l.ChatViewHeight = availableHeight - l.ComposerHeight - l.ChatViewHeaderHeight

	l.ComposerWidth = l.ChatViewWidth

	return l
}

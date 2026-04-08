package layout

type Layout struct {
	Width  int
	Height int

	ChatListWidth  int
	ChatListHeight int

	ChatViewWidth  int
	ChatViewHeight int

	ComposerWidth  int
	ComposerHeight int

	StatusBarWidth  int
	StatusBarHeight int

	SinglePanel bool
}

func Compute(width, height, chatListWidthPercent int) Layout {
	l := Layout{
		Width:  width,
		Height: height,
	}

	if chatListWidthPercent <= 0 {
		chatListWidthPercent = 30
	}

	l.StatusBarWidth = width
	l.StatusBarHeight = 1

	// Composer: 3 lines + 2 border = 5
	l.ComposerHeight = 5

	availableHeight := height - l.StatusBarHeight

	if width < 80 {
		l.SinglePanel = true
		l.ChatListWidth = width
		l.ChatListHeight = availableHeight
		l.ChatViewWidth = width
		l.ChatViewHeight = availableHeight - l.ComposerHeight
		l.ComposerWidth = width
		return l
	}

	// Left panel includes its border (2 chars)
	l.ChatListWidth = width * chatListWidthPercent / 100
	if l.ChatListWidth < 28 {
		l.ChatListWidth = 28
	}
	if l.ChatListWidth > 60 {
		l.ChatListWidth = 60
	}
	l.ChatListHeight = availableHeight

	l.ChatViewWidth = width - l.ChatListWidth
	l.ChatViewHeight = availableHeight - l.ComposerHeight
	l.ComposerWidth = l.ChatViewWidth

	return l
}

package chatview

// ScrollToBottomMsg requests scrolling to the newest message.
type ScrollToBottomMsg struct{}

// LoadMoreHistoryMsg requests loading older messages.
type LoadMoreHistoryMsg struct {
	ChatID int64
}

// MessageActionMsg is emitted for message actions (reply, edit, delete, forward).
type MessageActionMsg struct {
	Action    string // "reply", "edit", "delete", "forward"
	ChatID    int64
	MessageID int64
}

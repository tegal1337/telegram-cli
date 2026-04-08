package chatlist

// ChatSelectedMsg is emitted when the user selects a chat.
type ChatSelectedMsg struct {
	ChatID int64
}

// ChatListFilteredMsg is emitted when the chat list filter changes.
type ChatListFilteredMsg struct {
	Query string
}

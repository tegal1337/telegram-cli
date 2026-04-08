package telegram

import (
	"github.com/zelenin/go-tdlib/client"
)

// Tea messages produced from TDLib updates.
// These are sent into the bubbletea program via p.Send().

// AuthStateMsg carries authorization state changes.
type AuthStateMsg struct {
	State client.AuthorizationState
}

// NewMessageMsg is sent when a new message arrives.
type NewMessageMsg struct {
	Message *client.Message
}

// MessageEditedMsg is sent when a message is edited.
type MessageEditedMsg struct {
	ChatID    int64
	MessageID int64
}

// MessageDeletedMsg is sent when messages are deleted.
type MessageDeletedMsg struct {
	ChatID     int64
	MessageIDs []int64
}

// ChatUpdateMsg is sent when chat metadata changes (title, photo, etc).
type ChatUpdateMsg struct {
	Chat *client.Chat
}

// ChatPositionMsg is sent when a chat's position in the list changes.
type ChatPositionMsg struct {
	ChatID    int64
	Positions []*client.ChatPosition
}

// ChatLastMessageMsg is sent when a chat's last message changes.
type ChatLastMessageMsg struct {
	ChatID      int64
	LastMessage *client.Message
	Positions   []*client.ChatPosition
}

// ChatReadInboxMsg is sent when the read inbox state changes.
type ChatReadInboxMsg struct {
	ChatID                int64
	LastReadInboxMessageID int64
	UnreadCount           int32
}

// ChatReadOutboxMsg is sent when the read outbox state changes.
type ChatReadOutboxMsg struct {
	ChatID                 int64
	LastReadOutboxMessageID int64
}

// UserStatusMsg is sent when a user's online status changes.
type UserStatusMsg struct {
	UserID int64
	Status client.UserStatus
}

// UserUpdateMsg is sent when user info changes.
type UserUpdateMsg struct {
	User *client.User
}

// FileUpdateMsg is sent when a file download/upload progress changes.
type FileUpdateMsg struct {
	File *client.File
}

// ChatActionMsg is sent when someone is typing or performing an action.
type ChatActionMsg struct {
	ChatID int64
	UserID int64
	Action client.ChatAction
}

// ConnectionStateMsg is sent when the network connection state changes.
type ConnectionStateMsg struct {
	State client.ConnectionState
}

// UnreadCountMsg is sent when global unread counts change.
type UnreadCountMsg struct {
	UnreadCount        int32
	UnreadUnmutedCount int32
}

// MessageSendSucceededMsg is sent when a message is successfully sent.
type MessageSendSucceededMsg struct {
	Message      *client.Message
	OldMessageID int64
}

// MessageSendFailedMsg is sent when a message fails to send.
type MessageSendFailedMsg struct {
	Message      *client.Message
	OldMessageID int64
	ErrorCode    int32
	ErrorMessage string
}

// SupergroupUpdateMsg is sent when supergroup info changes.
type SupergroupUpdateMsg struct {
	Supergroup *client.Supergroup
}

// BasicGroupUpdateMsg is sent when basic group info changes.
type BasicGroupUpdateMsg struct {
	BasicGroup *client.BasicGroup
}

// NotificationMsg is sent for new notifications.
type NotificationMsg struct {
	GroupID       int32
	Notifications []*client.Notification
}

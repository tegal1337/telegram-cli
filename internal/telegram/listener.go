package telegram

import (
	"log"

	tea "charm.land/bubbletea/v2"
	"github.com/zelenin/go-tdlib/client"
)

// Listener bridges TDLib's async updates into bubbletea's event loop.
type Listener struct {
	tdClient *client.Client
	program  *tea.Program
}

// NewListener creates a listener that converts TDLib updates to tea.Msg.
func NewListener(tdClient *client.Client, program *tea.Program) *Listener {
	return &Listener{
		tdClient: tdClient,
		program:  program,
	}
}

// Start begins listening for TDLib updates in a goroutine.
// It converts each update type to a corresponding tea.Msg and sends it
// into the bubbletea program via p.Send().
func (l *Listener) Start() {
	listener := l.tdClient.GetListener()

	go func() {
		defer listener.Close()

		for update := range listener.Updates {
			msg := l.convertUpdate(update)
			if msg != nil {
				l.program.Send(msg)
			}
		}
	}()
}

func (l *Listener) convertUpdate(update client.Type) tea.Msg {
	switch u := update.(type) {
	case *client.UpdateAuthorizationState:
		return AuthStateMsg{State: u.AuthorizationState}

	case *client.UpdateNewMessage:
		return NewMessageMsg{Message: u.Message}

	case *client.UpdateMessageEdited:
		return MessageEditedMsg{
			ChatID:    u.ChatID,
			MessageID: u.MessageID,
		}

	case *client.UpdateDeleteMessages:
		if !u.FromCache {
			return MessageDeletedMsg{
				ChatID:     u.ChatID,
				MessageIDs: u.MessageIDs,
			}
		}

	case *client.UpdateChatTitle:
		return ChatUpdateMsg{Chat: &client.Chat{ID: u.ChatID, Title: u.Title}}

	case *client.UpdateChatPosition:
		// Handled via ChatPositionMsg
		return ChatPositionMsg{
			ChatID:    u.ChatID,
			Positions: []*client.ChatPosition{u.Position},
		}

	case *client.UpdateChatLastMessage:
		return ChatLastMessageMsg{
			ChatID:      u.ChatID,
			LastMessage: u.LastMessage,
			Positions:   u.Positions,
		}

	case *client.UpdateChatReadInbox:
		return ChatReadInboxMsg{
			ChatID:                 u.ChatID,
			LastReadInboxMessageID: u.LastReadInboxMessageID,
			UnreadCount:            u.UnreadCount,
		}

	case *client.UpdateChatReadOutbox:
		return ChatReadOutboxMsg{
			ChatID:                  u.ChatID,
			LastReadOutboxMessageID: u.LastReadOutboxMessageID,
		}

	case *client.UpdateUserStatus:
		return UserStatusMsg{
			UserID: u.UserID,
			Status: u.Status,
		}

	case *client.UpdateUser:
		return UserUpdateMsg{User: u.User}

	case *client.UpdateFile:
		return FileUpdateMsg{File: u.File}

	case *client.UpdateChatAction:
		return ChatActionMsg{
			ChatID: u.ChatID,
			UserID: extractSenderUserID(u.SenderId),
			Action: u.Action,
		}

	case *client.UpdateConnectionState:
		return ConnectionStateMsg{State: u.State}

	case *client.UpdateUnreadMessageCount:
		return UnreadCountMsg{
			UnreadCount:        u.UnreadCount,
			UnreadUnmutedCount: u.UnreadUnmutedCount,
		}

	case *client.UpdateMessageSendSucceeded:
		return MessageSendSucceededMsg{
			Message:      u.Message,
			OldMessageID: u.OldMessageID,
		}

	case *client.UpdateMessageSendFailed:
		return MessageSendFailedMsg{
			Message:      u.Message,
			OldMessageID: u.OldMessageID,
			ErrorCode:    u.Error.Code,
			ErrorMessage: u.Error.Message,
		}

	case *client.UpdateSupergroup:
		return SupergroupUpdateMsg{Supergroup: u.Supergroup}

	case *client.UpdateBasicGroup:
		return BasicGroupUpdateMsg{BasicGroup: u.BasicGroup}

	case *client.UpdateNotificationGroup:
		return NotificationMsg{
			GroupID:       u.NotificationGroupID,
			Notifications: u.AddedNotifications,
		}

	default:
		log.Printf("unhandled update: %T", update)
	}

	return nil
}

func extractSenderUserID(sender client.MessageSender) int64 {
	switch s := sender.(type) {
	case *client.MessageSenderUser:
		return s.UserID
	default:
		return 0
	}
}

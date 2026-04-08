package telegram

import (
	tea "charm.land/bubbletea/v2"
	"github.com/zelenin/go-tdlib/client"
)

type Listener struct {
	tdClient *client.Client
	program  *tea.Program
}

func NewListener(tdClient *client.Client, program *tea.Program) *Listener {
	return &Listener{
		tdClient: tdClient,
		program:  program,
	}
}

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
			ChatId:    u.ChatId,
			MessageId: u.MessageId,
		}

	case *client.UpdateDeleteMessages:
		if !u.FromCache {
			return MessageDeletedMsg{
				ChatId:     u.ChatId,
				MessageIds: u.MessageIds,
			}
		}

	case *client.UpdateNewChat:
		return ChatUpdateMsg{Chat: u.Chat}

	case *client.UpdateChatTitle:
		return ChatUpdateMsg{Chat: &client.Chat{Id: u.ChatId, Title: u.Title}}

	case *client.UpdateChatPosition:
		return ChatPositionMsg{
			ChatId:    u.ChatId,
			Positions: []*client.ChatPosition{u.Position},
		}

	case *client.UpdateChatLastMessage:
		return ChatLastMessageMsg{
			ChatId:      u.ChatId,
			LastMessage: u.LastMessage,
			Positions:   u.Positions,
		}

	case *client.UpdateChatReadInbox:
		return ChatReadInboxMsg{
			ChatId:                 u.ChatId,
			LastReadInboxMessageId: u.LastReadInboxMessageId,
			UnreadCount:            u.UnreadCount,
		}

	case *client.UpdateChatReadOutbox:
		return ChatReadOutboxMsg{
			ChatId:                  u.ChatId,
			LastReadOutboxMessageId: u.LastReadOutboxMessageId,
		}

	case *client.UpdateUserStatus:
		return UserStatusMsg{
			UserId: u.UserId,
			Status: u.Status,
		}

	case *client.UpdateUser:
		return UserUpdateMsg{User: u.User}

	case *client.UpdateFile:
		return FileUpdateMsg{File: u.File}

	case *client.UpdateChatAction:
		return ChatActionMsg{
			ChatId: u.ChatId,
			UserId: extractSenderUserId(u.SenderId),
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
			OldMessageId: u.OldMessageId,
		}

	case *client.UpdateMessageSendFailed:
		return MessageSendFailedMsg{
			Message:      u.Message,
			OldMessageId: u.OldMessageId,
			ErrorCode:    u.Error.Code,
			ErrorMessage: u.Error.Message,
		}

	case *client.UpdateSupergroup:
		return SupergroupUpdateMsg{Supergroup: u.Supergroup}

	case *client.UpdateBasicGroup:
		return BasicGroupUpdateMsg{BasicGroup: u.BasicGroup}

	case *client.UpdateNotificationGroup:
		return NotificationMsg{
			GroupId:       u.NotificationGroupId,
			Notifications: u.AddedNotifications,
		}

	case *client.UpdateChatAddedToList,
		*client.UpdateChatRemovedFromList,
		*client.UpdateUserFullInfo,
		*client.UpdateSupergroupFullInfo,
		*client.UpdateBasicGroupFullInfo,
		*client.UpdateChatActionBar,
		*client.UpdateChatHasScheduledMessages,
		*client.UpdateChatIsMarkedAsUnread,
		*client.UpdateChatNotificationSettings,
		*client.UpdateChatUnreadMentionCount,
		*client.UpdateChatUnreadReactionCount,
		*client.UpdateChatDraftMessage,
		*client.UpdateChatPhoto,
		*client.UpdateChatPermissions,
		*client.UpdateChatTheme,
		*client.UpdateChatAvailableReactions,
		*client.UpdateOption,
		*client.UpdateAnimationSearchParameters,
		*client.UpdateScopeNotificationSettings,
		*client.UpdateHavePendingNotifications,
		*client.UpdateChatFolders,
		*client.UpdateChatOnlineMemberCount,
		*client.UpdateAttachmentMenuBots,
		*client.UpdateActiveEmojiReactions,
		*client.UpdateDefaultReactionType,
		*client.UpdateUnreadChatCount,
		*client.UpdateStoryStealthMode,
		*client.UpdateChatBlockList,
		*client.UpdateAccentColors,
		*client.UpdateProfileAccentColors,
		*client.UpdateSavedMessagesTags,
		*client.UpdateOwnedStarCount,
		*client.UpdateChatPendingJoinRequests:
	
	}

	return nil
}

func extractSenderUserId(sender client.MessageSender) int64 {
	switch s := sender.(type) {
	case *client.MessageSenderUser:
		return s.UserId
	default:
		return 0
	}
}

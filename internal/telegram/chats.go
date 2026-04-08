package telegram

import (
	"context"

	"github.com/zelenin/go-tdlib/client"
)

// LoadChats fetches a batch of chats from the chat list.
func (c *Client) LoadChats(ctx context.Context, chatList client.ChatList, limit int32) error {
	_, err := c.tdClient.LoadChats(ctx, &client.LoadChatsRequest{
		ChatList: chatList,
		Limit:    limit,
	})
	return err
}

// GetChat returns a chat by its ID.
func (c *Client) GetChat(ctx context.Context, chatID int64) (*client.Chat, error) {
	return c.tdClient.GetChat(ctx, &client.GetChatRequest{
		ChatID: chatID,
	})
}

// GetChats returns an ordered list of chat IDs.
func (c *Client) GetChats(ctx context.Context, chatList client.ChatList, limit int32) (*client.Chats, error) {
	return c.tdClient.GetChats(ctx, &client.GetChatsRequest{
		ChatList: chatList,
		Limit:    limit,
	})
}

// GetChatHistory returns messages from a chat.
func (c *Client) GetChatHistory(ctx context.Context, chatID int64, fromMessageID int64, offset int32, limit int32) (*client.Messages, error) {
	return c.tdClient.GetChatHistory(ctx, &client.GetChatHistoryRequest{
		ChatID:        chatID,
		FromMessageID: fromMessageID,
		Offset:        offset,
		Limit:         limit,
	})
}

// SearchChats searches for chats by title.
func (c *Client) SearchChats(ctx context.Context, query string, limit int32) (*client.Chats, error) {
	return c.tdClient.SearchChats(ctx, &client.SearchChatsRequest{
		Query: query,
		Limit: limit,
	})
}

// SearchMessages searches for messages across all chats.
func (c *Client) SearchMessages(ctx context.Context, chatID int64, query string, fromMessageID int64, limit int32) (*client.FoundMessages, error) {
	return c.tdClient.SearchMessages(ctx, &client.SearchMessagesRequest{
		Query:  query,
		Limit:  limit,
		Filter: nil,
	})
}

// SearchChatMessages searches messages within a specific chat.
func (c *Client) SearchChatMessages(ctx context.Context, chatID int64, query string, fromMessageID int64, limit int32) (*client.FoundChatMessages, error) {
	return c.tdClient.SearchChatMessages(ctx, &client.SearchChatMessagesRequest{
		ChatID:        chatID,
		Query:         query,
		FromMessageID: fromMessageID,
		Limit:         limit,
	})
}

// GetChatMember returns info about a chat member.
func (c *Client) GetChatMember(ctx context.Context, chatID int64, memberID client.MessageSender) (*client.ChatMember, error) {
	return c.tdClient.GetChatMember(ctx, &client.GetChatMemberRequest{
		ChatID:   chatID,
		MemberID: memberID,
	})
}

// OpenChat marks a chat as opened (updates read state).
func (c *Client) OpenChat(ctx context.Context, chatID int64) error {
	_, err := c.tdClient.OpenChat(ctx, &client.OpenChatRequest{
		ChatID: chatID,
	})
	return err
}

// CloseChat marks a chat as closed.
func (c *Client) CloseChat(ctx context.Context, chatID int64) error {
	_, err := c.tdClient.CloseChat(ctx, &client.CloseChatRequest{
		ChatID: chatID,
	})
	return err
}

// JoinChat joins a public chat or channel.
func (c *Client) JoinChat(ctx context.Context, chatID int64) error {
	_, err := c.tdClient.JoinChat(ctx, &client.JoinChatRequest{
		ChatID: chatID,
	})
	return err
}

// LeaveChat leaves a chat.
func (c *Client) LeaveChat(ctx context.Context, chatID int64) error {
	_, err := c.tdClient.LeaveChat(ctx, &client.LeaveChatRequest{
		ChatID: chatID,
	})
	return err
}

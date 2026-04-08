package telegram

import (
	"github.com/zelenin/go-tdlib/client"
)

func (c *Client) LoadChats(chatList client.ChatList, limit int32) error {
	_, err := c.tdClient.LoadChats(&client.LoadChatsRequest{
		ChatList: chatList,
		Limit:    limit,
	})
	return err
}

func (c *Client) GetChat(chatId int64) (*client.Chat, error) {
	return c.tdClient.GetChat(&client.GetChatRequest{
		ChatId: chatId,
	})
}

func (c *Client) GetChats(chatList client.ChatList, limit int32) (*client.Chats, error) {
	return c.tdClient.GetChats(&client.GetChatsRequest{
		ChatList: chatList,
		Limit:    limit,
	})
}

func (c *Client) GetChatHistory(chatId int64, fromMessageId int64, offset int32, limit int32) (*client.Messages, error) {
	return c.tdClient.GetChatHistory(&client.GetChatHistoryRequest{
		ChatId:        chatId,
		FromMessageId: fromMessageId,
		Offset:        offset,
		Limit:         limit,
	})
}

func (c *Client) SearchChats(query string, limit int32) (*client.Chats, error) {
	return c.tdClient.SearchChats(&client.SearchChatsRequest{
		Query: query,
		Limit: limit,
	})
}

func (c *Client) SearchMessages(query string, limit int32) (*client.FoundMessages, error) {
	return c.tdClient.SearchMessages(&client.SearchMessagesRequest{
		Query: query,
		Limit: limit,
	})
}

func (c *Client) SearchChatMessages(chatId int64, query string, fromMessageId int64, limit int32) (*client.FoundChatMessages, error) {
	return c.tdClient.SearchChatMessages(&client.SearchChatMessagesRequest{
		ChatId:        chatId,
		Query:         query,
		FromMessageId: fromMessageId,
		Limit:         limit,
	})
}

func (c *Client) GetChatMember(chatId int64, memberId client.MessageSender) (*client.ChatMember, error) {
	return c.tdClient.GetChatMember(&client.GetChatMemberRequest{
		ChatId:   chatId,
		MemberId: memberId,
	})
}

func (c *Client) OpenChat(chatId int64) error {
	_, err := c.tdClient.OpenChat(&client.OpenChatRequest{
		ChatId: chatId,
	})
	return err
}

func (c *Client) CloseChat(chatId int64) error {
	_, err := c.tdClient.CloseChat(&client.CloseChatRequest{
		ChatId: chatId,
	})
	return err
}

func (c *Client) JoinChat(chatId int64) error {
	_, err := c.tdClient.JoinChat(&client.JoinChatRequest{
		ChatId: chatId,
	})
	return err
}

func (c *Client) LeaveChat(chatId int64) error {
	_, err := c.tdClient.LeaveChat(&client.LeaveChatRequest{
		ChatId: chatId,
	})
	return err
}

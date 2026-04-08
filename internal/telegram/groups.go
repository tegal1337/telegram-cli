package telegram

import (
	"context"

	"github.com/zelenin/go-tdlib/client"
)

// GetSupergroupInfo returns supergroup/channel metadata.
func (c *Client) GetSupergroupInfo(ctx context.Context, supergroupID int64) (*client.Supergroup, error) {
	return c.tdClient.GetSupergroup(ctx, &client.GetSupergroupRequest{
		SupergroupID: supergroupID,
	})
}

// GetSupergroupFullInfo returns full supergroup/channel info.
func (c *Client) GetSupergroupFullInfo(ctx context.Context, supergroupID int64) (*client.SupergroupFullInfo, error) {
	return c.tdClient.GetSupergroupFullInfo(ctx, &client.GetSupergroupFullInfoRequest{
		SupergroupID: supergroupID,
	})
}

// GetSupergroupMembers returns members of a supergroup/channel.
func (c *Client) GetSupergroupMembers(ctx context.Context, supergroupID int64, offset int32, limit int32) (*client.ChatMembers, error) {
	return c.tdClient.GetSupergroupMembers(ctx, &client.GetSupergroupMembersRequest{
		SupergroupID: supergroupID,
		Offset:       offset,
		Limit:        limit,
	})
}

// GetBasicGroupInfo returns basic group metadata.
func (c *Client) GetBasicGroupInfo(ctx context.Context, groupID int64) (*client.BasicGroup, error) {
	return c.tdClient.GetBasicGroup(ctx, &client.GetBasicGroupRequest{
		BasicGroupID: groupID,
	})
}

// GetBasicGroupFullInfo returns full basic group info.
func (c *Client) GetBasicGroupFullInfo(ctx context.Context, groupID int64) (*client.BasicGroupFullInfo, error) {
	return c.tdClient.GetBasicGroupFullInfo(ctx, &client.GetBasicGroupFullInfoRequest{
		BasicGroupID: groupID,
	})
}

// SetChatTitle sets the title of a chat.
func (c *Client) SetChatTitle(ctx context.Context, chatID int64, title string) error {
	_, err := c.tdClient.SetChatTitle(ctx, &client.SetChatTitleRequest{
		ChatID: chatID,
		Title:  title,
	})
	return err
}

// SetChatDescription sets the description of a chat.
func (c *Client) SetChatDescription(ctx context.Context, chatID int64, description string) error {
	_, err := c.tdClient.SetChatDescription(ctx, &client.SetChatDescriptionRequest{
		ChatID:      chatID,
		Description: description,
	})
	return err
}

// BanChatMember bans a member from a chat.
func (c *Client) BanChatMember(ctx context.Context, chatID int64, memberID client.MessageSender) error {
	_, err := c.tdClient.BanChatMember(ctx, &client.BanChatMemberRequest{
		ChatID:   chatID,
		MemberID: memberID,
	})
	return err
}

// CreatePrivateChat creates a private chat with a user.
func (c *Client) CreatePrivateChat(ctx context.Context, userID int64) (*client.Chat, error) {
	return c.tdClient.CreatePrivateChat(ctx, &client.CreatePrivateChatRequest{
		UserID: userID,
	})
}

// CreateBasicGroupChat creates a new basic group.
func (c *Client) CreateBasicGroupChat(ctx context.Context, title string, userIDs []int64) (*client.Chat, error) {
	return c.tdClient.CreateNewBasicGroupChat(ctx, &client.CreateNewBasicGroupChatRequest{
		Title:   title,
		UserIDs: userIDs,
	})
}

// CreateSupergroupChat creates a new supergroup or channel.
func (c *Client) CreateSupergroupChat(ctx context.Context, title string, isChannel bool, description string) (*client.Chat, error) {
	return c.tdClient.CreateNewSupergroupChat(ctx, &client.CreateNewSupergroupChatRequest{
		Title:       title,
		IsChannel:   isChannel,
		Description: description,
	})
}

// PinMessage pins a message in a chat.
func (c *Client) PinMessage(ctx context.Context, chatID int64, messageID int64) error {
	_, err := c.tdClient.PinChatMessage(ctx, &client.PinChatMessageRequest{
		ChatID:    chatID,
		MessageID: messageID,
	})
	return err
}

// UnpinMessage unpins a message in a chat.
func (c *Client) UnpinMessage(ctx context.Context, chatID int64, messageID int64) error {
	_, err := c.tdClient.UnpinChatMessage(ctx, &client.UnpinChatMessageRequest{
		ChatID:    chatID,
		MessageID: messageID,
	})
	return err
}

package telegram

import (
	"github.com/zelenin/go-tdlib/client"
)

func (c *Client) GetSupergroupInfo(supergroupId int64) (*client.Supergroup, error) {
	return c.tdClient.GetSupergroup(&client.GetSupergroupRequest{
		SupergroupId: supergroupId,
	})
}

func (c *Client) GetSupergroupFullInfo(supergroupId int64) (*client.SupergroupFullInfo, error) {
	return c.tdClient.GetSupergroupFullInfo(&client.GetSupergroupFullInfoRequest{
		SupergroupId: supergroupId,
	})
}

func (c *Client) GetSupergroupMembers(supergroupId int64, offset int32, limit int32) (*client.ChatMembers, error) {
	return c.tdClient.GetSupergroupMembers(&client.GetSupergroupMembersRequest{
		SupergroupId: supergroupId,
		Offset:       offset,
		Limit:        limit,
	})
}

func (c *Client) GetBasicGroupInfo(groupId int64) (*client.BasicGroup, error) {
	return c.tdClient.GetBasicGroup(&client.GetBasicGroupRequest{
		BasicGroupId: groupId,
	})
}

func (c *Client) GetBasicGroupFullInfo(groupId int64) (*client.BasicGroupFullInfo, error) {
	return c.tdClient.GetBasicGroupFullInfo(&client.GetBasicGroupFullInfoRequest{
		BasicGroupId: groupId,
	})
}

func (c *Client) SetChatTitle(chatId int64, title string) error {
	_, err := c.tdClient.SetChatTitle(&client.SetChatTitleRequest{
		ChatId: chatId,
		Title:  title,
	})
	return err
}

func (c *Client) SetChatDescription(chatId int64, description string) error {
	_, err := c.tdClient.SetChatDescription(&client.SetChatDescriptionRequest{
		ChatId:      chatId,
		Description: description,
	})
	return err
}

func (c *Client) BanChatMember(chatId int64, memberId client.MessageSender) error {
	_, err := c.tdClient.BanChatMember(&client.BanChatMemberRequest{
		ChatId:   chatId,
		MemberId: memberId,
	})
	return err
}

func (c *Client) CreatePrivateChat(userId int64) (*client.Chat, error) {
	return c.tdClient.CreatePrivateChat(&client.CreatePrivateChatRequest{
		UserId: userId,
	})
}

func (c *Client) CreateBasicGroupChat(title string, userIds []int64) (*client.CreatedBasicGroupChat, error) {
	return c.tdClient.CreateNewBasicGroupChat(&client.CreateNewBasicGroupChatRequest{
		Title:   title,
		UserIds: userIds,
	})
}

func (c *Client) CreateSupergroupChat(title string, isChannel bool, description string) (*client.Chat, error) {
	return c.tdClient.CreateNewSupergroupChat(&client.CreateNewSupergroupChatRequest{
		Title:       title,
		IsChannel:   isChannel,
		Description: description,
	})
}

func (c *Client) PinMessage(chatId int64, messageId int64) error {
	_, err := c.tdClient.PinChatMessage(&client.PinChatMessageRequest{
		ChatId:    chatId,
		MessageId: messageId,
	})
	return err
}

func (c *Client) UnpinMessage(chatId int64, messageId int64) error {
	_, err := c.tdClient.UnpinChatMessage(&client.UnpinChatMessageRequest{
		ChatId:    chatId,
		MessageId: messageId,
	})
	return err
}

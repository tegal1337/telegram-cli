package telegram

import (
	"github.com/zelenin/go-tdlib/client"
)

func (c *Client) GetUser(userId int64) (*client.User, error) {
	return c.tdClient.GetUser(&client.GetUserRequest{
		UserId: userId,
	})
}

func (c *Client) GetUserFullInfo(userId int64) (*client.UserFullInfo, error) {
	return c.tdClient.GetUserFullInfo(&client.GetUserFullInfoRequest{
		UserId: userId,
	})
}

func (c *Client) GetContacts() (*client.Users, error) {
	return c.tdClient.GetContacts()
}

func (c *Client) SearchContacts(query string, limit int32) (*client.Users, error) {
	return c.tdClient.SearchContacts(&client.SearchContactsRequest{
		Query: query,
		Limit: limit,
	})
}

func (c *Client) GetUserProfilePhotos(userId int64, offset int32, limit int32) (*client.ChatPhotos, error) {
	return c.tdClient.GetUserProfilePhotos(&client.GetUserProfilePhotosRequest{
		UserId: userId,
		Offset: offset,
		Limit:  limit,
	})
}

func (c *Client) BlockUser(senderId client.MessageSender) error {
	_, err := c.tdClient.SetMessageSenderBlockList(&client.SetMessageSenderBlockListRequest{
		SenderId:  senderId,
		BlockList: &client.BlockListMain{},
	})
	return err
}

func (c *Client) UnblockUser(senderId client.MessageSender) error {
	_, err := c.tdClient.SetMessageSenderBlockList(&client.SetMessageSenderBlockListRequest{
		SenderId:  senderId,
		BlockList: nil,
	})
	return err
}

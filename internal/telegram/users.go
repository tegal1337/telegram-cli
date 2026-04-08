package telegram

import (
	"context"

	"github.com/zelenin/go-tdlib/client"
)

// GetUser returns user info by ID.
func (c *Client) GetUser(ctx context.Context, userID int64) (*client.User, error) {
	return c.tdClient.GetUser(ctx, &client.GetUserRequest{
		UserID: userID,
	})
}

// GetUserFullInfo returns full user info.
func (c *Client) GetUserFullInfo(ctx context.Context, userID int64) (*client.UserFullInfo, error) {
	return c.tdClient.GetUserFullInfo(ctx, &client.GetUserFullInfoRequest{
		UserID: userID,
	})
}

// GetContacts returns the contact list.
func (c *Client) GetContacts(ctx context.Context) (*client.Users, error) {
	return c.tdClient.GetContacts(ctx)
}

// SearchContacts searches contacts by name.
func (c *Client) SearchContacts(ctx context.Context, query string, limit int32) (*client.Users, error) {
	return c.tdClient.SearchContacts(ctx, &client.SearchContactsRequest{
		Query: query,
		Limit: limit,
	})
}

// GetUserProfilePhotos returns a user's profile photos.
func (c *Client) GetUserProfilePhotos(ctx context.Context, userID int64, offset int32, limit int32) (*client.ChatPhotos, error) {
	return c.tdClient.GetUserProfilePhotos(ctx, &client.GetUserProfilePhotosRequest{
		UserID: userID,
		Offset: offset,
		Limit:  limit,
	})
}

// BlockUser blocks a user.
func (c *Client) BlockUser(ctx context.Context, senderID client.MessageSender) error {
	_, err := c.tdClient.ToggleMessageSenderIsBlocked(ctx, &client.ToggleMessageSenderIsBlockedRequest{
		SenderID:  senderID,
		IsBlocked: true,
	})
	return err
}

// UnblockUser unblocks a user.
func (c *Client) UnblockUser(ctx context.Context, senderID client.MessageSender) error {
	_, err := c.tdClient.ToggleMessageSenderIsBlocked(ctx, &client.ToggleMessageSenderIsBlockedRequest{
		SenderID:  senderID,
		IsBlocked: false,
	})
	return err
}

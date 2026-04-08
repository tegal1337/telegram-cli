package telegram

import (
	"context"

	"github.com/zelenin/go-tdlib/client"
)

// SendTextMessage sends a text message to a chat.
func (c *Client) SendTextMessage(ctx context.Context, chatID int64, text string, replyToMessageID int64) (*client.Message, error) {
	var replyTo client.InputMessageReplyTo
	if replyToMessageID != 0 {
		replyTo = &client.InputMessageReplyToMessage{
			MessageID: replyToMessageID,
		}
	}

	return c.tdClient.SendMessage(ctx, &client.SendMessageRequest{
		ChatID:  chatID,
		ReplyTo: replyTo,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: text,
			},
		},
	})
}

// EditTextMessage edits a text message.
func (c *Client) EditTextMessage(ctx context.Context, chatID int64, messageID int64, text string) (*client.Message, error) {
	return c.tdClient.EditMessageText(ctx, &client.EditMessageTextRequest{
		ChatID:    chatID,
		MessageID: messageID,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: text,
			},
		},
	})
}

// DeleteMessages deletes messages from a chat.
func (c *Client) DeleteMessages(ctx context.Context, chatID int64, messageIDs []int64, revoke bool) error {
	_, err := c.tdClient.DeleteMessages(ctx, &client.DeleteMessagesRequest{
		ChatID:     chatID,
		MessageIDs: messageIDs,
		Revoke:     revoke,
	})
	return err
}

// ForwardMessages forwards messages to another chat.
func (c *Client) ForwardMessages(ctx context.Context, chatID int64, fromChatID int64, messageIDs []int64) (*client.Messages, error) {
	return c.tdClient.ForwardMessages(ctx, &client.ForwardMessagesRequest{
		ChatID:     chatID,
		FromChatID: fromChatID,
		MessageIDs: messageIDs,
	})
}

// SendPhoto sends a photo message.
func (c *Client) SendPhoto(ctx context.Context, chatID int64, photoPath string, caption string) (*client.Message, error) {
	return c.tdClient.SendMessage(ctx, &client.SendMessageRequest{
		ChatID: chatID,
		InputMessageContent: &client.InputMessagePhoto{
			Photo: &client.InputFileLocal{
				Path: photoPath,
			},
			Caption: &client.FormattedText{
				Text: caption,
			},
		},
	})
}

// SendDocument sends a document/file message.
func (c *Client) SendDocument(ctx context.Context, chatID int64, filePath string, caption string) (*client.Message, error) {
	return c.tdClient.SendMessage(ctx, &client.SendMessageRequest{
		ChatID: chatID,
		InputMessageContent: &client.InputMessageDocument{
			Document: &client.InputFileLocal{
				Path: filePath,
			},
			Caption: &client.FormattedText{
				Text: caption,
			},
		},
	})
}

// SendVoiceNote sends a voice message.
func (c *Client) SendVoiceNote(ctx context.Context, chatID int64, voicePath string, duration int32) (*client.Message, error) {
	return c.tdClient.SendMessage(ctx, &client.SendMessageRequest{
		ChatID: chatID,
		InputMessageContent: &client.InputMessageVoiceNote{
			VoiceNote: &client.InputFileLocal{
				Path: voicePath,
			},
			Duration: duration,
		},
	})
}

// SendVideoNote sends a video note (round video).
func (c *Client) SendVideoNote(ctx context.Context, chatID int64, videoPath string, duration int32) (*client.Message, error) {
	return c.tdClient.SendMessage(ctx, &client.SendMessageRequest{
		ChatID: chatID,
		InputMessageContent: &client.InputMessageVideoNote{
			VideoNote: &client.InputFileLocal{
				Path: videoPath,
			},
			Duration: duration,
		},
	})
}

// SendVideo sends a video message.
func (c *Client) SendVideo(ctx context.Context, chatID int64, videoPath string, caption string) (*client.Message, error) {
	return c.tdClient.SendMessage(ctx, &client.SendMessageRequest{
		ChatID: chatID,
		InputMessageContent: &client.InputMessageVideo{
			Video: &client.InputFileLocal{
				Path: videoPath,
			},
			Caption: &client.FormattedText{
				Text: caption,
			},
		},
	})
}

// GetMessage retrieves a single message by ID.
func (c *Client) GetMessage(ctx context.Context, chatID int64, messageID int64) (*client.Message, error) {
	return c.tdClient.GetMessage(ctx, &client.GetMessageRequest{
		ChatID:    chatID,
		MessageID: messageID,
	})
}

// ViewMessages marks messages as read.
func (c *Client) ViewMessages(ctx context.Context, chatID int64, messageIDs []int64) error {
	_, err := c.tdClient.ViewMessages(ctx, &client.ViewMessagesRequest{
		ChatID:     chatID,
		MessageIDs: messageIDs,
	})
	return err
}

// SendSticker sends a sticker message.
func (c *Client) SendSticker(ctx context.Context, chatID int64, stickerPath string) (*client.Message, error) {
	return c.tdClient.SendMessage(ctx, &client.SendMessageRequest{
		ChatID: chatID,
		InputMessageContent: &client.InputMessageSticker{
			Sticker: &client.InputFileLocal{
				Path: stickerPath,
			},
		},
	})
}

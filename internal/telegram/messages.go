package telegram

import (
	"github.com/zelenin/go-tdlib/client"
)

func (c *Client) SendTextMessage(chatId int64, text string, replyToMessageId int64) (*client.Message, error) {
	var replyTo client.InputMessageReplyTo
	if replyToMessageId != 0 {
		replyTo = &client.InputMessageReplyToMessage{
			MessageId: replyToMessageId,
		}
	}

	return c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId:  chatId,
		ReplyTo: replyTo,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: text,
			},
		},
	})
}

func (c *Client) EditTextMessage(chatId int64, messageId int64, text string) (*client.Message, error) {
	return c.tdClient.EditMessageText(&client.EditMessageTextRequest{
		ChatId:    chatId,
		MessageId: messageId,
		InputMessageContent: &client.InputMessageText{
			Text: &client.FormattedText{
				Text: text,
			},
		},
	})
}

func (c *Client) DeleteMessages(chatId int64, messageIds []int64, revoke bool) error {
	_, err := c.tdClient.DeleteMessages(&client.DeleteMessagesRequest{
		ChatId:     chatId,
		MessageIds: messageIds,
		Revoke:     revoke,
	})
	return err
}

func (c *Client) ForwardMessages(chatId int64, fromChatId int64, messageIds []int64) (*client.Messages, error) {
	return c.tdClient.ForwardMessages(&client.ForwardMessagesRequest{
		ChatId:     chatId,
		FromChatId: fromChatId,
		MessageIds: messageIds,
	})
}

func (c *Client) SendPhoto(chatId int64, photoPath string, caption string) (*client.Message, error) {
	return c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId: chatId,
		InputMessageContent: &client.InputMessagePhoto{
			Photo:   &client.InputFileLocal{Path: photoPath},
			Caption: &client.FormattedText{Text: caption},
		},
	})
}

func (c *Client) SendDocument(chatId int64, filePath string, caption string) (*client.Message, error) {
	return c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId: chatId,
		InputMessageContent: &client.InputMessageDocument{
			Document: &client.InputFileLocal{Path: filePath},
			Caption:  &client.FormattedText{Text: caption},
		},
	})
}

func (c *Client) SendVoiceNote(chatId int64, voicePath string, duration int32) (*client.Message, error) {
	return c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId: chatId,
		InputMessageContent: &client.InputMessageVoiceNote{
			VoiceNote: &client.InputFileLocal{Path: voicePath},
			Duration:  duration,
		},
	})
}

func (c *Client) SendVideoNote(chatId int64, videoPath string, duration int32) (*client.Message, error) {
	return c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId: chatId,
		InputMessageContent: &client.InputMessageVideoNote{
			VideoNote: &client.InputFileLocal{Path: videoPath},
			Duration:  duration,
		},
	})
}

func (c *Client) SendVideo(chatId int64, videoPath string, caption string) (*client.Message, error) {
	return c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId: chatId,
		InputMessageContent: &client.InputMessageVideo{
			Video:   &client.InputFileLocal{Path: videoPath},
			Caption: &client.FormattedText{Text: caption},
		},
	})
}

func (c *Client) GetMessage(chatId int64, messageId int64) (*client.Message, error) {
	return c.tdClient.GetMessage(&client.GetMessageRequest{
		ChatId:    chatId,
		MessageId: messageId,
	})
}

func (c *Client) ViewMessages(chatId int64, messageIds []int64) error {
	_, err := c.tdClient.ViewMessages(&client.ViewMessagesRequest{
		ChatId:     chatId,
		MessageIds: messageIds,
	})
	return err
}

func (c *Client) SendSticker(chatId int64, stickerPath string) (*client.Message, error) {
	return c.tdClient.SendMessage(&client.SendMessageRequest{
		ChatId: chatId,
		InputMessageContent: &client.InputMessageSticker{
			Sticker: &client.InputFileLocal{Path: stickerPath},
		},
	})
}

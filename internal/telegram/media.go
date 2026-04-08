package telegram

import (
	"github.com/zelenin/go-tdlib/client"
)

func (c *Client) DownloadPhoto(photo *client.Photo) (*client.File, error) {
	if photo == nil || len(photo.Sizes) == 0 {
		return nil, nil
	}
	best := photo.Sizes[len(photo.Sizes)-1]
	return c.DownloadFile(best.Photo.Id, 16)
}

func (c *Client) DownloadPhotoThumbnail(photo *client.Photo) (*client.File, error) {
	if photo == nil || len(photo.Sizes) == 0 {
		return nil, nil
	}
	smallest := photo.Sizes[0]
	return c.DownloadFile(smallest.Photo.Id, 8)
}

func (c *Client) DownloadVoiceNote(voice *client.VoiceNote) (*client.File, error) {
	if voice == nil {
		return nil, nil
	}
	return c.DownloadFile(voice.Voice.Id, 16)
}

func (c *Client) DownloadVideoNote(videoNote *client.VideoNote) (*client.File, error) {
	if videoNote == nil {
		return nil, nil
	}
	return c.DownloadFile(videoNote.Video.Id, 16)
}

func (c *Client) DownloadVideo(video *client.Video) (*client.File, error) {
	if video == nil {
		return nil, nil
	}
	return c.DownloadFile(video.Video.Id, 8)
}

func (c *Client) DownloadDocument(doc *client.Document) (*client.File, error) {
	if doc == nil {
		return nil, nil
	}
	return c.DownloadFile(doc.Document.Id, 8)
}

func (c *Client) DownloadSticker(sticker *client.Sticker) (*client.File, error) {
	if sticker == nil {
		return nil, nil
	}
	return c.DownloadFile(sticker.Sticker.Id, 16)
}

func (c *Client) DownloadAnimation(anim *client.Animation) (*client.File, error) {
	if anim == nil {
		return nil, nil
	}
	return c.DownloadFile(anim.Animation.Id, 8)
}

func (c *Client) DownloadAudio(audio *client.Audio) (*client.File, error) {
	if audio == nil {
		return nil, nil
	}
	return c.DownloadFile(audio.Audio.Id, 8)
}

func (c *Client) GetVideoThumbnail(video *client.Video) (*client.File, error) {
	if video == nil || video.Thumbnail == nil {
		return nil, nil
	}
	return c.DownloadFile(video.Thumbnail.File.Id, 8)
}

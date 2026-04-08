package telegram

import (
	"context"

	"github.com/zelenin/go-tdlib/client"
)

// DownloadPhoto downloads the best available photo size.
func (c *Client) DownloadPhoto(ctx context.Context, photo *client.Photo) (*client.File, error) {
	if photo == nil || len(photo.Sizes) == 0 {
		return nil, nil
	}

	// Pick the largest available size for display.
	best := photo.Sizes[len(photo.Sizes)-1]
	return c.DownloadFile(ctx, best.Photo.ID, 16)
}

// DownloadPhotoThumbnail downloads the smallest photo size (thumbnail).
func (c *Client) DownloadPhotoThumbnail(ctx context.Context, photo *client.Photo) (*client.File, error) {
	if photo == nil || len(photo.Sizes) == 0 {
		return nil, nil
	}

	smallest := photo.Sizes[0]
	return c.DownloadFile(ctx, smallest.Photo.ID, 8)
}

// DownloadVoiceNote downloads a voice note file.
func (c *Client) DownloadVoiceNote(ctx context.Context, voice *client.VoiceNote) (*client.File, error) {
	if voice == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, voice.Voice.ID, 16)
}

// DownloadVideoNote downloads a video note file.
func (c *Client) DownloadVideoNote(ctx context.Context, videoNote *client.VideoNote) (*client.File, error) {
	if videoNote == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, videoNote.Video.ID, 16)
}

// DownloadVideo downloads a video file.
func (c *Client) DownloadVideo(ctx context.Context, video *client.Video) (*client.File, error) {
	if video == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, video.Video.ID, 8)
}

// DownloadDocument downloads a document file.
func (c *Client) DownloadDocument(ctx context.Context, doc *client.Document) (*client.File, error) {
	if doc == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, doc.Document.ID, 8)
}

// DownloadSticker downloads a sticker file.
func (c *Client) DownloadSticker(ctx context.Context, sticker *client.Sticker) (*client.File, error) {
	if sticker == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, sticker.Sticker.ID, 16)
}

// DownloadAnimation downloads an animation/GIF file.
func (c *Client) DownloadAnimation(ctx context.Context, anim *client.Animation) (*client.File, error) {
	if anim == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, anim.Animation.ID, 8)
}

// DownloadAudio downloads an audio file.
func (c *Client) DownloadAudio(ctx context.Context, audio *client.Audio) (*client.File, error) {
	if audio == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, audio.Audio.ID, 8)
}

// GetVideoThumbnail downloads the thumbnail for a video.
func (c *Client) GetVideoThumbnail(ctx context.Context, video *client.Video) (*client.File, error) {
	if video == nil || video.Thumbnail == nil {
		return nil, nil
	}
	return c.DownloadFile(ctx, video.Thumbnail.File.ID, 8)
}

package telegram

import (
	"context"

	"github.com/zelenin/go-tdlib/client"
)

// DownloadFile starts downloading a file by its file ID.
func (c *Client) DownloadFile(ctx context.Context, fileID int32, priority int32) (*client.File, error) {
	return c.tdClient.DownloadFile(ctx, &client.DownloadFileRequest{
		FileID:      fileID,
		Priority:    priority,
		Synchronous: false,
	})
}

// DownloadFileSync downloads a file synchronously (blocks until complete).
func (c *Client) DownloadFileSync(ctx context.Context, fileID int32) (*client.File, error) {
	return c.tdClient.DownloadFile(ctx, &client.DownloadFileRequest{
		FileID:      fileID,
		Priority:    32,
		Synchronous: true,
	})
}

// CancelDownloadFile cancels an active file download.
func (c *Client) CancelDownloadFile(ctx context.Context, fileID int32) error {
	_, err := c.tdClient.CancelDownloadFile(ctx, &client.CancelDownloadFileRequest{
		FileID: fileID,
	})
	return err
}

// GetFile returns file info by its ID.
func (c *Client) GetFile(ctx context.Context, fileID int32) (*client.File, error) {
	return c.tdClient.GetFile(ctx, &client.GetFileRequest{
		FileID: fileID,
	})
}

// GetRemoteFile returns file info by its remote ID.
func (c *Client) GetRemoteFile(ctx context.Context, remoteFileID string) (*client.File, error) {
	return c.tdClient.GetRemoteFile(ctx, &client.GetRemoteFileRequest{
		RemoteFileID: remoteFileID,
	})
}

// ReadFilePart reads a chunk of a local file managed by TDLib.
func (c *Client) ReadFilePart(ctx context.Context, fileID int32, offset int64, count int64) (*client.FilePart, error) {
	return c.tdClient.ReadFilePart(ctx, &client.ReadFilePartRequest{
		FileID: fileID,
		Offset: offset,
		Count:  count,
	})
}

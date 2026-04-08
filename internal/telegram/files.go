package telegram

import (
	"github.com/zelenin/go-tdlib/client"
)

func (c *Client) DownloadFile(fileId int32, priority int32) (*client.File, error) {
	return c.tdClient.DownloadFile(&client.DownloadFileRequest{
		FileId:      fileId,
		Priority:    priority,
		Synchronous: false,
	})
}

func (c *Client) DownloadFileSync(fileId int32) (*client.File, error) {
	return c.tdClient.DownloadFile(&client.DownloadFileRequest{
		FileId:      fileId,
		Priority:    32,
		Synchronous: true,
	})
}

func (c *Client) CancelDownloadFile(fileId int32) error {
	_, err := c.tdClient.CancelDownloadFile(&client.CancelDownloadFileRequest{
		FileId: fileId,
	})
	return err
}

func (c *Client) GetFile(fileId int32) (*client.File, error) {
	return c.tdClient.GetFile(&client.GetFileRequest{
		FileId: fileId,
	})
}

func (c *Client) GetRemoteFile(remoteFileId string) (*client.File, error) {
	return c.tdClient.GetRemoteFile(&client.GetRemoteFileRequest{
		RemoteFileId: remoteFileId,
	})
}

func (c *Client) ReadFilePart(fileId int32, offset int64, count int64) (*client.FilePart, error) {
	return c.tdClient.ReadFilePart(&client.ReadFilePartRequest{
		FileId: fileId,
		Offset: offset,
		Count:  count,
	})
}

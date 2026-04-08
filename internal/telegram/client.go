package telegram

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/tegal1337/telegram-cli/internal/config"
	"github.com/zelenin/go-tdlib/client"
)

// Client wraps the TDLib client with application-specific functionality.
type Client struct {
	tdClient *client.Client
	config   *config.Config
}

// NewClient creates a new TDLib client with the given configuration.
// The authorizer handles the authentication flow (phone, QR, bot token).
func NewClient(cfg *config.Config, authorizer client.AuthorizationStateHandler) (*Client, error) {
	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Printf("SetLogVerbosityLevel error: %s", err)
	}

	os.MkdirAll(cfg.Storage.DatabaseDir, 0o755)
	os.MkdirAll(cfg.Storage.FilesDir, 0o755)

	tdClient, err := client.NewClient(authorizer)
	if err != nil {
		return nil, fmt.Errorf("creating TDLib client: %w", err)
	}

	return &Client{
		tdClient: tdClient,
		config:   cfg,
	}, nil
}

// TdlibParameters returns the TDLib parameters from the config.
func TdlibParameters(cfg *config.Config) *client.SetTdlibParametersRequest {
	return &client.SetTdlibParametersRequest{
		UseTestDc:           false,
		DatabaseDirectory:   cfg.Storage.DatabaseDir,
		FilesDirectory:      cfg.Storage.FilesDir,
		UseFileDatabase:     true,
		UseChatInfoDatabase: true,
		UseMessageDatabase:  true,
		UseSecretChats:      true,
		ApiId:               cfg.Telegram.APIID,
		ApiHash:             cfg.Telegram.APIHash,
		SystemLanguageCode:  "en",
		DeviceModel:         "Tele-TUI",
		SystemVersion:       "1.0.0",
		ApplicationVersion:  "0.1.0",
	}
}

// Close gracefully shuts down the TDLib client.
func (c *Client) Close() {
	if c.tdClient != nil {
		c.tdClient.Close(context.Background())
	}
}

// TD returns the underlying TDLib client for direct API access.
func (c *Client) TD() *client.Client {
	return c.tdClient
}

// GetMe returns the current authorized user.
func (c *Client) GetMe(ctx context.Context) (*client.User, error) {
	return c.tdClient.GetMe(ctx)
}

// DataDir returns the configured data directory base path.
func (c *Client) DataDir() string {
	return filepath.Dir(c.config.Storage.DatabaseDir)
}

package telegram

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/tegal1337/telegram-cli/internal/config"
	"github.com/zelenin/go-tdlib/client"
)

type Client struct {
	mu       sync.RWMutex
	tdClient *client.Client
	config   *config.Config
	ready    chan struct{}
}

// NewClientAsync starts TDLib client creation in the background.
// The client blocks on authorization — call this before starting the TUI
// so the auth UI can feed credentials via the authorizer channels.
func NewClientAsync(cfg *config.Config, authorizer client.AuthorizationStateHandler) *Client {
	_, err := client.SetLogVerbosityLevel(&client.SetLogVerbosityLevelRequest{
		NewVerbosityLevel: 1,
	})
	if err != nil {
		log.Printf("SetLogVerbosityLevel error: %s", err)
	}

	os.MkdirAll(cfg.Storage.DatabaseDir, 0o755)
	os.MkdirAll(cfg.Storage.FilesDir, 0o755)

	c := &Client{
		config: cfg,
		ready:  make(chan struct{}),
	}

	go func() {
		tdClient, err := client.NewClient(authorizer)
		if err != nil {
			log.Printf("NewClient error: %s", err)
			return
		}
		c.mu.Lock()
		c.tdClient = tdClient
		c.mu.Unlock()
		close(c.ready)
	}()

	return c
}

// WaitReady blocks until the client is authorized and ready.
func (c *Client) WaitReady() {
	<-c.ready
}

// IsReady returns true if the client is authorized.
func (c *Client) IsReady() bool {
	select {
	case <-c.ready:
		return true
	default:
		return false
	}
}

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

func (c *Client) Close() {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.tdClient != nil {
		c.tdClient.Close()
	}
}

func (c *Client) TD() *client.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.tdClient
}

func (c *Client) GetMe() (*client.User, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.tdClient == nil {
		return nil, fmt.Errorf("client not ready")
	}
	return c.tdClient.GetMe()
}

func (c *Client) DataDir() string {
	return filepath.Dir(c.config.Storage.DatabaseDir)
}

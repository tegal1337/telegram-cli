package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml/v2"
)

type Config struct {
	Telegram      TelegramConfig      `toml:"telegram"`
	Storage       StorageConfig       `toml:"storage"`
	UI            UIConfig            `toml:"ui"`
	Media         MediaConfig         `toml:"media"`
	Notifications NotificationConfig  `toml:"notifications"`
	Keys          KeyConfig           `toml:"keys"`
}

type TelegramConfig struct {
	APIID   int32  `toml:"api_id"`
	APIHash string `toml:"api_hash"`
	Phone   string `toml:"phone"`
}

type StorageConfig struct {
	DatabaseDir string `toml:"database_dir"`
	FilesDir    string `toml:"files_dir"`
}

type UIConfig struct {
	Theme           string `toml:"theme"`
	ChatListWidth   int    `toml:"chat_list_width"`
	ShowAvatars     bool   `toml:"show_avatars"`
	TimestampFormat string `toml:"timestamp_format"`
	DateFormat      string `toml:"date_format"`
}

type MediaConfig struct {
	ImageProtocol      string `toml:"image_protocol"`
	MaxImageWidth      int    `toml:"max_image_width"`
	MaxImageHeight     int    `toml:"max_image_height"`
	VoicePlayer        string `toml:"voice_player"`
	VideoPlayer        string `toml:"video_player"`
	AutoDownloadPhotos bool   `toml:"auto_download_photos"`
	AutoDownloadVoice  bool   `toml:"auto_download_voice"`
	AutoDownloadLimitMB int   `toml:"auto_download_limit_mb"`
}

type NotificationConfig struct {
	Enabled     bool `toml:"enabled"`
	Sound       bool `toml:"sound"`
	ShowPreview bool `toml:"show_preview"`
}

type KeyConfig struct {
	Quit          string `toml:"quit"`
	FocusChatList string `toml:"focus_chat_list"`
	FocusChatView string `toml:"focus_chat_view"`
	FocusComposer string `toml:"focus_composer"`
	Search        string `toml:"search"`
	Contacts      string `toml:"contacts"`
	NextChat      string `toml:"next_chat"`
	PrevChat      string `toml:"prev_chat"`
	Reply         string `toml:"reply"`
	EditMessage   string `toml:"edit_message"`
	DeleteMessage string `toml:"delete_message"`
	Forward       string `toml:"forward"`
	ScrollUp      string `toml:"scroll_up"`
	ScrollDown    string `toml:"scroll_down"`
	PageUp        string `toml:"page_up"`
	PageDown      string `toml:"page_down"`
}

func Load() (*Config, error) {
	cfg := defaultConfig()

	configPath := findConfigPath()
	if configPath == "" {
		return cfg, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("reading config: %w", err)
	}

	if err := toml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}

	cfg.Storage.DatabaseDir = expandPath(cfg.Storage.DatabaseDir)
	cfg.Storage.FilesDir = expandPath(cfg.Storage.FilesDir)

	return cfg, nil
}

func defaultConfig() *Config {
	return &Config{
		Storage: StorageConfig{
			DatabaseDir: expandPath("~/.local/share/tele-tui/database"),
			FilesDir:    expandPath("~/.local/share/tele-tui/files"),
		},
		UI: UIConfig{
			Theme:           "dark",
			ChatListWidth:   30,
			ShowAvatars:     true,
			TimestampFormat: "15:04",
			DateFormat:      "2006-01-02",
		},
		Media: MediaConfig{
			ImageProtocol:       "auto",
			MaxImageWidth:       40,
			MaxImageHeight:      20,
			VoicePlayer:         "mpv",
			VideoPlayer:         "mpv",
			AutoDownloadPhotos:  true,
			AutoDownloadVoice:   true,
			AutoDownloadLimitMB: 10,
		},
		Notifications: NotificationConfig{
			Enabled:     true,
			Sound:       false,
			ShowPreview: true,
		},
		Keys: KeyConfig{
			Quit:          "ctrl+c",
			FocusChatList: "ctrl+1",
			FocusChatView: "ctrl+2",
			FocusComposer: "ctrl+3",
			Search:        "/",
			Contacts:      "ctrl+k",
			NextChat:      "ctrl+j",
			PrevChat:      "ctrl+k",
			Reply:         "r",
			EditMessage:   "e",
			DeleteMessage: "d",
			Forward:       "f",
			ScrollUp:      "k",
			ScrollDown:    "j",
			PageUp:        "ctrl+u",
			PageDown:      "ctrl+d",
		},
	}
}

func findConfigPath() string {
	if p := os.Getenv("TELETUI_CONFIG"); p != "" {
		return p
	}

	xdgConfig := os.Getenv("XDG_CONFIG_HOME")
	if xdgConfig == "" {
		home, _ := os.UserHomeDir()
		xdgConfig = filepath.Join(home, ".config")
	}

	path := filepath.Join(xdgConfig, "tele-tui", "config.toml")
	if _, err := os.Stat(path); err == nil {
		return path
	}

	return ""
}

func expandPath(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}

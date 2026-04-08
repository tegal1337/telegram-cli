# Telegram CLI

A full-featured Telegram client for the terminal, built with Go.

```
 _____ ___ _    ___   _____ _   _ ___
|_   _| __| |  | __| |_   _| | | |_ _|
  | | | _|| |__| _|    | | | |_| || |
  |_| |___|____|___|   |_|  \___/|___|
```

## Features

- **Chat Management** — Private chats, groups, supergroups, channels, secret chats
- **Rich Message Display** — Message bubbles with sender colors, timestamps, read status
- **Document Rendering** — Markdown/code rendering via [Glamour](https://github.com/charmbracelet/glamour)
- **Image Rendering** — Kitty graphics protocol, Sixel, Unicode half-block fallback
- **Voice Messages** — Inline playback via mpv/ffplay/paplay
- **Video** — External player launch (mpv/vlc/xdg-open)
- **File Transfer** — Download/upload with progress bars
- **Search** — Search chats, messages, and global Telegram directory
- **Contacts** — Contact list with online status indicators
- **Group Management** — Member list, admin actions, info panel
- **Authentication** — Phone/SMS code, 2FA password, QR code login
- **Notifications** — Desktop notifications via libnotify/osascript
- **Vim-style Navigation** — `j/k` scroll, `g/G` jump, `/` search
- **Responsive Layout** — Dual-panel (wide) or single-panel (narrow) mode
- **Theming** — Tokyo Night dark/light themes

## Architecture

```
┌─────────────────────────────────────────────────────────┐
│                    Bubbletea v2                          │
│  ┌──────────┐  ┌──────────────┐  ┌───────────────────┐  │
│  │ Chat List │  │  Chat View   │  │    Composer        │  │
│  │          │  │  (messages)  │  │  (input + reply)  │  │
│  └──────────┘  └──────────────┘  └───────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐│
│  │                  Status Bar                          ││
│  └──────────────────────────────────────────────────────┘│
├─────────────────────────────────────────────────────────┤
│                  Store (in-memory cache)                 │
│      Chats · Messages · Users · Files                   │
├─────────────────────────────────────────────────────────┤
│              TDLib via go-tdlib                          │
│   Listener → tea.Msg bridge (async → bubbletea loop)    │
└─────────────────────────────────────────────────────────┘
```

## Prerequisites

- **Go 1.23+**
- **TDLib** (libtdjson) — [installation guide](https://github.com/zelenin/go-tdlib#installation)
- **Telegram API credentials** — from [my.telegram.org](https://my.telegram.org/apps)

### Install TDLib (Ubuntu/Debian)

```bash
sudo apt install -y build-essential cmake gperf zlib1g-dev libssl-dev
git clone --depth 1 https://github.com/tdlib/td.git ~/td/td-src
cd ~/td/td-src && mkdir build && cd build
cmake -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX=~/td/tdlib ..
cmake --build . -j$(nproc) && cmake --install .
```

## Setup

```bash
# 1. Clone
git clone https://github.com/tegal1337/telegram-cli.git
cd telegram-cli

# 2. Configure
mkdir -p ~/.config/tele-tui
cp config.example.toml ~/.config/tele-tui/config.toml
# Edit config.toml with your api_id and api_hash

# 3. Build
make build

# 4. Run
make run
```

## Keybindings

| Key | Action |
|-----|--------|
| `Ctrl+1` | Focus chat list |
| `Ctrl+2` | Focus chat view |
| `Ctrl+3` | Focus composer |
| `j/k` | Navigate up/down |
| `Enter` | Select / Send message |
| `/` | Search |
| `Ctrl+K` | Contacts |
| `r` | Reply to message |
| `e` | Edit own message |
| `d` | Delete message |
| `f` | Forward message |
| `Ctrl+U/D` | Page up/down |
| `Ctrl+C` | Quit |

## Configuration

See [`config.example.toml`](config.example.toml) for all options:

- Telegram API credentials
- UI theme (dark/light)
- Media player preferences
- Image rendering protocol (auto/kitty/sixel/blocks)
- Notification settings
- Custom keybindings

## Project Structure

```
cmd/teletui/          Entry point
internal/
  app/                Root bubbletea model
  config/             TOML config loader
  telegram/           TDLib client wrapper + update listener
  ui/
    theme/            Lipgloss styles (Tokyo Night)
    layout/           Responsive panel layout
    widgets/          Reusable widgets (list, tabs, spinner, etc.)
    components/       UI components (chatlist, chatview, composer, auth, ...)
  media/              Image rendering (kitty/sixel/blocks), voice, video
  render/             Message rendering (markdown, entities, timestamps)
  notification/       Desktop notifications + sounds
  store/              Thread-safe in-memory caches
pkg/utils/            String/time/sanitize utilities
```

## License

MIT

<p align="center">
  <img src="https://upload.wikimedia.org/wikipedia/commons/8/82/Telegram_logo.svg" width="80" />
</p>

<h1 align="center">Telegram CLI</h1>

<p align="center">
  <strong>A full-featured Telegram client for the terminal</strong>
</p>

<p align="center">
  <a href="https://github.com/tegal1337/telegram-cli/actions"><img src="https://github.com/tegal1337/telegram-cli/actions/workflows/build.yml/badge.svg" alt="Build"></a>
  <a href="https://github.com/tegal1337/telegram-cli/releases"><img src="https://img.shields.io/github/v/release/tegal1337/telegram-cli?include_prereleases" alt="Release"></a>
  <a href="https://github.com/tegal1337/telegram-cli/blob/main/LICENSE"><img src="https://img.shields.io/github/license/tegal1337/telegram-cli" alt="License"></a>
  <img src="https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white" alt="Go">
</p>

---

## Features

- **Chat Management** — Private chats, groups, supergroups, channels, secret chats
- **Message Bubbles** — Rounded bordered bubbles, own messages right-aligned, read status indicators
- **Profile Avatars** — Colored initials or rendered profile photos in chat list
- **Markdown Rendering** — Code blocks, bold, italic, links via [Glamour](https://github.com/charmbracelet/glamour)
- **Image Rendering** — Kitty graphics protocol, Sixel, Unicode half-block fallback with CatmullRom scaling
- **Voice/Audio Playback** — Play voice messages and audio inline via `mpv` / `ffplay`
- **Video** — Open videos in external player (`mpv` / `vlc` / `xdg-open`)
- **File Transfer** — Download with `s`, open with `Enter`, progress bar during sync
- **Search** — Search chats, messages, and global Telegram directory
- **Contacts** — Contact list with online status indicators
- **Group Info** — Member list, admin roles, group description
- **Authentication** — Phone/SMS code, 2FA password, QR code login
- **First-Run Wizard** — Prompts for API credentials and saves config automatically
- **Notifications** — Desktop notifications via `notify-send` / `osascript`
- **Responsive Layout** — Dual-panel (wide) or single-panel (narrow terminals)
- **Theming** — Dark and light themes with 256-color support

## Screenshot

```
╭─ Chat List ─────────────╮╭─ Messages ──────────────────────────────────╮
│ DA  Dadang Jordan  08:15 ││                                             │
│     tes lim              ││                      ╭─────────────────────╮ │
│ SK  SKY API        13:24 ││                      │ naon we             │ │
│     sudah aman        2  ││                      │ 15:20 ✓✓            │ │
│ TG  Telegram       08:03 ││                      ╰─────────────────────╯ │
│     Login code: 90969... ││ ╭──────────────────╮                        │
│ AP  Api MX         14:38 ││ │ Dadang Jordan    │                        │
│     okesiap koo      81  ││ │ tah              │                        │
│                          ││ │ 15:22            │                        │
│                          ││ ╰──────────────────╯                        │
╰──────────────────────────╯╰─────────────────────────────────────────────╯
╭─ Compose ───────────────────────────────────────────────────────────────╮
│ █                                                                       │
│ Enter: send | Esc: cancel                                               │
╰─────────────────────────────────────────────────────────────────────────╯
● Connected  IMTAQIN    Tab:switch │ Esc:back │ /:search │ Alt+C:contacts
```

## Quick Start

```bash
# Clone
git clone https://github.com/tegal1337/telegram-cli.git
cd telegram-cli

# Auto-install TDLib + dependencies (one command)
make setup

# Build & run — first run prompts for API credentials
make run
```

`make setup` automatically handles everything:
- **Linux**: installs build deps via apt/dnf/pacman, builds TDLib, registers library path
- **macOS**: installs via `brew install tdlib`
- **Windows**: use MSYS2 (see below)

### Prerequisites

- **Go 1.23+**
- **mpv** (optional) — for voice/audio/video playback (`sudo apt install mpv`)
- **Telegram API credentials** — from [my.telegram.org/apps](https://my.telegram.org/apps)

### Windows (MSYS2)

```bash
pacman -S mingw-w64-x86_64-gcc mingw-w64-x86_64-cmake mingw-w64-x86_64-gperf
# Then build TDLib: https://github.com/tdlib/td#building
make build
```

On first run, you'll be prompted:

```
╔══════════════════════════════════════════╗
║         Telegram CLI - First Run         ║
╚══════════════════════════════════════════╝

Get your API credentials from:
https://my.telegram.org/apps

Enter API ID: xxxxxxx
Enter API Hash: xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
Enter phone number (optional): +628xxxxxxxxxx

Config saved! Starting Telegram CLI...
```

## Keybindings

### Navigation

| Key | Action |
|-----|--------|
| `Tab` / `Shift+Tab` | Cycle between panels |
| `Esc` | Go back / close overlay |
| `F1` / `Alt+1` | Focus chat list |
| `F2` / `Alt+2` | Focus messages |
| `F3` / `Alt+3` | Focus composer |
| `i` | Start composing (from chat view) |
| `j` / `k` | Scroll up/down |
| `g` / `G` | Jump to top/bottom |
| `PgUp` / `PgDn` | Page scroll |

### Actions

| Key | Action |
|-----|--------|
| `Enter` | Select chat / Send message / Play media |
| `o` | Open/play media |
| `s` | Save/download file |
| `/` | Search |
| `Alt+C` | Toggle contacts |
| `r` | Reply to message |
| `e` | Edit own message |
| `d` | Delete message |
| `Ctrl+Q` / `Ctrl+C` | Quit |

### Composer

| Key | Action |
|-----|--------|
| `Enter` | Send message |
| `Esc` | Cancel reply/edit, or leave composer |
| `Ctrl+W` | Delete word |
| `Ctrl+U` | Clear line before cursor |
| `Ctrl+K` | Clear line after cursor |

## Configuration

Config is stored at `~/.config/tele-tui/config.toml`. See [`config.example.toml`](config.example.toml) for all options:

```toml
[telegram]
api_id = 12345678
api_hash = "your_api_hash"

[ui]
theme = "dark"           # "dark" or "light"

[media]
image_protocol = "auto"  # "auto", "kitty", "sixel", "blocks"
voice_player = "mpv"     # "mpv", "ffplay"
video_player = "mpv"     # "mpv", "vlc", "xdg-open"
```

## Architecture

```
┌──────────────────────────────────────────────────────┐
│                   Bubbletea v2                        │
│  ╭────────╮  ╭──────────────╮  ╭──────────────────╮  │
│  │  Chat  │  │   Messages   │  │    Composer       │  │
│  │  List  │  │   (bubbles)  │  │  (text input)    │  │
│  ╰────────╯  ╰──────────────╯  ╰──────────────────╯  │
│  ╭──────────────────────────────────────────────────╮ │
│  │              Status Bar + Help                   │ │
│  ╰──────────────────────────────────────────────────╯ │
├──────────────────────────────────────────────────────┤
│              Store (thread-safe cache)                │
│         Chats · Messages · Users · Files              │
├──────────────────────────────────────────────────────┤
│           TDLib via go-tdlib (async bridge)            │
│      Listener goroutine → p.Send(tea.Msg)             │
└──────────────────────────────────────────────────────┘
```

## Project Structure

```
cmd/teletui/              Entry point + first-run wizard
internal/
  app/                    Root bubbletea model, key routing, layout
  config/                 TOML config loader + auto-save
  telegram/               TDLib client wrapper (async)
    auth.go               Phone/code/2FA/QR auth flow
    listener.go           TDLib update → tea.Msg bridge
    chats.go              Chat list, history, search
    messages.go           Send/edit/delete/forward
    media.go              Photo/voice/video download
  ui/
    theme/                256-color dark/light themes
    layout/               Responsive panel sizing
    widgets/              List, textarea, spinner, tabs, progress bar
    components/
      chatlist/           Chat list with avatars + unread badges
      chatview/           Message bubbles + media playback
      composer/           Text input with reply/edit modes
      auth/               Auth flow screens
      search/             Tabbed search overlay
      contacts/           Contact list
      groupinfo/          Group/channel info panel
      statusbar/          Connection status + typing indicators
      dialog/             Modal dialogs
  media/                  Image rendering (kitty/sixel/blocks)
  render/                 Message content → terminal output
  notification/           Desktop notifications
  store/                  Thread-safe in-memory caches
pkg/utils/                String/time/sanitize utilities
```

## Building from Source

```bash
make setup    # auto-install TDLib (Linux/macOS)
make build    # compile binary → bin/tele-tui
make run      # build + run (auto-detects TDLib path)
make clean    # remove build artifacts
```

The Makefile auto-detects TDLib in: `~/td/tdlib`, `/usr/local`, `/usr`, `/opt/homebrew/opt/tdlib`.

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/awesome`)
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

## License

MIT License - see [LICENSE](LICENSE) for details.

## Credits

- [Bubbletea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [Glamour](https://github.com/charmbracelet/glamour) — Markdown rendering
- [go-tdlib](https://github.com/zelenin/go-tdlib) — TDLib Go bindings

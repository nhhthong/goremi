<div align="center">
  <img src="resources/logo.png" alt="GoReMi logo" width="360" />

  <p><strong>Terminal music kit. Search, select, listen.</strong></p>

  <p>
    <a href="https://github.com/nhhthong/goremi"><img src="https://img.shields.io/github/go-mod/go-version/nhhthong/goremi?style=flat-square&logo=go&logoColor=white&color=00ADD8" alt="Go version"></a>
    <a href="https://github.com/nhhthong/goremi/commits"><img src="https://img.shields.io/github/last-commit/nhhthong/goremi?style=flat-square" alt="Last commit"></a>
    <img src="https://img.shields.io/badge/status-early%20development-orange?style=flat-square" alt="Status: early development">
    <img src="https://img.shields.io/badge/platform-Linux%20%7C%20macOS%20%7C%20Windows-informational?style=flat-square" alt="Platform: Linux, macOS, Windows">
    <a href="https://github.com/charmbracelet/bubbletea"><img src="https://img.shields.io/badge/TUI-Bubble%20Tea-FF5F87?style=flat-square" alt="TUI: Bubble Tea"></a>
    <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square" alt="License: MIT">
  </p>
</div>

---

## Overview

GoReMi (**Go** + **Do Re Mi**) is a minimal music player that runs in your terminal. Type a song, pick a result, and it plays, with album art and a live spectrum, all from the keyboard.

> **Status:** early development. Search, the interface, playback, artwork, the spectrum, themes and the install scripts are built. No release has been published yet, so the install scripts only work after the first release.

## Features

- **Search** YouTube from the terminal, 10 results at a time with `load more...`
- **Playback** through [mpv](https://mpv.io/), with play, pause, seek, previous and next track
- **Realtime spectrum analyzer** drawn with Unicode blocks
- **Album artwork** in the terminal (Kitty Graphics Protocol or Sixel), hidden when unsupported; shown only when the spectrum is off
- **Seven themes:** `default`, `catppuccin`, `dracula`, `gruvbox`, `nord`, `rosepine` and `tokyonight`, chosen with `goremi theme`
- **Keyboard-first**; clicking the player controls is optional
- **Linux, macOS and Windows**

## Install

GoReMi needs [mpv](https://mpv.io/) and [yt-dlp](https://github.com/yt-dlp/yt-dlp) on your `PATH`. The install scripts only report a missing one; they do not install it.

**Linux and macOS** (once a release exists):

```sh
curl -fsSL https://raw.githubusercontent.com/nhhthong/goremi/main/install.sh | sh
```

It installs `goremi` to `~/.local/bin` (no `sudo`). On macOS a missing dependency is reported with its `brew install` command. Make sure `~/.local/bin` is on your `PATH`.

**Windows** (amd64 only): download `install.ps1` from the repository and run it in PowerShell. It installs `goremi.exe` to `%LocalAppData%\Programs\goremi`; add that directory to your `PATH`.

**From source** (Go 1.26 or newer):

```sh
make build    # writes bin/goremi
make test     # go test, go vet, gofmt
```

## Usage

```text
goremi          open the player
goremi theme    choose a theme
goremi help     show this help
```

Type a query and press `Enter`. When the results list has focus:

| Key | Action |
|---|---|
| `↑` `↓` | select a track |
| `Enter` | play the selected track, or load the next page on `load more...` |
| `Space` or `k` | play or pause |
| `←` or `j` | seek back 10 seconds |
| `→` or `l` | seek forward 10 seconds |
| `p` / `n` | previous / next track |
| `Tab` or `Esc` | back to the search input |
| `q` | quit |

`Ctrl+C` quits from anywhere.

## Configuration

An optional TOML file in the user config directory (`~/.config/goremi/config.toml` on Linux, `~/Library/Application Support/goremi/config.toml` on macOS, `%AppData%\goremi\config.toml` on Windows):

```toml
[ui]
theme = "default"       # written by `goremi theme`
show_spectrum = true    # the spectrum on top of the player panel
show_artwork = true     # shown only when show_spectrum is false
mouse = true            # click the player controls; selecting text then needs Shift (Option on macOS)
```

A missing file or an unknown value never stops the player: it opens with the `default` theme. Errors from `mpv` and `yt-dlp` go to `goremi.log` in the user cache directory (`~/.cache/goremi/` on Linux).

## Interface

One screen: search on top, results on the left, the player on the right. With the spectrum on (the default) it fills the top of the player panel; with `show_spectrum = false` the album art takes its place.

```text
┌──────────────────────────────────────────────────────────┐
│ Search: daft punk_                          Ctrl+C: quit │
├───────────────────────────┬──────────────────────────────┤
│ ▶ Get Lucky               │    ▃ ▆ █ ▇ █ ▅ ▃ ▅ ▇ ▄       │
│   One More Time           │    ▂ ▅ █ ▇ █ ▅ ▃ ▇ █ ▆       │
│   Instant Crush           │    Daft Punk                 │
│   load more...            │    Get Lucky                 │
│                           │    ━━━━━●────── 2:31 / 4:08  │
└───────────────────────────┴──────────────────────────────┘
```

## Built with

- [Go](https://go.dev/)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lip Gloss](https://github.com/charmbracelet/lipgloss) for the terminal UI
- [mpv](https://mpv.io/) for audio playback
- [yt-dlp](https://github.com/yt-dlp/yt-dlp) as the YouTube source

## Design principles

- **Simple:** run `goremi` and start searching.
- **Minimal:** a player, not a music management system.
- **Graceful:** no artwork or spectrum, and music still plays.
- **Provider-independent:** YouTube is only the first music source.
- **Theme-independent:** components use semantic colors, never fixed ones.

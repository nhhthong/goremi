<div align="center">
  <img src="resources/logo.png" alt="GoReMi logo" width="360" />

  <p><strong>Terminal music kit. Search, select, listen.</strong></p>

  <p>
    <a href="https://github.com/nhhthong/goremi"><img src="https://img.shields.io/github/go-mod/go-version/nhhthong/goremi?style=flat-square&logo=go&logoColor=white&color=00ADD8" alt="Go version"></a>
    <a href="https://github.com/nhhthong/goremi/commits"><img src="https://img.shields.io/github/last-commit/nhhthong/goremi?style=flat-square" alt="Last commit"></a>
    <img src="https://img.shields.io/badge/status-early%20development-orange?style=flat-square" alt="Status: early development">
    <img src="https://img.shields.io/badge/platform-Linux-informational?style=flat-square&logo=linux&logoColor=white" alt="Platform: Linux">
    <a href="https://github.com/charmbracelet/bubbletea"><img src="https://img.shields.io/badge/TUI-Bubble%20Tea-FF5F87?style=flat-square" alt="TUI: Bubble Tea"></a>
    <img src="https://img.shields.io/badge/license-MIT-blue.svg?style=flat-square" alt="License: MIT">
  </p>
</div>

---

## Overview

GoReMi (**Go** + **Do Re Mi**) is a minimal music player that runs in your terminal. Type a song, pick a result, and it plays, with album art and a live spectrum, all from the keyboard.

> **Status:** early development. The search layer is built; the interface, playback, themes and installer are still in progress.

## Features

- **Search** YouTube from the terminal, 10 results at a time with `load more...`
- **Playback** through [mpv](https://mpv.io/), with play, pause and seek
- **Album artwork** in the terminal (Kitty Graphics Protocol or Sixel), hidden when unsupported
- **Realtime spectrum analyzer** drawn with Unicode blocks
- **Three themes:** Light, Dark (default) and Cyberpunk
- **Keyboard-first**, no mouse needed

## Interface

One screen: search on top, results on the left, the player on the right.

```text
┌──────────────────────────────────────────────────────────┐
│ Search: daft punk_                          Ctrl+C: quit │
├───────────────────────────┬──────────────────────────────┤
│ ▶ Get Lucky               │         [album art]          │
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

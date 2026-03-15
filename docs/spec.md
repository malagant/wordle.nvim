# wordle.nvim — Product Specification

**Author:** Sarah (PM) | **Date:** 2026-03-15 | **Status:** Approved

## Overview

wordle.nvim is a TUI-based Wordle game written in Go using bubbletea/lipgloss. It runs standalone in the terminal AND as a Neovim plugin (`:Wordle` command).

## User Stories

### US-1: Play Wordle
As a player, I want to guess a 5-letter word in up to 6 attempts, so I can enjoy the classic Wordle experience in my terminal.

**Acceptance Criteria:**
- [ ] Player can type a 5-letter word and submit with Enter
- [ ] Invalid words (not in word list) are rejected with feedback
- [ ] After each guess, letters are colored: green (correct position), yellow (wrong position), gray (not in word)
- [ ] Game ends after correct guess (win) or 6 failed attempts (loss)
- [ ] Solution is revealed on loss

### US-2: Language Selection
As a player, I want to choose between English and German at startup, so I can play in my preferred language.

**Acceptance Criteria:**
- [ ] CLI flag `--lang en|de` (default: en)
- [ ] Interactive language picker if no flag provided
- [ ] German words support Umlauts (ä, ö, ü) — treated as single characters
- [ ] Each language has its own curated word list (embedded in binary)

### US-3: Keyboard Display
As a player, I want to see an on-screen keyboard showing which letters I've used, so I can make better guesses.

**Acceptance Criteria:**
- [ ] QWERTY layout displayed below the game grid
- [ ] Letters colored based on best known status (green > yellow > gray)
- [ ] Unused letters shown in default color
- [ ] German keyboard includes Ä, Ö, Ü row

### US-4: Daily Word & Random Mode
As a player, I want a daily word (same for everyone) and a random mode for unlimited play.

**Acceptance Criteria:**
- [ ] Default mode: daily word derived from current date (deterministic hash)
- [ ] `--random` flag for random word each game
- [ ] Daily mode prevents replaying the same day (stats stored locally)

### US-5: Statistics
As a player, I want to track my wins, losses, and streak.

**Acceptance Criteria:**
- [ ] Stats persisted to `~/.local/share/wordle-nvim/stats.json`
- [ ] Track: games played, wins, losses, current streak, max streak, guess distribution
- [ ] Stats screen shown after each game
- [ ] Separate stats per language

### US-6: Neovim Integration
As a Neovim user, I want to run `:Wordle` to play directly in my editor.

**Acceptance Criteria:**
- [ ] Lua plugin provides `:Wordle` command
- [ ] Opens a terminal buffer running the wordle binary
- [ ] Supports `:Wordle de` and `:Wordle random` arguments
- [ ] Works with lazy.nvim package manager

## Technical Constraints

- **Language:** Go 1.21+
- **TUI:** charmbracelet/bubbletea + charmbracelet/lipgloss
- **Word lists:** Embedded via `//go:embed`
- **Binary name:** `wordle-nvim`
- **Platforms:** Linux (amd64, arm64), macOS (amd64, arm64)

## Out of Scope (v1)

- Multiplayer
- Hard mode
- Custom word lists
- Web UI
- Share results (emoji grid)

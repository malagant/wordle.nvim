# 🟩 wordle.nvim

A TUI-based Wordle game that runs standalone in the terminal **and** as a Neovim plugin.

Built with Go, [bubbletea](https://github.com/charmbracelet/bubbletea), and [lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- 🎮 Classic Wordle gameplay (5 letters, 6 attempts)
- 🎨 Colored letter hints (green/yellow/gray)
- ⌨️ On-screen keyboard with letter status
- 🌍 English and German (with Umlaut support: ä, ö, ü)
- 📅 Daily word (date-based) + random mode
- 📊 Persistent statistics (wins, losses, streaks)
- 🔌 Neovim plugin (`:Wordle` command)

## Installation

### Standalone (Go)

```bash
go install github.com/malagant/wordle-nvim/cmd/wordle@latest
```

### Neovim Plugin (lazy.nvim)

```lua
{
  "malagant/wordle.nvim",
  build = "go build -o wordle-nvim ./cmd/wordle/",
  config = function()
    -- Plugin auto-registers :Wordle command
  end,
}
```

## Usage

### Terminal

```bash
# English daily word (default)
wordle-nvim

# German daily word
wordle-nvim --lang de

# Random mode
wordle-nvim --random

# German random
wordle-nvim --lang de --random
```

### Neovim

```vim
:Wordle          " English daily
:Wordle de       " German daily
:Wordle random   " English random
:Wordle de random " German random
```

## Screenshots

```
🟩 WORDLE.NVIM — English (Daily)

 C   R   A   N   E
 S   T   O   R   M
 · · · · ·
 · · · · ·
 · · · · ·
 · · · · ·

  Q  W  E  R  T  Y  U  I  O  P
   A  S  D  F  G  H  J  K  L
    Z  X  C  V  B  N  M
```

## Development

```bash
# Build
go build -o wordle-nvim ./cmd/wordle/

# Test
go test ./...

# Run
./wordle-nvim
```

## Stats

Statistics are saved to `~/.local/share/wordle-nvim/stats-{en|de}.json`.

## License

MIT

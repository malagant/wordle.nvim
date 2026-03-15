# 🟩 wordle.nvim

Play [Wordle](https://www.nytimes.com/games/wordle/) in your terminal — or directly inside Neovim!

Built with [bubbletea](https://github.com/charmbracelet/bubbletea) and [lipgloss](https://github.com/charmbracelet/lipgloss).

## Features

- 🎮 Classic Wordle gameplay (5 letters, 6 attempts)
- 🎨 Color-coded hints (🟩 correct, 🟨 wrong position, ⬜ not in word)
- ⌨️  On-screen keyboard showing letter status
- 🌍 English (~800 words) and German (~350 words, with Umlaute)
- 📅 Daily word (date-based) + random mode
- 📊 Persistent statistics (wins, losses, streaks)
- 🔌 Neovim plugin (`:Wordle` command)
- 🚀 Cross-platform releases (Linux + macOS, amd64/arm64)

## Installation

### Go Install

```bash
go install github.com/malagant/wordle-nvim/cmd/wordle@latest
```

### Download Binary

Grab a release from the [Releases page](https://github.com/malagant/wordle.nvim/releases).

### Neovim Plugin

**lazy.nvim:**
```lua
{
  "malagant/wordle.nvim",
  build = "go build -o wordle-nvim ./cmd/wordle/",
  config = function() end,
}
```

**Manual:**
```bash
cp plugin/wordle.lua ~/.local/share/nvim/site/plugin/
```

## Usage

### Terminal

```bash
wordle              # English, daily word
wordle --lang de    # German
wordle --random     # Random mode
wordle de random    # German + random (positional args)
```

### Neovim

```vim
:Wordle             " English, daily word
:Wordle de          " German
:Wordle random      " Random mode
:Wordle de random   " German + random
```

## Screenshots

```
┌─────────────────────────────────┐
│         W O R D L E             │
│                                 │
│   ┌───┬───┬───┬───┬───┐        │
│   │ S │ T │ A │ R │ E │        │
│   └───┴───┴───┴───┴───┘        │
│   ┌───┬───┬───┬───┬───┐        │
│   │ C │ R │ A │ N │ E │        │
│   └───┴───┴───┴───┴───┘        │
│                                 │
│  Q W E R T Y U I O P           │
│   A S D F G H J K L            │
│     Z X C V B N M              │
└─────────────────────────────────┘
```

## Development

```bash
# Build
go build -o wordle ./cmd/wordle/

# Test
go test -v -race ./...

# Run
./wordle
```

## License

MIT

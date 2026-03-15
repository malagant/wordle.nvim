# Decision 001: Tech Stack

**Date:** 2026-03-15 | **Decider:** Alex (CEO)

## Decision
- Go with bubbletea/lipgloss for TUI
- Word lists embedded via `//go:embed`
- Stats persisted as JSON in `~/.local/share/wordle-nvim/`
- Neovim plugin as thin Lua wrapper that opens terminal buffer

## Rationale
- bubbletea is the de-facto Go TUI framework, well-maintained
- Embedded word lists = single binary, no external dependencies
- JSON stats = simple, human-readable, debuggable
- Lua wrapper keeps plugin minimal, binary does all the work

## Alternatives Considered
- Rust + ratatui: Good option, but Go is in the project brief
- ncurses: Too low-level, no advantage over bubbletea
- Vimscript plugin: Lua is modern standard for Neovim plugins

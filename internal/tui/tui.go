package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/malagant/wordle-nvim/internal/game"
	"github.com/malagant/wordle-nvim/internal/words"
)

// ── Color Palette ──────────────────────────────────────────────────────

var (
	green  = lipgloss.Color("#538d4e")
	yellow = lipgloss.Color("#b59f3b")
	gray   = lipgloss.Color("#3a3a3c")
	dark   = lipgloss.Color("#121213")
	white  = lipgloss.Color("#ffffff")
	dim    = lipgloss.Color("#565758")
	accent = lipgloss.Color("#818384")
	bg     = lipgloss.Color("#0a0a0b")
)

// ── Tile Styles (large, 3 lines per tile) ──────────────────────────────

func tileStyle(bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(bg).
		Width(7).
		Height(3).
		Align(lipgloss.Center, lipgloss.Center)
}

func keyStyle(bgColor lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(white).
		Background(bgColor).
		Width(5).
		Height(3).
		Align(lipgloss.Center, lipgloss.Center)
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Background(lipgloss.Color("#1a1a2e")).
			Padding(1, 4).
			BorderStyle(lipgloss.DoubleBorder()).
			BorderForeground(green).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(accent).
			Align(lipgloss.Center)

	messageStyle = lipgloss.NewStyle().
			Foreground(yellow).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 0)

	winMessageStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true).
			Align(lipgloss.Center).
			Padding(1, 0)

	loseMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#e74c3c")).
				Bold(true).
				Align(lipgloss.Center).
				Padding(1, 0)

	statsBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(accent).
			Padding(1, 3).
			Align(lipgloss.Center)

	statNumberStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(white).
			Align(lipgloss.Center)

	statLabelStyle = lipgloss.NewStyle().
			Foreground(accent).
			Align(lipgloss.Center)

	barFilledStyle = lipgloss.NewStyle().
			Background(green).
			Foreground(white).
			Bold(true)

	barEmptyStyle = lipgloss.NewStyle().
			Background(gray).
			Foreground(white)

	helpStyle = lipgloss.NewStyle().
			Foreground(dim).
			Align(lipgloss.Center)

	separatorStyle = lipgloss.NewStyle().
			Foreground(dim)
)

// ── Config & Model ─────────────────────────────────────────────────────

// Config holds the TUI configuration
type Config struct {
	Language words.Language
	Random   bool
}

// Model is the bubbletea model
type Model struct {
	game    *game.Game
	words   *words.WordList
	input   []rune
	message string
	msgType int // 0=info, 1=win, 2=lose
	stats   *game.Stats
	config  Config
	width   int
	height  int
}

// NewModel creates a new TUI model
func NewModel(cfg Config) Model {
	wl := words.Load(cfg.Language)

	var target string
	if cfg.Random {
		target = wl.RandomWord()
	} else {
		target = wl.DailyWord()
	}

	stats := game.LoadStats(string(cfg.Language))

	if !cfg.Random && stats.HasPlayedToday() {
		return Model{
			game:    game.New(target),
			words:   wl,
			input:   make([]rune, 0, 5),
			message: "Already played today's Wordle! Use --random for a new game.",
			stats:   stats,
			config:  cfg,
			width:   80,
			height:  40,
		}
	}

	return Model{
		game:   game.New(target),
		words:  wl,
		input:  make([]rune, 0, 5),
		stats:  stats,
		config: cfg,
		width:  80,
		height: 40,
	}
}

// ── Bubbletea Interface ────────────────────────────────────────────────

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyBackspace:
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
				m.message = ""
			}

		case tea.KeyEnter:
			if m.game.State != game.StatePlaying {
				return m, tea.Quit
			}
			if utf8.RuneCountInString(string(m.input)) != 5 {
				m.message = "⚠ Word must be 5 letters!"
				m.msgType = 0
				return m, nil
			}
			word := strings.ToLower(string(m.input))
			if !m.words.IsValid(word) {
				m.message = "⚠ Not in word list!"
				m.msgType = 0
				return m, nil
			}
			m.game.Guess(word)
			m.input = m.input[:0]
			m.message = ""

			if m.game.State != game.StatePlaying {
				won := m.game.State == game.StateWon
				m.stats.Record(won, len(m.game.Guesses))
				if !m.config.Random {
					m.stats.MarkDailyPlayed()
				}
				_ = m.stats.Save()

				if won {
					attempts := len(m.game.Guesses)
					emoji := []string{"", "🏆", "🎯", "🔥", "👍", "😅", "😰"}
					e := emoji[attempts]
					m.message = fmt.Sprintf("%s Solved in %d/6!  Press Enter or Esc to exit", e, attempts)
					m.msgType = 1
				} else {
					m.message = fmt.Sprintf("The word was: %s  —  Press Enter or Esc to exit", strings.ToUpper(m.game.Target))
					m.msgType = 2
				}
			}

		default:
			if m.game.State != game.StatePlaying {
				return m, nil
			}
			if msg.Type == tea.KeyRunes {
				for _, r := range msg.Runes {
					if utf8.RuneCountInString(string(m.input)) < 5 {
						m.input = append(m.input, r)
					}
				}
			}
		}
	}

	return m, nil
}

// ── View ───────────────────────────────────────────────────────────────

func (m Model) View() string {
	var sections []string

	// ─ Title ─
	langLabel := "English"
	if m.config.Language == words.German {
		langLabel = "Deutsch"
	}
	mode := "Daily"
	if m.config.Random {
		mode = "Random"
	}
	title := titleStyle.Render(fmt.Sprintf("W O R D L E\n%s · %s", langLabel, mode))
	sections = append(sections, title)

	// ─ Attempts remaining ─
	remaining := m.game.RemainingGuesses() - (m.game.MaxGuesses - len(m.game.Guesses) - func() int {
		if m.game.State == game.StatePlaying {
			return 0
		}
		return 0
	}())
	_ = remaining
	attemptInfo := subtitleStyle.Render(
		fmt.Sprintf("Attempt %d / %d", len(m.game.Guesses)+1, m.game.MaxGuesses))
	if m.game.State != game.StatePlaying {
		attemptInfo = subtitleStyle.Render(
			fmt.Sprintf("Finished in %d / %d", len(m.game.Guesses), m.game.MaxGuesses))
	}
	sections = append(sections, attemptInfo)
	sections = append(sections, "")

	// ─ Game Grid (large tiles) ─
	grid := m.renderGrid()
	sections = append(sections, grid)

	// ─ Message ─
	if m.message != "" {
		var styledMsg string
		switch m.msgType {
		case 1:
			styledMsg = winMessageStyle.Render(m.message)
		case 2:
			styledMsg = loseMessageStyle.Render(m.message)
		default:
			styledMsg = messageStyle.Render(m.message)
		}
		sections = append(sections, styledMsg)
	} else {
		sections = append(sections, "")
	}

	// ─ Keyboard ─
	keyboard := m.renderKeyboard()
	sections = append(sections, keyboard)

	// ─ Live Stats Panel ─
	statsPanel := m.renderStats()
	sections = append(sections, statsPanel)

	// ─ Help ─
	sections = append(sections, "")
	sections = append(sections, helpStyle.Render("Type a word · Enter to submit · Backspace to delete · Esc to quit"))

	// Join all sections
	content := lipgloss.JoinVertical(lipgloss.Center, sections...)

	// Center on screen
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, content,
		lipgloss.WithWhitespaceBackground(bg))
}

// ── Render: Grid ───────────────────────────────────────────────────────

func (m Model) renderGrid() string {
	var rows []string

	for i := 0; i < m.game.MaxGuesses; i++ {
		var tiles []string

		if i < len(m.game.Hints) {
			// Completed guess row
			for _, h := range m.game.Hints[i] {
				bgColor := gray
				switch h.Status {
				case game.StatusCorrect:
					bgColor = green
				case game.StatusPresent:
					bgColor = yellow
				}
				tiles = append(tiles, tileStyle(bgColor).Render(strings.ToUpper(string(h.Letter))))
			}
		} else if i == len(m.game.Guesses) && m.game.State == game.StatePlaying {
			// Current input row
			inputRunes := []rune(string(m.input))
			for j := 0; j < 5; j++ {
				if j < len(inputRunes) {
					tile := tileStyle(dark).
						BorderStyle(lipgloss.NormalBorder()).
						BorderForeground(accent).
						Width(5).Height(1).
						Render(strings.ToUpper(string(inputRunes[j])))
					tiles = append(tiles, tile)
				} else {
					tile := tileStyle(dark).
						BorderStyle(lipgloss.NormalBorder()).
						BorderForeground(dim).
						Width(5).Height(1).
						Render("·")
					tiles = append(tiles, tile)
				}
			}
		} else {
			// Empty row
			for j := 0; j < 5; j++ {
				tile := tileStyle(dark).
					BorderStyle(lipgloss.NormalBorder()).
					BorderForeground(dim).
					Width(5).Height(1).
					Render(" ")
				tiles = append(tiles, tile)
			}
		}

		row := lipgloss.JoinHorizontal(lipgloss.Center, tiles...)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

// ── Render: Keyboard ───────────────────────────────────────────────────

func (m Model) renderKeyboard() string {
	rows := []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
	if m.config.Language == words.German {
		rows = []string{"qwertzuiopü", "asdfghjklöä", "yxcvbnm"}
	}

	var keyboardRows []string
	for _, row := range rows {
		var keys []string
		for _, r := range row {
			bgColor := gray
			if status, ok := m.game.KeyStatus[r]; ok {
				switch status {
				case game.StatusCorrect:
					bgColor = green
				case game.StatusPresent:
					bgColor = yellow
				case game.StatusAbsent:
					bgColor = dark
				}
			}
			keys = append(keys, keyStyle(bgColor).Render(strings.ToUpper(string(r))))
		}
		keyboardRows = append(keyboardRows, lipgloss.JoinHorizontal(lipgloss.Center, keys...))
	}

	return lipgloss.JoinVertical(lipgloss.Center, keyboardRows...)
}

// ── Render: Stats Panel ────────────────────────────────────────────────

func (m Model) renderStats() string {
	s := m.stats

	// Top stats row
	statItems := []string{
		renderStatItem(fmt.Sprintf("%d", s.Played), "Played"),
		renderStatItem(fmt.Sprintf("%.0f%%", s.WinRate()), "Win %"),
		renderStatItem(fmt.Sprintf("%d", s.CurrentStreak), "Streak"),
		renderStatItem(fmt.Sprintf("%d", s.MaxStreak), "Best"),
	}
	statsRow := lipgloss.JoinHorizontal(lipgloss.Center, statItems...)

	// Distribution bars
	var distRows []string
	distRows = append(distRows, "")
	distRows = append(distRows, lipgloss.NewStyle().Foreground(white).Bold(true).Render("GUESS DISTRIBUTION"))
	distRows = append(distRows, "")

	maxCount := 0
	for i := 1; i <= 6; i++ {
		if c, ok := s.Distribution[i]; ok && c > maxCount {
			maxCount = c
		}
	}

	barMaxWidth := 30
	for i := 1; i <= 6; i++ {
		count := s.Distribution[i]
		barWidth := 1
		if maxCount > 0 && count > 0 {
			barWidth = (count * barMaxWidth) / maxCount
			if barWidth < 1 {
				barWidth = 1
			}
		}

		label := fmt.Sprintf(" %d ", i)
		isCurrentGuess := m.game.State == game.StateWon && len(m.game.Guesses) == i

		var bar string
		if count > 0 {
			style := barEmptyStyle
			if isCurrentGuess {
				style = barFilledStyle
			}
			bar = style.Render(fmt.Sprintf("%-*s", barWidth, fmt.Sprintf(" %d", count)))
		} else {
			bar = barEmptyStyle.Render(" 0")
		}

		distRows = append(distRows,
			lipgloss.NewStyle().Foreground(white).Render(label)+" "+bar)
	}

	distribution := lipgloss.JoinVertical(lipgloss.Left, distRows...)

	// Combine
	inner := lipgloss.JoinVertical(lipgloss.Center, statsRow, distribution)
	return statsBoxStyle.Render(inner)
}

func renderStatItem(value, label string) string {
	v := statNumberStyle.Width(8).Render(value)
	l := statLabelStyle.Width(8).Render(label)
	return lipgloss.JoinVertical(lipgloss.Center, v, l)
}

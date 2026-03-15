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

// ── Colors ─────────────────────────────────────────────────────────────

var (
	colorGreen  = lipgloss.Color("#538d4e")
	colorYellow = lipgloss.Color("#b59f3b")
	colorGray   = lipgloss.Color("#3a3a3c")
	colorDark   = lipgloss.Color("#272729")
	colorWhite  = lipgloss.Color("#d7dadc")
	colorDim    = lipgloss.Color("#565758")
	colorAccent = lipgloss.Color("#818384")
	colorBg     = lipgloss.Color("#121213")
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

	// If daily already played, auto-switch to random
	if !cfg.Random && stats.HasPlayedToday() {
		cfg.Random = true
		target = wl.RandomWord()
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
				m.message = "⚠  Word must be 5 letters!"
				m.msgType = 0
				return m, nil
			}
			word := strings.ToLower(string(m.input))
			if !m.words.IsValid(word) {
				m.message = "⚠  Not in word list!"
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
					n := len(m.game.Guesses)
					emoji := []string{"", "🏆", "🎯", "🔥", "👍", "😅", "😰"}
					m.message = fmt.Sprintf("%s Solved in %d/6!", emoji[n], n)
					m.msgType = 1
				} else {
					m.message = fmt.Sprintf("The word was: %s", strings.ToUpper(m.game.Target))
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
	langLabel := "EN"
	if m.config.Language == words.German {
		langLabel = "DE"
	}
	mode := "Daily"
	if m.config.Random {
		mode = "Random"
	}

	// All sections as simple strings, centered later
	var lines []string

	// Title
	lines = append(lines, "")
	lines = append(lines,
		lipgloss.NewStyle().Bold(true).Foreground(colorWhite).Render(
			"W  O  R  D  L  E"))
	lines = append(lines,
		lipgloss.NewStyle().Foreground(colorAccent).Render(
			fmt.Sprintf("%s · %s", langLabel, mode)))
	lines = append(lines, "")

	// Attempt info
	if m.game.State == game.StatePlaying {
		lines = append(lines,
			lipgloss.NewStyle().Foreground(colorDim).Render(
				fmt.Sprintf("Attempt %d of %d", len(m.game.Guesses)+1, m.game.MaxGuesses)))
	} else {
		lines = append(lines,
			lipgloss.NewStyle().Foreground(colorDim).Render(
				fmt.Sprintf("Finished — %d of %d", len(m.game.Guesses), m.game.MaxGuesses)))
	}
	lines = append(lines, "")

	// Grid
	for i := 0; i < m.game.MaxGuesses; i++ {
		lines = append(lines, m.renderRow(i))
	}
	lines = append(lines, "")

	// Message
	if m.message != "" {
		style := lipgloss.NewStyle().Bold(true)
		switch m.msgType {
		case 1:
			style = style.Foreground(colorGreen)
		case 2:
			style = style.Foreground(lipgloss.Color("#e74c3c"))
		default:
			style = style.Foreground(colorYellow)
		}
		lines = append(lines, style.Render(m.message))
	}
	lines = append(lines, "")

	// Keyboard
	kbLines := m.renderKeyboard()
	lines = append(lines, kbLines...)
	lines = append(lines, "")

	// Stats bar
	lines = append(lines, m.renderStatsBar())
	lines = append(lines, "")

	// Distribution (compact, single line per guess)
	lines = append(lines, m.renderDistribution()...)
	lines = append(lines, "")

	// Help
	if m.game.State != game.StatePlaying {
		lines = append(lines,
			lipgloss.NewStyle().Foreground(colorDim).Render("Press Enter or Esc to exit"))
	} else {
		lines = append(lines,
			lipgloss.NewStyle().Foreground(colorDim).Render("Type · Enter · Backspace · Esc"))
	}

	// Center each line
	content := lipgloss.JoinVertical(lipgloss.Center, lines...)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
		lipgloss.WithWhitespaceBackground(colorBg))
}

// ── Render: Single Grid Row ────────────────────────────────────────────

func (m Model) renderRow(i int) string {
	// Each tile: [ X ] with background color, fixed 5 chars wide
	tileWidth := 5
	gap := " "

	var tiles []string

	if i < len(m.game.Hints) {
		for _, h := range m.game.Hints[i] {
			bg := colorGray
			switch h.Status {
			case game.StatusCorrect:
				bg = colorGreen
			case game.StatusPresent:
				bg = colorYellow
			}
			t := lipgloss.NewStyle().
				Background(bg).Foreground(colorWhite).Bold(true).
				Width(tileWidth).Align(lipgloss.Center).
				Render(strings.ToUpper(string(h.Letter)))
			tiles = append(tiles, t)
		}
	} else if i == len(m.game.Guesses) && m.game.State == game.StatePlaying {
		inputRunes := []rune(string(m.input))
		for j := 0; j < 5; j++ {
			var t string
			if j < len(inputRunes) {
				t = lipgloss.NewStyle().
					Background(colorDark).Foreground(colorWhite).Bold(true).
					Width(tileWidth).Align(lipgloss.Center).
					Render(strings.ToUpper(string(inputRunes[j])))
			} else {
				t = lipgloss.NewStyle().
					Background(colorDark).Foreground(colorDim).
					Width(tileWidth).Align(lipgloss.Center).
					Render("·")
			}
			tiles = append(tiles, t)
		}
	} else {
		for j := 0; j < 5; j++ {
			t := lipgloss.NewStyle().
				Background(colorDark).Foreground(colorDim).
				Width(tileWidth).Align(lipgloss.Center).
				Render("·")
			tiles = append(tiles, t)
		}
	}

	return strings.Join(tiles, gap)
}

// ── Render: Keyboard ───────────────────────────────────────────────────

func (m Model) renderKeyboard() []string {
	rows := []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
	if m.config.Language == words.German {
		rows = []string{"qwertzuiopü", "asdfghjklöä", "yxcvbnm"}
	}

	var result []string
	for _, row := range rows {
		var keys []string
		for _, r := range row {
			bg := colorGray
			if status, ok := m.game.KeyStatus[r]; ok {
				switch status {
				case game.StatusCorrect:
					bg = colorGreen
				case game.StatusPresent:
					bg = colorYellow
				case game.StatusAbsent:
					bg = colorDark
				}
			}
			k := lipgloss.NewStyle().
				Background(bg).Foreground(colorWhite).Bold(true).
				Width(3).Align(lipgloss.Center).
				Render(strings.ToUpper(string(r)))
			keys = append(keys, k)
		}
		result = append(result, strings.Join(keys, " "))
	}
	return result
}

// ── Render: Stats Bar ──────────────────────────────────────────────────

func (m Model) renderStatsBar() string {
	s := m.stats
	num := lipgloss.NewStyle().Bold(true).Foreground(colorWhite)
	lbl := lipgloss.NewStyle().Foreground(colorAccent)

	parts := []string{
		num.Render(fmt.Sprintf("%d", s.Played)) + lbl.Render(" played"),
		num.Render(fmt.Sprintf("%.0f%%", s.WinRate())) + lbl.Render(" win"),
		num.Render(fmt.Sprintf("%d", s.CurrentStreak)) + lbl.Render(" streak"),
		num.Render(fmt.Sprintf("%d", s.MaxStreak)) + lbl.Render(" best"),
	}

	sep := lipgloss.NewStyle().Foreground(colorDim).Render("  │  ")
	return strings.Join(parts, sep)
}

// ── Render: Distribution ───────────────────────────────────────────────

func (m Model) renderDistribution() []string {
	s := m.stats
	maxCount := 0
	for i := 1; i <= 6; i++ {
		if c := s.Distribution[i]; c > maxCount {
			maxCount = c
		}
	}

	maxBarWidth := 20
	var lines []string

	for i := 1; i <= 6; i++ {
		count := s.Distribution[i]
		barWidth := 0
		if maxCount > 0 && count > 0 {
			barWidth = (count * maxBarWidth) / maxCount
			if barWidth < 1 {
				barWidth = 1
			}
		}

		isHighlight := m.game.State == game.StateWon && len(m.game.Guesses) == i

		label := lipgloss.NewStyle().Foreground(colorWhite).Width(2).Align(lipgloss.Right).
			Render(fmt.Sprintf("%d", i))

		countStr := fmt.Sprintf(" %d ", count)
		var bar string
		if isHighlight {
			bar = lipgloss.NewStyle().Background(colorGreen).Foreground(colorWhite).Bold(true).
				Render(countStr + strings.Repeat("▓", barWidth))
		} else if count > 0 {
			bar = lipgloss.NewStyle().Background(colorGray).Foreground(colorWhite).
				Render(countStr + strings.Repeat("▓", barWidth))
		} else {
			bar = lipgloss.NewStyle().Background(colorGray).Foreground(colorWhite).
				Render(countStr)
		}

		lines = append(lines, label+" "+bar)
	}

	return lines
}

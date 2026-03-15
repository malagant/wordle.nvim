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
	colorDark   = lipgloss.Color("#121213")
	colorWhite  = lipgloss.Color("#d7dadc")
	colorDim    = lipgloss.Color("#565758")
	colorAccent = lipgloss.Color("#818384")
	colorBg     = lipgloss.Color("#121213")
)

// ── Styles ─────────────────────────────────────────────────────────────

func makeTileStyle(bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Background(bg).
		Padding(0, 1).
		Align(lipgloss.Center)
}

func makeEmptyTile() lipgloss.Style {
	return lipgloss.NewStyle().
		Foreground(colorDim).
		Padding(0, 1).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorDim)
}

func makeInputTile() lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Padding(0, 1).
		Align(lipgloss.Center).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorAccent)
}

func makeKeyStyle(bg lipgloss.Color) lipgloss.Style {
	return lipgloss.NewStyle().
		Bold(true).
		Foreground(colorWhite).
		Background(bg).
		Padding(0, 1).
		Align(lipgloss.Center)
}

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Align(lipgloss.Center)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Align(lipgloss.Center)

	messageStyle = lipgloss.NewStyle().
			Foreground(colorYellow).
			Bold(true).
			Align(lipgloss.Center)

	winMsgStyle = lipgloss.NewStyle().
			Foreground(colorGreen).
			Bold(true).
			Align(lipgloss.Center)

	loseMsgStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#e74c3c")).
			Bold(true).
			Align(lipgloss.Center)

	statsBoxStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorDim).
			Padding(1, 2).
			Align(lipgloss.Center)

	statNumStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorWhite).
			Width(10).
			Align(lipgloss.Center)

	statLblStyle = lipgloss.NewStyle().
			Foreground(colorAccent).
			Width(10).
			Align(lipgloss.Center)

	helpStyle = lipgloss.NewStyle().
			Foreground(colorDim).
			Align(lipgloss.Center)
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
					m.message = fmt.Sprintf("%s  Solved in %d/6!  Press Enter or Esc to exit", emoji[n], n)
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

	// Title
	langLabel := "English"
	if m.config.Language == words.German {
		langLabel = "Deutsch"
	}
	mode := "Daily"
	if m.config.Random {
		mode = "Random"
	}

	sections = append(sections,
		titleStyle.Render("╔══════════════════════════════════╗"))
	sections = append(sections,
		titleStyle.Render("║        W  O  R  D  L  E         ║"))
	sections = append(sections,
		titleStyle.Render(fmt.Sprintf("║      %s · %-7s           ║", langLabel, mode)))
	sections = append(sections,
		titleStyle.Render("╚══════════════════════════════════╝"))
	sections = append(sections, "")

	// Attempt counter
	if m.game.State == game.StatePlaying {
		sections = append(sections,
			subtitleStyle.Render(fmt.Sprintf("Attempt %d / %d", len(m.game.Guesses)+1, m.game.MaxGuesses)))
	} else {
		sections = append(sections,
			subtitleStyle.Render(fmt.Sprintf("Finished — %d / %d", len(m.game.Guesses), m.game.MaxGuesses)))
	}
	sections = append(sections, "")

	// Game Grid
	sections = append(sections, m.renderGrid())
	sections = append(sections, "")

	// Message
	if m.message != "" {
		switch m.msgType {
		case 1:
			sections = append(sections, winMsgStyle.Render(m.message))
		case 2:
			sections = append(sections, loseMsgStyle.Render(m.message))
		default:
			sections = append(sections, messageStyle.Render(m.message))
		}
	}
	sections = append(sections, "")

	// Keyboard
	sections = append(sections, m.renderKeyboard())
	sections = append(sections, "")

	// Stats
	sections = append(sections, m.renderStats())
	sections = append(sections, "")

	// Help
	sections = append(sections,
		helpStyle.Render("Type a word · Enter to submit · Backspace to delete · Esc to quit"))

	content := lipgloss.JoinVertical(lipgloss.Center, sections...)

	// Center the whole thing on screen
	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
		lipgloss.WithWhitespaceBackground(colorBg))
}

// ── Grid ───────────────────────────────────────────────────────────────

func (m Model) renderGrid() string {
	var rows []string

	for i := 0; i < m.game.MaxGuesses; i++ {
		var cells []string

		if i < len(m.game.Hints) {
			// Completed guess
			for _, h := range m.game.Hints[i] {
				bg := colorGray
				switch h.Status {
				case game.StatusCorrect:
					bg = colorGreen
				case game.StatusPresent:
					bg = colorYellow
				}
				cell := makeTileStyle(bg).Render(fmt.Sprintf(" %s ", strings.ToUpper(string(h.Letter))))
				cells = append(cells, cell)
			}
		} else if i == len(m.game.Guesses) && m.game.State == game.StatePlaying {
			// Current input
			inputRunes := []rune(string(m.input))
			for j := 0; j < 5; j++ {
				if j < len(inputRunes) {
					cell := makeInputTile().Render(fmt.Sprintf(" %s ", strings.ToUpper(string(inputRunes[j]))))
					cells = append(cells, cell)
				} else {
					cell := makeEmptyTile().Render("   ")
					cells = append(cells, cell)
				}
			}
		} else {
			// Empty row
			for j := 0; j < 5; j++ {
				cell := makeEmptyTile().Render("   ")
				cells = append(cells, cell)
			}
		}

		row := lipgloss.JoinHorizontal(lipgloss.Top, cells...)
		rows = append(rows, row)
	}

	return lipgloss.JoinVertical(lipgloss.Center, rows...)
}

// ── Keyboard ───────────────────────────────────────────────────────────

func (m Model) renderKeyboard() string {
	rows := []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
	if m.config.Language == words.German {
		rows = []string{"qwertzuiopü", "asdfghjklöä", "yxcvbnm"}
	}

	var kbRows []string
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
			key := makeKeyStyle(bg).Render(fmt.Sprintf(" %s ", strings.ToUpper(string(r))))
			keys = append(keys, key)
		}
		kbRows = append(kbRows, lipgloss.JoinHorizontal(lipgloss.Top, keys...))
	}

	return lipgloss.JoinVertical(lipgloss.Center, kbRows...)
}

// ── Stats ──────────────────────────────────────────────────────────────

func (m Model) renderStats() string {
	s := m.stats

	// Stat numbers row
	items := []string{
		lipgloss.JoinVertical(lipgloss.Center,
			statNumStyle.Render(fmt.Sprintf("%d", s.Played)),
			statLblStyle.Render("Played")),
		lipgloss.JoinVertical(lipgloss.Center,
			statNumStyle.Render(fmt.Sprintf("%.0f%%", s.WinRate())),
			statLblStyle.Render("Win %")),
		lipgloss.JoinVertical(lipgloss.Center,
			statNumStyle.Render(fmt.Sprintf("%d", s.CurrentStreak)),
			statLblStyle.Render("Streak")),
		lipgloss.JoinVertical(lipgloss.Center,
			statNumStyle.Render(fmt.Sprintf("%d", s.MaxStreak)),
			statLblStyle.Render("Best")),
	}
	statsRow := lipgloss.JoinHorizontal(lipgloss.Top, items...)

	// Distribution
	var distLines []string
	distLines = append(distLines, "")
	distLines = append(distLines,
		lipgloss.NewStyle().Bold(true).Foreground(colorWhite).Render("GUESS DISTRIBUTION"))
	distLines = append(distLines, "")

	maxCount := 0
	for i := 1; i <= 6; i++ {
		if c := s.Distribution[i]; c > maxCount {
			maxCount = c
		}
	}

	maxBarWidth := 25
	for i := 1; i <= 6; i++ {
		count := s.Distribution[i]
		barWidth := 1
		if maxCount > 0 && count > 0 {
			barWidth = (count * maxBarWidth) / maxCount
			if barWidth < 1 {
				barWidth = 1
			}
		}

		label := fmt.Sprintf(" %d ", i)
		isHighlight := m.game.State == game.StateWon && len(m.game.Guesses) == i

		barText := fmt.Sprintf(" %d ", count)
		var bar string
		if isHighlight {
			bar = lipgloss.NewStyle().
				Background(colorGreen).Foreground(colorWhite).Bold(true).
				Width(barWidth + len(barText)).
				Render(barText)
		} else if count > 0 {
			bar = lipgloss.NewStyle().
				Background(colorGray).Foreground(colorWhite).
				Width(barWidth + len(barText)).
				Render(barText)
		} else {
			bar = lipgloss.NewStyle().
				Background(colorGray).Foreground(colorWhite).
				Render(barText)
		}

		distLines = append(distLines,
			lipgloss.NewStyle().Foreground(colorWhite).Render(label)+" "+bar)
	}

	dist := lipgloss.JoinVertical(lipgloss.Left, distLines...)
	inner := lipgloss.JoinVertical(lipgloss.Center, statsRow, dist)

	return statsBoxStyle.Render(inner)
}

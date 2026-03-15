package tui

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/malagant/wordle-nvim/internal/game"
	"github.com/malagant/wordle-nvim/internal/words"
)

// Colors
var (
	correctStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#538d4e"))

	presentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#b59f3b"))

	absentStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#3a3a3c"))

	emptyStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			Background(lipgloss.Color("#121213"))

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#ffffff")).
			MarginBottom(1)

	messageStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#818384")).
			MarginTop(1)

	statsStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ffffff")).
			Bold(true)
)

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
	stats   *game.Stats
	config  Config
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

	// Check if daily already played
	if !cfg.Random && stats.HasPlayedToday() {
		return Model{
			game:    game.New(target),
			words:   wl,
			input:   make([]rune, 0, 5),
			message: "Du hast das tägliche Wordle heute schon gespielt! Nutze --random für ein neues Spiel.",
			stats:   stats,
			config:  cfg,
		}
	}

	return Model{
		game:   game.New(target),
		words:  wl,
		input:  make([]rune, 0, 5),
		stats:  stats,
		config: cfg,
	}
}

// Init implements bubbletea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements bubbletea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
				m.message = "Wort muss 5 Buchstaben haben!"
				return m, nil
			}
			word := strings.ToLower(string(m.input))
			if !m.words.IsValid(word) {
				m.message = "Unbekanntes Wort!"
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
					m.message = fmt.Sprintf("🎉 Gewonnen in %d Versuch(en)! [Enter/Esc zum Beenden]", len(m.game.Guesses))
				} else {
					m.message = fmt.Sprintf("💀 Verloren! Das Wort war: %s [Enter/Esc zum Beenden]", strings.ToUpper(m.game.Target))
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

// View implements bubbletea.Model
func (m Model) View() string {
	var b strings.Builder

	// Title
	langLabel := "English"
	if m.config.Language == words.German {
		langLabel = "Deutsch"
	}
	mode := "Daily"
	if m.config.Random {
		mode = "Random"
	}
	b.WriteString(titleStyle.Render(fmt.Sprintf("🟩 WORDLE.NVIM — %s (%s)", langLabel, mode)))
	b.WriteString("\n\n")

	// Game grid
	for i := 0; i < m.game.MaxGuesses; i++ {
		if i < len(m.game.Hints) {
			// Rendered guess
			for _, h := range m.game.Hints[i] {
				style := absentStyle
				switch h.Status {
				case game.StatusCorrect:
					style = correctStyle
				case game.StatusPresent:
					style = presentStyle
				}
				b.WriteString(style.Render(fmt.Sprintf(" %s ", strings.ToUpper(string(h.Letter)))))
				b.WriteString(" ")
			}
		} else if i == len(m.game.Guesses) && m.game.State == game.StatePlaying {
			// Current input row
			inputRunes := []rune(string(m.input))
			for j := 0; j < 5; j++ {
				if j < len(inputRunes) {
					b.WriteString(emptyStyle.Render(fmt.Sprintf(" %s ", strings.ToUpper(string(inputRunes[j])))))
				} else {
					b.WriteString(emptyStyle.Render(" · "))
				}
				b.WriteString(" ")
			}
		} else {
			// Empty row
			for j := 0; j < 5; j++ {
				b.WriteString(emptyStyle.Render(" · "))
				b.WriteString(" ")
			}
		}
		b.WriteString("\n")
	}

	// Keyboard
	b.WriteString("\n")
	rows := []string{"qwertyuiop", "asdfghjkl", "zxcvbnm"}
	if m.config.Language == words.German {
		rows = []string{"qwertzuiopü", "asdfghjklöä", "yxcvbnm"}
	}

	for _, row := range rows {
		b.WriteString("  ")
		for _, r := range row {
			style := emptyStyle
			if status, ok := m.game.KeyStatus[r]; ok {
				switch status {
				case game.StatusCorrect:
					style = correctStyle
				case game.StatusPresent:
					style = presentStyle
				case game.StatusAbsent:
					style = absentStyle
				}
			}
			b.WriteString(style.Render(fmt.Sprintf(" %s ", strings.ToUpper(string(r)))))
			b.WriteString(" ")
		}
		b.WriteString("\n")
	}

	// Message
	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(messageStyle.Render(m.message))
	}

	// Stats after game over
	if m.game.State != game.StatePlaying {
		b.WriteString("\n\n")
		b.WriteString(statsStyle.Render("📊 Statistiken"))
		b.WriteString("\n")
		b.WriteString(fmt.Sprintf("  Gespielt: %d | Gewonnen: %d | Verloren: %d\n", m.stats.Played, m.stats.Won, m.stats.Lost))
		b.WriteString(fmt.Sprintf("  Gewinnrate: %.0f%% | Streak: %d | Max Streak: %d\n", m.stats.WinRate(), m.stats.CurrentStreak, m.stats.MaxStreak))
	}

	b.WriteString("\n")
	return b.String()
}

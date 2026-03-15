package game

import (
	"strings"
	"unicode/utf8"
)

// HintStatus represents the status of a letter in a guess
type HintStatus int

const (
	StatusAbsent  HintStatus = iota // gray — not in word
	StatusPresent                   // yellow — wrong position
	StatusCorrect                   // green — correct position
)

// LetterHint represents a single letter and its hint status
type LetterHint struct {
	Letter rune
	Status HintStatus
}

// GameState represents the current state of the game
type GameState int

const (
	StatePlaying GameState = iota
	StateWon
	StateLost
)

// Game holds the state of a Wordle game
type Game struct {
	Target     string
	Guesses    []string
	Hints      [][]LetterHint
	MaxGuesses int
	State      GameState
	KeyStatus  map[rune]HintStatus
}

// New creates a new Wordle game with the given target word
func New(target string) *Game {
	return &Game{
		Target:     strings.ToLower(target),
		Guesses:    make([]string, 0, 6),
		Hints:      make([][]LetterHint, 0, 6),
		MaxGuesses: 6,
		State:      StatePlaying,
		KeyStatus:  make(map[rune]HintStatus),
	}
}

// Guess processes a guess and returns the hints
func (g *Game) Guess(word string) []LetterHint {
	word = strings.ToLower(word)
	g.Guesses = append(g.Guesses, word)

	targetRunes := []rune(g.Target)
	guessRunes := []rune(word)
	hints := make([]LetterHint, len(guessRunes))

	// Track which target letters have been matched
	matched := make([]bool, len(targetRunes))

	// First pass: find correct positions (green)
	for i, r := range guessRunes {
		if i < len(targetRunes) && r == targetRunes[i] {
			hints[i] = LetterHint{Letter: r, Status: StatusCorrect}
			matched[i] = true
		} else {
			hints[i] = LetterHint{Letter: r, Status: StatusAbsent}
		}
	}

	// Second pass: find present but wrong position (yellow)
	for i, h := range hints {
		if h.Status == StatusCorrect {
			continue
		}
		for j, tr := range targetRunes {
			if !matched[j] && guessRunes[i] == tr {
				hints[i].Status = StatusPresent
				matched[j] = true
				break
			}
		}
	}

	g.Hints = append(g.Hints, hints)

	// Update keyboard status
	for _, h := range hints {
		existing, ok := g.KeyStatus[h.Letter]
		if !ok || h.Status > existing {
			g.KeyStatus[h.Letter] = h.Status
		}
	}

	// Check game state
	if word == g.Target {
		g.State = StateWon
	} else if len(g.Guesses) >= g.MaxGuesses {
		g.State = StateLost
	}

	return hints
}

// RemainingGuesses returns how many guesses are left
func (g *Game) RemainingGuesses() int {
	return g.MaxGuesses - len(g.Guesses)
}

// TargetLength returns the rune count of the target word
func (g *Game) TargetLength() int {
	return utf8.RuneCountInString(g.Target)
}

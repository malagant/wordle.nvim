package game

import (
	"testing"
)

func TestNewGame(t *testing.T) {
	g := New("hello")
	if g.Target != "hello" {
		t.Errorf("expected target 'hello', got '%s'", g.Target)
	}
	if g.State != StatePlaying {
		t.Errorf("expected StatePlaying, got %d", g.State)
	}
	if g.MaxGuesses != 6 {
		t.Errorf("expected 6 max guesses, got %d", g.MaxGuesses)
	}
	if len(g.Guesses) != 0 {
		t.Errorf("expected 0 guesses, got %d", len(g.Guesses))
	}
}

func TestGuessAllCorrect(t *testing.T) {
	g := New("hello")
	hints := g.Guess("hello")

	for i, h := range hints {
		if h.Status != StatusCorrect {
			t.Errorf("hint[%d]: expected StatusCorrect, got %d", i, h.Status)
		}
	}
	if g.State != StateWon {
		t.Errorf("expected StateWon, got %d", g.State)
	}
}

func TestGuessAllAbsent(t *testing.T) {
	g := New("hello")
	hints := g.Guess("fudgy")

	for i, h := range hints {
		if h.Status != StatusAbsent {
			t.Errorf("hint[%d]: expected StatusAbsent for '%c', got %d", i, h.Letter, h.Status)
		}
	}
	if g.State != StatePlaying {
		t.Errorf("expected StatePlaying, got %d", g.State)
	}
}

func TestGuessPresent(t *testing.T) {
	g := New("hello")
	hints := g.Guess("olelh")

	// o is present (in hello at pos 4, guessed at pos 0)
	if hints[0].Status != StatusPresent {
		t.Errorf("hint[0] 'o': expected StatusPresent, got %d", hints[0].Status)
	}
	// l at pos 1: l exists in hello at pos 2,3 — present
	if hints[1].Status != StatusPresent {
		t.Errorf("hint[1] 'l': expected StatusPresent, got %d", hints[1].Status)
	}
	// e at pos 2: e exists in hello at pos 1 — present
	if hints[2].Status != StatusPresent {
		t.Errorf("hint[2] 'e': expected StatusPresent, got %d", hints[2].Status)
	}
	// l at pos 3: correct position
	if hints[3].Status != StatusCorrect {
		t.Errorf("hint[3] 'l': expected StatusCorrect, got %d", hints[3].Status)
	}
	// h at pos 4: h exists at pos 0 — present
	if hints[4].Status != StatusPresent {
		t.Errorf("hint[4] 'h': expected StatusPresent, got %d", hints[4].Status)
	}
}

func TestGuessDoubleLetters(t *testing.T) {
	// Target has one 'a', guess has two 'a's
	g := New("crane")
	hints := g.Guess("salad")

	// s at 0: absent
	if hints[0].Status != StatusAbsent {
		t.Errorf("hint[0] 's': expected StatusAbsent, got %d", hints[0].Status)
	}
	// a at 1: present (a is at pos 2 in crane)
	if hints[1].Status != StatusPresent {
		t.Errorf("hint[1] 'a': expected StatusPresent, got %d", hints[1].Status)
	}
	// l at 2: absent
	if hints[2].Status != StatusAbsent {
		t.Errorf("hint[2] 'l': expected StatusAbsent, got %d", hints[2].Status)
	}
	// a at 3: absent (only one 'a' in target, already matched)
	if hints[3].Status != StatusAbsent {
		t.Errorf("hint[3] 'a': expected StatusAbsent, got %d", hints[3].Status)
	}
	// d at 4: absent
	if hints[4].Status != StatusAbsent {
		t.Errorf("hint[4] 'd': expected StatusAbsent, got %d", hints[4].Status)
	}
}

func TestGameLostAfterSixGuesses(t *testing.T) {
	g := New("hello")
	wrongGuesses := []string{"world", "about", "crane", "fudgy", "blimp", "stark"}

	for i, guess := range wrongGuesses {
		g.Guess(guess)
		if i < 5 && g.State != StatePlaying {
			t.Errorf("after guess %d: expected StatePlaying", i)
		}
	}
	if g.State != StateLost {
		t.Errorf("after 6 wrong guesses: expected StateLost, got %d", g.State)
	}
}

func TestRemainingGuesses(t *testing.T) {
	g := New("hello")
	if g.RemainingGuesses() != 6 {
		t.Errorf("expected 6 remaining, got %d", g.RemainingGuesses())
	}
	g.Guess("world")
	if g.RemainingGuesses() != 5 {
		t.Errorf("expected 5 remaining, got %d", g.RemainingGuesses())
	}
}

func TestKeyboardStatus(t *testing.T) {
	g := New("hello")
	g.Guess("helps")

	// h: correct (pos 0)
	if g.KeyStatus['h'] != StatusCorrect {
		t.Errorf("key 'h': expected StatusCorrect, got %d", g.KeyStatus['h'])
	}
	// e: correct (pos 1)
	if g.KeyStatus['e'] != StatusCorrect {
		t.Errorf("key 'e': expected StatusCorrect, got %d", g.KeyStatus['e'])
	}
	// l: correct (pos 2)
	if g.KeyStatus['l'] != StatusCorrect {
		t.Errorf("key 'l': expected StatusCorrect, got %d", g.KeyStatus['l'])
	}
	// p: absent
	if g.KeyStatus['p'] != StatusAbsent {
		t.Errorf("key 'p': expected StatusAbsent, got %d", g.KeyStatus['p'])
	}
	// s: absent
	if g.KeyStatus['s'] != StatusAbsent {
		t.Errorf("key 's': expected StatusAbsent, got %d", g.KeyStatus['s'])
	}
}

func TestKeyboardStatusUpgradesOnly(t *testing.T) {
	g := New("hello")
	// First guess: 'o' is present (yellow)
	g.Guess("going")
	if g.KeyStatus['o'] != StatusPresent {
		t.Errorf("key 'o' after first guess: expected StatusPresent, got %d", g.KeyStatus['o'])
	}

	// Second guess: 'o' is correct (green) — should upgrade
	g.Guess("hallo")
	if g.KeyStatus['o'] != StatusCorrect {
		t.Errorf("key 'o' after second guess: expected StatusCorrect, got %d", g.KeyStatus['o'])
	}
}

func TestGermanUmlauts(t *testing.T) {
	g := New("bäder")
	if g.TargetLength() != 5 {
		t.Errorf("expected target length 5, got %d", g.TargetLength())
	}

	hints := g.Guess("bäder")
	for i, h := range hints {
		if h.Status != StatusCorrect {
			t.Errorf("hint[%d]: expected StatusCorrect for '%c', got %d", i, h.Letter, h.Status)
		}
	}
	if g.State != StateWon {
		t.Errorf("expected StateWon")
	}
}

func TestTargetIsCaseInsensitive(t *testing.T) {
	g := New("HELLO")
	hints := g.Guess("hello")
	for i, h := range hints {
		if h.Status != StatusCorrect {
			t.Errorf("hint[%d]: expected StatusCorrect, got %d", i, h.Status)
		}
	}
}

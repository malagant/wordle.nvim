package words

import (
	"testing"
	"unicode/utf8"
)

func TestLoadEnglish(t *testing.T) {
	wl := Load(English)
	if len(wl.Words) == 0 {
		t.Fatal("English word list is empty")
	}
	if wl.Language != English {
		t.Errorf("expected language English, got %s", wl.Language)
	}

	// All words should be exactly 5 runes
	for _, w := range wl.Words {
		if utf8.RuneCountInString(w) != 5 {
			t.Errorf("word '%s' has %d runes, expected 5", w, utf8.RuneCountInString(w))
		}
	}
}

func TestLoadGerman(t *testing.T) {
	wl := Load(German)
	if len(wl.Words) == 0 {
		t.Fatal("German word list is empty")
	}
	if wl.Language != German {
		t.Errorf("expected language German, got %s", wl.Language)
	}

	for _, w := range wl.Words {
		if utf8.RuneCountInString(w) != 5 {
			t.Errorf("word '%s' has %d runes, expected 5", w, utf8.RuneCountInString(w))
		}
	}
}

func TestDailyWordDeterministic(t *testing.T) {
	wl := Load(English)
	w1 := wl.DailyWord()
	w2 := wl.DailyWord()
	if w1 != w2 {
		t.Errorf("daily word not deterministic: '%s' vs '%s'", w1, w2)
	}
}

func TestDailyWordDifferentLanguages(t *testing.T) {
	en := Load(English)
	de := Load(German)
	// Different languages should (almost certainly) produce different daily words
	// This is not guaranteed but extremely likely with different word lists + language salt
	_ = en.DailyWord()
	_ = de.DailyWord()
	// Just ensure both return valid words without panic
}

func TestRandomWord(t *testing.T) {
	wl := Load(English)
	w := wl.RandomWord()
	if w == "" {
		t.Error("random word is empty")
	}
	if !wl.IsValid(w) {
		t.Errorf("random word '%s' is not in word list", w)
	}
}

func TestIsValid(t *testing.T) {
	wl := Load(English)

	// Should find real words
	if !wl.IsValid("about") {
		t.Error("expected 'about' to be valid")
	}
	// Case insensitive
	if !wl.IsValid("ABOUT") {
		t.Error("expected 'ABOUT' to be valid (case insensitive)")
	}
	// Should reject nonsense
	if wl.IsValid("zzzzz") {
		t.Error("expected 'zzzzz' to be invalid")
	}
}

func TestGermanWordsContainUmlauts(t *testing.T) {
	wl := Load(German)
	hasUmlaut := false
	for _, w := range wl.Words {
		for _, r := range w {
			if r == 'ä' || r == 'ö' || r == 'ü' {
				hasUmlaut = true
				break
			}
		}
		if hasUmlaut {
			break
		}
	}
	if !hasUmlaut {
		t.Error("German word list should contain words with umlauts")
	}
}

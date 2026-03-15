package words

import (
	"crypto/sha256"
	_ "embed"
	"encoding/binary"
	"math/rand"
	"strings"
	"time"
	"unicode/utf8"
)

//go:embed en.txt
var enFile string

//go:embed de.txt
var deFile string

// Language represents a supported language
type Language string

const (
	English Language = "en"
	German  Language = "de"
)

// WordList holds the words for a language
type WordList struct {
	Words    []string
	Language Language
}

// Load returns the word list for the given language
func Load(lang Language) *WordList {
	var raw string
	switch lang {
	case German:
		raw = deFile
	default:
		raw = enFile
	}

	lines := strings.Split(strings.TrimSpace(raw), "\n")
	words := make([]string, 0, len(lines))
	for _, line := range lines {
		w := strings.TrimSpace(strings.ToLower(line))
		if utf8.RuneCountInString(w) == 5 {
			words = append(words, w)
		}
	}

	return &WordList{Words: words, Language: lang}
}

// DailyWord returns the word of the day (deterministic based on date)
func (wl *WordList) DailyWord() string {
	now := time.Now()
	dateStr := now.Format("2006-01-02")
	hash := sha256.Sum256([]byte(dateStr + string(wl.Language)))
	idx := binary.BigEndian.Uint64(hash[:8]) % uint64(len(wl.Words))
	return wl.Words[idx]
}

// RandomWord returns a random word from the list
func (wl *WordList) RandomWord() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return wl.Words[r.Intn(len(wl.Words))]
}

// IsValid checks if a word is in the word list
func (wl *WordList) IsValid(word string) bool {
	w := strings.ToLower(word)
	for _, candidate := range wl.Words {
		if candidate == w {
			return true
		}
	}
	return false
}

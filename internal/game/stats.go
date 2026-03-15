package game

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Stats tracks player statistics
type Stats struct {
	Played        int            `json:"played"`
	Won           int            `json:"won"`
	Lost          int            `json:"lost"`
	CurrentStreak int            `json:"current_streak"`
	MaxStreak     int            `json:"max_streak"`
	Distribution  map[int]int    `json:"distribution"` // guesses → count
	LastDaily     string         `json:"last_daily"`    // date of last daily play
	Language      string         `json:"language"`
}

// StatsPath returns the path to the stats file for a language
func StatsPath(lang string) string {
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, _ := os.UserHomeDir()
		dataDir = filepath.Join(home, ".local", "share")
	}
	return filepath.Join(dataDir, "wordle-nvim", "stats-"+lang+".json")
}

// LoadStats loads stats from disk
func LoadStats(lang string) *Stats {
	path := StatsPath(lang)
	data, err := os.ReadFile(path)
	if err != nil {
		return &Stats{
			Distribution: make(map[int]int),
			Language:     lang,
		}
	}
	var s Stats
	if err := json.Unmarshal(data, &s); err != nil {
		return &Stats{
			Distribution: make(map[int]int),
			Language:     lang,
		}
	}
	if s.Distribution == nil {
		s.Distribution = make(map[int]int)
	}
	return &s
}

// Save persists stats to disk
func (s *Stats) Save() error {
	path := StatsPath(s.Language)
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// Record records the result of a game
func (s *Stats) Record(won bool, guesses int) {
	s.Played++
	if won {
		s.Won++
		s.CurrentStreak++
		if s.CurrentStreak > s.MaxStreak {
			s.MaxStreak = s.CurrentStreak
		}
		s.Distribution[guesses]++
	} else {
		s.Lost++
		s.CurrentStreak = 0
	}
}

// HasPlayedToday checks if the daily word was already played today
func (s *Stats) HasPlayedToday() bool {
	today := time.Now().Format("2006-01-02")
	return s.LastDaily == today
}

// MarkDailyPlayed marks today's daily as played
func (s *Stats) MarkDailyPlayed() {
	s.LastDaily = time.Now().Format("2006-01-02")
}

// WinRate returns the win percentage
func (s *Stats) WinRate() float64 {
	if s.Played == 0 {
		return 0
	}
	return float64(s.Won) / float64(s.Played) * 100
}

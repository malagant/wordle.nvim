package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/malagant/wordle-nvim/internal/tui"
	"github.com/malagant/wordle-nvim/internal/words"
)

var version = "dev"

func main() {
	lang := flag.String("lang", "en", "Language: en or de")
	random := flag.Bool("random", false, "Random mode (no daily word)")
	showVersion := flag.Bool("version", false, "Show version")
	flag.Parse()

	if *showVersion {
		fmt.Printf("wordle-nvim %s\n", version)
		os.Exit(0)
	}

	// Also accept positional args for Neovim integration: wordle de, wordle random
	for _, arg := range flag.Args() {
		switch arg {
		case "de", "deutsch", "german":
			*lang = "de"
		case "en", "english":
			*lang = "en"
		case "random":
			*random = true
		}
	}

	var language words.Language
	switch *lang {
	case "de":
		language = words.German
	default:
		language = words.English
	}

	cfg := tui.Config{
		Language: language,
		Random:   *random,
	}

	m := tui.NewModel(cfg)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

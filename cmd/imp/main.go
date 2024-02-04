package main

import (
	"encoding/json"
	"flag"
	"github.com/karlovskiy/vocabacov/internal/database"
	"github.com/karlovskiy/vocabacov/internal/translate"
	"log/slog"
	"os"
)

var format = flag.String("format", "internal", "data format: internal or flashcards")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		slog.Error("arg with path to file to import not found")
		os.Exit(1)
	}
	if *format != "internal" && *format != "flashcards" {
		slog.Error("bad format", "format", *format)
		os.Exit(1)
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		slog.Error("import file read error", "err", err)
		os.Exit(2)
	}
	var phrases map[string][]Phrase
	if *format == "internal" {
		if err := json.Unmarshal(data, &phrases); err != nil {
			slog.Error("json unmarshall error", "format", *format, "err", err)
			os.Exit(3)
		}
	} else {
		var collection SimpleFlashcardsCollection
		if err := json.Unmarshal(data, &collection); err != nil {
			slog.Error("json unmarshall error", "format", *format, "err", err)
			os.Exit(3)
		}
		cards := make([]Phrase, 0, len(collection.Cards))
		for _, c := range collection.Cards {
			cards = append(cards, Phrase{
				Phrase:    c.FrontText,
				Translate: c.BackText,
			})
		}
		phrases = make(map[string][]Phrase, 1)
		phrases[collection.Name] = cards
	}
	db, err := database.OpenDb(true)
	if err != nil {
		slog.Error("db open error", "err", err)
		os.Exit(4)
	}
	for lang, langPhrases := range phrases {
		for _, p := range langPhrases {
			phrase := &translate.Phrase{
				Lang:        lang,
				Phrase:      p.Phrase,
				Translation: p.Translate,
			}
			if err := database.SavePhrase(db, phrase); err != nil {
				slog.Error("save phrase error", "err", err)
				os.Exit(5)
			}
		}
	}
}

type Phrase struct {
	Phrase    string `json:"p"`
	Translate string `json:"t"`
}

type SimpleFlashcardsCollection struct {
	Name  string      `json:"name"`
	Cards []Flashcard `json:"cards"`
}

type Flashcard struct {
	FrontText string `json:"frontText"`
	BackText  string `json:"backText"`
}

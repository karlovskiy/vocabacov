package main

import (
	"encoding/json"
	"flag"
	"github.com/karlovskiy/vocabacov/internal/database"
	"github.com/karlovskiy/vocabacov/internal/translate"
	"log/slog"
	"os"
)

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		slog.Error("arg with path to file to import not found")
		os.Exit(1)
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		slog.Error("import file read error", "err", err)
		os.Exit(2)
	}
	var phrases map[string][]Phrase
	if err := json.Unmarshal(data, &phrases); err != nil {
		slog.Error("json unmarshall error", "err", err)
		os.Exit(3)
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

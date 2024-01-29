package main

import (
	"encoding/json"
	"fmt"
	"github.com/karlovskiy/vocabacov/internal/database"
	"log/slog"
	"os"
)

func main() {
	db, err := database.OpenDb(false)
	if err != nil {
		slog.Error("db open error", "err", err)
		os.Exit(1)
	}
	dbPhrases, err := database.LoadAllPhrases(db)
	if err != nil {
		slog.Error("phrases load error", "err", err)
		os.Exit(2)
	}
	phrases := make(map[string][]Phrase)
	for _, p := range dbPhrases {
		langPhrases, ok := phrases[p.Lang]
		if !ok {
			langPhrases = make([]Phrase, 0, 100)
		}
		phrases[p.Lang] = append(langPhrases, Phrase{
			Phrase:    p.Phrase,
			Translate: p.Translation,
		})
	}
	data, err := json.Marshal(phrases)
	if err != nil {
		slog.Error("json marshal error", "err", err)
		os.Exit(3)
	}
	fmt.Printf("exported phrases:\n%s\n", data)
}

type Phrase struct {
	Phrase    string `json:"p"`
	Translate string `json:"t"`
}

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/karlovskiy/vocabacov/internal/database"
	"github.com/karlovskiy/vocabacov/internal/translate"
	"log/slog"
	"os"
)

var format = flag.String("format", "internal", "data format: internal or anki")

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		slog.Error("arg with path to file to export not found")
		os.Exit(1)
	}
	if *format != "internal" && *format != "anki" {
		slog.Error("bad format", "format", *format)
		os.Exit(1)
	}
	db, err := database.OpenDb(false)
	if err != nil {
		slog.Error("db open error", "err", err)
		os.Exit(2)
	}
	dbPhrases, err := database.LoadAllPhrases(db)
	if err != nil {
		slog.Error("phrases load error", "err", err)
		os.Exit(2)
	}
	if *format == "internal" {
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
			slog.Error("json marshal error", "format", *format, "err", err)
			os.Exit(3)
		}
		if err := os.WriteFile(args[0], data, 0666); err != nil {
			slog.Error("write file error", "format", *format, "err", err)
			os.Exit(4)
		}
		fmt.Printf("%d phrases exported in %s format to:%s\n", len(dbPhrases), *format, args[0])
	} else {
		phrases := make(map[string][]translate.Phrase)
		for _, p := range dbPhrases {
			langPhrases, ok := phrases[p.Lang]
			if !ok {
				langPhrases = make([]translate.Phrase, 0, 100)
			}
			phrases[p.Lang] = append(langPhrases, p)
		}
		for lang, ph := range phrases {
			data := translate.ExportAnki(ph)
			file := args[0] + "_" + lang + ".txt"
			if err := os.WriteFile(file, data, 0666); err != nil {
				slog.Error("write file error", "format", *format, "err", err)
				os.Exit(4)
			}
			fmt.Printf("%d %s phrases exported in %s format to:%s\n", len(dbPhrases), lang, *format, file)
		}

	}
}

type Phrase struct {
	Phrase    string `json:"p"`
	Translate string `json:"t"`
}

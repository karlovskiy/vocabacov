package main

import (
	"github.com/karlovskiy/vocabacov"
	"log"
)

func main() {
	bot, err := vocabacov.NewBot()
	if err != nil {
		log.Fatalf("bot creation error: %v", err)
	}
	db, err := vocabacov.NewDb()
	if err != nil {
		log.Fatalf("db error: %v", err)
	}
	if err := bot.Start(db); err != nil {
		log.Fatalf("bot error: %v", err)
	}
}

package main

import (
	"log/slog"
	"os"

	botapi "github.com/karlovskiy/vocabacov/internal/bot"
	"github.com/karlovskiy/vocabacov/internal/database"
)

func main() {
	bot, err := botapi.NewBot()
	if err != nil {
		slog.Error("bot creation error", "err", err)
		os.Exit(1)
	}
	db, err := database.OpenDb(true)
	if err != nil {
		slog.Error("db open error", "err", err)
		os.Exit(2)
	}
	defer func() {
		if err := db.Close(); err != nil {
			slog.Info("db close error", "err", err)
		}
	}()
	bot.Start(db)
}

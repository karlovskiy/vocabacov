package bot

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/karlovskiy/vocabacov/internal/database"
	"github.com/karlovskiy/vocabacov/internal/translate"
)

const (
	EnvToken   = "VOCABACOV_TOKEN"
	EnvChannel = "VOCABACOV_CHANNELS"
	EnvDebug   = "VOCABACOV_DEBUG"
	EnvTimeout = "VOCABACOV_TIMEOUT"
)

func NewBot() (*Bot, error) {
	token := os.Getenv(EnvToken)
	if token == "" {
		return nil, fmt.Errorf("token not found in environment variable %s", EnvToken)
	}
	channelIdsValue := os.Getenv(EnvChannel)
	if channelIdsValue == "" {
		return nil, fmt.Errorf("channelId not found in environment variable %s", EnvChannel)
	}
	channels, err := ParseChannels(channelIdsValue)
	if err != nil {
		return nil, err
	}
	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	api.Debug = "true" == os.Getenv(EnvDebug)
	return &Bot{
		api:      api,
		channels: channels,
	}, nil
}

func (b *Bot) Start(db *sql.DB) {
	updateConfig := tgbotapi.NewUpdate(0)
	timeoutValue := os.Getenv(EnvTimeout)
	if timeoutValue != "" {
		t, _ := strconv.Atoi(timeoutValue)
		updateConfig.Timeout = t
	}
	if updateConfig.Timeout == 0 {
		updateConfig.Timeout = 30
	}
	updates := b.api.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		_, ok := b.channels[update.Message.Chat.ID]
		if !ok {
			continue
		}
		cmd, err := translate.FindCommand(update.Message.Text)
		if err != nil {
			slog.Error("find command error", "err", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "command has incorrect syntax")
			msg.ReplyToMessageID = update.Message.MessageID
			b.sendMessage(&msg)
			continue
		}
		if cmd != nil {
			if cmd.Name == "export" {
				lang := cmd.Args[0]
				phrases, err := database.LoadActivePhrases(db, lang)
				if err != nil {
					slog.Error("load phrases error", "lang", lang, "err", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "error load phrases")
					msg.ReplyToMessageID = update.Message.MessageID
					b.sendMessage(&msg)
					continue
				}
				if len(phrases) == 0 {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "no active phrases exist to export")
					msg.ReplyToMessageID = update.Message.MessageID
					b.sendMessage(&msg)
					continue
				}
				ankiData := translate.ExportAnki(phrases)
				docName := lang + "-" + time.Now().Format("20060102150405") + ".txt"
				doc := tgbotapi.NewDocument(update.Message.Chat.ID, tgbotapi.FileBytes{
					Name:  docName,
					Bytes: ankiData,
				})
				sent, err := b.api.Send(doc)
				if err != nil {
					slog.Error("error sending document", "chatId", doc.ChatID, "docName", docName, "err", err)
					continue
				}
				slog.Info("document sent", "chatId", doc.ChatID, "doc", doc, "sent", sent)
				if err := database.SetPhrasesStatus(db, lang, "ARCHIVED"); err != nil {
					slog.Error("archive phrases error", "lang", lang, "err", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "error archive phrases after export")
					msg.ReplyToMessageID = update.Message.MessageID
					b.sendMessage(&msg)
				}
			} else if cmd.Name == "reset" {
				lang := cmd.Args[0]
				if err := database.SetPhrasesStatus(db, lang, "ACTIVE"); err != nil {
					slog.Error("reset phrases error", "lang", lang, "err", err)
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "error reset phrases")
					msg.ReplyToMessageID = update.Message.MessageID
					b.sendMessage(&msg)
					continue
				}
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "done")
				msg.ReplyToMessageID = update.Message.MessageID
				b.sendMessage(&msg)
			} else {
				slog.Info("unknown command", "name", cmd.Name)
			}
			continue
		}
		phrase, err := translate.FindPhrase(update.Message.Text)
		if err != nil {
			slog.Error("find phase error", "err", err)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "phrase has incorrect syntax")
			msg.ReplyToMessageID = update.Message.MessageID
			b.sendMessage(&msg)
			continue
		}
		if phrase != nil {
			text := "done"
			if err := database.SavePhrase(db, phrase); err != nil {
				slog.Error("db save error", "err", err)
				text = "db save error"
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			msg.ReplyToMessageID = update.Message.MessageID
			b.sendMessage(&msg)
		}
	}
}

func ParseChannels(channelIdsValue string) (map[int64]struct{}, error) {
	channelIds := strings.Split(channelIdsValue, ",")
	channels := map[int64]struct{}{}
	for _, channelId := range channelIds {
		channelId = strings.TrimSpace(channelId)
		if channelId == "" {
			continue
		}
		id, err := strconv.ParseInt(channelId, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("channel %q int parsing error %w", channelId, err)
		}
		channels[id] = struct{}{}
	}
	if len(channels) == 0 {
		return nil, errors.New("channels not found")
	}
	return channels, nil
}

func (b *Bot) sendMessage(msg *tgbotapi.MessageConfig) {
	sent, err := b.api.Send(msg)
	if err != nil {
		slog.Error("error sending message", "chatId", msg.ChatID, "replyId", msg.ReplyToMessageID,
			"msg", msg, "err", err)
		return
	}
	slog.Info("message sent", "chatId", msg.ChatID, "replyId", msg.ReplyToMessageID, "msg", msg, "sent", sent)
}

type Bot struct {
	api      *tgbotapi.BotAPI
	channels map[int64]struct{}
}

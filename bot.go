package vocabacov

import (
	"database/sql"
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"os"
	"strconv"
	"strings"
)

const (
	EnvToken    = "VOCABACOV_TOKEN"
	EnvChannel  = "VOCABACOV_CHANNELS"
	EnvDebug    = "VOCABACOV_DEBUG"
	EnvTimeout  = "VOCABACOV_TIMEOUT"
	EnvDatabase = "VOCABACOV_DB_PATH"
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
	channels, err := parseChannels(channelIdsValue)
	if err != nil {
		return nil, err
	}
	return &Bot{
		token:    token,
		channels: channels,
	}, nil
}

func (b *Bot) Start(db *sql.DB) error {
	bot, err := tgbotapi.NewBotAPI(b.token)
	if err != nil {
		return err
	}
	bot.Debug = "true" == os.Getenv(EnvDebug)
	updateConfig := tgbotapi.NewUpdate(0)
	timeoutValue := os.Getenv(EnvTimeout)
	if timeoutValue != "" {
		t, _ := strconv.Atoi(timeoutValue)
		updateConfig.Timeout = t
	}
	if updateConfig.Timeout == 0 {
		updateConfig.Timeout = 30
	}
	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		_, ok := b.channels[update.Message.Chat.ID]
		if !ok {
			continue
		}
		lang, phrase := findPhrase(update.Message.Text)
		if lang != "" && phrase != "" {
			if err := savePhrase(db, lang, phrase); err != nil {
				return err
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "done")
			msg.ReplyToMessageID = update.Message.MessageID
			if _, err := bot.Send(msg); err != nil {
				return err
			}
		}
	}
	return nil
}

func parseChannels(channelIdsValue string) (map[int64]struct{}, error) {
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

func findPhrase(text string) (string, string) {
	if text != "" {
		text = strings.TrimLeft(text, "/")
		tokens := strings.SplitN(text, " ", 2)
		if len(tokens) > 1 {
			return tokens[0], tokens[1]
		}
	}
	return "", ""
}

type Bot struct {
	token    string
	channels map[int64]struct{}
}

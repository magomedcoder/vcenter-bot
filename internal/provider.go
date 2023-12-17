package internal

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

func NewBotAPI(conf *Config) *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	return bot
}

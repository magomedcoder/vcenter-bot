package internal

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Start(bot *tgbotapi.BotAPI) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := b.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message != nil {
			userId := update.FromChat().ID
			fmt.Println(userId)
		}
	}
}

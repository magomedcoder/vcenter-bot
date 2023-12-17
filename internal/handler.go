package internal

import (
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	BotAPI         *tgbotapi.BotAPI
	VCenterApiCall *VCenterApiCall
}

func NewBotHandler(botAPI *tgbotapi.BotAPI, vcenterApiCall *VCenterApiCall) *Bot {
	return &Bot{BotAPI: botAPI, VCenterApiCall: vcenterApiCall}
}

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		if update.Message == nil {
			continue
		}
		if update.Message.IsCommand() {
			switch update.Message.Command() {
			case "login":
				b.VCenterApiCall.session()
			}
		}

	}
}

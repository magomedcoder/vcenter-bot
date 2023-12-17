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
			case "vm":
				items, err := b.VCenterApiCall.getListVM()
				if err != nil {
					continue
				}
				if items != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)

					var keyboard [][]tgbotapi.InlineKeyboardButton
					for _, item := range items {
						keyboard = append(keyboard, tgbotapi.NewInlineKeyboardRow(
							tgbotapi.NewInlineKeyboardButtonData(item.Name, item.Name),
						))
					}

					msg.ReplyMarkup = tgbotapi.InlineKeyboardMarkup{
						InlineKeyboard: keyboard,
					}

					b.BotAPI.Send(msg)
				}
			}
		}

	}
}

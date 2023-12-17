package internal

import (
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
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
		if update.CallbackQuery != nil {
			data := update.CallbackQuery
			chatID := data.Message.Chat.ID
			messageID := data.Message.MessageID
			vm := strings.Split(data.Data, ":")
			item, err := b.VCenterApiCall.getVM(vm[1])
			if err != nil {
				continue
			}
			switch vm[0] {
			case "vm":
				var buttons []tgbotapi.InlineKeyboardButton
				if item.PowerState == "POWERED_ON" {
					buttons = append(buttons,
						tgbotapi.NewInlineKeyboardButtonData("Выключить", "vmReboot:"+vm[1]),
						tgbotapi.NewInlineKeyboardButtonData("Перезагрузить", "vmOff:"+vm[1]),
					)
				}
				if item.PowerState == "POWERED_OFF" {
					buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData("Включить", "vmOn:"+vm[1]))
				}
				msg := tgbotapi.NewMessage(data.From.ID, fmt.Sprintf("%s", item.Name))
				msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(tgbotapi.NewInlineKeyboardRow(buttons...))
				b.BotAPI.Send(msg)
				break
			case "vmOn":
				msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s\n", item.Name))
				b.BotAPI.Send(msg)
				break
			case "vmOff":
				msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s\n", item.Name))
				b.BotAPI.Send(msg)
				break
			case "vmReboot":
				msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s\n", item.Name))
				b.BotAPI.Send(msg)
				break
			}
		}
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
							tgbotapi.NewInlineKeyboardButtonData(item.Name, "vm:"+item.Id),
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

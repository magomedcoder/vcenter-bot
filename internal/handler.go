package internal

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"regexp"
	"strings"
)

type Bot struct {
	Conf           *Config
	BotAPI         *tgbotapi.BotAPI
	VCenterApiCall *VCenterApiCall
	Db             *sql.DB
}

func NewBotHandler(conf *Config, botAPI *tgbotapi.BotAPI, vcenterApiCall *VCenterApiCall, db *sql.DB) *Bot {
	return &Bot{Conf: conf, BotAPI: botAPI, VCenterApiCall: vcenterApiCall, Db: db}
}

func parseLoginPassword(input string) []string {
	// Используем регулярное выражение для извлечения логина и пароля
	re := regexp.MustCompile(`"([^"]+)":"([^"]+)"`)
	matches := re.FindStringSubmatch(input)

	if len(matches) != 3 {
		return nil
	}

	return matches
}

func (b *Bot) CallbackQuery(userId int64, data *tgbotapi.CallbackQuery) {
	chatID := data.Message.Chat.ID
	messageID := data.Message.MessageID
	vm := strings.Split(data.Data, ":")
	item, err := b.VCenterApiCall.getVM(userId, vm[1])
	if err != nil {
		return
	}
	switch vm[0] {
	case "vm":
		var buttons []tgbotapi.InlineKeyboardButton
		if item.PowerState == "POWERED_ON" {
			buttons = append(buttons,
				tgbotapi.NewInlineKeyboardButtonData("Выключить", "vmOff:"+vm[1]),
				tgbotapi.NewInlineKeyboardButtonData("Перезагрузить", "vmReboot:"+vm[1]),
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
		if b.VCenterApiCall.StartVM(userId, vm[1]) {
			msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s\n", item.Name))
			b.BotAPI.Send(msg)
		}
		break
	case "vmOff":
		if b.VCenterApiCall.StopVM(userId, vm[1]) {
			msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s\n", item.Name))
			b.BotAPI.Send(msg)
		}
		break
	case "vmReboot":
		if b.VCenterApiCall.RebootVM(userId, vm[1]) {
			msg := tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%s\n", item.Name))
			b.BotAPI.Send(msg)
		}
		break
	}
}

func (b *Bot) Command(userId int64, message *tgbotapi.Message) {
	switch message.Command() {
	case "vm":
		items, err := b.VCenterApiCall.getListVM(userId)
		if err != nil {
			return
		}
		if items != nil {
			msg := tgbotapi.NewMessage(message.Chat.ID, message.Text)
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

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		userId := update.FromChat().ID

		login := parseLoginPassword(update.Message.Text)

		if login != nil && login[1] != "" && login[2] != "" {
			_, err := b.Db.Exec("DELETE FROM users WHERE user_id = ?", userId)
			if err != nil {
				log.Println(err)
			}

			_, _err := b.Db.Exec("INSERT INTO users (user_id, username, password) VALUES (?, ?, ?)", userId, login[1], login[2])
			if _err != nil {
				log.Println(err)
			}
		}

		if update.CallbackQuery != nil {
			b.CallbackQuery(userId, update.CallbackQuery)
		}

		if update.Message == nil {
			continue
		}

		if update.Message.IsCommand() {
			b.Command(userId, update.Message)
		}
	}
}

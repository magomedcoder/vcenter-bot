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
	infoVM := fmt.Sprintf("VM: *%s*\n - CPU: *%d*\n - RAM: *%d*\n", item.Name, item.Cpu, item.Ram)
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
		msg := tgbotapi.NewMessage(data.From.ID, infoVM)
		msg.ParseMode = tgbotapi.ModeMarkdown
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
		break
	case "logout":
		_, err := b.Db.Exec("DELETE FROM users WHERE user_id = ?", userId)
		if err != nil {
			log.Println(err)
		}
		break
	}
}

var numericKeyboard = tgbotapi.NewReplyKeyboard(
	tgbotapi.NewKeyboardButtonRow(
		tgbotapi.NewKeyboardButton("/vm"),
		tgbotapi.NewKeyboardButton("/logout"),
	),
)

func (b *Bot) Start() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates := b.BotAPI.GetUpdatesChan(u)
	for update := range updates {
		userId := update.FromChat().ID

		if update.Message != nil {
			login := parseLoginPassword(update.Message.Text)
			if login != nil {
				_, err := b.Db.Exec("DELETE FROM users WHERE user_id = ?", userId)
				if err != nil {
					log.Println(err)
					continue
				}

				_, _err := b.Db.Exec("INSERT INTO users (user_id, username, password) VALUES (?, ?, ?)", userId, login[1], login[2])
				if _err != nil {
					log.Println(err)
					continue
				}
				if b.VCenterApiCall.session(userId) == false {
					_, err := b.Db.Exec("DELETE FROM users WHERE user_id = ?", userId)
					if err != nil {
						log.Println(err)
						continue
					}
					msg := tgbotapi.NewMessage(userId, "Ошибка авторизации.")
					msg.ReplyMarkup = numericKeyboard
					b.BotAPI.Send(msg)

					continue
				}
				msg := tgbotapi.NewMessage(userId, "Добро пожаловать.")
				msg.ReplyMarkup = numericKeyboard
				b.BotAPI.Send(msg)
			} else {

				var sessionId string
				err := b.Db.QueryRow("SELECT user_id FROM users WHERE user_id = ?", userId).Scan(&sessionId)
				if err != nil {

				}

				if sessionId == "" {
					msg := tgbotapi.NewMessage(userId, "Для входа в систему введите ваш логин и пароль в следующем формате:\n *\"логин\":\"пароль\"*")
					msg.ReplyMarkup = numericKeyboard
					msg.ParseMode = tgbotapi.ModeMarkdown
					b.BotAPI.Send(msg)
					continue
				}
			}

			if update.Message.IsCommand() {

				msg := tgbotapi.NewMessage(userId, "Пожалуйста, подождите.")
				msg.ReplyMarkup = numericKeyboard
				b.BotAPI.Send(msg)

				b.Command(userId, update.Message)
			}

		} else if update.CallbackQuery != nil {

			msg := tgbotapi.NewMessage(userId, "Пожалуйста, подождите.")
			msg.ReplyMarkup = numericKeyboard
			b.BotAPI.Send(msg)

			b.CallbackQuery(userId, update.CallbackQuery)
		}
	}
}

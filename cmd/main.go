package main

import (
	"flag"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"vcenter-bot/internal"
)

func main() {
	env := flag.String("c", "./config.yaml", "Требуется конфиг")
	flag.Parse()

	conf, err := internal.ReadConfig(*env)
	if err != nil {
		log.Fatal(err)
	}

	bot, err := tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	internal.Start(bot)
}

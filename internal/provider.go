package internal

import (
	"database/sql"
	"fmt"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/proxy"
	"log"
	"net/http"
)

func NewBotAPI(conf *Config) *tgbotapi.BotAPI {
	if conf.Telegram.ProxySocks5 != nil {
		dialSocks5, err := proxy.SOCKS5(
			"tcp",
			fmt.Sprintf("%s:%s", conf.Telegram.ProxySocks5.Host, conf.Telegram.ProxySocks5.Port),
			&proxy.Auth{
				User:     conf.Telegram.ProxySocks5.Username,
				Password: conf.Telegram.ProxySocks5.Password,
			},
			proxy.Direct)
		if err != nil {
			log.Panicf("SOCKS5-%s", err)
		}

		transport := &http.Transport{Dial: dialSocks5.Dial}
		httpClient := &http.Client{}
		httpClient.Transport = transport

		bot, err := tgbotapi.NewBotAPIWithClient(conf.Telegram.Token, "", httpClient)
		if err != nil {
			log.Panicf("NewBotAPIWithClient-%s", err)
		}

		bot.Debug = false
		return bot
	} else {
		bot, err := tgbotapi.NewBotAPI(conf.Telegram.Token)
		if err != nil {
			log.Panicf("NewBotAPI-%s", err)
		}

		bot.Debug = false
		return bot
	}
}

func NewDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "/var/lib/vcenter-bot/vcenter-bot.db")
	if err != nil {
		log.Panicf("NewDatabase-sql-%s", err)
	}

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER,
			username VARCHAR,
			password VARCHAR,
			session_id VARCHAR
		)
	`)
	if err != nil {
		log.Panicf("NewDatabase-%s", err)
	}

	return db
}

package internal

import (
	"database/sql"
	"github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/mattn/go-sqlite3"
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

func NewDatabase() *sql.DB {
	db, err := sql.Open("sqlite3", "vcenter-bot.db")
	if err != nil {
		log.Println(err)
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
		log.Println(err)
	}

	return db
}

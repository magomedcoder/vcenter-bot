package main

import (
	"flag"
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

	app := Initialize(conf)
	app.Bot.Start()
}

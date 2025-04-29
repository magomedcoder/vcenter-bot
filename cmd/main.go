package main

import (
	"flag"
	"io"
	"log"
	"os"
	"vcenter-bot/internal"
)

func main() {
	const logPath = "/var/log/vcenter-bot.log"
	logFile, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("не удалось открыть файл лога %s: %v", logPath, err)
	}
	defer logFile.Close()

	mw := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(mw)

	log.SetPrefix("vcenter-bot: ")
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile)

	env := flag.String("c", "/etc/vcenter-bot/config.yaml", "путь до конфига")
	flag.Parse()

	conf, err := internal.ReadConfig(*env)
	if err != nil {
		log.Panicf("main-ReadConfig: %v", err)
	}

	app := Initialize(conf)
	app.Bot.Start()
}

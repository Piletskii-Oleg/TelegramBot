package main

import (
	"flag"
	"log"
	tgBot "telegramBot/bot"
	"telegramBot/consumer"
	"telegramBot/events/telegram"
	"telegramBot/storage/files"
)

const host = "api.telegram.org"

const storage = "storage"

const batchSize = 100

func main() {
	bot := tgBot.NewBot(token(), host)

	processor := telegram.NewTelegramProcessor(bot, files.NewFileStorage(storage))

	log.Print("service started")

	eventConsumer := consumer.NewEventConsumer(processor, processor, batchSize)

	if err := eventConsumer.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}

func token() string {
	token := flag.String("token", "", "Telegram API token")

	flag.Parse()
	if *token == "" {
		log.Fatal("No token provided")
	}

	return *token
}

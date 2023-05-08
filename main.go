package main

import (
	"flag"
	"log"
	tgBot "telegramBot/bot"
	"telegramBot/consumer"
	telegram_storage "telegramBot/storage"
	"telegramBot/storage/files"
)

const host = "api.telegram.org"
const storageFolder = "storage"
const batchSize = 100

func main() {
	bot := tgBot.NewBot(token(), host)

	processor := telegram_storage.NewStorageProcessor(bot, files.NewFileStorage(storageFolder))

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

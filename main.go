package main

import (
	"flag"
	"log"
	tgBot "telegramBot/bot"
	"telegramBot/consumer"
	sunny_day "telegramBot/events/telegram/sunny-day"
	weather2 "telegramBot/weather"
)

const host = "api.telegram.org"
const storageFolder = "storage"
const batchSize = 100

func main() {
	tgToken := token("tg", "Telegram API token")
	bot := tgBot.NewBot(tgToken, host)

	//processor := telegram_storage.NewStorageProcessor(bot, files.NewFileStorage(storageFolder))

	weather := weather2.NewClient(weatherToken)
	weatherProcessor := sunny_day.NewWeatherProcessor(bot, *weather)
	log.Print("service started")

	eventConsumerWeather := consumer.NewEventConsumer(weatherProcessor, weatherProcessor, batchSize)
	//eventConsumer := consumer.NewEventConsumer(processor, processor, batchSize)

	if err := eventConsumerWeather.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}

func token(name, usage string) string {
	token := flag.String(name, "", usage)

	flag.Parse()
	if *token == "" {
		log.Fatal("No token provided")
	}

	return *token
}

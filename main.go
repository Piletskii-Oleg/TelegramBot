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
	info := map[string]string{
		"tg":  "Telegram API token",
		"omw": "Open weather map API token",
	}
	tokens := processTokens(info)
	tgToken := tokens["tg"]
	weatherToken := tokens["omw"]

	bot := tgBot.NewBot(*tgToken, host)

	//processor := telegram_storage.NewStorageProcessor(bot, files.NewFileStorage(storageFolder))

	weather := weather2.NewClient(*weatherToken)
	weatherProcessor := sunny_day.NewWeatherProcessor(bot, *weather)
	log.Print("service started")

	eventConsumerWeather := consumer.NewEventConsumer(weatherProcessor, weatherProcessor, batchSize)
	//eventConsumer := consumer.NewEventConsumer(processor, processor, batchSize)

	if err := eventConsumerWeather.Start(); err != nil {
		log.Fatal("service stopped", err)
	}
}

func processTokens(info map[string]string) map[string]*string {
	var tokens = make(map[string]*string, len(info))
	for name, usage := range info {
		tokens[name] = flag.String(name, "", usage)
	}

	flag.Parse()
	for key, token := range tokens {
		if *token == "" {
			log.Fatal("token is not provided: ", key, " - ", info[key])
		}
	}

	return tokens
}

type Token struct {
	name, usage string
}

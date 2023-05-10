package main

import (
	"flag"
	"log"
	tgBot "telegramBot/bot"
	"telegramBot/consumer"
	sunny_day "telegramBot/events/telegram/sunny-day"
	"telegramBot/files"
	"telegramBot/weather"
)

const host = "api.telegram.org"
const batchSize = 100

func main() {
	info := map[string]string{
		"tg":  "Telegram API token",
		"owm": "Open weather map API token",
	}
	tokens := processTokens(info)

	tgToken := tokens["tg"]
	weatherToken := tokens["owm"]

	bot := tgBot.NewBot(*tgToken, host)

	fileStorage := files.NewFileStorage("location_files")

	weatherClient := weather.NewClient(*weatherToken)
	weatherProcessor := sunny_day.NewWeatherProcessor(bot, *weatherClient, *fileStorage)

	eventConsumerWeather := consumer.NewEventConsumer(weatherProcessor, weatherProcessor, batchSize)
	if err := eventConsumerWeather.Start(); err != nil {
		log.Fatal("service cannot be started", err)
	}
	log.Print("service started")
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

package sunny_day

import (
	"log"
	"strings"
	"telegramBot/events"
	"telegramBot/weather"
	"time"
)

const (
	hello    = "Hello"
	help     = "help"
	location = "location"
)

var (
	defaultButtons  = []string{"location", "hello", "help"}
	locationButtons = []string{"Moscow", "Saint Petersburg", "Novosibirsk"}
)

func (p *Processor) DoCmd(event events.Event, meta Meta) error {
	text := event.Text
	text = strings.TrimSpace(text)

	username := meta.Username
	chatID := meta.ChatID

	log.Printf("got command '%s' from '%s'", text, username)

	switch text {
	case location:
		return p.requestLocation(chatID)
	case help:
		return p.sendHelp(chatID)
	case hello:
		return p.sendHello(chatID)
	default:
		return p.Processor.Bot.SendMessage(chatID, msgUnknownCommand, defaultButtons)
	}
}

func (p *Processor) requestLocation(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgEnterLocation, locationButtons)
}

func (p *Processor) sendLocationInfo(response weather.Response, chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, time.Unix(response.Time, 0).String(), defaultButtons)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgHelp, defaultButtons)
}

func (p *Processor) sendHello(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgHello, defaultButtons)
}

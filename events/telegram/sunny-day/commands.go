package sunny_day

import (
	"log"
	"strings"
	"telegramBot/events"
	"telegramBot/storage"
)

const (
	hello        = "/start"
	help         = "help"
	getLocation  = "get location info"
	saveLocation = "save location"
	savedInfo    = "show saved"
)

var (
	defaultButtons  = []string{help, getLocation, saveLocation, savedInfo}
	locationButtons = []string{"Moscow", "Saint Petersburg", "Novosibirsk"}
)

func (p *Processor) DoCmd(event events.Event, meta Meta) error {
	text := event.Text
	text = strings.TrimSpace(text)

	username := meta.Username
	chatID := meta.ChatID

	log.Printf("got command '%s' from '%s'", text, username)

	switch text {
	case getLocation:
		return p.requestLocation(chatID)
	case help:
		return p.sendHelp(chatID)
	case hello:
		return p.sendHello(chatID)
	case saveLocation:
		return p.requestLocation(chatID)
	case savedInfo:
		return p.sendSavedLocationInfo(chatID, username)
	default:
		return p.Processor.Bot.SendMessage(chatID, msgUnknownCommand, defaultButtons)
	}
}

func (p *Processor) requestLocation(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgEnterLocation, locationButtons)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgHelp, defaultButtons)
}

func (p *Processor) sendHello(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgHello, defaultButtons)
}

func (p *Processor) confirmSaveLocation(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgSavedLocation, defaultButtons)
}

func (p *Processor) sendSavedLocationInfo(chatID int, username string) error {
	loc, err := p.Storage.Location(username)
	if err == storage.ErrNoSavedLocation {
		return p.Processor.Bot.SendMessage(chatID, msgNoSavedLocation, defaultButtons)
	} else if err != nil {
		return err
	}

	response, err := p.weather.MakeRequest(*loc)
	if err != nil {
		return err
	}

	return p.sendLocationInfo(response, chatID)
}

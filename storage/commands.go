package storage

import (
	errors2 "errors"
	"log"
	url2 "net/url"
	"strings"
	"telegramBot/bot"
	"telegramBot/bot/errors"
	"telegramBot/events"
	"telegramBot/events/telegram"
	"time"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) DoCmd(event events.Event, meta telegram.Meta) error {
	text := event.Text
	text = strings.TrimSpace(text)

	username := meta.Username
	chatID := meta.ChatID

	log.Printf("got command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.Processor.Bot.SendMessage(chatID, msgUnknownCommand)
	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() {
		err = errors.WrapIfError("can't do command: save page", err)
	}()

	sendMsg := NewMessageSender(chatID, p.Bot)

	page := &Page{
		URL:      pageURL,
		Username: username,
		Created:  time.Now(),
	}

	isExists, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isExists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() {
		err = errors.WrapIfError("can't do command: save page", err)
	}()

	sendMsg := NewMessageSender(chatID, p.Processor.Bot)

	page, err := p.storage.PickRandom(username)
	if err != nil && !errors2.Is(err, ErrNoSavedPages) {
		return err
	}
	if errors2.Is(err, ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgHelp)
}

func (p *Processor) sendHello(chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, msgHello)
}

func NewMessageSender(chatID int, bot *bot.Bot) func(string) error {
	return func(message string) error {
		return bot.SendMessage(chatID, message)
	}
}

func isAddCmd(text string) bool {
	return isUrl(text)
}

func isUrl(text string) bool {
	url, err := url2.Parse(text)
	return err == nil && url.Host != ""
}

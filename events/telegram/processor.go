package telegram

import (
	errors2 "errors"
	"telegramBot/bot"
	"telegramBot/bot/errors"
	"telegramBot/events"
	"telegramBot/storage"
)

type Processor struct {
	bot     *bot.Bot
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors2.New("unknown event type")
	ErrUnknownMetaType  = errors2.New("unknown meta type")
)

func NewTelegramProcessor(bot *bot.Bot, storage storage.Storage) *Processor {
	return &Processor{
		bot:     bot,
		storage: storage}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.bot.Updates(p.offset, limit)
	if err != nil {
		return nil, errors.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	events := make([]events.Event, 0, len(updates))
	for _, upd := range updates {
		events = append(events, event(upd))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return events, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.Wrap("unknown event type", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errors.Wrap("can't process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return errors.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (Meta, error) {
	result, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.Wrap("unable to get meta", ErrUnknownMetaType)
	}

	return result, nil
}

func event(upd bot.Update) events.Event {
	updateType := fetchType(upd)
	res := events.Event{
		Type: updateType,
		Text: fetchText(upd),
	}

	if updateType == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

func fetchText(upd bot.Update) string {
	if upd.Message == nil {
		return ""
	}
	return upd.Message.Text
}

func fetchType(upd bot.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}

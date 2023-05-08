package telegram

import (
	errors2 "errors"
	"telegramBot/bot"
	"telegramBot/bot/errors"
	"telegramBot/events"
)

type CmdProcessor interface {
	Processor
	DoCmd(event events.Event, meta Meta) error
}

type Processor struct {
	Bot    *bot.Bot
	Offset int
}

type Meta struct {
	ChatID   int
	Username string
}

var (
	ErrUnknownEventType = errors2.New("unknown event type")
	ErrUnknownMetaType  = errors2.New("unknown meta type")
)

func NewTelegramProcessor(bot *bot.Bot) *Processor {
	return &Processor{
		Bot: bot,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.Bot.Updates(p.Offset, limit)
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

	p.Offset = updates[len(updates)-1].ID + 1

	return events, nil
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

package storage

import (
	"telegramBot/bot"
	"telegramBot/bot/errors"
	"telegramBot/events"
	"telegramBot/events/telegram"
)

type Processor struct {
	*telegram.Processor
	storage Storage
}

func NewStorageProcessor(bot *bot.Bot, storage Storage) *Processor {
	return &Processor{Processor: telegram.NewTelegramProcessor(bot), storage: storage}
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return errors.Wrap("unknown event type", telegram.ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errors.Wrap("can't process message", err)
	}

	if err := p.DoCmd(event, meta); err != nil {
		return errors.Wrap("can't process message", err)
	}

	return nil
}

func meta(event events.Event) (telegram.Meta, error) {
	result, ok := event.Meta.(telegram.Meta)
	if !ok {
		return telegram.Meta{}, errors.Wrap("unable to get meta", telegram.ErrUnknownMetaType)
	}

	return result, nil
}

package sunny_day

import (
	"telegramBot/bot"
	"telegramBot/bot/errors"
	"telegramBot/events"
	"telegramBot/events/telegram"
	"telegramBot/storage/files"
	"telegramBot/weather"
)

type Processor struct {
	*telegram.Processor
	weather       weather.Client
	States        States
	SavedLocation string

	Storage files.FileStorage
}

type States struct {
	AskedForLocation           bool
	AskedForPersistentLocation bool
}

type Meta struct {
	ChatID   int
	Username string
}

func NewWeatherProcessor(bot *bot.Bot, weather weather.Client, storage files.FileStorage) *Processor {
	return &Processor{
		Processor: telegram.NewTelegramProcessor(bot),
		weather:   weather,
		States: States{
			AskedForLocation:           false,
			AskedForPersistentLocation: false,
		},
		Storage: storage,
	}
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		if p.States.AskedForLocation {
			p.States.AskedForLocation = false
			return p.processLocation(event)
		} else if p.States.AskedForPersistentLocation {
			p.States.AskedForPersistentLocation = false
			return p.saveLocation(event)
		} else {
			return p.processMessage(event)
		}
	case events.AskForLocation:
		p.States.AskedForLocation = true
		return p.processMessage(event)
	case events.ReceiveLocation:
		return p.processLocation(event)
	case events.AskForPersistentLocation:
		p.States.AskedForPersistentLocation = true
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

func (p *Processor) processLocation(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errors.Wrap("can't process message", err)
	}

	locationName := event.Text

	response, err := p.weather.MakeRequest(locationName)
	if err != nil {
		return errors.Wrap("can't make weather request", err)
	}

	if err := p.sendLocationInfo(response, meta.ChatID); err != nil {
		return errors.Wrap("can't process location", err)
	}

	return err
}

func (p *Processor) sendLocationInfo(response *weather.Response, chatID int) error {
	return p.Processor.Bot.SendMessage(chatID, locationInfo(response), defaultButtons)
}

func (p *Processor) saveLocation(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return errors.Wrap("can't process message", err)
	}

	if err := p.Storage.SaveLocation(meta.Username, event.Text); err != nil {
		return errors.Wrap("can't save location", err)
	}

	if err := p.confirmSaveLocation(meta.ChatID); err != nil {
		return errors.Wrap("can't process location", err)
	}

	return err
}

func meta(event events.Event) (Meta, error) {
	result, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, errors.Wrap("unable to get meta", telegram.ErrUnknownMetaType)
	}

	return result, nil
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.Bot.Updates(p.Offset, limit)
	if err != nil {
		return nil, errors.Wrap("can't get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	receivedEvents := make([]events.Event, 0, len(updates))
	for _, upd := range updates {
		receivedEvents = append(receivedEvents, event(upd))
	}

	p.Offset = updates[len(updates)-1].ID + 1

	return receivedEvents, nil
}

func event(upd bot.Update) events.Event {
	updateType := fetchType(upd)
	res := events.Event{
		Type: updateType,
		Text: fetchText(upd),
	}

	if updateType != events.Unknown {
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
	} else if upd.Message.Text == getLocation {
		return events.AskForLocation
	} else if upd.Message.Text == saveLocation {
		return events.AskForPersistentLocation
	}

	return events.Message
}

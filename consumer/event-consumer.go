package consumer

import (
	"log"
	"telegramBot/events"
	"time"
)

type EventConsumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func NewEventConsumer(fetcher events.Fetcher, processor events.Processor, batchSize int) *EventConsumer {
	return &EventConsumer{fetcher: fetcher, processor: processor, batchSize: batchSize}
}

func (c *EventConsumer) Start() error {
	for {
		receivedEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err.Error())
			continue
		}

		if len(receivedEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}

		if err := c.handleEvents(receivedEvents); err != nil {
			log.Print(err)

			continue
		}

	}
}

func (c *EventConsumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("got new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("can't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}

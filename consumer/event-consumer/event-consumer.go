package event_consumer

import (
	"TelegramBot/events"
	"log"
	"time"
)

type Consumer struct {
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int) *Consumer {
	return &Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (consumer *Consumer) Start() error {
	for {
		gotEvents, err := consumer.fetcher.Fetch(consumer.batchSize)
		if err != nil {
			log.Println("[ERROR] consumer", err.Error())
			continue
		}
		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)
			continue
		}
		if err := consumer.handleEvents(gotEvents); err != nil {
			log.Println("[ERROR] consumer", err.Error())
			continue
		}
	}
}

func (consumer *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("event: %#v", event)
		if err := consumer.processor.Process(event); err != nil {
			log.Println("[ERROR] processor", err.Error())
		}

	}
	return nil
}

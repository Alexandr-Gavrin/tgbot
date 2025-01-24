package main

import (
	tgClient "TelegramBot/clients/telegram"
	event_consumer "TelegramBot/consumer/event-consumer"
	"TelegramBot/events/telegram"
	"TelegramBot/storage/files"
	"flag"
	"log"
)

const (
	tgBotHost = "api.telegram.org"
	batchSize = 100
)

var (
	storagePath = "storage"
)

func main() {

	evProcessor := telegram.New(tgClient.New(tgBotHost, mustToken()), files.New(storagePath))
	log.Printf("Service started")
	consumer := event_consumer.New(evProcessor, evProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal(err)
	}

	// fetcher := fetcher.New()

	// processor := processor.New()

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String(
		"token-bot-token",
		"",
		"Telegram Bot Token",
	)
	flag.Parse()
	if *token == "" {
		log.Fatal("Token is required")
	}
	return *token
}

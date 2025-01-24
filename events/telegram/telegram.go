package telegram

import (
	"TelegramBot/clients/telegram"
	"TelegramBot/events"
	"TelegramBot/lib/e"
	"TelegramBot/storage"
	"errors"
)

var (
	ErrUnknownEventType = errors.New("Unknown event type")
	ErrUnknownMetaType  = errors.New("Unknown meta type")
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}

type Meta struct {
	ChatID   int
	Username string
}

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{tg: client, offset: 0, storage: storage}
}

func fetchText(upd telegram.Update) string {
	if upd.Message != nil {
		return upd.Message.Text
	}
	return ""
}

func fetchType(upd telegram.Update) events.Type {
	if upd.Message != nil {
		return events.Message
	}
	return events.Unknown
}

func event(upd telegram.Update) events.Event {
	t := fetchType(upd)
	res := events.Event{
		Type: t,
		Text: fetchText(upd),
	}
	if t == events.Message {
		res.Meta = Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}
	return res

}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	update, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, e.Wrap("Cant get events", err)
	}
	if len(update) == 0 {
		return nil, nil
	}
	res := make([]events.Event, 0, len(update))
	for _, u := range update {
		res = append(res, event(u))
	}
	p.offset = update[len(update)-1].ID + 1
	return res, nil
}

func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		e.Wrap("Cant process message", err)
	}
	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("cant process msg", err)
	}
	return nil

}

func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("Cant get meta", ErrUnknownEventType)
	}
	return res, nil
}

func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("Can't process message", ErrUnknownEventType)
	}
}

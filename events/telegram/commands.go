package telegram

import (
	"TelegramBot/clients/telegram"
	"TelegramBot/lib/e"
	"TelegramBot/storage"
	"errors"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"
	HelpCmd  = "/help"
	StartCmd = "/start"
)

func (p *Processor) doCmd(text string, chatID int, username string) error {
	text = strings.TrimSpace(text)
	log.Printf("Got new command %s from user %s\n", text, chatID)
	// add page : http://..
	// rnd page: /rnd
	// help: /help
	// start: /start

	if isAddCmd(text) {
		p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, MsgUnknown)

	}
}

func (p *Processor) savePage(chatID int, pageURL string, username string) (err error) {
	defer func() { err = e.Wrap("Cant do this command: save page", err) }()
	sendMsg := NewMessageSender(chatID, p.tg)
	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}
	isEx, err := p.storage.IsExists(page)
	if err != nil {
		return err
	}
	if isEx {
		return sendMsg(MsgAlreadyExists)
		//p.tg.SendMessage(chatID, MsgAlreadyExists)
	}
	if err := p.storage.Save(page); err != nil {
		return err
	}
	if err := p.tg.SendMessage(chatID, MsgSaved); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendRandom(chatID int, username string) (err error) {
	defer func() { err = e.Wrap("Cant do this command: pick random", err) }()
	page, err := p.storage.PickRandom(username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return p.tg.SendMessage(chatID, MsgNoSavedPage)
	}
	if err := p.tg.SendMessage(chatID, page.URL); err != nil {
		return err
	}
	return p.storage.Remove(page)
}

func (p *Processor) sendHelp(chatID int) error {
	if err := p.tg.SendMessage(chatID, MsgHelp); err != nil {
		return err
	}
	return nil
}

func (p *Processor) sendHello(chatID int) error {
	if err := p.tg.SendMessage(chatID, MsgHello); err != nil {
		return err
	}
	return nil
}

func NewMessageSender(chatID int, tg *telegram.Client) func(string) error {
	return func(s string) error {
		return tg.SendMessage(chatID, s)
	}
}

func isAddCmd(text string) bool {
	return isURL(text)
}

func isURL(text string) bool {
	u, err := url.Parse(text)
	return err == nil && u.Host != ""
}

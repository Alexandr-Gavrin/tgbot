package storage

import (
	"TelegramBot/lib/e"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
)

var ErrNoSavedPages = errors.New("No saved pages")

type Storage interface {
	Save(p *Page) error
	PickRandom(userName string) (*Page, error)
	Remove(p *Page) error
	IsExists(p *Page) (bool, error)
}

type Page struct {
	URL      string
	UserName string
	// Created time.Time
}

func (p Page) Hash() (string, error) {
	h := sha1.New()
	if _, err := io.WriteString(h, p.URL); err != nil {
		return "", e.Wrap("cant make hash", err)
	}
	if _, err := io.WriteString(h, p.UserName); err != nil {
		return "", e.Wrap("cant make hash", err)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

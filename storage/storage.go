package storage

import (
	"crypto/sha256"
	errors2 "errors"
	"fmt"
	"io"
	"telegramBot/bot/errors"
	"time"
)

var ErrNoSavedPages = errors2.New("no saved pages")

type Storage interface {
	Save(page *Page) error
	PickRandom(username string) (*Page, error)
	Remove(page *Page) error
	IsExists(page *Page) (bool, error)
}

type Page struct {
	URL      string
	Username string
	Created  time.Time
}

func (page *Page) Hash() (string, error) {
	hash := sha256.New()

	if _, err := io.WriteString(hash, page.URL); err != nil {
		return "", errors.Wrap("unable to calculate hash", err)
	}

	if _, err := io.WriteString(hash, page.Username); err != nil {
		return "", errors.Wrap("unable to calculate hash", err)
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

package bs

import (
	"fmt"
	"net/url"

	"github.com/dgraph-io/badger/v3"
	"github.com/gocolly/colly/v2/storage"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var _ storage.Storage = (*BadgerStorage)(nil)

type BadgerStorage struct {
	DB *badger.DB
}

func (b *BadgerStorage) Init() error {
	if b.DB == nil {
		return errors.New("badger db is not initialized")
	}
	return nil
}

func (b *BadgerStorage) requestIDKey(requestID uint64) []byte {
	return []byte(fmt.Sprintf("RequestID:%d", requestID))
}

func (b *BadgerStorage) cookieKey(u *url.URL) []byte {
	return []byte(fmt.Sprintf("Cookie:%s", u.String()))
}

func (b *BadgerStorage) Visited(requestID uint64) error {
	return b.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(b.requestIDKey(requestID), []byte{})
	})
}

func (b *BadgerStorage) IsVisited(requestID uint64) (out bool, err error) {
	err = b.DB.View(func(txn *badger.Txn) error {
		_, err = txn.Get(b.requestIDKey(requestID))
		if err == badger.ErrKeyNotFound {
			err = nil
		} else if err == nil {
			out = true
		}
		return err
	})
	return
}

func (b *BadgerStorage) Cookies(u *url.URL) (out string) {
	err := b.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get(b.cookieKey(u))
		if err == nil {
			err = item.Value(func(val []byte) error {
				out = string(val)
				return nil
			})
		}
		return err
	})
	if err != nil {
		log.Err(err).Msg("failed to get cookies")
	}
	return
}

func (b *BadgerStorage) SetCookies(u *url.URL, cookies string) {
	err := b.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(b.cookieKey(u), []byte(cookies))
	})
	if err != nil {
		log.Err(err).Msg("failed to set cookies")
	}
}

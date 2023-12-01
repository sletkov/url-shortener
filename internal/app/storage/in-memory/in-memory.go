package inmemory

import (
	"urlshortener/internal/app/storage"

	"github.com/patrickmn/go-cache"
)

type Storage struct {
	Db *cache.Cache
}

func New() (*Storage, error) {
	db := cache.New(cache.NoExpiration, cache.NoExpiration)

	return &Storage{Db: db}, nil
}

func (s *Storage) SaveURL(url string, alias string) error {
	err := s.Db.Add(alias, url, cache.NoExpiration)
	if err != nil {
		return storage.ErrURLExists
	}
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	url, ok := s.Db.Get(alias)
	if !ok {
		return "", storage.ErrURLNotFound
	}
	return url.(string), nil
}

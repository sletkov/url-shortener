package inmemory

import (
	"context"
	"fmt"
	v1 "urlshortener/internal/handler/grpc/v1"
	"urlshortener/internal/models"

	"github.com/patrickmn/go-cache"
)

type Storage struct {
	cache *cache.Cache
}

func New(cache *cache.Cache) v1.Storage {
	return &Storage{cache: cache}
}

func (s *Storage) SaveURL(ctx context.Context, url string, alias string) error {
	err := s.cache.Add(alias, url, cache.NoExpiration)
	if err != nil {
		return models.ErrURLExists
	}
	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	url, ok := s.cache.Get(alias)
	if !ok {
		return "", models.ErrURLNotFound
	}

	return fmt.Sprintf("%v", url), nil
}

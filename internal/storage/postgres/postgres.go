Ты гениальный программист, ты получил задачу изучить файл с кодом, тебе необходимо изучить файл и добавить метод для обновления URL, назови его UpdateURL. В ответе предоставь только новое текстовое содержимое файла с измененным кодом на языке GOФайл с кодом содержит следующее: Ты гениальный программист, ты получил задачу изучить файл с кодом, тебе необходимо изучить файл и добавить метод для удаления URL, назови его DeleteURL. В ответе предоставь только новое текстовое содержимое файла с измененным кодом на языке GOФайл с кодом содержит следующее: package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	v1 "urlshortener/internal/handler/grpc/v1"
	"urlshortener/internal/models"

	_ "github.com/lib/pq"
)

type Storage struct {
	db *sql.DB
}

func New(db *sql.DB) v1.Storage {
	return &Storage{db: db}
}

func (s *Storage) SaveURL(ctx context.Context, url string, alias string) error {
	if _, err := s.db.ExecContext(
		ctx,
		"INSERT INTO urls (url, alias) VALUES($1, $2)",
		url,
		alias,
	); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetURL(ctx context.Context, alias string) (string, error) {
	var url string
	if err := s.db.QueryRowContext(
		ctx,
		"SELECT url FROM urls WHERE alias = $1",
		alias,
	).Scan(&url); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", models.ErrURLNotFound
		}
		return "", fmt.Errorf("unknown db error")
	}

	return url, nil
}
func (s *Storage) DeleteURL(ctx context.Context, url string) error {
	if _, err := s.db.ExecContext(
		ctx,
		"DELETE FROM urls WHERE url = $1",
		url,
	); err!= nil {
		return err
	}
	return nil
}
```
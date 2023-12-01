package postgres

import (
	"database/sql"
	"fmt"
	"urlshortener/internal/app/storage"

	_ "github.com/lib/pq"
)

type Storage struct {
	Db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("postgres", storagePath)
	if err != nil {
		return nil, fmt.Errorf("can't open connection to postgres: %w", err)
	}

	return &Storage{Db: db}, nil
}

func (s *Storage) SaveURL(url string, alias string) error {
	if _, err := s.Db.Exec("INSERT INTO urls (url, alias) VALUES($1, $2)",
		url,
		alias,
	); err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	var url string
	if err := s.Db.QueryRow(
		"SELECT url FROM urls WHERE alias = $1",
		alias,
	).Scan(&url); err != nil {
		if err == sql.ErrNoRows {
			return "", storage.ErrURLNotFound
		}
	}

	return url, nil
}

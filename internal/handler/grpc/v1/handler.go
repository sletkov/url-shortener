package v1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	api "urlshortener/api/proto"
	random "urlshortener/internal/pkg"
)

type Storage interface {
	SaveURL(ctx context.Context, url string, alias string) error
	GetURL(ctx context.Context, alias string) (string, error)
}

type GRPCServer struct {
	api.UnimplementedURLShortenerServer
	Storage Storage
	Logger  *slog.Logger
}

func New(logger *slog.Logger, storage Storage) *GRPCServer {
	return &GRPCServer{
		Logger:  logger,
		Storage: storage,
	}
}

func (s *GRPCServer) SaveURL(ctx context.Context, req *api.SaveURLRequest) (*api.SaveURLResponse, error) {
	s.Logger.InfoContext(ctx, fmt.Sprintf("saving url %s", req.GetUrl()))
	alias := random.MakeRandomString(10)

	err := validateUrl(req.GetUrl())
	if err != nil {
		s.Logger.ErrorContext(ctx, "invalid url")
		return nil, err
	}

	err = s.Storage.SaveURL(ctx, req.GetUrl(), alias)
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to save url")
		return nil, err
	}

	s.Logger.InfoContext(ctx, "url saved")

	return &api.SaveURLResponse{
		Alias: alias,
	}, nil
}

func (s *GRPCServer) GetURL(ctx context.Context, req *api.GetURLRequest) (*api.GetURLResponse, error) {
	s.Logger.InfoContext(ctx, "getting url")

	err := validateAlias(req.GetAlias())
	if err != nil {
		s.Logger.ErrorContext(ctx, "invalid alias")
		return nil, err
	}

	url, err := s.Storage.GetURL(ctx, req.GetAlias())
	if err != nil {
		s.Logger.ErrorContext(ctx, "failed to get url")
		return nil, err
	}

	s.Logger.InfoContext(ctx, fmt.Sprintf("got url %s", url))

	return &api.GetURLResponse{
		Url: url,
	}, nil
}

func validateUrl(urlToValidate string) error {
	if urlToValidate == "" {
		return errors.New("url required")
	}

	_, err := url.ParseRequestURI(urlToValidate)

	if err != nil {
		return err
	}

	return nil
}

func validateAlias(alias string) error {
	if alias == "" {
		return errors.New("alias required")
	}

	if len(alias) != 10 {
		return errors.New("invalid alias length")
	}

	return nil
}

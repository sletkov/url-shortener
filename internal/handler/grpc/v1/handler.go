package v1

import (
	"context"
	"fmt"
	"log/slog"
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
	err := s.Storage.SaveURL(ctx, req.GetUrl(), alias)
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

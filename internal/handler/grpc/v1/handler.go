package v1

import (
	"context"
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
}

func New(storage Storage) *GRPCServer {
	return &GRPCServer{
		Storage: storage,
	}
}

func (s *GRPCServer) SaveURL(ctx context.Context, req *api.SaveURLRequest) (*api.SaveURLResponse, error) {
	alias := random.MakeRandomString(10)
	err := s.Storage.SaveURL(ctx, req.GetUrl(), alias)
	if err != nil {
		return nil, err
	}
	return &api.SaveURLResponse{
		Alias: alias,
	}, nil
}

func (s *GRPCServer) GetURL(ctx context.Context, req *api.GetURLRequest) (*api.GetURLResponse, error) {
	url, err := s.Storage.GetURL(ctx, req.GetAlias())
	if err != nil {
		return nil, err
	}
	return &api.GetURLResponse{
		Url: url,
	}, nil
}

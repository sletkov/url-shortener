package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"
	api "urlshortener/api/proto"
	random "urlshortener/internal/app/lib"
	inmemory "urlshortener/internal/app/storage/in-memory"
	"urlshortener/internal/app/storage/postgres"

	"google.golang.org/grpc"
)

const (
	inMemoryStorage = "in-memory"
	postgresStorage = "postgres"
)

type Storage interface {
	SaveURL(url string, alias string) error
	GetURL(alias string) (string, error)
}

type GRPCServer struct {
	api.UnimplementedURLShortenerServer
	Storage Storage
}

func NewServer(storage Storage) *GRPCServer {
	return &GRPCServer{
		Storage: storage,
	}
}

func Start(config *Config, storageType string) error {
	// TODO: init logger

	// Connect to db
	var srv *GRPCServer

	switch storageType {
	case postgresStorage:
		storage, err := postgres.New(config.StoragePath)
		if err != nil {
			log.Fatal(err)
		}
		srv = NewServer(storage)
		defer storage.Db.Close()
	case inMemoryStorage:
		storage, err := inmemory.New()
		if err != nil {
			log.Fatal(err)
		}
		srv = NewServer(storage)
	default:
		log.Fatal("invalid type of storage")
	}

	// Starting server
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost%s", config.BindAddr))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	// ...
	grpcServer := grpc.NewServer(opts...)
	api.RegisterURLShortenerServer(grpcServer, srv)
	println("gRPC server started...")
	return grpcServer.Serve(lis)
}

func (s *GRPCServer) SaveURL(ctx context.Context, req *api.SaveURLRequest) (*api.SaveURLResponse, error) {
	alias := random.MakeRandomString(10)
	err := s.Storage.SaveURL(req.Url, alias)
	if err != nil {
		return nil, err
	}
	return &api.SaveURLResponse{
		Alias: alias,
	}, nil
}

func (s *GRPCServer) GetURL(ctx context.Context, req *api.GetURLRequest) (*api.GetURLResponse, error) {
	url, err := s.Storage.GetURL(req.Alias)
	if err != nil {
		return nil, err
	}
	return &api.GetURLResponse{Url: url}, nil
}

package main

import (
	"database/sql"
	"flag"
	"log"
	"log/slog"
	"net"
	"os"
	api "urlshortener/api/proto"
	"urlshortener/internal/config"
	v1 "urlshortener/internal/handler/grpc/v1"
	inmemory "urlshortener/internal/storage/in-memory"
	"urlshortener/internal/storage/postgres"

	"github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

const (
	inMemoryStorage = "in-memory"
	postgresStorage = "postgres"
)

var (
	configPath  string
	storageType string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/config.yaml", "path to config file")
	flag.StringVar(&storageType, "storage-type", "postgres", "type of storage")
}

func main() {
	flag.Parse()

	config := config.NewConfig()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	var storage v1.Storage

	switch storageType {
	case postgresStorage:
		db, err := sql.Open("postgres", config.StoragePath)
		if err != nil {
			log.Fatal(err)
		}
		storage = postgres.New(db)
	case inMemoryStorage:
		cache := cache.New(cache.NoExpiration, cache.NoExpiration)
		storage = inmemory.New(cache)
	default:
		log.Fatal("unknown storage type")
	}

	grpcHandler := v1.New(logger, storage)

	grpcServer := grpc.NewServer()
	api.RegisterURLShortenerServer(grpcServer, grpcHandler)

	lis, err := net.Listen("tcp", net.JoinHostPort(config.Host, config.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	logger.Info("grpc server started")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

}

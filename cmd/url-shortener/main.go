package main

import (
	"flag"
	"log"
	"urlshortener/internal/app/grpcserver"
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
	config := grpcserver.NewConfig()
	if err := grpcserver.Start(config, storageType); err != nil {
		log.Fatal(err)
	}
}

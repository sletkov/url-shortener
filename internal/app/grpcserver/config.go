package grpcserver

type Config struct {
	BindAddr    string
	StoragePath string
}

func NewConfig() *Config {
	return &Config{
		BindAddr:    ":8082",
		StoragePath: "host=localhost user=sletkov password=postgres dbname=urlsDB sslmode=disable",
	}
}

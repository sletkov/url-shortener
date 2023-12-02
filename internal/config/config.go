package config

type Config struct {
	Host        string
	Port        string
	StoragePath string
}

func NewConfig() *Config {
	return &Config{
		Host:        "localhost",
		Port:        "8082",
		StoragePath: "host=localhost user=sletkov password=postgres dbname=urlsDB sslmode=disable",
	}
}

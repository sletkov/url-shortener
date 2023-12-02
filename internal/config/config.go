package config

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

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

func (c *Config) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Host, validation.Required),
		validation.Field(&c.Port, validation.Required),
		validation.Field(&c.StoragePath, validation.Required),
	)
}

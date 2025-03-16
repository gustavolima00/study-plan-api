package config

import (
	"log"

	env "github.com/caarlos0/env/v11"
)

// Config defines the application env vars
type Config struct {
	ConnStr string `env:"CONN_STR" envDefault:""`
	Port    string `env:"PORT" envDefault:"8080"`
}

// New will parse the necessary env vars to
// struct Config
func New() *Config {
	c := new(Config)

	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}

	return c
}

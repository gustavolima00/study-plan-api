package config

import (
	"log"

	env "github.com/caarlos0/env/v11"
)

// Config defines the application env vars
type Config struct {
	Port string `env:"PORT" envDefault:"8080"`

	// Keycloak
	KeycloakBaseURL      string `env:"KEYCLOAK_BASE_URL" envDefault:"http://localhost:8088"`
	KeycloakRealm        string `env:"KEYCLOAK_REALM" envDefault:"myrealm"`
	KeycloakClientID     string `env:"KEYCLOAK_CLIENT_ID" envDefault:"myclient"`
	KeycloakClientSecret string `env:"KEYCLOAK_CLIENT_SECRET" envDefault:"mysecret"`
	KeycloakTimoutMS     int    `env:"KEYCLOAK_TIMEOUT_MS" envDefault:"10000"`

	// Database
	PostgresConnectionString string `env:"POSTGRES_CONNECTION_STRING" envDefault:"postgres://myuser:mypassword@localhost:5432/mydb?sslmode=disable"`
	PostgresMaxOpenConns     int    `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"10"`
	PostgresMaxIdleConns     int    `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"10"`
	PostgresConnMaxLifetime  int    `env:"POSTGRES_CONN_MAX_LIFETIME" envDefault:"300"` // in seconds
	PostgresMaxIdleTime      int    `env:"POSTGRES_MAX_IDLE_TIME" envDefault:"300"`     // in seconds
	PostgresMaxLifetime      int    `env:"POSTGRES_MAX_LIFETIME" envDefault:"300"`      // in seconds
}

// NewConfig will parse the necessary env vars to
// struct Config
func NewConfig() *Config {
	c := new(Config)

	if err := env.Parse(c); err != nil {
		log.Fatal(err)
	}

	return c
}

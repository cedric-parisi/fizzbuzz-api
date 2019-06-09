package config

import (
	"log"

	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
)

// Config holds the config parameters.
type Config struct {
	AppPort string `env:"APP_PORT"`

	DbHost     string `env:"DB_HOST"`
	DbPort     string `env:"DB_PORT"`
	DbUser     string `env:"DB_USER"`
	DbName     string `env:"DB_NAME"`
	DbPassword string `env:"DB_PASSWORD"`
	DbTimeout  int    `env:"DB_TIMEOUT"`
}

// NewConfig read config from .env or env
// and return a config struct.
func NewConfig() (Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("could not load .env file, read from env")
	}

	var cfg Config
	if err := env.Parse(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

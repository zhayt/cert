package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	AppPort    string `env:"APP_PORT" envDefault:"8000"`
	DBHost     string `env:"DB_HOST" envDefault:"localhost"`
	DBPort     string `env:"DB_PORT" envDefault:"5432"`
	DBUser     string `env:"DB_USER" envDefault:"cert"`
	DBName     string `env:"DB_NAME" envDefault:"cert_db"`
	DBPassword string `env:"DB_PASSWORD"`
	TZ         string `env:"TZ" envDefault:"Asia/Almaty"`
}

func NewConfig() (*Config, error) {
	var cfg Config

	if err := env.Parse(&cfg); err != nil {
		return nil, fmt.Errorf("cannot parse env: %w", err)
	}

	return &cfg, nil
}

func PrepareENV() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
}

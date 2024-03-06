package config

import (
	"fmt"
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"log"
)

type Config struct {
	RedisDB       int    `env:"REDIS_DB" envDefault:"0"`
	AppPort       string `env:"APP_PORT" envDefault:"8000"`
	DBHost        string `env:"DB_HOST" envDefault:"localhost"`
	DBPort        string `env:"DB_PORT" envDefault:"5432"`
	DBUser        string `env:"DB_USER" envDefault:"cert"`
	DBName        string `env:"DB_NAME" envDefault:"cert_db"`
	DBPassword    string `env:"DB_PASSWORD"`
	TZ            string `env:"TZ" envDefault:"Asia/Almaty"`
	RedisHost     string `env:"REDIS_HOST" envDefault:"localhost"`
	RedisPort     string `env:"REDIS_PORT" envDefault:"6379"`
	RedisPassword string `env:"REDIS_PASSWORD"`
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

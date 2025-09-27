package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	HTTPPort      string
	RabbitMQURL   string
	RedisURL      string
	TelegramToken string
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		HTTPPort:      os.Getenv("HTTP_PORT"),
		RabbitMQURL:   os.Getenv("RABBITMQ_URL"),
		RedisURL:      os.Getenv("REDIS_URL"),
		TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
	}
}
